package main

import (
	"errors"
	"fmt"
	"github.com/zetamatta/experimental/argf"
	"os"
	"regexp"
)

func main1() error {
	if len(os.Args) < 2 {
		return errors.New("Usage: grep.exe REGEXP Files...")
	}
	rx, err := regexp.Compile(os.Args[1])
	if err != nil {
		return err
	}
	r := argf.NewFiles(os.Args[2:])
	for r.Scan() {
		if rx.MatchString(r.Text()) {
			fmt.Println(r.Text())
		}
	}
	return r.Err()
}

func main() {
	if err := main1(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
