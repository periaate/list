package main

import (
	_ "embed"
	"fmt"
	"io/fs"
	"list/inf"
	"list/sorting"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// TODO:
// Symlinks flag
// Archive support
// Proper slicing (start, end)
// Custom pattern files
// Pass in patterns as string
// Recurse depth
// Exlusive recurse depth (inverse depth)
// Support piping in paths

//go:embed patterns.yml
var patternsFile []byte

var conf patterns // Configuration variable.

var (
	inclusionMap = map[string]bool{}
	exclusionMap = map[string]bool{}

	highestTime int64                 // Highest found unix timestamp.
	lowestTime  int64 = math.MaxInt64 // Lowest found unix timestamp. Uses MaxInt64 to make comparisons work.

	fp = "."
)

type patterns struct {
	Extensions map[string][]string `yaml:"extensions"`
}

func main() {

	if inf.Opts.Query != "" {
		if inf.Opts.Ngram != 0 {
			sorting.N = inf.Opts.Ngram
		}
		sorting.QueryGram = sorting.GenNgram(inf.Opts.Query, sorting.N)
		sorting.Query = strings.ToLower(inf.Opts.Query)
	}

	files := sorting.SortableFiles{}

	// If no arguments are given, use the current directory.
	if len(inf.Args) == 0 {
		inf.Args = append(inf.Args, fp)
	}

	for _, fp = range inf.Args {
		// Validate filepath
		stat, err := os.Stat(fp)
		if err != nil || !stat.IsDir() {
			log.Fatalln(err)
		}

		err = yaml.Unmarshal(patternsFile, &conf)
		if err != nil {
			log.Fatalln(err)
		}
		parsePatterns()

		getFiles(&files)

		if inf.Opts.Date {
			files = sorting.CountingSort(files, lowestTime, highestTime)
		}

		if !inf.Opts.Combine {
			printResults(files)
			files = sorting.SortableFiles{}
		}
	}

	if !inf.Opts.Combine {
		return
	}

	printResults(files)
}

func printResults(files sorting.SortableFiles) {
	switch {
	case inf.Opts.Query != "":
		sorting.SortByScore(files)
		if inf.Opts.Prune != 0.0 {
			if inf.Opts.Prune == -1.0 {
				inf.Opts.Prune = 0.0
			}
			files = sorting.Prune(files, inf.Opts.Prune)
		}
	default:
		sort.Sort(files)
	}
	defaulPrint(files)
}

func defaulPrint(files sorting.SortableFiles) {
	top := -1
	if inf.Opts.Top > 0 && inf.Opts.Top < len(files) {
		top = inf.Opts.Top
	}

	if top != -1 {
		if inf.Opts.Invert {
			top = len(files) - top
			files = files[top:]
		} else {
			files = files[:top]
		}
	}

	for i := range files {
		k := i
		if inf.Opts.Invert {
			k = len(files) - 1 - i
		}

		printResult(files[k])
	}
}

func printResult(sf *sorting.SortableFile) {
	fp := filepath.ToSlash(sf.Fp)
	if inf.Opts.Absolute {
		fp, _ = filepath.Abs(sf.Fp)
		fp = filepath.ToSlash(fp)
	}

	if inf.Opts.Score {
		fmt.Printf("%f\t%s\n", sf.Score, sf.Fp)
		return
	}
	fmt.Println(fp)
}

// getFiles attempts to populate the files array using the existing configurations.
func getFiles(files *sorting.SortableFiles) {
	// perEntry is ran on each file to construct sorting.SortableFile and check if it matches any patterns.
	perEntry := func(pre, base string, d fs.DirEntry) {
		fp := filepath.ToSlash(path.Join(pre, base))
		file := &sorting.SortableFile{
			Fp:           fp,
			SortableName: strings.ToLower(fp),
		}

		if inf.Opts.Query != "" {
			// file.Score = sorting.CalculateScore(file.SortableName, inf.Opts.Query)
			// file.Ngram = sorting.GenNgram(file.SortableName)
			file.Score = sorting.CalculateMatchScore(file.SortableName, sorting.N)
		}

		if inf.Opts.Date {
			inf, err := d.Info()
			if err != nil {
				return
			}
			unixTime := inf.ModTime().Unix()
			updateTime(unixTime)
			file.Value = unixTime
		}

		// Check files extension and match against patterns.
		ext := filepath.Ext(file.Fp)
		incl := true
		if len(inclusionMap) != 0 {
			if _, ok := inclusionMap[ext]; !ok {
				incl = false
			}
		}
		if len(exclusionMap) != 0 {
			if _, ok := exclusionMap[ext]; ok {
				incl = false
			}
		}

		if incl {
			*files = append(*files, file)
		}
	}

	// Get files
	if inf.Opts.ExclusiveRecursion {
		res, err := os.ReadDir(fp)
		if err != nil {
			log.Fatalln(err)
		}
		dirs := []fs.DirEntry{}
		for _, d := range res {
			if d.IsDir() {
				dirs = append(dirs, d)
			}
		}
		fn := func(pre string) fs.WalkDirFunc {
			return func(path string, d fs.DirEntry, err error) error {
				if d == nil || err != nil {
					return err
				}
				perEntry(pre, path, d)
				return nil
			}
		}

		for _, dir := range dirs {
			recurse(filepath.Join(fp, dir.Name()), fn)
		}

	} else if !inf.Opts.Recurse {
		res, err := os.ReadDir(fp)
		if err != nil {
			log.Fatalln(err)
		}
		for _, d := range res {
			perEntry(fp, d.Name(), d)
		}
	} else {
		fn := func(pre string) fs.WalkDirFunc {
			return func(path string, d fs.DirEntry, err error) error {
				if d == nil || err != nil {
					return err
				}
				perEntry(pre, path, d)
				return nil
			}
		}

		recurse(fp, fn)
	}
}

func recurse(fp string, fn func(string) fs.WalkDirFunc) {
	absfp, err := filepath.Abs(fp)
	if err != nil {
		log.Fatalln(err)
	}
	pre := fp

	err = fs.WalkDir(os.DirFS(absfp), ".", fn(pre))
	if err != nil {
		log.Fatalln(err)
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

// parsePatterns checks for flags and builds pattern maps accordingly which will be used to
// include or exclude files as per the pattern.
func parsePatterns() {
	if inf.Opts.Include != "" {
		// Split flag arguments by comma for multiple items.
		incstr := strings.Split(inf.Opts.Include, ",")

		for _, key := range incstr {
			// Find matching extensions and add each of their elements to the inclusion map.
			if v, ok := conf.Extensions[key]; ok {
				for _, val := range v {
					inclusionMap[val] = ok
				}
			}
		}
	}

	if inf.Opts.Exclude != "" {
		// Split flag arguments by comma for multiple items.
		excstr := strings.Split(inf.Opts.Exclude, ",")

		for _, key := range excstr {
			// Find matching extensions and add each of their elements to the exclusion map.
			if v, ok := conf.Extensions[key]; ok {
				for _, val := range v {
					exclusionMap[val] = ok
				}
			}
		}
	}
}
