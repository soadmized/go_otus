package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_envValue(t *testing.T) {
	tests := []struct {
		name string
		path string
		want EnvValue
	}{
		{
			name: "positive",
			path: "testdata/env/BAR",
			want: EnvValue{
				Value:      "bar",
				NeedRemove: false,
			},
		},
		{
			name: "empty file",
			path: "testdata/env/EMPTY",
			want: EnvValue{
				Value:      "",
				NeedRemove: false,
			},
		},
		{
			name: "unset env",
			path: "testdata/env/UNSET",
			want: EnvValue{NeedRemove: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := envValue(tt.path)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReadDir(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		want    Environment
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "positive",
			dir:  "testdata/env",
			want: Environment{
				"BAR":   EnvValue{Value: "bar"},
				"EMPTY": EnvValue{Value: ""},
				"FOO":   EnvValue{Value: "   foo\nwith new line"},
				"HELLO": EnvValue{Value: "\"hello\""},
				"UNSET": EnvValue{NeedRemove: true},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "cant read directory",
			dir:     "testdata/foo",
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadDir(tt.dir)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
