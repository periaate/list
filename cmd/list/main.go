package main

import (
	"os"

	"github.com/periaate/list"
)

func main() {
	rest := list.Program.EvalOnly(os.Args[1:], []string{"quiet", "help", "debug"})
	if list.Program.PrintedHelp() {
		return
	}

	list.PrintWithBuf(list.Parse(rest))
}
