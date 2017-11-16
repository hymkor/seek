// +build !windows

package argf

import (
	"path/filepath"
)

func glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}
