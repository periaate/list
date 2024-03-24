package main

import (
	"sort"

	"github.com/facette/natsort"
)

func ProcessList(res *result, fns []process) {
	for _, fn := range fns {
		res.files = fn(res.files)
	}
}

func reverse(filenames []*finfo) []*finfo {
	for i := 0; i < len(filenames)/2; i++ {
		j := len(filenames) - i - 1
		filenames[i], filenames[j] = filenames[j], filenames[i]
	}
	return filenames
}

func collectProcess() []process {
	var fns []process

	switch {
	case len(Opts.Query) > 0:
		fns = append(fns, queryProcess)
		if Opts.Ascending {
			fns = append(fns, reverse)
		}
	case Opts.Ascending || len(Opts.Sort) != 0:
		sorting := strToSortBy(Opts.Sort)

		if sorting == byNone {
			break
		}

		order := toDesc
		if Opts.Ascending {
			order = toAsc
		}
		fns = append(fns, sortProcess(sorting, order))
	}

	if len(Opts.Select) >= len("[0]") {
		fns = append(fns, sliceProcess(Opts.Select))
	}
	return fns
}

func sortProcess(sorting sortBy, ordering orderTo) process {
	return func(filenames []*finfo) []*finfo {
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

func sliceProcess(pattern string) process {
	return func(filenames []*finfo) []*finfo {
		return sliceArray(pattern, filenames)
	}
}
