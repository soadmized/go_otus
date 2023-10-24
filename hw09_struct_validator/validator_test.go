package hw09structvalidator

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) { //nolint: funlen
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "positive",
			in: Response{
				Code: 200,
				Body: "payload",
			},
			expectedErr: nil,
		},
		{
			name:        "interface is not a struct",
			in:          42,
			expectedErr: ErrNotStruct,
		},
		{
			name: "no validate tags",
			in: Token{
				Header:    []byte{1},
				Payload:   []byte{2},
				Signature: []byte{3},
			},
			expectedErr: nil,
		},
		{
			name: "invalid string len",
			in:   App{Version: "1234"},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrIncorrectStringLen,
				},
			},
		},
		{
			name: "invalid min value",
			in: User{
				ID:     "c133763a-8bbd-4087-94ba-dfa9d5e8e805",
				Age:    17,
				Email:  "smth@mail.com",
				Role:   "admin",
				Phones: []string{"89003215467"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   ErrIntLessThanMin,
				},
			},
		},
		{
			name: "invalid max value",
			in: User{
				ID:     "c133763a-8bbd-4087-94ba-dfa9d5e8e805",
				Age:    55,
				Email:  "smth@mail.com",
				Role:   "admin",
				Phones: []string{"89003215467"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   ErrIntMoreThanMax,
				},
			},
		},
		{
			name: "not match to regexp",
			in: User{
				ID:     "c133763a-8bbd-4087-94ba-dfa9d5e8e805",
				Age:    20,
				Email:  "smth",
				Role:   "admin",
				Phones: []string{"89003215467"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Email",
					Err:   ErrRegexString,
				},
			},
		},
		{
			name: "invalid string slice",
			in: User{
				ID:     "c133763a-8bbd-4087-94ba-dfa9d5e8e805",
				Age:    20,
				Email:  "smth@mail.com",
				Role:   "admin",
				Phones: []string{"89003215"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Phones",
					Err:   ErrIncorrectStringLen,
				},
			},
		},
		{
			name: "invalid int slice",
			in: struct {
				integers []int `validate:"min:25"`
			}{
				integers: []int{20},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "integers",
					Err:   ErrIntLessThanMin,
				},
			},
		},
		{
			name: "invalid int tag",
			in: struct {
				integers []int `validate:"min:25.5"`
			}{
				integers: []int{20},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "integers",
					Err:   ErrTagValueMustBeInt,
				},
			},
		},
		{
			name: "invalid regexp tag",
			in: struct {
				strs []string `validate:"regexp:(qwerty]"`
			}{
				strs: []string{"123"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "strs",
					Err:   ErrIncorrectTagRegexPattern,
				},
			},
		},
		{
			name: "invalid | tag",
			in: struct {
				strs []string `validate:"|len:22"`
			}{
				strs: []string{"123"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "strs",
					Err:   ErrIncorrectRule,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
