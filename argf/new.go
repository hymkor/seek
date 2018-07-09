package argf

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type scanner struct {
	*bufio.Scanner
	files   []string
	n       int
	fd      io.ReadCloser
	err     error
	nr      int
	fnr     int
	OnError func(error) error
}

func NewFiles(files []string) *scanner {
	if len(files) < 1 {
		files = []string{"-"}
	}
	_files := make([]string, 0, len(files)*2)
	for _, file1 := range files {
		if file1 == "-" {
			_files = append(_files, file1)
		} else if matches, err := glob(file1); err == nil && matches != nil {
			for _, m := range matches {
				stat1, err := os.Stat(m)
				if err == nil && !stat1.IsDir() {
					_files = append(_files, m)
				}
			}
		} else {
			_files = append(_files, file1)
		}
	}
	return &scanner{
		Scanner: nil,
		files:   _files,
		n:       -1,
		fd:      nil,
		OnError: func(err error) error { return err },
	}
}

func New() *scanner {
	return NewFiles(os.Args[1:])
}

func (this *scanner) Err() error {
	return this.err
}

func (this *scanner) Scan() bool {
	if this.err != nil {
		return false
	}
	this.nr++
	for {
		if this.Scanner == nil {
			this.n++
			this.fnr = 0
			if this.n >= len(this.files) {
				return false
			}
			if this.files[this.n] == "-" {
				this.fd = ioutil.NopCloser(os.Stdin)
			} else {
				fd, err := os.Open(this.files[this.n])
				if err != nil {
					err = errors.Wrap(err, this.files[this.n])
					if err := this.OnError(err); err != nil {
						this.err = err
						return false
					} else {
						continue
					}
				}
				this.fd = fd
			}
			this.Scanner = bufio.NewScanner(this.fd)
		}
		this.fnr++
		if this.Scanner.Scan() {
			return true
		}
		if err := this.Scanner.Err(); err != nil {
			err = errors.Wrap(err, this.files[this.n])
			if err := this.OnError(err); err != nil {
				this.fd.Close()
				this.err = err
				return false
			}
		}
		if err := this.fd.Close(); err != nil {
			err = errors.Wrap(err, this.files[this.n])
			if err := this.OnError(err); err != nil {
				this.err = err
				return false
			}
		}
		this.Scanner = nil
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
