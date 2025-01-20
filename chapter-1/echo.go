package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("asdf")
	// variables initialize to empty by default
	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Printf("Echoing: %v\n", s)

	s_1, sep_1 := "", ""
	for _, arg := range os.Args[1:] {
		s_1 += sep_1 + arg
		sep_1 = " "
	}
	fmt.Printf("Echoing Twice: %v\n", s_1)
}
