package hw09structvalidator

import (
	"strings"

	"github.com/pkg/errors"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	strBuilder := strings.Builder{}

	for _, ve := range v {
		errMessage := ve.Err.Error()
		if strBuilder.Len() != 0 {
			errMessage += "\n\t"
		}

		strBuilder.WriteString(errMessage)
	}

	return strBuilder.String()
}

var (
	ErrNotStruct                = errors.New("value is not a struct")
	ErrIncorrectRule            = errors.New("unable to parse rule")
	ErrIncorrectStringLen       = errors.New("incorrect string length")
	ErrRegexString              = errors.New("string do not satisfy regexp")
	ErrNoMatchingElementInSlice = errors.New("no matching element")
	ErrIntLessThanMin           = errors.New("int value is less than min")
	ErrIntMoreThanMax           = errors.New("int value is more than max")
	ErrUnsupportedType          = errors.New("type is not supported")
	ErrTagValueMustBeInt        = errors.New("tag value must be integer")
	ErrIncorrectTagRegexPattern = errors.New("regexp pattern is invalid")
)

func (v ValidationErrors) Is(target error) bool {
	var tErr ValidationErrors

	if !errors.As(target, &tErr) {
		return false
	}

	if len(v) != len(tErr) {
		return false
	}

	for i := 0; i < len(tErr); i++ {
		areFieldsEqual := v[i].Field == tErr[i].Field
		areErrsEqual := errors.Is(v[i].Err, tErr[i].Err)
		if !areFieldsEqual || !areErrsEqual {
			return false
		}
	}

	return true
}
