package list

import (
	"log/slog"
	"math/rand"
	"sort"

	"github.com/facette/natsort"
	"github.com/periaate/slice"
)

type Process func(filenames []*Element) []*Element

const (
	ByNone SortBy = iota
	ByMod
	BySize
	ByCreation
	ByName
)

type SortBy uint8

func StrToSortBy(s string) SortBy {
	switch s {
	case "date", "mod", "time", "t":
		return ByMod
	case "creation", "c":
		return ByCreation
	case "size", "s":
		return BySize
	case "name", "n":
		return ByName
	case "none":
		fallthrough
	default:
		return ByNone
	}
}

func ProcessList(res *Result, fns []Process) {
	for _, fn := range fns {
		res.Files = fn(res.Files)
	}
}

func CollectProcess(opts *Options) []Process {
	var fns []Process

	switch {
	// case len(opts.Query) > 0:
	// 	fns = append(fns, QueryProcess(opts))
	case opts.Ascending || len(opts.Sort) != 0:
		sorting := StrToSortBy(opts.Sort)

		if sorting == ByNone {
			break
		}

		fns = append(fns, SortProcess(sorting))
	}

	if opts.Shuffle {
		source := rand.NewSource(rand.Int63())
		if opts.Seed != -1 {
			source = rand.New(rand.NewSource(opts.Seed))
		}
		fns = append(fns, ShuffleProcess(source))
	}

	if opts.Ascending {
		fns = append(fns, Reverse[*Element])
	}

	if len(opts.Select) > 0 {
		fns = append(fns, SliceProcess(opts.Select))
	}
	return fns
}

func Reverse[T any](filenames []T) (res []T) {
	res = make([]T, 0, len(filenames))
	for i := len(filenames) - 1; i >= 0; i-- {
		res = append(res, filenames[i])
	}
	return
}

func ShuffleProcess(src rand.Source) Process {
	return func(filenames []*Element) []*Element {
		for i := range filenames {
			j := src.Int63() % int64(len(filenames))
			filenames[i], filenames[j] = filenames[j], filenames[i]
		}
		return filenames
	}
}

func SortProcess(sorting SortBy) Process {
	return func(filenames []*Element) []*Element {
		if sorting == ByName {
			sort.Slice(filenames, func(i, j int) bool {
				return natsort.Compare(filenames[i].Name, filenames[j].Name)
			})

			return filenames
		}

		sort.Slice(filenames, func(i, j int) bool {
			return filenames[j].Vany < filenames[i].Vany
		})

		return filenames
	}
}

func SliceProcess(patterns []string) Process {
	exp := slice.NewExpression[*Element]()
	for _, pattern := range patterns {
		exp.Parse(pattern)
	}

	return func(filenames []*Element) (res []*Element) {
		res, err := exp.Eval(filenames)
		if err != nil {
			slog.Error("error in Slice", "error", err)
		}
		return
	}
}

func LookBehind() Process {
	return func(els []*Element) (res []*Element) {
		res = make([]*Element, 0, len(els))

		return
	}
}
