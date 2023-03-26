package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	var recurse bool
	var fp string
	switch len(os.Args) {
	case 0:
		panic("no arguments")
	case 1:
		fp = "."
	case 2:
		recurse = false
	default:
		recurse = os.Args[2] == "-r" || os.Args[2] == "--recurse" || os.Args[2] == "-R"
	}

	if fp == "" {
		fp = os.Args[1]
	}
	fp = fullpath(fp)

	stat, err := os.Stat(fp)
	if err != nil {
		fmt.Println("No such file or directory", err)
		os.Exit(1)
	}

	if !stat.IsDir() {
		fmt.Println("Not a directory")
		os.Exit(1)
	}

	if recurse {
		filepath.WalkDir(fp, func(path string, d fs.DirEntry, err error) error {
			fmt.Println(fullpath(path))
			return nil
		})
		return
	}

	files, err := os.ReadDir(fp)
	if err != nil {
		fmt.Println("No such file or directory", err)
		os.Exit(1)
	}

	for _, file := range files {
		fmt.Println(fullpath(file.Name()))
	}

}

func fullpath(str string) string {
	if str[0] == '/' {
		return str
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.ToSlash(filepath.Join(wd, str))
}
