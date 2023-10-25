package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type (
	validationRules map[string]string
	validatorFn     func(reflect.Value, validationRules) error
)

func Validate(v interface{}) error {
	var errs ValidationErrors

	refType := reflect.TypeOf(v)
	refValue := reflect.ValueOf(v)

	if refType.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	for i := 0; i < refValue.NumField(); i++ {
		field := refType.Field(i)
		value := refValue.Field(i)

		tag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}

		rules, err := processTag(tag)
		if err != nil {
			return append(errs, ValidationError{
				Field: field.Name,
				Err:   err,
			})
		}

		errs = append(errs, validate(field, value, rules)...)
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func validate(field reflect.StructField, value reflect.Value, rules validationRules) ValidationErrors {
	var errs ValidationErrors

	kind := field.Type.Kind()

	switch kind { //nolint: exhaustive
	case reflect.Int:
		if err := validateInt(value, rules); err != nil {
			errs = append(errs, ValidationError{
				Field: field.Name,
				Err:   err,
			})
		}
	case reflect.String:
		if err := validateString(value, rules); err != nil {
			errs = append(errs, ValidationError{
				Field: field.Name,
				Err:   err,
			})
		}
	case reflect.Slice:
		elKind := field.Type.Elem().Kind()
		switch elKind { //nolint: exhaustive
		case reflect.Int:
			for _, err := range validateSlice(validateInt, value, rules) {
				errs = append(errs, ValidationError{
					Field: field.Name,
					Err:   err,
				})
			}
		case reflect.String:
			for _, err := range validateSlice(validateString, value, rules) {
				errs = append(errs, ValidationError{
					Field: field.Name,
					Err:   err,
				})
			}
		default:
			errs = append(errs, ValidationError{
				Field: field.Name,
				Err:   ErrUnsupportedType,
			})
		}
	default:
		errs = append(errs, ValidationError{
			Field: field.Name,
			Err:   ErrUnsupportedType,
		})
	}

	return errs
}

func processTag(tag string) (validationRules, error) {
	vm := make(validationRules)

	for _, rule := range strings.Split(tag, "|") {
		fieldRule := strings.Split(rule, ":")

		if len(fieldRule) != 2 || (fieldRule[0] == "") || fieldRule[1] == "" {
			return nil, fmt.Errorf("%w: rule: %s", ErrIncorrectRule, tag)
		}

		vm[fieldRule[0]] = fieldRule[1]
	}

	return vm, nil
}

func validateInt(value reflect.Value, rules validationRules) error {
	for k, v := range rules {
		switch k {
		case "min":
			minVal, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: invalid min value: %s", ErrTagValueMustBeInt, v)
			}

			if value.Int() < minVal {
				return ErrIntLessThanMin
			}
		case "max":
			maxVal, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: invalid max value: %s", ErrTagValueMustBeInt, v)
			}

			if value.Int() > maxVal {
				return ErrIntMoreThanMax
			}
		case "in":
			for _, integer := range strings.Split(v, ",") {
				expected, err := strconv.ParseInt(integer, 10, 64)
				if err != nil {
					return fmt.Errorf("%w: invalid in slice: %s", ErrTagValueMustBeInt, integer)
				}

				if value.Int() == expected {
					return nil
				}
			}

			return ErrNoMatchingElementInSlice
		}
	}

	return nil
}

func validateString(value reflect.Value, rules validationRules) error {
	for k, v := range rules {
		switch k {
		case "len":
			expected, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("%w: invalid len value: %s", err, v)
			}

			if value.Len() != expected {
				return ErrIncorrectStringLen
			}
		case "regexp":
			rgxp, err := regexp.Compile(v)
			if err != nil {
				return fmt.Errorf("%w; invalid regexp: %s", ErrIncorrectTagRegexPattern, v)
			}

			if !rgxp.MatchString(value.String()) {
				return ErrRegexString
			}
		case "in":
			for _, word := range strings.Split(v, ",") {
				if value.String() == word {
					return nil
				}
			}

			return ErrNoMatchingElementInSlice
		}
	}

	return nil
}

func validateSlice(fn validatorFn, values reflect.Value, rules validationRules) []error {
	var errs []error

	for i := 0; i < values.Len(); i++ {
		value := values.Index(i)

		if err := fn(value, rules); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
