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

type ArgfScanner struct {
	Files    []string
	N        int
	Fd       *os.File
	Reader   *bufio.Scanner
	ErrValue error
}

func NewFiles(files []string) Scanner {
	if len(files) < 1 {
		return bufio.NewScanner(os.Stdin)
	}
	fd, err := os.Open(files[0])
	if err != nil {
		return &ArgfScanner{ErrValue: err}
	}
	return &ArgfScanner{
		Files:  files,
		N:      0,
		Fd:     fd,
		Reader: bufio.NewScanner(fd),
	}
}

func New() Scanner {
	return NewFiles(os.Args[1:])
}

func (this *ArgfScanner) Err() error {
	return this.ErrValue
}

func (this *ArgfScanner) Text() string {
	return this.Reader.Text()
}

func (this *ArgfScanner) Bytes() []byte {
	return this.Reader.Bytes()
}

func (this *ArgfScanner) Scan() bool {
	for {
		if this.ErrValue != nil {
			return false
		}
		if this.Reader.Scan() {
			return true
		}
		this.ErrValue = this.Reader.Err()
		if this.ErrValue != nil {
			this.Fd.Close()
			return false
		}
		this.ErrValue = this.Fd.Close()
		if this.ErrValue != nil {
			return false
		}
		this.N++
		if this.N >= len(this.Files) {
			return false
		}
		this.Fd, this.ErrValue = os.Open(this.Files[this.N])
		if this.ErrValue != nil {
			return false
		}
		this.Reader = bufio.NewScanner(this.Fd)
	}
}
