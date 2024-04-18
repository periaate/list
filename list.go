package list

import (
	"bufio"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/periaate/common"
)

func Run(opts *Options) []*Element {
	res := []*Element{}
	rfn := GetRfn(CollectFilters(opts), res)
	process := CollectProcess(opts)
	yfn := GetYieldFs(opts)
	Traverse(opts, yfn, rfn)

	res = process(res)
	slog.Debug("final result", "len", len(res))
	return res
}

func Initialize(opts *Options) ([]*Element, Filter, Process) {
	res := []*Element{}
	filters := CollectFilters(opts)
	processes := CollectProcess(opts)

	return res, filters, processes
}

func Do(args ...string) []*Element {
	opts := Parse(args)
	return Run(opts)
}

func Exec(res []*Element, opts *Options) {
	opts.ExecArgs = append(opts.ExecArgs, common.Collect(res, func(e *Element) string { return e.Name })...)

	cmd := exec.Command(opts.ExecArgs[0], opts.ExecArgs[1:]...)
	cmd.Stdout = os.Stdout
	if opts.Debug {
		cmd.Stderr = os.Stderr
	}
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		slog.Debug("error running command", "err", err)
	}
}

const bufLength = 500

func PrintWithBuf(els []*Element, opts *Options) {
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
