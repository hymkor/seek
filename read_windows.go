package main

import (
	"github.com/zetamatta/go-mbcs"
	"unicode/utf8"
)

func readline(sc scanner) string {
	line := sc.Bytes()

	var text string
	if utf8.Valid(line) {
		text = string(line)
	} else {
		var err error
		text, err = mbcs.AtoU(line)
		if err != nil {
			text = err.Error()
		}
	}
	return text
}
