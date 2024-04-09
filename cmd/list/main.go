package main

import (
	"os"

	"github.com/periaate/list"
	"github.com/periaate/list/internal"
)

func main() {
	opts := list.Parse(os.Args[1:])

	pipedValues := internal.ReadPipe()
	if len(pipedValues) != 0 {
		opts.Args = append(opts.Args, pipedValues...)
	}
	if len(opts.Args) == 0 {
		opts.Args = append(opts.Args, "./")
	}

	res := list.Run(opts)
	list.PrintWithBuf(res.Files, opts)
}
