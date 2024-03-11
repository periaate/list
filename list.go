package main

import (
	"fmt"
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// filters
// include file types
// exclude file types
// ignore patterns
// strict search

// selection
// slices -S

// sorting & ordering
// sort by date -d
// sort by filename -f (natural sorting)
// sort in ascending order -a

// onhold
// dir only
// file only
// fuzzy search with score

var (
	highestTime int64                 // Highest found unix timestamp.
	lowestTime  int64 = math.MaxInt64 // Lowest found unix timestamp. Uses MaxInt64 to make comparisons work.

	fp = "./"
)

type finfo struct {
	name string
	path string // includes name, relative path to cwd
	mod  int64  // unix timestamp
}

func List() {
	filters := collectFilters()
	processes := collectProcess()

	res := collectFiles(filters)

	ProcessList(res, processes)
	printWithBuf(res.files)
}
func collectFiles(filters []filter) *result {
	res := &result{[]*finfo{}}

	if len(Args) > 0 {
		fp = Args[0]
	}

	if Opts.Recurse {
		var walkFn fs.WalkDirFunc
		if Opts.Date {
			walkFn = filterFnRecursive(filters, res, optsHasDate)
		} else {
			walkFn = filterFnRecursive(filters, res)
		}

		stat, err := os.Stat(fp)
		if err != nil || !stat.IsDir() {
			log.Fatalln(err)
		}
		err = filepath.WalkDir(fp, walkFn)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		var fn func(string, fs.DirEntry) error
		if Opts.Date {
			fn = filterFn(filters, res, optsHasDate)
		} else {
			fn = filterFn(filters, res)
		}

		stat, err := os.Stat(fp)
		if err != nil || !stat.IsDir() {
			log.Fatalln(err)
		}
		entries, err := os.ReadDir(fp)
		if err != nil {
			log.Fatalln(err)
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			filep := filepath.Join(fp, entry.Name())
			err = fn(filep, entry)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	return res
}

// updateTime updates the unix timestamp boundaries.
func updateTime(t int64) {
	if t > highestTime {
		highestTime = t
	}
	if t < lowestTime {
		lowestTime = t
	}
}

func collectFilters() []filter {
	var fns []filter
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

func collectProcess() []process {
	var fns []process
	sorting := byName
	order := toDesc
	if Opts.Ascending {
		order = toAsc
	}
	if Opts.Date {
		sorting = byDate
	}
	fns = append(fns, sortProcess(sorting, order))
	// 4 is the minimum length of a slice pattern, as [:n] or [n:] are the smallest possible patterns.
	if len(Opts.Slice) >= 4 {
		fns = append(fns, sliceProcess(Opts.Slice))
	}
	return fns
}

type dirEntryfns func(*finfo, fs.DirEntry)

func optsHasDate(fi *finfo, d fs.DirEntry) {
	fileinfo, err := d.Info()
	if err != nil || fileinfo == nil {
		return
	}

	info, err := d.Info()
	if err != nil {
		return
	}
	unixTime := info.ModTime().Unix()
	updateTime(unixTime)
	fi.mod = unixTime
}

type result struct {
	files []*finfo
}

func filterFnRecursive(fns []filter, res *result, dirEntryfns ...dirEntryfns) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if d == nil || err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		fi := &finfo{name: d.Name(), path: path}
		for _, fn := range dirEntryfns {
			fn(fi, d)
		}
		for _, fn := range fns {
			res := fn(fi)
			if !res {
				return nil
			}
		}
		res.files = append(res.files, fi)
		return nil
	}
}

func filterFn(fns []filter, res *result, dirEntryfns ...dirEntryfns) func(string, fs.DirEntry) error {
	return func(path string, d fs.DirEntry) error {
		if d.IsDir() {
			return nil
		}
		fi := &finfo{name: d.Name(), path: path}
		for _, fn := range dirEntryfns {
			fn(fi, d)
		}
		for _, fn := range fns {
			res := fn(fi)
			if !res {
				return nil
			}
		}
		res.files = append(res.files, fi)
		return nil
	}
}

type filter func(*finfo) bool

func filterList(include []contentType, exclude []contentType, ignore []string, search []string) filter {
	return func(fi *finfo) bool {
		for _, s := range search {
			if !strings.Contains(fi.name, s) {
				return false
			}
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

type process func(filenames []*finfo) []*finfo

func ProcessList(res *result, fns []process) {
	for _, fn := range fns {
		res.files = fn(res.files)
	}
}

type sortBy uint8

const (
	byDef sortBy = iota
	byDate
	byName
)

type orderTo uint8

const (
	toDesc orderTo = iota
	toAsc
)

func sortProcess(sorting sortBy, ordering orderTo) process {
	return func(filenames []*finfo) []*finfo {
		fmt.Println("SORTING")
		switch sorting {
		case byDate:
			if ordering == toAsc {
				return countingSortAsc(filenames, lowestTime, highestTime)
			}
			return countingSortDesc(filenames, lowestTime, highestTime)
		case byName:
			if ordering == toAsc {
				sort.Slice(filenames, func(i, j int) bool {
					return naturalAsc(filenames[i].name, filenames[j].name)
				})
			} else {
				sort.Slice(filenames, func(i, j int) bool {
					return naturalDesc(filenames[j].name, filenames[i].name)
				})
			}
		}

		return filenames
	}
}

func sliceProcess(pattern string) process {
	return func(filenames []*finfo) []*finfo {
		fmt.Println("SLICING")
		return sliceArray(pattern, filenames)
	}
}

// sliceArray takes a string pattern and a generic slice, then returns a slice according to the pattern.
func sliceArray[T any](pattern string, input []T) []T {
	// Default slice indices
	start, end := 0, len(input)

	// Remove brackets and split by colon
	trimPattern := strings.Trim(pattern, "[]")
	parts := strings.Split(trimPattern, ":")

	// Parse start index if it exists
	if parts[0] != "" {
		parsedStart, err := strconv.Atoi(parts[0])
		if err == nil {
			start = parsedStart
		}

		// If the start index is negative, adjust it to be relative to the end of the slice
		if parsedStart < 0 {
			start = len(input) + parsedStart
		}
	}

	// Parse end index if it exists
	if len(parts) > 1 && parts[1] != "" {
		parsedEnd, err := strconv.Atoi(parts[1])
		if err == nil {
			end = parsedEnd
		}

		// If the end index is negative, adjust it to be relative to the end of the slice
		if parsedEnd < 0 {
			end = len(input) + parsedEnd
		}
	}

	// Adjust indices to prevent out-of-bounds slicing
	if start < 0 {
		start = 0
	}
	if end > len(input) {
		end = len(input)
	}
	if start > end {
		start = end
	}

	// Return the sliced input
	return input[start:end]
}
