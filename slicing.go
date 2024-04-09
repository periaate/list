package list

import (
	"fmt"
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

const (
	f = 0
	t = 1

	pageToken = '='
)

func Slice[T any](pat string, input []T) (_ []T, err error) {
	var from, to int

	if len(input) == 0 {
		err = fmt.Errorf("input is empty")
		slog.Debug(err.Error())
		return
	}

	if len(pat) < 3 {
		err = fmt.Errorf("last argument is not long enough to be a slice")
		slog.Debug(err.Error())
		return
	}

	L := len(pat) - 1
	if pat[0] != '[' || pat[L] != ']' {
		err = fmt.Errorf("last argument does not match slice pattern, does not start and end with brackets")
		slog.Debug(err.Error())
		return
	}
	pat = pat[1:L]
	slog.Debug("slice pattern", "pattern", pat)

	for _, r := range pat {
		if !(r == '-' || r == '+' || r == ':' || r == pageToken || r >= '0' || r <= '9') {
			err = fmt.Errorf("slice pattern included non integer values")
			slog.Debug(err.Error())
			return
		}
	}

	pageSize := 1

	if ind := strings.Index(pat, string(pageToken)); ind != -1 {
		vl := pat[ind+1:]
		pageSize, err = strconv.Atoi(vl)
		if err != nil {
			slog.Debug(err.Error())
			return
		}

		pat = pat[:ind]
	}

	ind := strings.Index(pat, ":")

	if len(pat) == 1 || (ind == -1 && pat[0] == '-') {
		slog.Debug("slice pattern is only one character long")
		if pat[0] == ':' {
			return input, nil
		}

		if pat[0] == '-' {
			slog.Debug("negative single index")
			from, err = parseMinus(pat, len(input), pageSize)
			if err != nil {
				slog.Debug(err.Error())
				return
			}
			to = from + 1*pageSize
		} else {
			from, err = strconv.Atoi(pat)
			if err != nil {
				slog.Debug(err.Error())
				return
			}

			to = from + 1
			to *= pageSize
			from *= pageSize
		}

		slog.Debug("slice results", "from", from, "to", to, "pagesize", pageSize, "input length", len(input))
		from = Clamp(from, 0, len(input))
		to = Clamp(to, 0, len(input))
		from = Clamp(from, 0, to)
		slog.Debug("clamped results", "from", from, "to", to, "pagesize", pageSize, "input length", len(input))
		return input[from:to], nil
	}

	if ind == -1 {
		err = fmt.Errorf("slice pattern does not contain a colon")
		slog.Debug(err.Error())
		return
	}

	fromTo := []string{pat[:ind], pat[ind+1:]}

	if fromTo[f] == "" {
		from = 0
	} else {
		if fromTo[f][0] == '-' {
			from, err = parseMinus(fromTo[f], len(input), pageSize)
			if err != nil {
				slog.Debug(err.Error())
				return
			}
		} else {
			from, err = strconv.Atoi(fromTo[f])
			if err != nil {
				slog.Debug(err.Error())
				return
			}
			from *= pageSize
		}
	}

	if fromTo[t] == "" {
		to = len(input)
	} else {
		switch fromTo[t][0] {
		case '+':
			to, err = parsePlus(fromTo[t], from, pageSize)
			if err != nil {
				slog.Debug(err.Error())
				return
			}
		case '-':
			to, err = parseMinus(fromTo[t], len(input), pageSize)
			if err != nil {
				slog.Debug(err.Error())
				return
			}
		default:
			to, err = strconv.Atoi(fromTo[t])
			if err != nil {
				slog.Debug(err.Error())
				return
			}
			to *= pageSize
		}
	}

	slog.Debug("slice results", "from", from, "to", to, "pagesize", pageSize, "input length", len(input))
	from = Clamp(from, 0, len(input))
	to = Clamp(to, 0, len(input))
	from = Clamp(from, 0, to)
	slog.Debug("clamped results", "from", from, "to", to, "pagesize", pageSize, "input length", len(input))

	return input[from:to], nil
}

func parseMinus(s string, l, size int) (int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return l + val*size, nil
}

func parsePlus(s string, l, size int) (int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return l + val*size, nil
}
