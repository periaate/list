package main

import (
	"fmt"
	"os"

	"github.com/periaate/common"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	switch os.Args[1] {
	case "pipe":
		for i, r := range common.ReadPipe() {
			fmt.Println(i, r)
		}
	case "args":
		for i, r := range os.Args[2:] {
			fmt.Println(i, r)
		}
	}
}
