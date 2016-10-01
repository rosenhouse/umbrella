package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rosenhouse/umbrella/example/lib"
)

func main() {
	words := os.Args[1:]
	uppercase := os.Getenv("UPPERCASE")

	// use if-else here, instead of a bare if
	// so that no single spec can cover both branches
	if uppercase != "" {
		fmt.Println(strings.ToUpper(lib.Concat(words...)))
	} else {
		fmt.Println(lib.Concat(words...))
	}
}
