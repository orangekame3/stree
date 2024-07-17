// Package pkg provides the core functionality of the program.
package pkg

import (
	"strings"
)

// ProcessKeys returns the number of files and directories in the provided slice of keys
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
