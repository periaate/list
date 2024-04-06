package list

import "github.com/periaate/list/cfg"

func Run(opts *cfg.Options) *Result {
	res := &Result{Files: []*Finfo{}}
	filters := CollectFilters(opts)
	processes := CollectProcess(opts)
	wfn := BuildWalkDirFn(filters, res)
	Traverse(wfn, opts)

	ProcessList(res, processes)
	return res
}
