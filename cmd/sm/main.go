package main

// func main() {
// 	pipedValues := internal.ReadPipe()
// 	if len(pipedValues) == 0 {
// 		fmt.Println("USAGE: ... | sm [{from}:{to}]")
// 		os.Exit(0)
// 	}
// 	if len(os.Args) == 0 {
// 		fmt.Println("USAGE: ... | sm [{from}:{to}]")
// 		os.Exit(0)
// 	}

// 	for _, v := range pipedValues {
// 		res, _ := generic.Slice(os.Args[1], []rune(v))
// 		fmt.Println(string(res))
// 	}
// }
