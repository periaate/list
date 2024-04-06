package main

import (
	"bufio"
	"os"

	"github.com/periaate/list/cfg"

	"github.com/periaate/list"
)

func main() {
	opts := cfg.Parse(os.Args[1:])

	pipedValues := readPipe()
	if len(pipedValues) != 0 {
		opts.Args = append(opts.Args, pipedValues...)
	}
	if len(opts.Args) == 0 {
		opts.Args = append(opts.Args, "./")
	}

	res := list.Run(opts)
	list.PrintWithBuf(res.Files, opts)
}

func readPipe() (res []string) {
	fileInfo, _ := os.Stdin.Stat()
	if (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			res = append(res, scanner.Text())
		}
	}
	return
}
