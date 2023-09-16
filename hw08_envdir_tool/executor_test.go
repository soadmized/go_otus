package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setEnvs(t *testing.T) {
	tests := []struct {
		name    string
		envs    Environment
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "positive - set and unset variables",
			envs: Environment{
				"FOO": EnvValue{Value: "foo"},
				"BAR": EnvValue{Value: "bar"},
				"FOOBAR": EnvValue{
					Value:      "42",
					NeedRemove: true,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, setEnvs(tt.envs))
		})
	}
}

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name     string
		cmd      []string
		env      Environment
		positive bool
		want     int
	}{
		{
			name:     "positive",
			cmd:      []string{"ls", "-h"},
			env:      Environment{"FOO": EnvValue{Value: "foo"}, "BAR": EnvValue{Value: "bar"}},
			positive: true,
			want:     0,
		},
		{
			name: "error - not found in path",
			cmd:  []string{"ll", "-h"},
			env:  Environment{"FOO": EnvValue{Value: "foo"}, "BAR": EnvValue{Value: "bar"}},
			want: 0,
		},
		{
			name: "error - exec command",
			cmd:  []string{"ls", "-12"},
			env:  Environment{"FOO": EnvValue{Value: "foo"}, "BAR": EnvValue{Value: "bar"}},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.positive {
				assert.Equal(t, tt.want, RunCmd(tt.cmd, tt.env))
			} else {
				assert.NotEqual(t, tt.want, RunCmd(tt.cmd, tt.env))
			}
		})
	}
}
