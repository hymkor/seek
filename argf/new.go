package argf

import (
	"bufio"
	"os"
)

type Scanner interface {
	Scan() bool
	Text() string
	Bytes() []byte
	Err() error
	NR() int
	FNR() int
	Filename() string
}

type stdinScanner struct {
	*bufio.Scanner
	nr int
}

func (this *stdinScanner) Scan() bool {
	this.nr++
	return this.Scanner.Scan()
}

func (this *stdinScanner) NR() int {
	return this.nr
}

func (this *stdinScanner) FNR() int {
	return this.nr
}

func (this *stdinScanner) Filename() string {
	return "-"
}

type argfScanner struct {
	*bufio.Scanner
	files []string
	n     int
	fd    *os.File
	err   error
	nr    int
	fnr   int
}

func NewFiles(files []string) Scanner {
	if len(files) < 1 {
		return &stdinScanner{Scanner: bufio.NewScanner(os.Stdin)}
	}
	fd, err := os.Open(files[0])
	if err != nil {
		return &argfScanner{err: err}
	}
	return &argfScanner{
		Scanner: bufio.NewScanner(fd),
		files:   files,
		n:       0,
		fd:      fd,
	}
}

func New() Scanner {
	return NewFiles(os.Args[1:])
}

func (this *argfScanner) Err() error {
	return this.err
}

func (this *argfScanner) Scan() bool {
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
			this.fd.Close()
			return false
		}
		this.err = this.fd.Close()
		if this.err != nil {
			return false
		}
		this.n++
		this.fnr = 0
		if this.n >= len(this.files) {
			return false
		}
		this.fd, this.err = os.Open(this.files[this.n])
		if this.err != nil {
			return false
		}
		this.Scanner = bufio.NewScanner(this.fd)
	}
}

func (this *argfScanner) NR() int {
	return this.nr
}

func (this *argfScanner) FNR() int {
	return this.fnr
}

func (this *argfScanner) Filename() string {
	return this.files[this.n]
}
