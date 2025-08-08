package store

import (
	"tasks/errorcodes"
)

type StoreError struct {
	Code         errorcodes.Code
	Msg          string
	WrappedError error
}

func (e *StoreError) Error() string {
	return e.Msg
}

func (e *StoreError) Unwrap() error {
	return e.WrappedError
}
