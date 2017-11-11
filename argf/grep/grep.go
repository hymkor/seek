package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/mattn/go-colorable"
	"github.com/zetamatta/experimental/argf"
)

func main1() error {
	if len(os.Args) < 2 {
		return errors.New("Usage: grep.exe REGEXP Files...")
	}
	rx, err := regexp.Compile(os.Args[1])
	if err != nil {
		return err
	}
	out := colorable.NewColorableStdout()
	r := argf.NewFiles(os.Args[2:])
	for r.Scan() {
		text := r.Text()
		m := rx.FindAllStringIndex(text, -1)
		if m != nil {
			fmt.Fprintf(out, "\x1B[35;1m%s:\x1B[32;1m%d\x1B[36;1m:\x1B[37;1m", r.Filename(), r.FNR())
			last := 0
			for i := 0; i < len(m); i++ {
				fmt.Fprintf(out, "%s\x1B[31;1m%s\x1B[37;1m",
					text[last:m[i][0]],
					text[m[i][0]:m[i][1]])
				last = m[i][1]
			}
			fmt.Fprintln(out, text[last:])
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
