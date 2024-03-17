package main

import (
	"log"
	"log/slog"
	"os"

	gf "github.com/jessevdk/go-flags"
)

var Opts Options

var Args []string

type Options struct {
	Absolute bool `short:"A" long:"absolute" description:"Format paths to be absolute. Relative by default."`
	Recurse  bool `short:"r" long:"recurse" description:"Recursively list files in subdirectories"`

	ToDepth   int `short:"T" long:"todepth" description:"List files to a certain depth."`
	FromDepth int `short:"F" long:"fromdepth" description:"List files from a certain depth."`

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

	Sort bool `short:"n" long:"sort" description:"Sort the result. Files are ordered by filename by default."`

	Debug bool `short:"D" long:"debug" description:"Debug flag enables debug logging."`
	Quiet bool `short:"Q" long:"quiet" description:"Quiet flag disables printing results."`
}

func main() {
	Opts = Options{}
	args, err := gf.Parse(&Opts)
	if err != nil {
		if gf.WroteHelp(err) {
			os.Exit(0)
		}
		log.Fatalln("Error parsing flags:", err)
	}
	Args = args

	if Opts.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	implicitSlice()

	List()
}

func implicitSlice() {
	if Opts.Slice != "" {
		slog.Debug("slice is already set. ignoring implicit slice.")
		return
	}

	if len(Args) == 0 {
		slog.Debug("implicit slice found no args")
		return
	}

	L := len(Args) - 1

	if _, _, ok := parseSlice(Args[L]); ok {
		Opts.Slice = Args[L]
		Args = Args[:L]
	}
}
