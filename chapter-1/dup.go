package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		// Treat input as file
		countLines(os.Stdin, counts)
	} else {
		// Open files if they are specified
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				// Error opening file
				fmt.Fprintf(os.Stderr, "dup: %v\n", err)
			}
			countLines(f, counts)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

func countLines(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
	input.Err()
	for input.Scan() {
		counts[input.Text()]++
	}
}
