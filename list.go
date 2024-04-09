package list

func Run(opts *Options) *Result {
	res := &Result{Files: []*Finfo{}}
	filters := CollectFilters(opts)
	processes := CollectProcess(opts)
	wfn := BuildWalkDirFn(filters, res)

	Traverse(wfn, opts)
	ProcessList(res, processes)
	return res
}

func Initialize(opts *Options) (*Result, []Filter, []Process) {
	res := &Result{Files: []*Finfo{}}
	filters := CollectFilters(opts)
	processes := CollectProcess(opts)

	return res, filters, processes
}
