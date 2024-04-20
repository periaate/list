package list

import (
	"log"
	"log/slog"
	"math"
	"os"

	gf "github.com/jessevdk/go-flags"
)

type ListingOpts struct {
	Recurse  bool `short:"r" long:"recurse" description:"Recursively list files in subdirectories. Directory traversal is done iteratively and breadth first."`
	Archive  bool `short:"z" description:"Treat zip archives as directories."`
	NoHide   bool `short:"h" long:"hide" description:"Toggle of hiding of commonly unwanted files."`
	MaxLimit int  `short:"m" long:"max" description:"Maximum number of elements traversed in a single directory. Unlimited by default."`
}

type FilterOpts struct {
	Search []string `short:"s" long:"search" description:""`

	DirOnly   bool `long:"dirs" description:"Only include directories in the result."`
	OnlyFiles bool `long:"files" description:"Only include files in the result."`
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
}

// Parse is used to parse command line or string arguments to the Options struct.
func Parse(args []string) *Options {
	opts := &Options{
		Filters:   []func(*Element) bool{NoneFilter},
		Processes: []func(els []*Element) []*Element{NoneProcess},
		ToDepth:   0,
		FromDepth: -1,
	}
	opts.MaxLimit = math.MaxInt64

	rest, err := gf.ParseArgs(opts, args)
	if err != nil {
		if gf.WroteHelp(err) {
			os.Exit(0)
		}
		log.Fatalln("Error parsing flags:", err)
	}
	opts.Args = rest

	Implicit(opts)
	opts.Args = ApplyFlags(opts.Args, opts)
	if opts.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	bef := len(opts.Args)
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
