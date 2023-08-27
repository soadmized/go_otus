package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fromPath string
		toPath   string
		offset   int64
		limit    int64
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name:     "positive, limit = 0",
			fromPath: "testdata/input.txt",
			toPath:   "new_file.txt",
			offset:   0,
			limit:    0,
			wantErr:  assert.NoError,
		},
		{
			name:     "positive, limit = 100000000",
			fromPath: "testdata/input.txt",
			toPath:   "new_file.txt",
			offset:   0,
			limit:    100000000,
			wantErr:  assert.NoError,
		},
		{
			name:     "positive, with offset",
			fromPath: "testdata/input.txt",
			toPath:   "new_file.txt",
			offset:   100,
			limit:    0,
			wantErr:  assert.NoError,
		},
		{
			name:     "offset exceeds file size",
			fromPath: "testdata/input.txt",
			toPath:   "new_file.txt",
			offset:   1000000000,
			limit:    0,
			wantErr:  assert.Error,
		},
		{
			name:     "unsupported file",
			fromPath: "/dev/urandom",
			toPath:   "new_file.txt",
			offset:   1,
			limit:    0,
			wantErr:  assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(tt.fromPath, tt.toPath, tt.offset, tt.limit)
			tt.wantErr(t, err)
		})
		os.Remove(tt.toPath)
	}
}
