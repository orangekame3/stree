// Package: pkg is a package that contains the business logic for stree
package pkg

import (
	"strings"
)

func ProcessKeys(keys [][]string) (int, int) {
	var fileCount int
	var uniqueDirs = map[string]struct{}{}

	for _, key := range keys {
		// Count directories
		for i := 1; i < len(key); i++ {
			uniqueDirs[strings.Join(key[:i], "/")] = struct{}{}
		}

		// Count files
		if len(key) == 1 || key[len(key)-1] != "" {
			fileCount++
		}
	}
	return fileCount, len(uniqueDirs)
}
