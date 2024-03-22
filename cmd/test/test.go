package main

import (
	"fmt"
	"log"
	"os"

	gf "github.com/jessevdk/go-flags"
)

type Opts struct {
	A bool `short:"a"`
}

func main() {
	opts := Opts{}
	rest, err := gf.Parse(&opts)
	if err != nil {
		if gf.WroteHelp(err) {
			os.Exit(0)
		}
		log.Fatalln("Error parsing flags:", err)
	}
	args := rest

	fmt.Println(args)
}
