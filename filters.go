package main

import (
	"io/fs"
	"strings"
)

func collectFilters() []filter {
	var fns []filter

	switch {
	case Opts.DirOnly:
		fns = append(fns, func(_ *finfo, d fs.DirEntry) bool {
			return d.IsDir()
		})
	case Opts.FileOnly:
		fns = append(fns, func(_ *finfo, d fs.DirEntry) bool {
			return !d.IsDir()
		})
	}

	switch strToSortBy(Opts.Sort) {
	case byMod:
		fns = append(fns, addModT)
	case bySize:
		fns = append(fns, addSize)
	}

	if (len(Opts.Search) + len(Opts.Include) + len(Opts.Exclude) + len(Opts.Ignore)) > 0 {
		var include []contentType
		var exclude []contentType
		for _, inc := range Opts.Include {
			include = append(include, stringToContentType(inc))
		}
		for _, exc := range Opts.Exclude {
			exclude = append(exclude, stringToContentType(exc))
		}

		fns = append(fns, filterList(include, exclude, Opts.Ignore, Opts.Search))
	}
	return fns
}

func addModT(fi *finfo, d fs.DirEntry) bool {
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

func addSize(fi *finfo, d fs.DirEntry) bool {
	fileinfo, err := d.Info()
	if err != nil || fileinfo == nil {
		return false
	}

	fi.vany = fileinfo.Size()
	return true
}

func filterList(include []contentType, exclude []contentType, ignore []string, search []string) filter {
	// to avoid checking flags for every element.
	var searchFn func(string) bool
	if Opts.SearchAnd {
		searchFn = func(str string) bool {
			for _, sub := range search {
				if !strings.Contains(str, sub) {
					return false
				}
			}
			return true
		}
	} else {
		searchFn = func(str string) bool {
			for _, sub := range search {
				if strings.Contains(str, sub) {
					return true
				}
			}
			return false
		}
	}

	return func(fi *finfo, _ fs.DirEntry) bool {
		any := searchFn(fi.name)
		if len(search) > 0 && !any {
			return false
		}

		for _, inc := range include {
			if inc != getContentType(fi.name) {
				return false
			}
		}

		for _, ign := range ignore {
			if strings.Contains(fi.path, ign) {
				return false
			}
		}

		for _, exc := range exclude {
			if exc == getContentType(fi.name) {
				return false
			}
		}

		return true
	}
}
