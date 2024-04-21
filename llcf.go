package list

import (
	"log/slog"

	"github.com/periaate/common"
	"github.com/periaate/ls"
	"github.com/periaate/ls/files"
	"github.com/periaate/ls/lfs"
)

type Runner struct {
	processes []ls.Process
	options   []ls.Option
	paths     []string
	direct    bool
}

func (r *Runner) Recurse(args []string) {
	if len(args) == 0 {
		r.options = append(r.options, ls.Recurse)
		return
	}
	r.options = append(r.options, ls.DepthPattern(args[0]))
}

func (r *Runner) Reverse(args []string) {
	r.processes = append(r.processes, func(inp []*lfs.Element) (out []*lfs.Element) {
		return common.ReverseC(inp)
	})
}

func (r *Runner) Is(args []string) {
	slog.Debug("is", "args", args)
	var masks []uint32
	for _, arg := range args {
		masks = append(masks, files.StrToMask(arg))
	}

	r.options = append(r.options, ls.Masks(ls.Include, masks...))
}

func (r *Runner) Not(args []string) {
	slog.Debug("not", "args", args)
	var masks []uint32
	for _, arg := range args {
		masks = append(masks, files.StrToMask(arg))
	}

	r.options = append(r.options, ls.Masks(ls.Exclude, masks...))
}

func (r *Runner) Slice(args []string) {
	slog.Debug("slice", "pattern", args[0])
	r.processes = append(r.processes, ls.SliceProcess(args[0]))
}

func (r *Runner) Sort(args []string) {
	slog.Debug("sort", "by", args[0])
	r.processes = append(r.processes, ls.Sort(StrToSortBy(args[0])))
}

func (r *Runner) Search(args []string) {
	r.options = append(r.options, ls.Search(args...))
}

func (r *Runner) Dir(args ...string) {
	r.paths = append(r.paths, args...)
}

func (r *Runner) Eval() (rr []*lfs.Element) {
	if !r.direct {
		if len(r.paths) == 0 {
			r.paths = append(r.paths, "./")
		}
		r.options = append(r.options, ls.Paths(r.paths...))
		r.options = common.ReverseC(r.options)
		if len(r.options) != 0 {
			ar := []ls.Process{ls.Combine(ls.Dir(r.options...))}
			r.processes = append(ar, r.processes...)
		}
		slog.Debug("evaluating", "paths", r.paths, "options", len(r.options), "processes", len(r.processes))
		rr = ls.Do(r.processes...)
	} else {
		rr = res
		slog.Debug("direct eval", "results", len(res))
		for _, proc := range r.processes {
			rr = proc(rr)
		}
	}
	r.processes = nil
	r.options = nil
	r.paths = nil
	return
}

func RunList(args []string) {
	sars := [][][]string{}
	thenR := SplitArgs(args, "then")
	if len(thenR) == 0 {
		sar := SplitArgs(args, "also")
		slog.Info("sar", "len", len(sar))
		if len(sar) == 0 {
			rr := RunExpr(args)
			res = rr
			return
		}

		RunAlso(sar)
		return
	}

	for _, args := range thenR {
		sars = append(sars, SplitArgs(args, "also"))
	}

	RunThen(sars)
}

func RunExpr(args []string) []*lfs.Element {
	rest, _ := Program.Eval(args)
	run.paths = append(run.paths, rest...)
	return run.Eval()
	// res = append(res, run.Eval()...)
}

func RunAlso(exprs [][]string) {
	var sum []*lfs.Element
	for i, args := range exprs {
		slog.Debug("run also loop", "iteration", i, "args", args)
		rr := RunExpr(args)
		sum = append(sum, rr...)
	}

	if len(sum) != 0 {
		res = sum
	}
}
func RunThen(flows [][][]string) []*lfs.Element {
	slog.Debug("FLOWS", "args", flows)
	for _, flow := range flows {
		slog.Debug("current flow", "args", flow)
		RunAlso(flow)
		run.direct = true
	}
	return res
}

func SplitArgs(args []string, key string) [][]string {
	var res [][]string
	var current []string
	for _, arg := range args {
		if arg == key {
			res = append(res, current)
			current = []string{}
			continue
		}
		current = append(current, arg)
	}
	if len(current) != 0 {
		res = append(res, current)
	}
	return res
}
