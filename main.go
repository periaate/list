// Copyright © 2023 Daniel S. (GitHub: periaate)
// All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package main provides a simple command-line utility for listing the files in
// a directory.
//
// Usage:
//
//	program [directory] [-r|--recurse|-R]
//
// Arguments:
//
//	directory     Path to the directory to list (default is current directory)
//
// Flags:
//
//	-r, --recurse   Recurse into subdirectories
//	-R              Recurse into subdirectories
//
// The program can be used to quickly and easily view the files in a directory
// or to generate a list of file paths for use in other programs or scripts. By
// default, the program lists only the files in the specified directory. If the -r
// or --recurse flag is provided, the program will also list files in all
// subdirectories of the specified directory.
//
// To use the program, you can either compile it using the "go build" command
// and then run it from the command line, or you can run it directly using the "go
// run" command. For example:
//
//	# Compile the program
//	$ go build
//
//	# List the files in the current directory
//	$ ./program
//
//	# List the files in a specific directory
//	$ ./program /path/to/directory
//
//	# Recurse into subdirectories and list all files
//	$ ./program /path/to/directory -r
//
//	# Run the program directly with "go run"
//	$ go run program.go [directory] [-r|--recurse|-R]
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	var recurse bool
	fp := "."

	if len(os.Args) > 1 {
		fp = os.Args[1]
	}
	if len(os.Args) > 2 {
		recurse = os.Args[2] == "-r" || os.Args[2] == "--recurse" || os.Args[2] == "-R"
	}

	stat, err := os.Stat(fp)
	if err != nil || !stat.IsDir() {
		fmt.Println("Invalid directory")
		os.Exit(1)
	}

	if recurse {
		filepath.WalkDir(fp, func(path string, d fs.DirEntry, err error) error {
			fmt.Println(fullpath(path))
			return nil
		})
		return
	}

	files, _ := os.ReadDir(fp)
	for _, file := range files {
		fmt.Println(fullpath(file.Name()))
	}
}

// fullpath is a function that takes a string input str representing a file or
// directory path and returns an absolute path string. It ensures that the path is
// represented with forward slashes ("/") regardless of the operating system the
// code is executed on. The returned path is derived based on various conditions
// described in great detail below, covering all possible branches of the code
// execution.
//
// When the input str is an empty string or its first character is a forward
// slash ("/"), the function performs the following steps:
//
// 1. If str starts with a forward slash ("/"), the function assumes that str
// is already an absolute path and directly returns the input str as-is without
// any further modification.
//
// When the input str does not start with a forward slash ("/"), the function
// performs the following steps:
//
// 2. The function calls the os.Getwd() function to retrieve the current
// working directory. os.Getwd() can return two possible outcomes: the current
// working directory wd as a string, or an error err.
//
// a. If os.Getwd() returns an error err, the function will immediately panic
// and terminate the program execution with the error message contained in err.
// This error may occur if, for example, the current working directory is deleted
// or its permissions are changed during the function execution, making it
// inaccessible to the program.
//
// b. If os.Getwd() returns the current working directory wd, the function
// proceeds to the next step.
//
// 3. The function calls the filepath.Join() function, passing the wd and str
// as arguments. filepath.Join() is responsible for concatenating the two input
// strings in a manner that is compatible with the operating system's path
// structure. It takes care of adding or removing any necessary path separators
// (e.g., slashes) between the wd and str arguments. The concatenated path is then
// returned as a new string.
//
// 4. Finally, the function calls the filepath.ToSlash() function, passing the
// concatenated path string from the previous step as an argument.
// filepath.ToSlash() converts the path string into a canonical form that uses
// forward slashes ("/") as path separators, regardless of the operating system.
// This canonical path string is then returned as the final output of the fullpath
// function.
//
// Note: The fullpath function does not validate the existence or accessibility
// of the input str or the final output path. It is the responsibility of the
// caller to ensure that the input str and the returned path are valid and
// accessible according to the specific requirements of their program or use case.
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
