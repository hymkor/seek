package argf

import (
	"bufio"
	"io"
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

type argsScanner struct {
	*bufio.Scanner
	files  []string
	n      int
	closer io.Closer
	err    error
	nr     int
	fnr    int
}

func NewFiles(files []string) Scanner {
	if len(files) < 1 {
		return &stdinScanner{Scanner: bufio.NewScanner(os.Stdin)}
	}
	if files[0] == "-" {
		return &argsScanner{
			Scanner: bufio.NewScanner(os.Stdin),
			files:   files,
			n:       0,
			closer:  nil,
		}
	}
	fd, err := os.Open(files[0])
	if err != nil {
		return &argsScanner{err: err}
	}
	return &argsScanner{
		Scanner: bufio.NewScanner(fd),
		files:   files,
		n:       0,
		closer:  fd,
	}
}

func New() Scanner {
	return NewFiles(os.Args[1:])
}

func (this *argsScanner) Err() error {
	return this.err
}

func (this *argsScanner) Scan() bool {
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

func (this *argsScanner) NR() int {
	return this.nr
}

func (this *argsScanner) FNR() int {
	return this.fnr
}

func (this *argsScanner) Filename() string {
	return this.files[this.n]
}
