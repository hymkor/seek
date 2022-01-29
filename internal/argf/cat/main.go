package main

import (
	"fmt"
	"os"

	"github.com/zetamatta/seek/internal/argf"
)

func main() {
	r := argf.New()
	r.OnError = func(err error) error {
		fmt.Fprintf(os.Stderr, err.Error())
		return nil
	}
	for r.Scan() {
		fmt.Println(r.Text())
	}
	if err := r.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
