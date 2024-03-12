package main

import (
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// 3 paths
// walkDir -r recursion
// readDir no recursion
// recursive readDir for -T --toDepth and -F --fromDepth

// dirEntryFns - add additional information to finfo if necessary
// including mod time in finfo -d
// including score in finfo for query -q --query

// filters - filters determine whether or not a given file should be included in the result
// include file types
// exclude file types
// ignore patterns
// strict search

// processes - processes are applied to the result after filtering. They cna manipulate the result in any way.
// fuzzy search with score -q --query
// sort by date -d
// sort by filename -f (natural sorting)
// sort in ascending order -a
// slices -S

// print - printing the result
// absolute paths -A

// onhold
// dir only
// file only

var (
	highestTime int64                 // Highest found unix timestamp.
	lowestTime  int64 = math.MaxInt64 // Lowest found unix timestamp. Uses MaxInt64 to make comparisons work.

	fp = "./"
)

type result struct {
	files []*finfo
}

type finfo struct {
	name string
	path string // includes name, relative path to cwd
	mod  int64  // unix timestamp
}

func List() {
	res := &result{[]*finfo{}}
	filters := collectFilters()
	processes := collectProcess()
	collectResult(buildWalkDirFn(filters, res))

	ProcessList(res, processes)
	printWithBuf(res.files)
}

func collectResult(fn fs.WalkDirFunc) {
	if len(Args) > 0 {
		fp = Args[0]
	}

	stat, err := os.Stat(fp)
	if err != nil || !stat.IsDir() {
		log.Fatalln(err)
	}

	switch {
	case Opts.Recurse:
		err = filepath.WalkDir(fp, fn)
		if err != nil {
			log.Fatalln(err)
		}
	case Opts.ToDepth > 0 || Opts.FromDepth > 0:
		if Opts.ToDepth == 0 {
			Opts.ToDepth = math.MaxInt64
		}
		if Opts.FromDepth == 0 {
			Opts.FromDepth = math.MinInt64
		}
		depthDir(Opts.ToDepth, Opts.FromDepth, fn)
	default:
		entries, err := os.ReadDir(fp)
		if err != nil {
			log.Fatalln(err)
		}
		for _, entry := range entries {
			filep := filepath.Join(fp, entry.Name())
			err = fn(filep, entry, nil)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

}

var (
	to   int
	from int
	rfn  fs.WalkDirFunc
)

func depthDir(To, From int, fn fs.WalkDirFunc) {
	to = To
	from = From
	rfn = fn
	dirFn(fp, 0)
}
func dirFn(path string, depth int) {
	if depth < from || depth > to {
		return
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	for _, entry := range entries {
		filep := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			dirFn(filep, depth+1)
		}
		err = rfn(filep, entry, nil)
		if err != nil {
			log.Fatalln(err)
		}
	}
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

func reverse(filenames []*finfo) []*finfo {
	for i := 0; i < len(filenames)/2; i++ {
		j := len(filenames) - i - 1
		filenames[i], filenames[j] = filenames[j], filenames[i]
	}
	return filenames
}

func collectProcess() []process {
	var fns []process

	switch {
	case len(Opts.Query) > 0:
		fns = append(fns, queryProcess)
		if Opts.Ascending {
			fns = append(fns, reverse)
		}
	case Opts.Ascending, Opts.Date:
		sorting := byName
		order := toDesc
		if Opts.Ascending {
			order = toAsc
		}
		if Opts.Date {
			sorting = byDate
		}
		fns = append(fns, sortProcess(sorting, order))
	}

	// 4 is the minimum length of a slice pattern, as [:n] or [n:] are the smallest possible patterns.
	if len(Opts.Slice) >= 4 {
		fns = append(fns, sliceProcess(Opts.Slice))
	}
	return fns
}

type filter func(*finfo, fs.DirEntry) bool

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
	fi.mod = unixTime
	return true
}

func buildWalkDirFn(fns []filter, res *result) func(string, fs.DirEntry, error) error {
	return func(path string, d fs.DirEntry, err error) error {
		if d == nil || err != nil {
			return nil
		}
		fi := &finfo{name: d.Name(), path: path}
		for _, fn := range fns {
			res := fn(fi, d)
			if !res {
				return nil
			}
		}
		res.files = append(res.files, fi)
		return nil
	}
}

func filterList(include []contentType, exclude []contentType, ignore []string, search []string) filter {
	return func(fi *finfo, _ fs.DirEntry) bool {
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
