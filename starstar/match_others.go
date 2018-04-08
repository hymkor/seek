// +build !windows

package starstar

import "path/filepath"

func match(pattern, path string) (bool, error) {
	return filepath.Match(pattern, path)
}
