package inf

import (
	"log"
	"os"

	gf "github.com/jessevdk/go-flags"
)

var Opts Options

var Args []string

// TODO: Custom pattern files
type Options struct {
	Absolute           bool `short:"A" long:"absolute" description:"Format paths to be absolute. Relative by default."`
	Recurse            bool `short:"r" long:"recurse" description:"Recursively list files in subdirectories"`
	ExclusiveRecursion bool `short:"x" long:"xrecurse" description:"Exclusively list files in subdirectories"`
	Invert             bool `short:"a" long:"invert" description:"Results will be ordered in inverted order. Files are ordered into descending order by default."`
	Date               bool `short:"d" long:"date" description:"Results will be ordered by their modified time. Files are ordered by filename by default"`

	Combine bool `short:"c" long:"combine" description:"If given multiple paths, will combine the results into one list."`

	Query    string  `short:"q" long:"query" description:"Returns items ordered by their similarity to the query."`
	QueryAll string  `short:"Q" long:"queryAll" description:"Takes comma separated queries and evaluates each query individually. Works with --combine."`
	Ngram    int     `short:"n" long:"ngram" description:"Defines the n-gram size for the query algorithm. Default is 3."`
	Top      int     `short:"t" long:"top" description:"Returns first n items."`
	Prune    float64 `short:"p" long:"prune" description:"Prunes items with a score lower than the given value. -1 to prune all items without a score. 0.0 is the default and will not prune any items."`
	Score    bool    `short:"s" long:"score" description:"Returns items with their score. Intended for debugging purposes."`

	Include string `short:"i" long:"include" description:"Given an existing extension pattern configuration target, will include only items fitting the pattern. Use ',' to define multiple patterns."`
	Exclude string `short:"e" long:"exclude" description:"Given an existing extension pattern configuration target, will exclude items fitting the pattern. Use ',' to define multiple patterns."`
}

func init() {
	Opts = Options{}
	args, err := gf.Parse(&Opts)
	if err != nil {
		if gf.WroteHelp(err) {
			os.Exit(0)
		}
		log.Fatalln("Error parsing flags:", err)
	}
	Args = args
}
