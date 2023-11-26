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
	"strconv"
	"strings"
	"unicode"

	gf "github.com/jessevdk/go-flags"
	"github.com/texttheater/golang-levenshtein/levenshtein"
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

const (
	// using parts to make the weights always add up to 1.0 as well as make it easier to see the ratio.
	distanceParts = 1.0
	scoreParts    = 1.0

	distanceWeight = distanceParts / (distanceParts + scoreParts)
	scoreWeight    = scoreParts / (distanceParts + scoreParts)
)

// SortableFiles is an array of sortableFiles which is sortable.
type SortableFiles []sortableFile

func (s SortableFiles) Len() int           { return len(s) }
func (s SortableFiles) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SortableFiles) Less(i, j int) bool { return naturalLess(s[i].sortableName, s[j].sortableName) }

type sortableFile struct {
	Fp           string  // Filepath
	sortableName string  // Filepath in lowercase for sorting
	Value        int64   // Countable value to be used by counting sort. Populated by unix timestamp.
	score        float64 // Score for how well the string matches the query
}

// TODO: Custom pattern files
type Options struct {
	Absolute           bool `short:"A" long:"absolute" description:"Format paths to be absolute. Relative by default."`
	Recurse            bool `short:"r" long:"recurse" description:"Recursively list files in subdirectories"`
	ExclusiveRecursion bool `short:"x" long:"xrecurse" description:"Exclusively list files in subdirectories"`
	Ascending          bool `short:"a" long:"ascending" description:"Results will be ordered in ascending order. Files are ordered into descending order by default."`
	Date               bool `short:"d" long:"date" description:"Results will be ordered by their modified time. Files are ordered by filename by default"`

	Select string `short:"S" long:"select" description:"Selects the item which matches the query the best."`
	Query  string `short:"Q" long:"query" description:"Returns items ordered by their similarity to the query."`

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

	switch {
	case opts.Select != "":
		printSelectResult(files)
	case opts.Query != "":
		sortByScore(files)
		fallthrough
	default:
		printResults(files)
	}
}

// getFiles attempts to populate the files array using the existing configurations.
func getFiles(files *SortableFiles) {
	// perEntry is ran on each file to construct sortableFile and check if it matches any patterns.
	perEntry := func(pre, base string, d fs.DirEntry) {
		fp := filepath.ToSlash(path.Join(pre, base))
		file := sortableFile{
			Fp:           fp,
			sortableName: strings.ToLower(fp),
		}

		if opts.Query != "" || opts.Select != "" {
			file.score = levenshtein.RatioForStrings([]rune(fp), []rune(opts.Query), levenshtein.DefaultOptions)*distanceWeight + calculateScore(fp, opts.Include)*scoreWeight
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

func sortByScore(files SortableFiles) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].score > files[j].score
	})
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

		printResult(files[k].Fp)
	}
}

func printSelectResult(files SortableFiles) {
	bestFile := files[0]

	for i := 1; i < len(files); i++ {
		file := files[i]
		if file.score > bestFile.score {
			bestFile = file
		}
	}

	printResult(bestFile.Fp)
}

func printResult(fp string) {
	if opts.Absolute {
		fp, _ = filepath.Abs(fp)
		fp = filepath.ToSlash(fp)
	}
	fmt.Println(fp)
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

// naturalLess compares two strings and returns true if a < b in natural order.
func naturalLess(a, b string) bool {
	var ai, bi int
	for ai < len(a) && bi < len(b) {
		ach, bch := rune(a[ai]), rune(b[bi])
		if unicode.IsDigit(ach) && unicode.IsDigit(bch) {
			var anum, bnum string
			for ; ai < len(a) && unicode.IsDigit(rune(a[ai])); ai++ {
				anum += string(a[ai])
			}
			for ; bi < len(b) && unicode.IsDigit(rune(b[bi])); bi++ {
				bnum += string(b[bi])
			}
			an, _ := strconv.Atoi(anum)
			bn, _ := strconv.Atoi(bnum)
			if an != bn {
				return an < bn
			}
		} else {
			if ach != bch {
				return ach < bch
			}
			ai++
			bi++
		}
	}
	return len(a) < len(b)
}

// Calculates a score for how well the string matches the query using subsequence matching
func calculateScore(str, query string) (score float64) {
	strIndex, queryIndex := 0, 0

	for strIndex < len(str) && queryIndex < len(query) {
		if str[strIndex] == query[queryIndex] {
			score += 1   // Increment score for each matching character
			queryIndex++ // Move to the next character in the query
		}
		strIndex++ // Always move to the next character in the string
	}

	if queryIndex == len(query) {
		return score / float64(len(str)) // Return the score divided by the length of the string
	}
	return 0 // Return 0 if the query is not a subsequence of the string
}
