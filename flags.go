package list

import (
	"log/slog"

	"github.com/periaate/clf"
	"github.com/periaate/ls/lfs"
)

var (
	version = "undefined"
	run     = &Runner{}
	res     []*lfs.Element
)

var program = clf.Register(
	clf.Info(
		clf.Name("list"),
		clf.Description("List is a filesystem listing CLI core utility."),
		clf.Author("Daniel Saury"),
		clf.Source("https://github.com/periaate/list"),
		clf.Copyright("Copyright (c) 2024 Daniel Saury"),
		clf.Version(version),
	),
	clf.Flags(flags),
)

var flags = []*clf.Flag{
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
	{Keys: []string{"reverse", "a", "asc"}, AtLeast: 1, Description: "Search by name.",
		Handler: run.Reverse},

	{Keys: []string{"absolute", "abs", "A"}, AtLeast: 1, Description: "Absolute paths.",
		Handler: func(_ []string) { Abs = true }},

	clf.Group([]*clf.Flag{
		{Keys: []string{"--debug", "-D"}, Exactly: -1,
			Description: "Set debug level.",
			Handler:     func(_ []string) { slog.SetLogLoggerLevel(slog.LevelDebug) }},
		{Keys: []string{"--quiet", "-Q"}, Exactly: -1,
			Description: "Don't print results.",
			Handler:     func(_ []string) { Quiet = true }},
	}),
}
