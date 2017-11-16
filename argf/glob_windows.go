// +build windows

package argf

import (
	"github.com/zetamatta/go-findfile"
)

// On Windows, ignore case of filename.

func glob(pattern string) ([]string, error) {
	return findfile.Glob(pattern)
}
