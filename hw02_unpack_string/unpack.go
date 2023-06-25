package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var res strings.Builder
	var prev, curr rune
	var digitsFlag bool

	if len(s) == 0 {
		return "", nil
	}

	_, err := strconv.Atoi(string(s[0]))
	if err == nil {
		return "", ErrInvalidString
	}

	for i, v := range s {
		curr = v
		count := int(v - '0')

		if count > 9 {
			res.WriteRune(curr)

			prev = curr
			digitsFlag = false

			continue
		}

		if digitsFlag {
			return "", ErrInvalidString
		}

		if count == 0 {
			str := []byte(res.String())
			str = append(str[:i-1], str[i:]...)

			res.Reset()
			res.Write(str)
		} else {
			str := strings.Repeat(string(prev), count-1)
			r := []byte(str)
			res.Write(r)

			digitsFlag = true
			prev = curr
		}
	}

	return res.String(), nil
}
