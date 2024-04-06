package list

import (
	"sort"

	"github.com/periaate/list/cfg"

	"github.com/facette/natsort"
)

func ProcessList(res *Result, fns []Process) {
	for _, fn := range fns {
		res.Files = fn(res.Files)
	}
}

func Reverse(filenames []*Finfo) []*Finfo {
	for i := 0; i < len(filenames)/2; i++ {
		j := len(filenames) - i - 1
		filenames[i], filenames[j] = filenames[j], filenames[i]
	}
	return filenames
}

func CollectProcess(opts *cfg.Options) []Process {
	var fns []Process

	switch {
	case len(opts.Query) > 0:
		fns = append(fns, QueryProcess(opts))
		if opts.Ascending {
			fns = append(fns, Reverse)
		}
	case opts.Ascending || len(opts.Sort) != 0:
		sorting := StrToSortBy(opts.Sort)

		if sorting == ByNone {
			break
		}

		order := ToDesc
		if opts.Ascending {
			order = ToAsc
		}
		fns = append(fns, SortProcess(sorting, order))
	}

	if len(opts.Select) >= len("[0]") {
		fns = append(fns, SliceProcess(opts.Select))
	}
	return fns
}

func SortProcess(sorting SortBy, ordering OrderTo) Process {
	return func(filenames []*Finfo) []*Finfo {
		if sorting == ByName {
			sort.Slice(filenames, func(i, j int) bool {
				return natsort.Compare(filenames[i].name, filenames[j].name)
				// return natural(filenames[j].name, filenames[i].name)
			})
			if ordering == ToAsc {
				return Reverse(filenames)
			}

			return filenames
		}

		sort.Slice(filenames, func(i, j int) bool {
			return filenames[j].vany < filenames[i].vany
		})
		if ordering == ToAsc {
			return Reverse(filenames)
		}

		return filenames
	}
}

func SliceProcess(pattern string) Process {
	return func(filenames []*Finfo) []*Finfo {
		return SliceArray(pattern, filenames)
	}
}
