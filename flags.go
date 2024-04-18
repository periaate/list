package list

import (
	"fmt"
	"log"
	"log/slog"
	"math"
	"os"

	gf "github.com/jessevdk/go-flags"
	"github.com/periaate/clf"
	"github.com/periaate/common"
	"github.com/periaate/slice"
)

type ListingOpts struct {
	Recurse  bool `short:"r" long:"recurse" description:"Recursively list files in subdirectories. Directory traversal is done iteratively and breadth first."`
	Archive  bool `short:"z" description:"Treat zip archives as directories."`
	NoHide   bool `short:"h" long:"hide" description:"Toggle of hiding of commonly unwanted files."`
	MaxLimit int  `short:"m" long:"max" description:"Maximum number of elements traversed in a single directory. Unlimited by default."`
}

type FilterOpts struct {
	Search []string `short:"s" long:"search" description:""`

	DirOnly  bool `long:"dirs" description:"Only include directories in the result."`
	FileOnly bool `long:"files" description:"Only include files in the result."`
}

type ProcessOpts struct {
	Ascending bool `short:"a" long:"ascending" description:"Results will be ordered in ascending order. Files are ordered into descending order by default."`
	Reverse   bool `short:"R" long:"reverse" description:"The exact same as the ascending flag."`

	Sort string `short:"S" long:"sort" description:"Sort the result by word." default:"none" choice:"none" choice:"name" choice:"n" choice:"mod" choice:"time" choice:"t" choice:"size" choice:"s" choice:"creation" choice:"c"`
}

type Printing struct {
	Absolute bool `short:"A" long:"absolute" description:"Format paths to be absolute. Relative by default."`
	Debug    bool `short:"D" long:"debug" description:"Debug flag enables debug logging."`
	Quiet    bool `short:"Q" long:"quiet" description:"Quiet flag disables printing results."`
}

type Options struct {
	ListingOpts `group:"Traversal options - Determines how the traversal is done."`
	FilterOpts  `group:"Filtering options - Applied while traversing, called on every entry found."`
	ProcessOpts `group:"Processing options - Applied after traversal, called on the final list of files."`
	Printing    `group:"Printing options - Determines how the results are printed."`

	ToDepth   int
	FromDepth int
	Queries   []Query

	ExecArgs []string
	Args     []string

	Filters   []func(*Element) bool
	Processes []func(els []*Element) []*Element

	Select []string
}

func Recurse(opts *Options) {
	opts.ToDepth = math.MaxInt64
}

func NoneFilter(*Element) bool              { return true }
func NoneProcess(els []*Element) []*Element { return els }

func Parse(args []string) *Options {
	if _, ind := common.First(args, func(s string) bool { return s == "-D" }); ind != -1 {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		if ind == len(args)-1 {
			args = args[:ind]
		} else {
			args = append(args[:ind], args[ind+1:]...)
		}
	}
	var execArgs []string

	if _, i := common.First(args, func(f string) bool { return f == "::" }); i != -1 {
		// drop the "::", everything after goes to execargs
		execArgs = args[i+1:]
		args = args[:i]
	}
	opts := &Options{
		ExecArgs:  execArgs,
		Filters:   []func(*Element) bool{NoneFilter},
		Processes: []func(els []*Element) []*Element{NoneProcess},
	}
	opts.MaxLimit = math.MaxInt64

	args = ApplyFlags(args, opts)
	fmt.Println(len(opts.Filters), len(opts.Processes))
	rest, err := gf.ParseArgs(opts, args)
	if err != nil {
		if gf.WroteHelp(err) {
			os.Exit(0)
		}
		log.Fatalln("Error parsing flags:", err)
	}
	fmt.Println(len(opts.Filters), len(opts.Processes))
	opts.Args = rest

	if opts.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	bef := len(opts.Args)
	Implicit(opts)
	if bef != len(opts.Args) {
		slog.Debug("Found implicit commands", "len", bef-len(opts.Args))
	}

	if len(opts.Args) == 0 {
		slog.Debug("No args found, defaulting to current directory")
		opts.Args = []string{"./"}
	}

	return opts
}

func Implicit(opts *Options) {
	if len(opts.Args) == 0 {
		slog.Debug("implicit slice found no Args")
		return
	}

	newArgs := make([]string, 0, len(opts.Args))
	for _, arg := range opts.Args {
		if len(arg) > 2 && arg[0] == '[' && arg[len(arg)-1] == ']' {
			opts.Select = append(opts.Select, arg)
			slog.Debug("implicitly found cmd", "type", "Slice", "arg", arg)
		} else {
			slog.Debug("implicit slice found no Args")
			newArgs = append(newArgs, arg)
		}
	}

	opts.Args = newArgs
}

var flags = []*clf.Flag{
	{
		Name:        "only",
		Greedy:      true, // Only capture at most 1 value
		Keys:        []string{"only"},
		Description: "Only files or only dirs",
		Exactly:     1,
	},
	{
		Name:        "recurse",
		Greedy:      true, // Only capture at most 1 value
		Keys:        []string{"r", "recurse"},
		Description: "Toggles on infinite depth, or can be given a depth slice",
		AtMost:      1,
	},
	{
		Name:        "is",
		Keys:        []string{"is", "?"},
		Description: "File type inference inclusion; e.g., image, video, etc.",
		AtLeast:     1,
	},
	{
		Name:        "not",
		Keys:        []string{"not", "!"},
		Description: "File type inference exclusion; e.g., image, video, etc.",
		AtLeast:     1,
	},
	{
		Name:        "where",
		Keys:        []string{"where", "s"},
		Description: "File type inference exclusion; e.g., image, video, etc.",
		AtLeast:     1,
	},
}

func ApplyFlags(args []string, opts *Options) []string {
	if len(args) == 0 {
		slog.Debug("clf flags: no args found")
		return args
	}
	op, err := clf.Parse(args, flags)
	if err != nil {
		slog.Debug("Error parsing clf flags", "error", err)
		return args
	}

	opts.Args = op.Rest
	for _, flag := range op.Yield() {
		switch flag.Name {
		case "only":
			switch flag.Values[0] {
			case "file", "f":
				opts.FileOnly = true
			case "dir", "d":
				opts.DirOnly = true
			}
		case "recurse":
			slog.Debug("found clf flag", "type", flag.Name)
			if len(flag.Values) > 0 {
				f, t, err := slice.ParsePattern(flag.Values[0], math.MaxInt)
				if err != nil {
					continue
				}
				slog.Debug("recurse range", "from", f, "to", t)
				opts.FromDepth = f
				opts.ToDepth = t
				continue
			}
			slog.Debug("recurse call")
			Recurse(opts)
		case "is":
			opts.Filters = append(opts.Filters, ParseKind(flag.Values, true))
		case "not":
			opts.Filters = append(opts.Filters, ParseKind(flag.Values, false))
		case "where":
			opts.Filters = append(opts.Filters, ParseSearch(flag.Values))
		}
	}

	if len(args) != len(op.Rest) {
		slog.Debug("found clf flags", "diff", len(args)-len(op.Rest))
	}
	return op.Rest
}
