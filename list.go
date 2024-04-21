package list

import (
	"bufio"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/periaate/ls/lfs"
)

func Parse(args []string) []*lfs.Element {
	RunList(args)
	return res
}

var (
	Quiet bool
	Abs   bool
)

func PrintWithBuf(els []*lfs.Element) {
	const bufLength = 500
	if len(els) == 0 {
		return
	}
	if Quiet {
		slog.Debug("quiet flag is set, returning from print function")
		return
	}

	w := bufio.NewWriterSize(os.Stdout, 4096*bufLength)

	for i, file := range els {
		fp := filepath.ToSlash(file.Path)
		if Abs {
			fp, _ = filepath.Abs(file.Path)
			fp = filepath.ToSlash(fp)
		}
		res := fp + "\n"

		w.WriteString(res)
		if i%bufLength == 0 {
			w.Flush()
		}
	}

	w.Flush()
}

func StrToSortBy(s string) lfs.SortBy {
	switch s {
	case "date", "mod", "time", "t":
		return lfs.ByMod
	case "creation", "c":
		return lfs.ByCreation
	case "size", "s":
		return lfs.BySize
	case "name", "n":
		return lfs.ByName
	case "none":
		fallthrough
	default:
		return lfs.ByNone
	}
}
