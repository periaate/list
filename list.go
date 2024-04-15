package list

import (
	"log"
	"log/slog"
	"math"
	"os"

	gf "github.com/jessevdk/go-flags"
	"github.com/periaate/common"
	"github.com/periaate/slice"
)

func Run(opts *Options) *Result {
	res := &Result{Files: []*Finfo{}}
	filters := InitFilters(CollectFilters(opts), res)
	processes := CollectProcess(opts)
	traverser := GetTraverser(opts)

	traverser(opts, filters)
	ProcessList(res, processes)
	return res
}

func Initialize(opts *Options) (*Result, []Filter, []Process) {
	res := &Result{Files: []*Finfo{}}
	filters := CollectFilters(opts)
	processes := CollectProcess(opts)

	return res, filters, processes
}

type ModeOpts struct {
	ArgMode  bool   `short:"l" long:"arg" description:"Skips listing, uses the input arguments as elements."`
	FileMode string `long:"file" description:"Reads the files given as arguments and uses either words or lines as elements." choice:"words" choice:"w" choice:"lines" choice:"l"`
}

type ListingOpts struct {
	Recurse   bool     `short:"r" long:"recurse" description:"Recursively list files in subdirectories. Directory traversal is done iteratively and breadth first."`
	Archive   bool     `short:"z" description:"Treat zip archives as directories."`
	ToDepth   int      `short:"T" long:"todepth" description:"List files to a certain depth." default:"0"`
	FromDepth int      `short:"F" long:"fromdepth" description:"List files from a certain depth." default:"-1"`
	DirSearch []string `short:"d" long:"dirsearch" description:"Only include directories which have search terms as substrings. Can be used multiple times. Multiple values are inclusive by default. (OR) Does not work within archives."`
}

type FilterOpts struct {
	Search    []string `short:"s" long:"search" description:"Only include items which have search terms as substrings. Can be used multiple times. Multiple values are inclusive by default. (OR)"`
	SearchAnd bool     `long:"AND" description:"Including this flag makes search with multiple values conjuctive, i.e., all search terms must be matched. (AND)"`
	Include   []string `short:"i" long:"include" description:"File type inclusion. Can be used multiple times."`
	Exclude   []string `short:"e" long:"exclude" description:"File type exclusion. Can be used multiple times."`
	Ignore    []string `short:"I" long:"ignore" description:"Ignores all paths which include any given strings."`

	DirOnly  bool `long:"dirs" description:"Only include directories in the result."`
	FileOnly bool `long:"files" description:"Only include files in the result."`
}

type ProcessOpts struct {
	Query     []string `short:"q" long:"query" description:"Fuzzy search query. Results will be ordered by their score."`
	Ascending bool     `short:"a" long:"ascending" description:"Results will be ordered in ascending order. Files are ordered into descending order by default."`

	Sort string `short:"S" long:"sort" description:"Sort the result by word." default:"none" choice:"none" choice:"name" choice:"n" choice:"mod" choice:"time" choice:"t" choice:"size" choice:"s" choice:"creation" choice:"c"`

	// Mod      bool `long:"mod" description:"Results will be ordered by their modified time."`
	// Size     bool `long:"size" description:"Results will be ordered by their size time."`
	// None     bool `long:"none" description:"Results will be ordered by their modified time."`

	Select []string `long:"select" description:"Select a single element or a range of elements. Usage: [{index}] [{from}:{to}] Supports negative indexing. Can be used without a flag as the last argument."`

	Shuffle bool  `long:"shuffle" description:"Randomly shuffle the result."`
	Seed    int64 `long:"seed" description:"Seed for the random shuffle." default:"-1"`
}

type Printing struct {
	Absolute bool `short:"A" long:"absolute" description:"Format paths to be absolute. Relative by default."`
	Debug    bool `short:"D" long:"debug" description:"Debug flag enables debug logging."`
	Quiet    bool `short:"Q" long:"quiet" description:"Quiet flag disables printing results."`
	Count    bool `short:"C" long:"count" description:"Print the number of results."`
	Tree     bool `long:"tree" description:"Prints as tree."`
}

type Options struct {
	ModeOpts    `group:"Mode options - Determines which mode list executes in, fs, string, file, etc."`
	ListingOpts `group:"Traversal options - Determines how the traversal is done."`
	FilterOpts  `group:"Filtering options - Applied while traversing, called on every entry found."`
	ProcessOpts `group:"Processing options - Applied after traversal, called on the final list of files."`
	Printing    `group:"Printing options - Determines how the results are printed."`

	ExecArgs []string
	Args     []string
}

func Parse(args []string) *Options {
	var execArgs []string

	if i, ok := common.Any(args, func(f string) bool { return f == "::" }); ok {
		// drop the "::", everything after goes to execargs
		execArgs = args[i+1:]
		args = args[:i]
	}

	opts := &Options{
		ExecArgs: execArgs,
	}
	rest, err := gf.ParseArgs(opts, args)
	if err != nil {
		if gf.WroteHelp(err) {
			os.Exit(0)
		}
		log.Fatalln("Error parsing flags:", err)
	}

	opts.Args = rest

	if opts.ToDepth == 0 && opts.Recurse {
		opts.ToDepth = math.MaxInt64
	}

	opts.ToDepth = slice.Clamp(opts.ToDepth, opts.FromDepth+1, math.MaxInt64)

	if opts.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	slog.Debug("args before generic.Slice", "len", len(opts.Args))
	Implicit(opts)
	slog.Debug("left after generic.Slice", "len", len(opts.Args))

	return opts
}

func Implicit(opts *Options) {
	if len(opts.Args) == 0 {
		slog.Debug("implicit slice found no Args")
		return
	}

	newArgs := make([]string, 0, len(opts.Args))
	for _, arg := range opts.Args {
		switch {
		case len(arg) > 2 && arg[0] == '[' && arg[len(arg)-1] == ']':
			opts.Select = append(opts.Select, arg)

			slog.Debug("implicit generic.Slice found", "generic.Slice", opts.Select, "Args left", len(opts.Args))
		case len(arg) > 1 && arg[0] == '?':
			QuickCommand(arg[1:], opts)
		default:
			newArgs = append(newArgs, arg)
		}
	}

	opts.Args = newArgs
}

func QuickCommand(arg string, opts *Options) {
	for _, r := range arg {
		if fn, ok := pairs[r]; ok {
			fn(opts)
		}
	}
}

func Do(args ...string) *Result {
	opts := Parse(args)
	return Run(opts)
}

var pairs = map[rune]func(*Options){
	'm': func(opts *Options) { opts.Include = append(opts.Include, Audio, Video, Image) },
	'a': func(opts *Options) { opts.Include = append(opts.Include, Audio) },
	'v': func(opts *Options) { opts.Include = append(opts.Include, Video) },
	'i': func(opts *Options) { opts.Include = append(opts.Include, Image) },
	'n': func(opts *Options) { opts.Sort = "name" },
	'c': func(opts *Options) { opts.Sort = "creation" },
	't': func(opts *Options) { opts.Sort = "time" },
	'f': func(opts *Options) { opts.FileOnly = true },
	'd': func(opts *Options) { opts.DirOnly = true },
	'r': func(opts *Options) { opts.Recurse = true },
	'z': func(opts *Options) { opts.Archive = true },
	'C': func(opts *Options) { opts.Count = true },
}
