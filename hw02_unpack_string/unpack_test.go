package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  assert.ErrorAssertionFunc
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde", wantErr: assert.NoError},
		{input: "abccd", expected: "abccd", wantErr: assert.NoError},
		{input: "", expected: "", wantErr: assert.NoError},
		{input: "aaa0b", expected: "aab", wantErr: assert.NoError},
		{input: "aaa今3日2b", expected: "aaa今今今日日b", wantErr: assert.NoError},
		{input: "今3日2", expected: "今今今日日", wantErr: assert.NoError},
		{input: "🏴3", expected: "🏴🏴🏴", wantErr: assert.NoError},
		{input: "3ab", expected: "", wantErr: assert.Error},
		{input: "a33b", expected: "", wantErr: assert.Error},

		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			tc.wantErr(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
