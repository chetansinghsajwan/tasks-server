package utils

import "strconv"

func ParseInt32(str string) (int32, error) {
	var value, err = strconv.ParseInt(str, 10, 32)

	if err != nil {
		return 0, nil
	}

	return int32(value), nil
}
