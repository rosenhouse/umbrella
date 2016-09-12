package main

import (
	"fmt"
	"os"

	"github.com/rosenhouse/umbrella/example/lib"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("usage: %s [some] [things] [to] [concat]\n", os.Args[0])
		os.Exit(1)
	}

	words := os.Args[1:]
	concatenated := lib.Concat(words...)

	fmt.Println(concatenated)
}
