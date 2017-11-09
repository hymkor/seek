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
	files  []string
	n      int
	fd     *os.File
	reader *bufio.Scanner
	err    error
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
		files:  files,
		n:      0,
		fd:     fd,
		reader: bufio.NewScanner(fd),
	}
}

func New() Scanner {
	return NewFiles(os.Args[1:])
}

func (this *argfScanner) Err() error {
	return this.err
}

func (this *argfScanner) Text() string {
	return this.reader.Text()
}

func (this *argfScanner) Bytes() []byte {
	return this.reader.Bytes()
}

func (this *argfScanner) Scan() bool {
	for {
		if this.err != nil {
			return false
		}
		if this.reader.Scan() {
			return true
		}
		this.err = this.reader.Err()
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
		this.reader = bufio.NewScanner(this.fd)
	}
}
