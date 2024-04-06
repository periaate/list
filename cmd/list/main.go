package main

import (
	"bufio"
	"os"

	"github.com/periaate/list/cfg"

	"github.com/periaate/list"
)

func main() {
	Opts := cfg.Parse(os.Args[1:])

	pipedValues := readPipe()
	if len(pipedValues) != 0 {
		cfg.Args = append(cfg.Args, pipedValues...)
	}
	if len(cfg.Args) == 0 {
		cfg.Args = append(cfg.Args, "./")
	}

	res := list.Run(Opts)
	list.PrintWithBuf(res.Files)
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
