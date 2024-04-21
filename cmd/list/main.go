package main

import (
	"os"

	"github.com/periaate/clf"
	"github.com/periaate/list"
)

func main() {
	rest := clf.Run(os.Args[1:], fl)

	list.PrintWithBuf(list.Parse(rest))
}
