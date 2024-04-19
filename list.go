package list

import (
	"bufio"
	"log/slog"
	"os"
	"path/filepath"
)

type Result struct {
	Files []*Element
}

func Run(opts *Options) []*Element {
	res := &Result{}
	rfn := GetRfn(CollectFilters(opts), res)
	process := CollectProcess(opts)
	yfn := GetYieldFs(opts)
	Traverse(opts, yfn, rfn)

	res.Files = process(res.Files)
	slog.Debug("final result", "len", len(res.Files))
	return res.Files
}

func Initialize(opts *Options) ([]*Element, Filter, Process) {
	res := []*Element{}
	filters := CollectFilters(opts)
	processes := CollectProcess(opts)

	return res, filters, processes
}

func PrintWithBuf(els []*Element, opts *Options) {
	const bufLength = 500
	if len(els) == 0 {
		return
	}
	if opts.Quiet {
		slog.Debug("quiet flag is set, returning from print function")
		return
	}

	w := bufio.NewWriterSize(os.Stdout, 4096*bufLength)

	for i, file := range els {
		fp := filepath.ToSlash(file.Path)
		if opts.Absolute {
			fp, _ = filepath.Abs(file.Path)
			fp = filepath.ToSlash(fp)
		}
		res := fp + "\n"

		w.WriteString(res)
		if i%bufLength == 0 {
			w.Flush()
		}
	}

	w.Flush()
}
