package utils

import "strconv"

type IntConstraint interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type UintConstraint interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func ParseInt[T IntConstraint](str string, bitSize int) (T, error) {

	value, err := strconv.ParseUint(str, 10, bitSize)

	if err != nil {
		var zero T
		return zero, err
	}

	return T(value), nil
}

func ParseUint[T UintConstraint](str string, bitSize int) (T, error) {

	value, err := strconv.ParseUint(str, 10, bitSize)

	if err != nil {
		var zero T
		return zero, err
	}

	return T(value), nil
}

func ParseInt8(str string) (int8, error) {
	return ParseInt[int8](str, 8)
}

func ParseInt16(str string) (int16, error) {
	return ParseInt[int16](str, 16)
}

func ParseInt32(str string) (int32, error) {
	return ParseInt[int32](str, 32)
}

func ParseInt64(str string) (int64, error) {
	return ParseInt[int64](str, 64)
}

func ParseUint8(str string) (uint8, error) {
	return ParseUint[uint8](str, 8)
}

func ParseUint16(str string) (uint16, error) {
	return ParseUint[uint16](str, 16)
}

func ParseUint32(str string) (uint32, error) {
	return ParseUint[uint32](str, 32)
}

func ParseUint64(str string) (uint64, error) {
	return ParseUint[uint64](str, 64)
}
