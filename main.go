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

type Options struct {

	// File listing
	Recurse   bool `short:"r" long:"recurse" description:"Recursively list files in subdirectories"`
	Archive   bool `short:"z" description:"Treat zip archives as directories."`
	ToDepth   int  `short:"T" long:"todepth" description:"List files to a certain depth." default:"-1"`
	FromDepth int  `short:"F" long:"fromdepth" description:"List files from a certain depth." default:"-1"`

	// filters
	Include []string `short:"i" long:"include" description:"File type inclusion: image, video, audio"`
	Exclude []string `short:"e" long:"exclude" description:"File type exclusion: image, video, audio."`
	Ignore  []string `short:"I" long:"ignore" description:"Ignores all paths which include any given strings."`
	Search  []string `short:"s" long:"search" description:"Only include paths which include any given strings."`

	DirOnly  bool `long:"dirs" description:"Only include directories in the result."`
	FileOnly bool `long:"files" description:"Only include files in the result."`

	// process
	Query     []string `short:"q" long:"query" description:"Fuzzy search query. Results will be ordered by their score."`
	Ascending bool     `short:"a" long:"ascending" description:"Results will be ordered in ascending order. Files are ordered into descending order by default."`
	Date      bool     `short:"d" long:"date" description:"Results will be ordered by their modified time. Files are ordered by filename by default"`
	Slice     string   `short:"S" long:"slice" description:"Slice [{from}:{to}]. Supports negative indexing. Can be used without a flag as the last argument."`
	Sort      bool     `short:"n" long:"sort" description:"Sort the result. Files are ordered by filename by default."`

	// Printing
	Absolute bool `short:"A" long:"absolute" description:"Format paths to be absolute. Relative by default."`
	Debug    bool `short:"D" long:"debug" description:"Debug flag enables debug logging."`
	Quiet    bool `short:"Q" long:"quiet" description:"Quiet flag disables printing results."`
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

	if Opts.ToDepth == -1 {
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
	if Opts.Slice != "" {
		slog.Debug("slice is already set. ignoring implicit slice.")
		return
	}

	if len(args) == 0 {
		slog.Debug("implicit slice found no args")
		return
	}

	L := len(args) - 1

	if _, _, ok := parseSlice(args[L]); ok {
		Opts.Slice = args[L]
		args = args[:L]
	}
}
