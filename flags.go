package list

import (
	"log/slog"
	"os"

	"github.com/periaate/clf"
	"github.com/periaate/common"
	"github.com/periaate/ls/lfs"
)

var (
	version = "undefined"
	run     = &Runner{}
	res     []*lfs.Element
	Log     *slog.Logger
)

var Program = clf.Register(
	clf.Info(clf.Meta{
		Name:        "list",
		Description: "List is a filesystem listing CLI core utility.",
		Author:      "Daniel Saury",
		Source:      "https://github.com/periaate/list",
		Copyright:   "Copyright (c) 2024 Daniel Saury",
		Version:     version,
	}),
	clf.Flags(flags),
)

var flags = common.Join(
	clf.Group("def", []*clf.Flag{
		// {Keys: []string{"from", "f"}, Handler: run.Dir, AtLeast: 1},
		{Keys: []string{"recurse", "r"}, AtMost: 1, Description: "Recurse dirs. Optional range slice.",
			Handler: run.Recurse,
		},
		{Keys: []string{"is", "="}, AtLeast: 1, Description: "File type inclusion.",
			Handler: run.Is},
		{Keys: []string{"not", "!"}, AtLeast: 1, Description: "File type exclusion.",
			Handler: run.Not},
		{Keys: []string{"slice", "]"}, Exactly: 1, Description: "Slice the list.",
			Handler: run.Slice},
		{Keys: []string{"sort", "$"}, Exactly: 1, Description: "Sort the list.",
			Handler: run.Sort},
		{Keys: []string{"search", "where", "?"}, AtLeast: 1, Description: "Search by name.",
			Handler: run.Search},
		{Keys: []string{"reverse"}, AtLeast: 1, Description: "Search by name.",
			Handler: run.Reverse},
		{Keys: []string{"all", "h"}, AtLeast: 1, Description: "Include everything.",
			Handler: run.All},

		{Keys: []string{"absolute", "abs", "A"}, AtLeast: 1, Description: "Absolute paths.",
			Handler: func(_ []string) { Abs = true }},
	}...),
	clf.Group("other", []*clf.Flag{
		{Name: "debug", Keys: []string{"--debug", "-D"}, Exactly: -1,
			Description: "Set debug level.",
			Handler: func(_ []string) {
				Log = common.NewClog(os.Stdout, slog.LevelDebug, common.MaxLen(30))
				clf.SetGlobalLogger(Log)
				slog.SetDefault(Log)
			}},
		{Name: "quiet", Keys: []string{"--quiet", "-Q"}, Exactly: -1,
			Description: "Don't print results.",
			Handler:     func(_ []string) { Quiet = true }},
		clf.DefaultHelp(nil),
	}...),
)
