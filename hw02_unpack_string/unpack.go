package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var res strings.Builder
	var prev rune

	for i, v := range str {
		if unicode.IsDigit(v) {
			if i == 0 || unicode.IsDigit(prev) {
				return "", ErrInvalidString
			}
			num, _ := strconv.Atoi(string(v))

			if num == 0 {
				tmp := res.String()
				res = strings.Builder{}
				res.WriteString(removeCharAt(tmp, i-1))
			} else {
				res.WriteString(strings.Repeat(string(prev), num-1))
			}
		} else {
			res.WriteRune(v)
		}
		prev = v
	}

	return res.String(), nil
}

func removeCharAt(str string, pos int) string {
	if pos < 0 || pos >= len(str) {
		return str
	}
	return str[:pos] + str[pos+1:]
}
