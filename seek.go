package main

import (
	"errors"
	"flag"
	"fmt"
	"html"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/mattn/go-zglob"

	"github.com/zetamatta/seek/argf"
)

const (
	UTF8BOM = "\xEF\xBB\xBF"

	MAGENTA = "\x1B[35;1m"
	GREEN   = "\x1B[32;1m"
	AQUA    = "\x1B[36;1m"
	WHITE   = "\x1B[37;1m"
	RED     = "\x1B[31;1m"

	RESET = "\x1B[0m"
)

type FlagStrings []string

func (f *FlagStrings) String() string {
	return strings.Join([]string(*f), ",")
}

func (f *FlagStrings) Set(value string) error {
	*f = append(*f, value)
	return nil
}

var multiPatterns FlagStrings

var ignoreCase = flag.Bool("i", false, "ignore case")
var recursive = flag.Bool("r", false, "recursive")
var outputHtml = flag.Bool("html", false, "output html")
var flagBefore = flag.Int("B", 0, "print N lines before matching lines")
var flagAfter = flag.Int("A", 0, "print N lines after matching lines")
var flagNoColor = flag.Bool("no-color", false, "no color")

func makeRegularExpression(pattern string) (*regexp.Regexp, error) {
	if *ignoreCase {
		pattern = "(?i)" + pattern
	}

	pattern = strings.Replace(pattern, `\<`, `\b`, -1)
	pattern = strings.Replace(pattern, `\>`, `\b`, -1)

	return regexp.Compile(pattern)
}

func main1() error {
	flag.Var(&multiPatterns, "m", "multi regular expression")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags...] REGEXP Files...\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, " Files: **/*.go ... find *.go files recursively.\n")
		return nil
	}

	var rxs []*regexp.Regexp
	if multiPatterns == nil || len(multiPatterns) <= 0 {
		rx, err := makeRegularExpression(args[0])
		if err != nil {
			return err
		}
		rxs = []*regexp.Regexp{rx}
		args = args[1:]
	} else {
		for _, pattern := range multiPatterns {
			rx, err := makeRegularExpression(pattern)
			if err != nil {
				return err
			}
			rxs = append(rxs, rx)
		}
	}
	var output func(fname string, line int, text string, m [][]int)
	if *outputHtml {
		fmt.Println("<html><body style=\"background-color:lightgray\">")
		output = func(fname string, line int, text string, m [][]int) {
			fmt.Printf(`<div><span style="color:magenta">%s</span>:<span style="color:green">%d</span><span style="color:aqua">:</span>`,
				html.EscapeString(fname), line)
			last := 0
			for i := 0; i < len(m); i++ {
				fmt.Printf(`%s<span style="color:red">%s</span>`,
					html.EscapeString(text[last:m[i][0]]),
					html.EscapeString(text[m[i][0]:m[i][1]]))
				last = m[i][1]
			}
			fmt.Printf("%s</div>\n",
				html.EscapeString(text[last:]))
		}
		defer func() {
			fmt.Println("</body></html>")
		}()
	} else {
		var out io.Writer
		if *flagNoColor {
			out = colorable.NewNonColorable(os.Stdout)
		} else if isatty.IsTerminal(os.Stdout.Fd()) {
			out = colorable.NewColorableStdout()
		} else {
			out = colorable.NewNonColorable(os.Stdout)
		}
		needReset := false
		output = func(fname string, line int, text string, m [][]int) {
			fmt.Fprintf(out, MAGENTA+"%s"+WHITE+":"+GREEN+"%d"+AQUA+":"+WHITE, fname, line)
			last := 0
			for i := 0; i < len(m); i++ {
				fmt.Fprintf(out, "%s"+RED+"%s"+WHITE,
					text[last:m[i][0]],
					text[m[i][0]:m[i][1]])
				last = m[i][1]
			}
			fmt.Fprintln(out, text[last:])
			needReset = true
		}
		defer func() {
			if needReset {
				fmt.Fprint(out, RESET)
			}
		}()
	}

	var files []string
	if *recursive {
		for _, arg1 := range args {
			stat1, err := os.Stat(arg1)
			if err == nil && stat1.IsDir() {
				filepath.Walk(arg1, func(path string, info os.FileInfo, err error) error {
					if !info.IsDir() {
						files = append(files, path)
					}
					return nil
				})
			} else {
				files = append(files, arg1)
			}
		}
	} else {
		for _, arg1 := range args {
			if addfiles, err := zglob.Glob(arg1); err == nil && addfiles != nil && len(addfiles) > 0 {
				for _, file1 := range addfiles {
					stat, err := os.Stat(file1)
					if err == nil && !stat.IsDir() {
						files = append(files, file1)
					}
				}
			} else if stat1, err := os.Stat(arg1); err == nil && stat1.IsDir() {
				fmt.Fprintf(os.Stderr, "%s is directory\n", arg1)
			} else {
				files = append(files, arg1)
			}
		}
	}
	r := argf.NewFiles(files)
	r.OnError = func(err error) error {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	found := false
	beforeBuffer := make([]string, 0, *flagBefore)
	afterCount := 0
	for r.Scan() {
		text := r.Text()
		text = strings.Replace(text, UTF8BOM, "", 1)

		m := (func() [][]int {
			for _, rx1 := range rxs {
				m1 := rx1.FindAllStringIndex(text, -1)
				if m1 != nil {
					return m1
				}
			}
			return nil
		})()
		if m != nil {
			found = true
			// for `-B n`
			for i, s := range beforeBuffer {
				output(r.Filename(), r.FNR()-len(beforeBuffer)+i, s, [][]int{})
			}
			beforeBuffer = beforeBuffer[:0]

			output(r.Filename(), r.FNR(), text, m)
			afterCount = *flagAfter
		} else if afterCount > 0 {
			afterCount--
			output(r.Filename(), r.FNR(), text, [][]int{})
		} else if *flagBefore > 0 {
			beforeBuffer = append(beforeBuffer, text)
			if len(beforeBuffer) > *flagBefore {
				beforeBuffer = beforeBuffer[1:]
			}
		}
	}
	if r.Err() != nil {
		return r.Err()
	}
	if found {
		return nil
	} else {
		return errors.New("Not found")
	}
}

func main() {
	if err := main1(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
