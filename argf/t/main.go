package main

import (
	"fmt"
	"github.com/zetamatta/experimental/argf"
)

func main() {
	argf.Read_(func(line string) error {
		fmt.Println(line)
		return nil
	})
}
