package list

import (
	"log"
	"log/slog"
	"math"
	"os"

	gf "github.com/jessevdk/go-flags"
)

type ListingOpts struct {
	Recurse   bool `short:"r" long:"recurse" description:"Recursively list files in subdirectories. Directory traversal is done iteratively and breadth first."`
	Archive   bool `short:"z" description:"Treat zip archives as directories."`
	ToDepth   int  `short:"T" long:"todepth" description:"List files to a certain depth." default:"0"`
	FromDepth int  `short:"F" long:"fromdepth" description:"List files from a certain depth." default:"-1"`
}

type FilterOpts struct {
	Search    []string `short:"s" long:"search" description:"Only include items which have search terms as substrings. Can be used multiple times. Multiple values are inclusive by default. (OR)"`
	SearchAnd bool     `long:"AND" description:"Including this flag makes search with multiple values conjuctive, i.e., all search terms must be matched. (AND)"`
	Include   []string `short:"i" long:"include" description:"File type inclusion: image, video, audio. Can be used multiple times."`
	Exclude   []string `short:"e" long:"exclude" description:"File type exclusion: image, video, audio. Can be used multiple times."`
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

	Select string `long:"select" description:"Select a single element or a range of elements. Usage: [{index}] [{from}:{to}] Supports negative indexing. Can be used without a flag as the last argument."`

	Shuffle bool  `long:"shuffle" description:"Randomly shuffle the result."`
	Seed    int64 `long:"seed" description:"Seed for the random shuffle." default:"-1"`
}

type Printing struct {
	Absolute  bool `short:"A" long:"absolute" description:"Format paths to be absolute. Relative by default."`
	Debug     bool `short:"D" long:"debug" description:"Debug flag enables debug logging."`
	Quiet     bool `short:"Q" long:"quiet" description:"Quiet flag disables printing results."`
	Clipboard bool `short:"c" long:"clipboard" description:"Copy the result to the clipboard."`
	Tree      bool `long:"tree" description:"Prints as tree."`
}

type Options struct {
	ListingOpts `group:"Traversal options - Determines how the traversal is done."`
	FilterOpts  `group:"Filtering options - Applied while traversing, called on every entry found."`
	ProcessOpts `group:"Processing options - Applied after traversal, called on the final list of files."`
	Printing    `group:"Printing options - Determines how the results are printed."`

	Args []string
}

func Parse(args []string) *Options {
	opts := &Options{}
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

	opts.ToDepth = Clamp(opts.ToDepth, opts.FromDepth+1, math.MaxInt64)

	if opts.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	slog.Debug("args before slice", "len", len(opts.Args))
	ImplicitSlice(opts)
	slog.Debug("left after slice", "len", len(opts.Args))

	if len(opts.Args) == 0 {
		opts.Args = append(opts.Args, "./")
	}

	return opts
}

func ImplicitSlice(opts *Options) {
	if opts.Select != "" {
		slog.Debug("slice is already set. ignoring implicit slice.")
		return
	}

	if len(opts.Args) == 0 {
		slog.Debug("implicit slice found no Args")
		return
	}

	for i, arg := range opts.Args {
		if len(arg) < 2 {
			continue
		}
		if arg[0] == '[' && arg[len(arg)-1] == ']' {
			opts.Select = arg
			opts.Args = append(opts.Args[:i], opts.Args[i+1:]...)
			slog.Debug("implicit slice found", "slice", opts.Select, "Args left", len(opts.Args))
			break
		}
	}
}
