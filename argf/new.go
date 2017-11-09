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
}

type argfScanner struct {
	*bufio.Scanner
	files []string
	n     int
	fd    *os.File
	err   error
}

func NewFiles(files []string) Scanner {
	if len(files) < 1 {
		return bufio.NewScanner(os.Stdin)
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
	for {
		if this.err != nil {
			return false
		}
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
