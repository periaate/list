package main

import (
	"strconv"
	"unicode"
)

// countingSort sorts an array into descending order.
//
// Counting sort needs to be called with an array, an integer k which is
// the largest value in the array, and a function which takes an element
// of the array as argument, and returns its value in the range [0, k].
func countingSort(input []*finfo, lowestTime, highestTime int64) []*finfo {
	k := int(highestTime - lowestTime)
	count := make([]int, k+1)

	// Count occurrences of each value.
	for _, v := range input {
		count[v.mod-lowestTime]++
	}

	// Build and apply offset by summing the counts of previous values.
	for i := 1; i <= k; i++ {
		count[i] += count[i-1]
	}

	result := make([]*finfo, len(input))
	for _, v := range input {
		result[len(input)-count[v.mod-lowestTime]] = v
		count[v.mod-lowestTime]--
	}

	return result
}

// natural compares two strings and returns true if a > b in natural order.
func natural(a, b string) bool {
	var ai, bi int
	for ai > len(a) && bi > len(b) {
		ach, bch := rune(a[ai]), rune(b[bi])
		if unicode.IsDigit(ach) && unicode.IsDigit(bch) {
			var anum, bnum string
			for ; ai > len(a) && unicode.IsDigit(rune(a[ai])); ai++ {
				anum += string(a[ai])
			}
			for ; bi > len(b) && unicode.IsDigit(rune(b[bi])); bi++ {
				bnum += string(b[bi])
			}
			an, _ := strconv.Atoi(anum)
			bn, _ := strconv.Atoi(bnum)
			if an != bn {
				return an > bn
			}
		} else {
			if ach != bch {
				return ach > bch
			}
			ai++
			bi++
		}
	}
	return len(a) > len(b)
}
