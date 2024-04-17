package main

import (
	"fmt"
	"log"
	"os"

	"github.com/periaate/common"
	"github.com/periaate/slice"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("No slice expression given\nUsage:\tslice [PATTERN]")
	}
	arg := os.Args[1]

	vals := common.ReadPipe()
	expr := slice.NewExpression[string]()
	expr.Parse(arg)

	res, err := expr.Eval(vals)
	if err != nil {
		log.Fatalln("error during slicing", err)
	}

	for _, s := range res {
		fmt.Println(s)
	}
}
