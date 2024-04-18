package list

import (
	"log/slog"
	"sort"

	"github.com/facette/natsort"
	"github.com/periaate/common"
	"github.com/periaate/slice"
)

type Process func(els []*Element) []*Element

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

func CollectProcess(opts *Options) Process {
	switch {
	case opts.Ascending || len(opts.Sort) != 0:
		sorting := StrToSortBy(opts.Sort)

		if sorting == ByNone {
			break
		}

		opts.Processes = append(opts.Processes, SortProcess(sorting))
	}

	if opts.Ascending {
		opts.Processes = append(opts.Processes, Reverse[*Element])
	}

	if len(opts.Select) > 0 {
		opts.Processes = append(opts.Processes, SliceProcess(opts.Select))
	}
	return common.Pipe(opts.Processes...)
}

func Reverse[T any](filenames []T) (res []T) {
	res = make([]T, 0, len(filenames))
	for i := len(filenames) - 1; i >= 0; i-- {
		res = append(res, filenames[i])
	}
	return
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
