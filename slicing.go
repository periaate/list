package list

import (
	"golang.org/x/exp/constraints"
)

func Clamp[T constraints.Ordered](val, lower, upper T) (res T) {
	if val >= upper {
		return upper
	}
	if val <= lower {
		return lower
	}
	return val
}

// SliceArray takes a string pattern and a generic slice, then returns a slice according to the pattern.
func SliceArray[T any](pattern string, input []T) []T {
	if len(input) == 0 {
		return input
	}
	// // Default slice indices
	iar, isSlice, ok := ParseSlice(pattern)
	if !ok {
		return input
	}

	for i, v := range iar {
		if v < 0 {
			iar[i] = len(input) + v
		}
		iar[i] = Clamp(iar[i], 0, len(input))
	}

	if !isSlice {
		return []T{input[iar[0]]}
	}
	if iar[1] == 0 {
		iar[1] = len(input) - 1
	}
	iar[1]++

	start := Clamp(iar[0], iar[0], len(input))
	end := Clamp(iar[1], iar[0], len(input))
	start = Clamp(start, 0, end)

	return input[start:end]
}
