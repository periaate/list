package sorting

import (
	"sort"
	"strconv"
	"strings"
	"unicode"
)

const (
	contains = 100.0 // the score for a perfect match
)

var (
	N         = 3
	QueryGram map[string]bool
	Query     string
)

// SortableFiles is an array of sortableFiles which is sortable.
type SortableFiles []*SortableFile

func (s SortableFiles) Len() int           { return len(s) }
func (s SortableFiles) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SortableFiles) Less(i, j int) bool { return naturalLess(s[i].SortableName, s[j].SortableName) }

type SortableFile struct {
	Fp           string  // Filepath
	SortableName string  // Filepath in lowercase for sorting
	Value        int64   // Countable value to be used by counting sort. Populated by unix timestamp.
	Score        float64 // Score for how well the string matches the query
	Ngram        map[string]bool
}

// naturalLess compares two strings and returns true if a < b in natural order.
func naturalLess(a, b string) bool {
	var ai, bi int
	for ai < len(a) && bi < len(b) {
		ach, bch := rune(a[ai]), rune(b[bi])
		if unicode.IsDigit(ach) && unicode.IsDigit(bch) {
			var anum, bnum string
			for ; ai < len(a) && unicode.IsDigit(rune(a[ai])); ai++ {
				anum += string(a[ai])
			}
			for ; bi < len(b) && unicode.IsDigit(rune(b[bi])); bi++ {
				bnum += string(b[bi])
			}
			an, _ := strconv.Atoi(anum)
			bn, _ := strconv.Atoi(bnum)
			if an != bn {
				return an < bn
			}
		} else {
			if ach != bch {
				return ach < bch
			}
			ai++
			bi++
		}
	}
	return len(a) < len(b)
}

func CalculateMatchScore(str string, n int) (score float64) {
	str = strings.ToLower(str)

	ngramLen := len(str) - n + 1
	ngrams := make(map[string]bool, ngramLen)

	for i := 0; i < len(str)-len(Query)+1; i++ {
		sub := str[i : i+len(Query)]
		if sub == Query {
			// Return as containing the query is considered a perfect match
			return contains
		}

		gram := str[i : i+n]
		ngrams[gram] = true
	}

	// check for ngrams
	for gram := range QueryGram {
		if _, ok := ngrams[gram]; ok {
			score += 1 / float64(len(Query))
		}
	}

	return score
}

func GenNgram(s string, n int) map[string]bool {
	l := len(s) - n + 1
	ngrams := make(map[string]bool, l)
	for i := 0; i < l; i++ {
		gram := s[i : i+n]
		ngrams[gram] = true
	}

	return ngrams
}

func Prune(files SortableFiles, score float64) (sf SortableFiles) {
	sf = make(SortableFiles, 0, len(files))
	for _, f := range files {
		if f.Score > score {
			sf = append(sf, f)
		}
	}
	return sf
}

func SortByScore(files SortableFiles) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Score > files[j].Score
	})
}

// CountingSort sorts an array into descending order.
//
// Counting sort needs to be called with an array, an integer k which is
// the largest value in the array, and a function which takes an element
// of the array as argument, and returns its value in the range [0, k].
func CountingSort(input []*SortableFile, lowestTime, highestTime int64) []*SortableFile {
	k := int(highestTime - lowestTime)
	count := make([]int, k+1)

	// Count occurrences of each value.
	for _, v := range input {
		count[v.Value-lowestTime]++
	}

	// Build and apply offset by summing the counts of previous values.
	for i := 1; i <= k; i++ {
		count[i] += count[i-1]
	}

	result := make([]*SortableFile, len(input))
	for _, v := range input {
		result[len(input)-count[v.Value-lowestTime]] = v
		count[v.Value-lowestTime]--
	}

	return result
}
