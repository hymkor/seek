// +build windows

package starstar

import (
	"path/filepath"
	"strings"
)

func match(pattern, path string) (bool, error) {
	return filepath.Match(strings.ToUpper(pattern), strings.ToUpper(path))
}
