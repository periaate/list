package list

import (
	"log/slog"
	"math"
	"strings"

	c "github.com/periaate/common"
	"github.com/periaate/slice"
)

const (
	SortR    = '#'
	SortS    = "#"
	SliceR   = '['
	SliceS   = "["
	SearchS  = "?["
	ExcludeS = "!["
	KindAddS = "+["
	KindSubS = "-["
)

var funcKeys = []string{SortS, SliceS, ExcludeS, SearchS, KindAddS, KindSubS}

func Strip(str string, left, right string) string {
	if len(str) > len(left) {
		if str[:len(left)] == left {
			str = str[len(left):]
		}
	}
	if len(str) > len(right) {
		if str[len(str)-len(right):] == right {
			str = str[:len(str)-len(right)]
		}
	}
	return str
}

const (
	traversal = "traversal"
	files     = "files"
	dirs      = "dirs"
)

type FPPair struct {
	Tar     string
	Filter  Filter  // Ran during traversal, regardless of pair used
	Process Process // Ran after traversal. Traversal can't be applied to traversal.
}

func TargetedSearch(opts *Options, pat string) (fpp *FPPair, ok bool) {
	if len(pat) == 0 {
		return
	}
	fpp = new(FPPair)
	switch {
	case pat[0] == 'r' || pat[0] == 't':
		slog.Debug("recursive traversal")
		if opts.ToDepth == 0 {
			opts.ToDepth = math.MaxInt
		}
		fpp.Tar = traversal
	case pat[0] == 'f':
		fpp.Tar = files
	case pat[0] == 'd':
		fpp.Tar = dirs
	case len(pat) < 2:
		return nil, false
	case pat[:2] == "!f" || pat[:2] == "?d":
		opts.DirOnly = true
		return nil, true
	case pat[:2] == "!d" || pat[:2] == "?f":
		opts.FileOnly = true
		return nil, true
	default:
		return
	}
	var fn func(string)

	if fpp.Tar == traversal {
		fn = func(s string) {
			f, t, err := slice.ParsePattern(s, math.MaxInt)
			if err != nil {
				slog.Error("error parsing pattern", "pattern", s, "error", err)
				return
			}
			opts.FromDepth = f
			opts.ToDepth = t
		}
	}

	if len(pat) > 1 {
		fpp.Filter, fpp.Process = Search(pat[1:], fpp.Tar, fn)
	}
	return fpp, true
}

func Search(past string, tar string, recCase ...func(string)) (Filter, Process) {
	split := slice.SplitWithAll(past, funcKeys...)
	sf := []Filter{}
	sp := []Process{}
	for _, pat := range split {
		switch {
		case len(pat) == 0:
			continue
		case pat[0] == SortR:
			sp = append(sp, GetSortAlg(pat[1:]))
		case pat[0] == SliceR:
			if len(recCase) != 0 && recCase[0] != nil {
				recCase[0](pat)
				continue
			}
			sp = append(sp, SliceProcess([]string{pat}))
		case len(pat) < 2:
			continue
		case pat[:2] == SearchS:
			pat = Strip(pat, "?[", "]")
			sf = append(sf, GetSearchFilter(pat, false))
		case pat[:2] == ExcludeS:
			pat = Strip(pat, "![", "]")
			sf = append(sf, GetSearchFilter(pat, true))
		case pat[:2] == KindAddS && tar == files:
			sf = append(sf, GetKindFilter(pat[2:], false))
		case pat[:2] == KindSubS && tar == files:
			sf = append(sf, GetKindFilter(pat[2:], true))
		}
	}

	if len(sf) == 0 {
		sf = append(sf, NoneFilter)
	}

	if len(sp) == 0 {
		sp = append(sp, NoneSort)
	}

	filterFuncs := make([]func(*Element) bool, 0, len(sf))
	for _, f := range sf {
		filterFuncs = append(filterFuncs, f)
	}

	processFuncs := make([]func([]*Element) []*Element, 0, len(sp))
	for _, p := range sp {
		processFuncs = append(processFuncs, p)
	}

	return c.All(true, filterFuncs...),
		c.Pipe(processFuncs...)
}

func GetKindFilter(pat string, exclude bool) Filter {
	pat = Strip(pat, "+[", "]")
	pat = Strip(pat, "-[", "]")
	pat = strings.ToLower(pat)
	var kindMask uint32
	pats := strings.Split(pat, ".")
	for _, p := range pats {
		kindMask |= StrToMask(p)
	}

	if exclude {
		return func(el *Element) bool {
			return el.Mask&kindMask == 0
		}
	}
	return func(el *Element) bool {
		return el.Mask&kindMask != 0
	}
}

func GetSearchFilter(pat string, exclude bool) (res Filter) {
	switch pat[0] {
	case '=':
		res = getMatch(pat[1:])
	// case '~':
	// 	res = QueryAsFilter(pat[1:])
	default:
		res = getSubstringMatch(pat[1:])
	}
	if exclude {
		res = c.Negate(res)
	}
	return
}

func getSubstringMatch(str string) Filter {
	str = strings.ToLower(str)
	return func(e *Element) bool {
		return strings.Contains(strings.ToLower(e.Name), str)
	}
}

func getMatch(str string) Filter {
	str = strings.ToLower(str)
	return func(e *Element) bool {
		return str == strings.ToLower(e.Name)
	}
}

func GetSortAlg(pat string) Process {
	var res Process
	alg, _ := c.First([]rune(pat), func(r rune) bool {
		return r == 'n' || r == 't' || r == 'c' || r == 's'
	})
	dir, _ := c.First([]rune(pat), func(r rune) bool {
		return r == 'a' || r == 'r'
	})

	switch alg {
	case 'n':
		res = SortProcess(ByName)
	case 't':
		res = SortProcess(ByMod)
	case 'c':
		res = SortProcess(ByCreation)
	case 's':
		res = SortProcess(BySize)
	default:
		res = NoneSort
	}

	if dir == 'r' || dir == 'a' {
		return func(filenames []*Element) []*Element {
			return Reverse(res(filenames))
		}
	}

	return res
}

func NoneSort(els []*Element) []*Element {
	return els
}

func NoneFilter(e *Element) bool {
	return true
}
