package main

import (
	"log"
	"log/slog"
	"math"
	"os"

	gf "github.com/jessevdk/go-flags"
)

var Opts Options

var args []string

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

	Sort string `long:"sort" description:"Sort the result by word. Choices: [name|date|mod|size|none]. Date and mod are the same." default:"none"`

	// Mod      bool `long:"mod" description:"Results will be ordered by their modified time."`
	// Size     bool `long:"size" description:"Results will be ordered by their size time."`
	// None     bool `long:"none" description:"Results will be ordered by their modified time."`

	Select string `long:"select" description:"Select a single element or a range of elements. Usage: [{index}] [{from}:{to}] Supports negative indexing. Can be used without a flag as the last argument."`
}

type Printing struct {
	Absolute  bool `short:"A" long:"absolute" description:"Format paths to be absolute. Relative by default."`
	Debug     bool `short:"D" long:"debug" description:"Debug flag enables debug logging."`
	Quiet     bool `short:"Q" long:"quiet" description:"Quiet flag disables printing results."`
	Clipboard bool `short:"c" long:"clipboard" description:"Copy the result to the clipboard."`
}

type Options struct {
	ListingOpts `group:"Traversal options - Determines how the traversal is done."`
	FilterOpts  `group:"Filtering options - Applied while traversing, called on every entry found."`
	ProcessOpts `group:"Processing options - Applied after traversal, called on the final list of files."`
	Printing    `group:"Printing options - Determines how the results are printed."`
}

func main() {
	Opts = Options{}
	rest, err := gf.Parse(&Opts)
	if err != nil {
		if gf.WroteHelp(err) {
			os.Exit(0)
		}
		log.Fatalln("Error parsing flags:", err)
	}
	args = rest

	if Opts.Recurse {
		Opts.ToDepth = math.MaxInt64
	}

	if Opts.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	implicitSlice()

	if len(args) == 0 {
		args = append(args, "./")
	}

	res := &result{[]*finfo{}}
	filters := collectFilters()
	processes := collectProcess()
	wfn := buildWalkDirFn(filters, res)
	Traverse(wfn)

	ProcessList(res, processes)
	printWithBuf(res.files)
}

func implicitSlice() {
	if Opts.Select != "" {
		slog.Debug("slice is already set. ignoring implicit slice.")
		return
	}

	if len(args) == 0 {
		slog.Debug("implicit slice found no args")
		return
	}

	L := len(args) - 1

	if _, _, ok := parseSlice(args[L]); ok {
		Opts.Select = args[L]
		args = args[:L]
	}
}
