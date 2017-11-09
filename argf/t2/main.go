package main

import (
	"fmt"
	"os"

	"github.com/zetamatta/experimental/argf"
)

func main() {
	r := argf.New()
	for r.Scan() {
		fmt.Println(r.Text())
	}
	if err := r.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
