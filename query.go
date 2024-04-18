package list

import (
	"log/slog"
	"sort"
	"strings"
)

const N = 3

type scored[T any] struct {
	t     T
	score float32
}

type ScoredFiles[T any] []scored[T]

// Items returns all T with a score above 0.
func (s ScoredFiles[T]) Items() []T {
	items := make([]T, 0, len(s))
	for i, v := range s {
		if v.score == 0 {
			continue
		}
		items = append(items, s[i].t)
	}
	return items
}

// func QueryProcess(opts *Options) Process {
// 	return func(filenames []*Element) []*Element {
// 		scorer := GetScoringFunction(opts.Query)
// 		scorable := ScoredFiles[*Element](make([]scored[*Element], len(filenames)))
// 		for i, file := range filenames {
// 			score := scorer(file.Name)
// 			scorable[i] = scored[*Element]{file, score}
// 		}

// 		SortByScore(scorable)
// 		return scorable.Items()
// 	}
// }

func GetScoringFunction(queries []string) func(string) float32 {
	queryGrams, _ := GenNgrams(queries, N)
	n := N
	return func(str string) (score float32) {
		str = strings.ToLower(str)

		ngramLen := len(str) - n + 1
		ngrams := make(map[string]int, ngramLen)
		for i := 0; i < len(str)-n+1; i++ {
			gram := str[i : i+n]
			ngrams[gram]++
		}

		// check for ngrams
		for gram, weight := range queryGrams {
			if v, ok := ngrams[gram]; ok {
				score += (float32(v) * float32(weight)) / float32(len(queryGrams))
			}
		}

		return score
	}
}

func GenNgrams(sar []string, n int) (map[string]int, int) {
	var qlen int
	l := len(sar) - n + 1
	ngrams := make(map[string]int, l)
	for _, query := range sar {
		qlen += len(query)
		query = strings.ToLower(query)
		for i := 0; i < len(query)-n+1; i++ {
			gram := query[i : i+n]
			ngrams[gram]++
		}
	}

	return ngrams, qlen
}

func SortByScore[T any](files ScoredFiles[T]) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].score > files[j].score
	})
}

func QueryAsFilter(qr []string) Filter {
	scorer := GetScoringFunction(qr)
	return func(e *Element) bool {
		score := scorer(e.Name)
		slog.Debug("query filter", "name", e.Name, "score", score)
		return score != 0
	}
}
