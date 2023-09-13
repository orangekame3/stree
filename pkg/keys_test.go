package pkg

import (
	"testing"
)

func TestProcessKeys(t *testing.T) {
	tests := []struct {
		name      string
		keys      [][]string
		wantFiles int
		wantDirs  int
	}{
		{
			name:      "No keys",
			keys:      [][]string{},
			wantFiles: 0,
			wantDirs:  0,
		},
		{
			name:      "Single file no directories",
			keys:      [][]string{{"file1.txt"}},
			wantFiles: 1,
			wantDirs:  0,
		},
		{
			name:      "Multiple files in root",
			keys:      [][]string{{"file1.txt"}, {"file2.txt"}},
			wantFiles: 2,
			wantDirs:  0,
		},
		{
			name:      "Single directory no files",
			keys:      [][]string{{"dir1", ""}},
			wantFiles: 0,
			wantDirs:  1,
		},
		{
			name:      "Nested directories with files",
			keys:      [][]string{{"dir1", "dir2", "file1.txt"}},
			wantFiles: 1,
			wantDirs:  2,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			files, dirs := ProcessKeys(tt.keys)
			if files != tt.wantFiles || dirs != tt.wantDirs {
				t.Fatalf("want %d files and %d dirs, got %d files and %d dirs", tt.wantFiles, tt.wantDirs, files, dirs)
			}
		})
	}
}
