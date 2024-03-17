package main

import (
	"io/fs"
	"strings"
)

// updateTime updates the unix timestamp boundaries.
func updateTime(t int64) {
}
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

	if Opts.Date {
		fns = append(fns, addDate)
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

func addDate(fi *finfo, d fs.DirEntry) bool {
	fileinfo, err := d.Info()
	if err != nil || fileinfo == nil {
		return false
	}

	info, err := d.Info()
	if err != nil {
		return false
	}
	unixTime := info.ModTime().Unix()
	updateTime(unixTime)

	if unixTime > highestTime {
		highestTime = unixTime
	}
	if unixTime < lowestTime {
		lowestTime = unixTime
	}
	fi.mod = unixTime
	return true
}

func filterList(include []contentType, exclude []contentType, ignore []string, search []string) filter {
	return func(fi *finfo, _ fs.DirEntry) bool {
		var any bool
		for _, s := range search {
			if strings.Contains(fi.name, s) {
				any = true
				break
			}
		}
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
