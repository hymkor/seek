package argf

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func ReadFd(reader io.Reader, f func(line string) error) error {
	br := bufio.NewScanner(reader)
	for br.Scan() {
		line := br.Text()
		err := f(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadFiles(args []string, f func(line string) error) error {
	if len(args) <= 0 {
		return ReadFd(os.Stdin, f)
	}
	for _, arg1 := range args {
		fd, err := os.Open(arg1)
		if err != nil {
			return err
		}
		if err = ReadFd(fd, f); err != nil {
			fd.Close()
			return err
		}
		if err = fd.Close(); err != nil {
			return err
		}
	}
	return nil
}

func Read(f func(line string) error) error {
	if len(os.Args) >= 2 {
		return ReadFiles(os.Args[1:], f)
	} else {
		return ReadFd(os.Stdin, f)
	}
}

func Read_(f func(line string) error) {
	if err := Read(f); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
