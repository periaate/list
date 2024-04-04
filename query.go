package list

import (
	"list/cfg"
	"sort"
	"strings"
)

const N = 3

type scored[T any] struct {
	t     T
	score float32
}

type scoredFiles[T any] []scored[T]

// Items returns all T with a score above 0.
func (s scoredFiles[T]) Items() []T {
	items := make([]T, 0, len(s))
	for i, v := range s {
		if v.score == 0 {
			continue
		}
		items = append(items, s[i].t)
	}
	return items
}

func QueryProcess(filenames []*Finfo) []*Finfo {
	scorer := getScoringFunction(cfg.Opts.Query)
	scorable := scoredFiles[*Finfo](make([]scored[*Finfo], len(filenames)))
	for i, file := range filenames {
		score := scorer(file.name)
		scorable[i] = scored[*Finfo]{file, score}
	}

	SortByScore(scorable)
	return scorable.Items()
}

func getScoringFunction(queries []string) func(string) float32 {
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

func SortByScore[T any](files scoredFiles[T]) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].score > files[j].score
	})
}
