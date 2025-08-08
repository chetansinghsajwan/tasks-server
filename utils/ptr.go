package utils

func Ptr[T any](v T) *T {
	return &v
}

func DoublePtr[T any](v T) **T {
	var vptr = &v
	return &vptr
}
