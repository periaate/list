package list

import (
	"log/slog"
	"math"

	"github.com/periaate/clf"
	"github.com/periaate/slice"
)

var flags = []*clf.Flag{
	{Keys: []string{"recurse", "r"}, AtMost: 1, Description: "Recurse dirs. Optional range slice."},
	{Keys: []string{"is", "?"}, AtLeast: 1, Description: "File type inclusion."},
	{Keys: []string{"not", "!"}, AtLeast: 1, Description: "File type exclusion."},
}

func ApplyFlags(args []string, opts *Options) []string {
	if len(args) == 0 {
		slog.Debug("clf flags: no args found")
		return args
	}
	op, err := clf.Parse(args, flags)
	if err != nil {
		slog.Debug("Error parsing clf flags", "error", err)
		return args
	}

	opts.Args = op.Rest
	for _, flag := range op.Yield() {
		switch flag.Name {
		case "only":
			switch flag.Values[0] {
			case "file", "f", "files":
				opts.OnlyFiles = true
			case "dir", "d", "dirs":
				opts.DirOnly = true
			}
		case "recurse":
			slog.Debug("found clf flag", "type", flag.Name)
			if len(flag.Values) > 0 {
				f, t, err := slice.ParsePattern(flag.Values[0], math.MaxInt)
				if err != nil {
					continue
				}
				slog.Debug("recurse range", "from", f, "to", t)
				opts.FromDepth = f
				opts.ToDepth = t
				continue
			}
			slog.Debug("recurse call")
			opts.ToDepth = math.MaxInt
		case "is":
			opts.Filters = append(opts.Filters, ParseKind(flag.Values, true))
		case "not":
			opts.Filters = append(opts.Filters, ParseKind(flag.Values, false))
		case "where":
			opts.Filters = append(opts.Filters, ParseSearch(flag.Values))
		}
	}

	if len(args) != len(op.Rest) {
		slog.Debug("found clf flags", "diff", len(args)-len(op.Rest))
	}
	return op.Rest
}
