package list

import (
	"log"
	"log/slog"
	"math"
	"os"

	gf "github.com/jessevdk/go-flags"
	"github.com/periaate/clf"
	"github.com/periaate/common"
)

func Run(opts *Options) *Result {
	res := &Result{Files: []*Element{}}
	filters := InitFilters(CollectFilters(opts), res)
	processes := CollectProcess(opts)
	traverser := GetTraverser(opts)

	traverser(opts, filters)
	ProcessList(res, processes)
	return res
}

func Initialize(opts *Options) (*Result, []Filter, []Process) {
	res := &Result{Files: []*Element{}}
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
	NoHide    bool     `short:"h" long:"hide" description:"Toggle of hiding of commonly unwanted files."`
	MaxLimit  int      `short:"m" long:"max" description:"Maximum number of elements traversed in a single directory. Unlimited by default."`
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
	Reverse   bool     `short:"R" long:"reverse" description:"The exact same as the ascending flag."`

	Sort string `short:"S" long:"sort" description:"Sort the result by word." default:"none" choice:"none" choice:"name" choice:"n" choice:"mod" choice:"time" choice:"t" choice:"size" choice:"s" choice:"creation" choice:"c"`

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

	traversalFpp *FPPair
	filesFpp     *FPPair
	dirsFpp      *FPPair
}

func Recurse(opts *Options) {
	opts.ToDepth = math.MaxInt64
}

func Parse(args []string) *Options {
	var execArgs []string

	if _, i := common.First(args, func(f string) bool { return f == "::" }); i != -1 {
		// drop the "::", everything after goes to execargs
		execArgs = args[i+1:]
		args = args[:i]
	}

	opts := &Options{
		ExecArgs: execArgs,
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

	if opts.ToDepth == 0 && opts.Recurse {
		Recurse(opts)
	}

	opts.ToDepth = common.Clamp(opts.ToDepth, opts.FromDepth+1, math.MaxInt64)

	if opts.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	bef := len(opts.Args)
	Implicit(opts)
	if bef != len(opts.Args) {
		slog.Debug("Found implicit commands", "len", bef-len(opts.Args))
	}

	d(opts.Args, opts)

	return opts
}

func Implicit(opts *Options) {
	if len(opts.Args) == 0 {
		slog.Debug("implicit slice found no Args")
		return
	}

	newArgs := make([]string, 0, len(opts.Args))
	for _, arg := range opts.Args {
		// 	fpp, ok := TargetedSearch(opts, arg)

		// 	if !ok {
		// 		goto Old
		// 	}

		// 	if fpp == nil {
		// 		continue
		// 	}
		// 	switch fpp.Tar {
		// 	case traversal:
		// 		if opts.traversalFpp != nil && opts.traversalFpp.Filter != nil {
		// 			slog.Debug("found traversal filter", "filter", fpp.Filter)
		// 			opts.traversalFpp.Filter = common.All(true, opts.traversalFpp.Filter, fpp.Filter)
		// 			continue
		// 		}
		// 		slog.Debug("didn't find traversal filter", "filter", fpp.Filter)
		// 		opts.traversalFpp = fpp
		// 		continue
		// 	case files:
		// 		if opts.filesFpp != nil && opts.filesFpp.Filter != nil && opts.filesFpp.Process != nil {
		// 			opts.filesFpp.Filter = common.All(true, opts.filesFpp.Filter, fpp.Filter)
		// 			opts.filesFpp.Process = common.Pipe(opts.filesFpp.Process, fpp.Process)
		// 			continue
		// 		}
		// 		opts.filesFpp = fpp
		// 		continue
		// 	case dirs:
		// 		if opts.dirsFpp != nil && opts.dirsFpp.Filter != nil && opts.dirsFpp.Process != nil {
		// 			opts.dirsFpp.Filter = common.All(true, opts.dirsFpp.Filter, fpp.Filter)
		// 			opts.dirsFpp.Process = common.Pipe(opts.dirsFpp.Process, fpp.Process)
		// 			continue
		// 		}
		// 		opts.dirsFpp = fpp
		// 		continue
		// 	}
		// Old:
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
	'h': func(opts *Options) { opts.NoHide = true },
	'n': func(opts *Options) { opts.Sort = "name" },
	'c': func(opts *Options) { opts.Sort = "creation" },
	't': func(opts *Options) { opts.Sort = "time" },
	'f': func(opts *Options) { opts.FileOnly = true },
	'd': func(opts *Options) { opts.DirOnly = true },
	'r': func(opts *Options) { Recurse(opts) },
	'z': func(opts *Options) { opts.Archive = true },
	'C': func(opts *Options) { opts.Count = true },
	'R': func(opts *Options) { opts.Ascending = true },
	'M': func(opts *Options) { opts.MaxLimit = 1000 },
}

var file = &clf.Flag{
	Keys: []string{"f"},
	Name: "file",
}

var recurse = &clf.Flag{
	Toggle: true,
	Keys:   []string{"r"},
	Name:   "recurse",
}

var incl = &clf.Flag{
	Keys: []string{"i"},
	Name: "include",
}

func d(args []string, opts *Options) {
	op, err := clf.Parse(args, []*clf.Flag{file, recurse, incl})
	if err != nil {
		return
	}

	opts.Args = op.Rest

	if op.Get("recurse").Present != 0 {
		Recurse(opts)
	}
	if len(op.Get("include").Values) != 0 {
		opts.Include = append(opts.Include, op.Get("include").Values...)
	}
	if len(op.Get("file").Values) != 0 {
		opts.Query = append(opts.Query, op.Get("file").Values...)
	}
}
