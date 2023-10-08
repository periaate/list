package main

import (
	"fmt"
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	gf "github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v3"
)

const configFile = "config.yml"

var conf extcfg
var opts Options

type Options struct {
	Recurse bool `short:"r" long:"recurse" description:"Recursively list files in subdirectories"`

	Ascending  bool `short:"A" long:"ascending" description:"Results will be ordered in ascending order."`
	Descending bool `short:"D" long:"descending" description:"Results will be ordered in descending order (default)."`

	Date bool `short:"d" long:"date" description:"Results will be ordered by their modified time."`
	Name bool `short:"n" long:"name" description:"Results will be ordered by their names (default)."`

	Include string `short:"i" long:"include" description:"Given an existing extension pattern configuration target, will include only items fitting the pattern. Use ',' to define multiple patterns."`
	Exclude string `short:"e" long:"exclude" description:"Given an existing extension pattern configuration target, will excldue items fitting the pattern. Use ',' to define multiple patterns."`
}

type extcfg struct {
	Extensions map[string][]string `yaml:"extensions"`
}

func main() {
	_, err := gf.Parse(&opts)
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		os.Exit(1)
	}

	fp := "."
	if len(os.Args) > 1 && os.Args[1][0] != '-' {
		fp = os.Args[1]
	}

	stat, err := os.Stat(fp)
	if err != nil || !stat.IsDir() {
		fmt.Println("Invalid directory")
		os.Exit(1)
	}

	inc := map[string]bool{}
	exc := map[string]bool{}

	if opts.Include != "" {

		incstr := strings.Split(opts.Include, ",")

		for _, key := range incstr {
			if v, ok := conf.Extensions[key]; ok {
				for _, val := range v {
					inc[val] = ok
				}
			}
		}
	}

	if opts.Exclude != "" {
		excstr := strings.Split(opts.Exclude, ",")
		for _, key := range excstr {
			if v, ok := conf.Extensions[key]; ok {
				for _, val := range v {
					exc[val] = ok
				}
			}
		}
	}

	ress := ByFp{}
	var L int64
	var S int64 = math.MaxInt64

	if !opts.Recurse {
		res, err := os.ReadDir(fp)
		if err != nil {
			log.Fatalln(err)
		}
		for _, e := range res {
			inf, err := e.Info()
			if err != nil {
				continue
			}
			dt := inf.ModTime().Unix()
			if dt > L {
				L = dt
			}
			if dt < S {
				S = dt
			}
			ress = append(ress, sortableFile{
				Fp:  fmt.Sprintf("%s/%s", fp, e.Name()),
				Val: dt,
			})
		}
	} else {
		vfs := os.DirFS(fp)
		err := fs.WalkDir(vfs, fp, func(path string, d fs.DirEntry, _ error) error {
			if d.IsDir() {
				return nil
			}
			inf, err := d.Info()
			if err != nil {
				return nil
			}
			dt := inf.ModTime().Unix()
			if dt > L {
				L = dt
			}
			if dt < S {
				S = dt
			}
			ress = append(ress, sortableFile{
				Fp:  fmt.Sprintf("%s/%s", path, d.Name()),
				Val: dt,
			})
			return nil
		})
		if err != nil {
			log.Fatalln(err)
		}
	}

	files := ByFp{}
	for _, e := range ress {
		ext := filepath.Ext(e.Fp)
		incl := true
		if len(inc) != 0 {
			if _, ok := inc[ext]; !ok {
				incl = false
			}
		}
		if len(exc) != 0 {
			if _, ok := exc[ext]; ok {
				incl = false
			}
		}

		if incl {
			e.Val -= S
			files = append(files, e)
		}
	}

	if opts.Date {
		files = CountingSort(files, int(L-S))
	} else {
		sort.Sort(files)
	}

	if opts.Ascending {
		for i := range files {
			f := files[len(files)-1-i]
			fmt.Println(f.Fp)
		}
	} else {
		for _, f := range files {
			fmt.Println(f.Fp)
		}
	}

}

type ByFp []sortableFile

func (s ByFp) Len() int {
	return len(s)
}

func (s ByFp) Less(i, j int) bool {
	return s[i].Fp < s[j].Fp
}

func (s ByFp) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type sortableFile struct {
	Fp  string
	Val int64
}

func init() {
	yml, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yml, &conf)
	if err != nil {
		log.Fatal(err)
	}
}

// CountingSort sorts an array into descending order.
//
// Counting sort needs to be called with an array, an integer k which is
// the largest value in the array, and a function which takes an element
// of the array as argument, and returns its value in the range [0, k].
func CountingSort(input []sortableFile, k int) []sortableFile {
	count := make([]int, k+1)

	// Count occurrences of each value.
	for _, v := range input {
		count[v.Val]++
	}

	// Build and apply offset by summing the counts of previous values.
	for i := 1; i <= k; i++ {
		count[i] += count[i-1]
	}

	result := make([]sortableFile, len(input))
	for _, v := range input {
		result[len(input)-count[v.Val]] = v
		count[v.Val]--
	}

	return result
}
