package main

import (
	"os"

	"github.com/periaate/common"
	"github.com/periaate/list"
)

func main() {
	opts := list.Parse(os.Args[1:])

	pipedValues := common.ReadPipe()
	if len(pipedValues) != 0 {
		opts.Args = append(opts.Args, pipedValues...)
	}

	res := list.Run(opts)

	if opts.ExecArgs != nil || len(opts.ExecArgs) != 0 {
		list.Exec(res, opts)
		return
	}

	list.PrintWithBuf(res, opts)
}
