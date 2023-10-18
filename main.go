package main

import (
	_ "embed"
	"fmt"
	"io/fs"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	gf "github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v3"
)

//go:embed patterns.yml
var patternsFile []byte

var conf patterns // Configuration variable.
var opts Options  // Flag options.

var (
	inclusionMap = map[string]bool{}
	exclusionMap = map[string]bool{}

	highestTime int64                 // Highest found unix timestamp.
	lowestTime  int64 = math.MaxInt64 // Lowest found unix timestamp. Uses MaxInt64 to make comparisons work.

	fp = "."
)

// SortableFiles is an array of sortableFiles which is sortable.
type SortableFiles []sortableFile

func (s SortableFiles) Len() int           { return len(s) }
func (s SortableFiles) Less(i, j int) bool { return s[i].sortable < s[j].sortable }
func (s SortableFiles) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type sortableFile struct {
	Fp       string // Filepath
	sortable string
	Value    int64 // Countable value to be used by counting sort. Populated by unix timestamp.
}

// Todo: Absolute path flag, Custom pattern files
type Options struct {
	Recurse            bool `short:"r" long:"recurse" description:"Recursively list files in subdirectories"`
	ExclusiveRecursion bool `short:"x" long:"xrecurse" description:"Exclusively list files in subdirectories"`
	Ascending          bool `short:"A" long:"ascending" description:"Results will be ordered in ascending order. Files are ordered into descending order by default."`
	Date               bool `short:"d" long:"date" description:"Results will be ordered by their modified time. Files are ordered by filename by default"`

	Include string `short:"i" long:"include" description:"Given an existing extension pattern configuration target, will include only items fitting the pattern. Use ',' to define multiple patterns."`
	Exclude string `short:"e" long:"exclude" description:"Given an existing extension pattern configuration target, will excldue items fitting the pattern. Use ',' to define multiple patterns."`
}

type patterns struct {
	Extensions map[string][]string `yaml:"extensions"`
}

func main() {
	args, err := gf.Parse(&opts)
	if err != nil {
		log.Fatalln("Error parsing flags:", err)
	}

	// Get filepath
	if len(args) > 0 {
		fp = args[0]
	}

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

	files := SortableFiles{}
	getFiles(&files)

	if opts.Date {
		files = countingSort(files)
	} else {
		sort.Sort(files)
	}

	printResults(files)
}

// getFiles attempts to populate the files array using the existing configurations.
func getFiles(files *SortableFiles) {
	// perEntry is ran on each file to construct sortableFile and check if it matches any patterns.
	perEntry := func(pre, base string, d fs.DirEntry) {
		fp := filepath.ToSlash(path.Join(pre, base))
		file := sortableFile{
			Fp:       fp,
			sortable: strings.ToLower(fp),
		}

		if opts.Date {
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
	if opts.ExclusiveRecursion {
		fmt.Println("hello")
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

	} else if !opts.Recurse {
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

// printResults prints the results to stdout.
func printResults(files SortableFiles) {
	for i := range files {
		k := i
		if opts.Ascending {
			k = len(files) - 1 - i
		}
		fmt.Println(files[k].Fp)
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

// countingSort sorts an array into descending order.
//
// Counting sort needs to be called with an array, an integer k which is
// the largest value in the array, and a function which takes an element
// of the array as argument, and returns its value in the range [0, k].
func countingSort(input []sortableFile) []sortableFile {
	k := int(highestTime - lowestTime)
	count := make([]int, k+1)

	// Count occurrences of each value.
	for _, v := range input {
		count[v.Value-lowestTime]++
	}

	// Build and apply offset by summing the counts of previous values.
	for i := 1; i <= k; i++ {
		count[i] += count[i-1]
	}

	result := make([]sortableFile, len(input))
	for _, v := range input {
		result[len(input)-count[v.Value-lowestTime]] = v
		count[v.Value-lowestTime]--
	}

	return result
}

// parsePatterns checks for flags and builds pattern maps accordingly which will be used to
// include or exclude files as per the pattern.
func parsePatterns() {
	if opts.Include != "" {
		// Split flag arguments by comma for multiple items.
		incstr := strings.Split(opts.Include, ",")

		for _, key := range incstr {
			// Find matching extensions and add each of their elements to the inclusion map.
			if v, ok := conf.Extensions[key]; ok {
				for _, val := range v {
					inclusionMap[val] = ok
				}
			}
		}
	}

	if opts.Exclude != "" {
		// Split flag arguments by comma for multiple items.
		excstr := strings.Split(opts.Exclude, ",")

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
