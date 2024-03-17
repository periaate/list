package main

import (
	"log/slog"
	"strconv"
	"strings"

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

// sliceArray takes a string pattern and a generic slice, then returns a slice according to the pattern.
func sliceArray[T any](pattern string, input []T) []T {
	if len(input) == 0 {
		return input
	}
	// // Default slice indices
	iar, isSlice, ok := parseSlice(pattern)
	if !ok {
		return input
	}

	for i, v := range iar {
		if v < 0 {
			iar[i] = len(input) - 1 + v
		}
		iar[i] = Clamp(iar[i], 0, len(input)-1)
	}

	if !isSlice {
		return []T{input[iar[0]]}
	}
	if iar[1] == 0 {
		iar[1] = len(input)
	}

	start := Clamp(iar[0], iar[0], len(input)-1)
	end := Clamp(iar[1], iar[0], len(input)-1)
	start = Clamp(start, 0, end)

	return input[start:end]
}

func parseSlice(inp string) (iar []int, isSlice bool, ok bool) {
	if len(inp) < 3 {
		slog.Debug("last argument is not long enough to be a slice")
		return
	}
	L := len(inp) - 1
	if inp[0] != '[' || inp[L] != ']' {
		slog.Debug("last argument does not match slice pattern, is not within brackets")
		return
	}

	slice := strings.Split(inp[1:L], ":")

	if len(slice) > 2 {
		slog.Debug("last argument does not match slice pattern, split returned too many items")
		return
	}

	for _, s := range slice {
		if len(s) == 0 {
			continue
		}
		if !isInt(s) {
			slog.Debug("slice pattern included non integer values")
			return
		}
	}

	iar = make([]int, len(slice))
	for i, s := range slice {
		iar[i], _ = strconv.Atoi(s)
	}

	return iar, len(iar) > 1, true
}

func isInt(s string) bool {
	rar := []rune(s)
	if len(rar) == 0 {
		return false
	}
	if a := rar[0]; a == '-' {
		if len(rar) == 1 {
			return false
		}
		rar = rar[1:]
	}

	for _, r := range rar {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
