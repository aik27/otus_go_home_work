package hw02unpackstring

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidString = errors.New("invalid string")
	ErrInvalidAtoi   = errors.New("invalid cast to int")
)

func Unpack(str string) (string, error) {
	var res strings.Builder
	var prev rune
	var isPrevSlash bool

	re := regexp.MustCompile(`^(([a-zA-Z\x{1F600}-\x{1F64F}])+\d?|\\+\d+)*$`)
	ok := re.Match([]byte(str))

	if !ok {
		return "", ErrInvalidString
	}

	for _, char := range str {
		if isPrevSlash {
			res.WriteRune(prev)
			isPrevSlash = false
			prev = char
			continue
		}

		if string(char) == `\` {
			isPrevSlash = true
			continue
		}

		if unicode.IsDigit(char) {
			num, err := strconv.Atoi(string(char))
			if err != nil {
				return "", ErrInvalidAtoi
			}

			res.WriteString(strings.Repeat(string(prev), num))
			prev = 0
		} else {
			if prev != 0 {
				res.WriteRune(prev)
			}
			prev = char
		}
	}

	if prev != 0 {
		res.WriteRune(prev)
	}

	return res.String(), nil
}
