package list

import (
	"io/fs"
	"strings"
)

type Filter func(*Finfo, fs.DirEntry) bool

func CollectFilters(opts *Options) []Filter {
	var fns []Filter

	switch {
	case opts.DirOnly:
		fns = append(fns, func(_ *Finfo, d fs.DirEntry) bool {
			return d.IsDir()
		})
	case opts.FileOnly:
		fns = append(fns, func(_ *Finfo, d fs.DirEntry) bool {
			return !d.IsDir()
		})
	}

	switch StrToSortBy(opts.Sort) {
	case ByMod:
		fns = append(fns, addModT)
	case BySize:
		fns = append(fns, addSize)
	case ByCreation:
		fns = append(fns, addCreationT)
	}

	if (len(opts.Search) + len(opts.Include) + len(opts.Exclude) + len(opts.Ignore)) > 0 {
		fns = append(fns, FilterList(opts))
	}
	return fns
}

func addModT(fi *Finfo, d fs.DirEntry) bool {
	info, err := d.Info()
	if err != nil || info == nil {
		return false
	}

	unixTime := info.ModTime().Unix()

	fi.Vany = unixTime
	return true
}

func addSize(fi *Finfo, d fs.DirEntry) bool {
	info, err := d.Info()
	if err != nil || info == nil {
		return false
	}

	fi.Vany = info.Size()
	return true
}

func FilterList(opts *Options) Filter {
	incMask := AsMask(opts.Include)
	excMask := AsMask(opts.Exclude)

	var searchFn func(string) bool
	if opts.SearchAnd {
		searchFn = func(str string) bool {
			for _, sub := range opts.Search {
				if !strings.Contains(str, sub) {
					return false
				}
			}
			return true
		}
	} else {
		searchFn = func(str string) bool {
			for _, sub := range opts.Search {
				if strings.Contains(str, sub) {
					return true
				}
			}
			return false
		}
	}

	return func(fi *Finfo, _ fs.DirEntry) bool {
		any := searchFn(fi.Name)
		if len(opts.Search) > 0 && !any {
			return false
		}

		if incMask > 0 && (incMask&fi.Mask) == 0 {
			return false
		}

		for _, ign := range opts.Ignore {
			if strings.Contains(fi.Path, ign) {
				return false
			}
		}

		if excMask > 0 && (excMask&fi.Mask) != 0 {
			return false
		}

		return true
	}
}
