package main

import (
	"math"
	"sort"
)

// Counting sort should only be used if there is a large number of elements or a low gap between high and low.
var (
	highestTime int64                 // Highest found unix timestamp.
	lowestTime  int64 = math.MaxInt64 // Lowest found unix timestamp. Uses MaxInt64 to make comparisons work.
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
	case Opts.Ascending || Opts.Date || Opts.Sort:
		sorting := byName
		order := toDesc
		if Opts.Ascending {
			order = toAsc
		}
		if Opts.Date {
			sorting = byDate
		}
		fns = append(fns, sortProcess(sorting, order))
	}

	// 4 is the minimum length of a slice pattern, as [:n] or [n:] are the smallest possible patterns.
	if len(Opts.Slice) >= 3 {
		fns = append(fns, sliceProcess(Opts.Slice))
	}
	return fns
}

func sortProcess(sorting sortBy, ordering orderTo) process {
	return func(filenames []*finfo) []*finfo {
		switch sorting {
		case byDate:
			filenames = countingSort(filenames, lowestTime, highestTime)
			if ordering == toAsc {
				return reverse(filenames)
			}
		case byName:
			sort.Slice(filenames, func(i, j int) bool {
				return natural(filenames[j].name, filenames[i].name)
			})
			if ordering == toAsc {
				return reverse(filenames)
			}
		}

		return filenames
	}
}

func sliceProcess(pattern string) process {
	return func(filenames []*finfo) []*finfo {
		return sliceArray(pattern, filenames)
	}
}
