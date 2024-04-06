package list

import (
	"io/fs"
	"strings"

	"github.com/periaate/list/cfg"
)

func CollectFilters(opts *cfg.Options) []filter {
	var fns []filter

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
	case byMod:
		fns = append(fns, addModT)
	case bySize:
		fns = append(fns, addSize)
	}

	if (len(opts.Search) + len(opts.Include) + len(opts.Exclude) + len(opts.Ignore)) > 0 {
		fns = append(fns, filterList(opts))
	}
	return fns
}

func addModT(fi *Finfo, d fs.DirEntry) bool {
	fileinfo, err := d.Info()
	if err != nil || fileinfo == nil {
		return false
	}

	info, err := d.Info()
	if err != nil {
		return false
	}
	unixTime := info.ModTime().Unix()

	fi.vany = unixTime
	return true
}

func addSize(fi *Finfo, d fs.DirEntry) bool {
	fileinfo, err := d.Info()
	if err != nil || fileinfo == nil {
		return false
	}

	fi.vany = fileinfo.Size()
	return true
}

func filterList(opts *cfg.Options) filter {
	// to avoid checking flags for every element.
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
		any := searchFn(fi.name)
		if len(opts.Search) > 0 && !any {
			return false
		}

		ext := GetContentTypes(fi.name)

		for _, inc := range opts.Include {
			if !ext.contains(inc) {
				return false
			}
		}

		for _, ign := range opts.Ignore {
			if strings.Contains(fi.path, ign) {
				return false
			}
		}

		for _, exc := range opts.Exclude {
			if ext.contains(exc) {
				return false
			}
		}

		return true
	}
}
