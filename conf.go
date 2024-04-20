package list

import "github.com/periaate/ls/lfs"

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
