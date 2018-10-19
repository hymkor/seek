package main

import (
	"fmt"
	"os"

	"github.com/zetamatta/seek/starstar"
)

func main() {
	for _, arg1 := range os.Args[1:] {
		files, err := starstar.Expand(arg1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", arg1, err)
			continue
		}
		for _, file1 := range files {
			fmt.Println(file1)
		}
	}
}
