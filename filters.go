package list

import (
	"io/fs"
	"list/cfg"
	"strings"
)

func CollectFilters() []filter {
	var fns []filter

	switch {
	case cfg.Opts.DirOnly:
		fns = append(fns, func(_ *Finfo, d fs.DirEntry) bool {
			return d.IsDir()
		})
	case cfg.Opts.FileOnly:
		fns = append(fns, func(_ *Finfo, d fs.DirEntry) bool {
			return !d.IsDir()
		})
	}

	switch StrToSortBy(cfg.Opts.Sort) {
	case byMod:
		fns = append(fns, addModT)
	case bySize:
		fns = append(fns, addSize)
	}

	if (len(cfg.Opts.Search) + len(cfg.Opts.Include) + len(cfg.Opts.Exclude) + len(cfg.Opts.Ignore)) > 0 {
		fns = append(fns, filterList(cfg.Opts.Include, cfg.Opts.Exclude, cfg.Opts.Ignore, cfg.Opts.Search))
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

func filterList(include []string, exclude []string, ignore []string, search []string) filter {
	// to avoid checking flags for every element.
	var searchFn func(string) bool
	if cfg.Opts.SearchAnd {
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

	return func(fi *Finfo, _ fs.DirEntry) bool {
		any := searchFn(fi.name)
		if len(search) > 0 && !any {
			return false
		}

		ext := GetContentTypes(fi.name)

		for _, inc := range include {
			if !ext.contains(inc) {
				return false
			}
		}

		for _, ign := range ignore {
			if strings.Contains(fi.path, ign) {
				return false
			}
		}

		for _, exc := range exclude {
			if ext.contains(exc) {
				return false
			}
		}

		return true
	}
}
