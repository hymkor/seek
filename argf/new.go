package argf

import (
	"bufio"
	"io"
	"os"
)

type scanner struct {
	*bufio.Scanner
	files  []string
	n      int
	closer io.Closer
	err    error
	nr     int
	fnr    int
}

func NewFiles(files []string) *scanner {
	if len(files) < 1 {
		return &scanner{
			Scanner: bufio.NewScanner(os.Stdin),
			files:   []string{"-"},
			n:       0,
			closer:  nil,
		}
	}
	files_ := make([]string, 0, len(files)*2)
	for _, file1 := range files {
		if matches, err := glob(file1); err == nil && matches != nil {
			files_ = append(files_, matches...)
		} else {
			files_ = append(files_, file1)
		}
	}
	files = files_
	if files[0] == "-" {
		return &scanner{
			Scanner: bufio.NewScanner(os.Stdin),
			files:   files,
			n:       0,
			closer:  nil,
		}
	}
	fd, err := os.Open(files[0])
	if err != nil {
		return &scanner{err: err}
	}
	return &scanner{
		Scanner: bufio.NewScanner(fd),
		files:   files,
		n:       0,
		closer:  fd,
	}
}

func New() *scanner {
	return NewFiles(os.Args[1:])
}

func (this *scanner) Err() error {
	return this.err
}

func (this *scanner) Scan() bool {
	this.nr++
	for {
		if this.err != nil {
			return false
		}
		this.fnr++
		if this.Scanner.Scan() {
			return true
		}
		this.err = this.Scanner.Err()
		if this.err != nil {
			if this.closer != nil {
				this.closer.Close()
			}
			return false
		}
		if this.closer != nil {
			this.err = this.closer.Close()
			this.closer = nil
			if this.err != nil {
				return false
			}
		}
		this.n++
		this.fnr = 0
		if this.n >= len(this.files) {
			return false
		}
		if this.files[this.n] == "-" {
			this.closer = nil
			this.Scanner = bufio.NewScanner(os.Stdin)
		} else {
			fd, err := os.Open(this.files[this.n])
			if err != nil {
				this.err = err
				return false
			}
			this.closer = fd
			this.Scanner = bufio.NewScanner(fd)
		}
	}
}

func (this *scanner) NR() int {
	return this.nr
}

func (this *scanner) FNR() int {
	return this.fnr
}

func (this *scanner) Filename() string {
	return this.files[this.n]
}
