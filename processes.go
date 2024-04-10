package list

import (
	"math/rand"
	"sort"

	"github.com/facette/natsort"
)

func ProcessList(res *Result, fns []Process) {
	for _, fn := range fns {
		res.Files = fn(res.Files)
	}
}

func Reverse[T any](filenames []T) []T {
	for i := 0; i < len(filenames)/2; i++ {
		j := len(filenames) - i - 1
		filenames[i], filenames[j] = filenames[j], filenames[i]
	}
	return filenames
}

func CollectProcess(opts *Options) []Process {
	var fns []Process

	switch {
	case len(opts.Query) > 0:
		fns = append(fns, QueryProcess(opts))
		if opts.Ascending {
			fns = append(fns, Reverse[*Finfo])
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

	if opts.Shuffle {
		source := rand.NewSource(rand.Int63())
		if opts.Seed != -1 {
			source = rand.New(rand.NewSource(opts.Seed))
		}
		fns = append(fns, ShuffleProcess(source))
	}
	return fns
}

func ShuffleProcess(src rand.Source) Process {
	return func(filenames []*Finfo) []*Finfo {
		for i := range filenames {
			j := src.Int63() % int64(len(filenames))
			filenames[i], filenames[j] = filenames[j], filenames[i]
		}
		return filenames
	}
}

func SortProcess(sorting SortBy, ordering OrderTo) Process {
	return func(filenames []*Finfo) []*Finfo {
		if sorting == ByName {
			sort.Slice(filenames, func(i, j int) bool {
				return natsort.Compare(filenames[i].Name, filenames[j].Name)
				// return natural(filenames[j].name, filenames[i].name)
			})
			if ordering == ToAsc {
				return Reverse(filenames)
			}

			return filenames
		}

		sort.Slice(filenames, func(i, j int) bool {
			return filenames[j].Vany < filenames[i].Vany
		})
		if ordering == ToAsc {
			return Reverse(filenames)
		}

		return filenames
	}
}

func SliceProcess(pattern string) Process {
	return func(filenames []*Finfo) []*Finfo {
		res, err := Slice(pattern, filenames)
		if err != nil {
			return filenames
		}
		return res
	}
}
