package list

import (
	"list/cfg"
	"sort"

	"github.com/facette/natsort"
)

func ProcessList(res *Result, fns []process) {
	for _, fn := range fns {
		res.Files = fn(res.Files)
	}
}

func reverse(filenames []*Finfo) []*Finfo {
	for i := 0; i < len(filenames)/2; i++ {
		j := len(filenames) - i - 1
		filenames[i], filenames[j] = filenames[j], filenames[i]
	}
	return filenames
}

func CollectProcess() []process {
	var fns []process

	switch {
	case len(cfg.Opts.Query) > 0:
		fns = append(fns, QueryProcess)
		if cfg.Opts.Ascending {
			fns = append(fns, reverse)
		}
	case cfg.Opts.Ascending || len(cfg.Opts.Sort) != 0:
		sorting := StrToSortBy(cfg.Opts.Sort)

		if sorting == byNone {
			break
		}

		order := toDesc
		if cfg.Opts.Ascending {
			order = toAsc
		}
		fns = append(fns, SortProcess(sorting, order))
	}

	if len(cfg.Opts.Select) >= len("[0]") {
		fns = append(fns, SliceProcess(cfg.Opts.Select))
	}
	return fns
}

func SortProcess(sorting sortBy, ordering orderTo) process {
	return func(filenames []*Finfo) []*Finfo {
		if sorting == byName {
			sort.Slice(filenames, func(i, j int) bool {
				return natsort.Compare(filenames[i].name, filenames[j].name)
				// return natural(filenames[j].name, filenames[i].name)
			})
			if ordering == toAsc {
				return reverse(filenames)
			}

			return filenames
		}

		sort.Slice(filenames, func(i, j int) bool {
			return filenames[j].vany < filenames[i].vany
		})
		if ordering == toAsc {
			return reverse(filenames)
		}

		return filenames
	}
}

func SliceProcess(pattern string) process {
	return func(filenames []*Finfo) []*Finfo {
		return SliceArray(pattern, filenames)
	}
}
