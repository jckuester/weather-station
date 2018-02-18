// Package binary implements some binary operations
// for decoding received pulse sequences.
package binary

import (
	"strconv"
)

// ToNumber converts the characters from position b to e (exclusive) of a string input,
// which must be only 0s and 1s, into a decimal number.
func ToNumber(input string, b int, e int) (int, error) {
	i, err := strconv.ParseInt(input[b:e], 2, 0)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

// ToSignedNumber converts the characters from position b to e (exclusive) of a string input,
// which must be only 0s and 1s, into a signed decimal number
// (i.e. the first bit of input is interpreted as the sign).
func ToSignedNumber(input string, b int, e int) (int, error) {
	i, err := strconv.ParseInt(string(input[b]), 2, 0)
	if err != nil {
		return 0, err
	}

	b++

	if i == 1 {
		return toSignedNumberMSBLSB(input, b, e)
	} else {
		return toNumberMSBLSB(input, b, e)
	}
}

func toSignedNumberMSBLSB(input string, b int, e int) (int, error) {
	number := ^0
	i := b

	for i <= e {
		s, err := strconv.ParseInt(string(input[i]), 2, 0)
		if err != nil {
			return 0, err
		}

		number <<= 1
		number |= int(s)
		i++
	}
	return number, nil
}

func toNumberMSBLSB(input string, b int, e int) (int, error) {
	number := 0
	i := b

	for i <= e {
		s, err := strconv.ParseInt(string(input[i]), 2, 0)
		if err != nil {
			return 0, err
		}

		number <<= 1
		number |= int(s)
		i++
	}
	return number, nil
}

// ToBoolean converts the character at position i of a string input,
// which must be only 0s and 1s, into a boolean (1 means true).
func ToBoolean(input string, i int) bool {
	return string(input[i]) == "1"
}
