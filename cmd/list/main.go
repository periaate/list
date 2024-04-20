package main

import (
	"os"

	"github.com/periaate/list"
)

func main() { list.Parse(os.Args[1:]).Run() }
