package binary

import "strconv"

// Package binary implements some binary operations
// needed for decoding received pulse sequences.

func ToNumber(data string, b int, e int) (int, error) {
	i, err := strconv.ParseInt(data[b:e], 2, 0)
	if err != nil {
		return 0, err
	}

	return int(i), nil

}

func ToSignedNumber(data string, b int, e int) (int, error) {
	signedPos := b
	b++

	i, err := strconv.ParseInt(string(data[signedPos]), 2, 0)
	if err != nil {
		return 0, err
	}

	if i == 1 {
		return ToSignedNumberMSBLSB(data, b, e)
	} else {
		return ToNumberMSBLSB(data, b, e)
	}
}

func ToSignedNumberMSBLSB(data string, b int, e int) (int, error) {
	number := ^0
	i := b

	for i <= e {
		s, err := strconv.ParseInt(string(data[i]), 2, 0)
		if err != nil {
			return 0, err
		}

		number <<= 1
		number |= int(s)
		i++
	}
	return number, nil
}

func ToNumberMSBLSB(data string, b int, e int) (int, error) {
	number := 0
	i := b

	for i <= e {
		s, err := strconv.ParseInt(string(data[i]), 2, 0)
		if err != nil {
			return 0, err
		}

		number <<= 1
		number |= int(s)
		i++
	}
	return number, nil
}

func ToBoolean(data string, i int) bool {
	return string(data[i]) == "1"
}
