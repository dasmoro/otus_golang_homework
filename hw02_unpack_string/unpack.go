package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strings"
	"unicode"
)

// ErrInvalidString - error when uncpacked stgring have errors.
var ErrInvalidString = errors.New("invalid string")

// Unpack is a simple string unpack.
func Unpack(s *string) (string, error) {
	var buf, result string
	var err error
	var isEscape bool

	for _, c := range *s {
		isEscape = !isEscape && buf == `\`

		if unicode.IsDigit(c) {
			switch {
			case isEscape:
				buf = string(c)
			case len(buf) > 0:
				result += strings.Repeat(buf, int(c-'0'))
				buf = ""
			default:
				err = ErrInvalidString
				return "", err
			}
		} else {
			if len(buf) > 0 && !isEscape {
				result += buf
			}
			buf = string(c)
		}
	}
	result += buf
	return result, err
}
