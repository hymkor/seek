package starstar

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

func expand(dir string, pattern string, recurse bool, callback func(string) error) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			if !recurse {
				continue
			}
			err := expand(filepath.Join(dir, f.Name()), pattern, recurse, callback)
			if err != nil {
				return err
			}
		} else {
			m, err := match(pattern, f.Name())
			if err != nil {
				return err
			}
			if m {
				err = callback(filepath.Join(dir, f.Name()))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func Enumerate(path1 string, callback func(string) error) error {
	pathSlash := filepath.ToSlash(path1)
	if strings.HasPrefix(pathSlash, "**/") {
		return expand(".", path1[3:], true, callback)
	} else if index := strings.Index(pathSlash, "/**/"); index >= 0 {
		return expand(path1[:index], path1[index+4:], true, callback)
	} else {
		dir := filepath.Dir(path1)
		name := filepath.Base(path1)
		if strings.ContainsAny(name, "*?") {
			return expand(dir, name, false, callback)
		} else {
			return callback(path1)
		}
	}
}

func Expand(path1 string) ([]string, error) {
	result := []string{}
	err := Enumerate(path1, func(fname string) error {
		result = append(result, fname)
		return nil
	})
	return result, err
}
