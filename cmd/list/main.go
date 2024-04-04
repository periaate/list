package main

import (
	"bufio"
	"list/cfg"
	"log"
	"log/slog"
	"math"
	"os"

	"list"

	gf "github.com/jessevdk/go-flags"
)

func main() {
	cfg.Opts = &cfg.Options{}
	Opts := cfg.Opts
	rest, err := gf.Parse(Opts)
	if err != nil {
		if gf.WroteHelp(err) {
			os.Exit(0)
		}
		log.Fatalln("Error parsing flags:", err)
	}
	cfg.Args = rest

	if cfg.Opts.ToDepth == 0 && cfg.Opts.Recurse {
		Opts.ToDepth = math.MaxInt64
	}

	if cfg.Opts.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	implicitSlice()
	pipedValues := readPipe()
	if len(pipedValues) != 0 {
		cfg.Args = append(cfg.Args, pipedValues...)
	}

	if len(cfg.Args) == 0 {
		cfg.Args = append(cfg.Args, "./")
	}

	res := &list.Result{Files: []*list.Finfo{}}
	filters := list.CollectFilters()
	processes := list.CollectProcess()
	wfn := list.BuildWalkDirFn(filters, res)
	list.Traverse(wfn)

	list.ProcessList(res, processes)
	list.PrintWithBuf(res.Files)
}

func implicitSlice() {
	if cfg.Opts.Select != "" {
		slog.Debug("slice is already set. ignoring implicit slice.")
		return
	}

	if len(cfg.Args) == 0 {
		slog.Debug("implicit slice found no cfg.Args")
		return
	}

	L := len(cfg.Args) - 1

	if _, _, ok := list.ParseSlice(cfg.Args[L]); ok {
		cfg.Opts.Select = cfg.Args[L]
		cfg.Args = cfg.Args[:L]
	}
}

func readPipe() (res []string) {
	fileInfo, _ := os.Stdin.Stat()
	if (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			res = append(res, scanner.Text())
		}
	}
	return
}
