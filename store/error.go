package store

import (
	"fmt"
	"tasks/errorcodes"
)

type StoreError struct {
	Code         errorcodes.Code
	Msg          string
	WrappedError error
}

func (e *StoreError) Error() string {
	return fmt.Sprintf("StoreError(%v, %s, %v)", e.Code, e.Msg, e.WrappedError)
}

func (e *StoreError) Unwrap() error {
	return e.WrappedError
}
