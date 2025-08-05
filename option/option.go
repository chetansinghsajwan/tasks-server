package option

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

// Option represents a value that may or may not be present.
type Option[T any] struct {
	value  T
	isSome bool
}

// Some returns an Option containing a value.
func Some[T any](v T) Option[T] {
	return Option[T]{value: v, isSome: true}
}

// None returns an Option with no value.
func None[T any]() Option[T] {
	var zero T
	return Option[T]{value: zero, isSome: false}
}

// FromPtr converts a pointer to an Option.
func FromPtr[T any](v *T) Option[T] {
	if v == nil {
		return None[T]()
	}
	return Some(*v)
}

func (o Option[_]) IsSome() bool {
	return o.isSome
}

func (o Option[_]) IsNone() bool {
	return !o.isSome
}

func (o Option[T]) Get() (T, bool) { return o.value, o.isSome }

func (o Option[T]) MustGet() T {
	if !o.isSome {
		panic("called Unwrap on None")
	}
	return o.value
}

func (o Option[T]) Ptr() *T {
	if !o.isSome {
		return nil
	}
	return &o.value
}

// String implements fmt.Stringer
func (o Option[T]) String() string {

	if o.IsSome() {
		return fmt.Sprintf("%v", o.value)
	}

	return ""
}

func (o Option[T]) Value() (driver.Value, error) {

	if !o.isSome {
		return nil, nil
	}

	return fmt.Sprintf("%v", o.value), nil
}

// Scan implements the sql.Scanner interface (pgx supports this).
func (o *Option[T]) Scan(src any) error {
	if src == nil {
		*o = None[T]()
		return nil
	}

	// Handle if already of type T
	if val, ok := src.(T); ok {
		*o = Some(val)
		return nil
	}

	// Try using database/sql scanner if T implements it
	var target T
	if scanner, ok := any(&target).(sql.Scanner); ok {
		if err := scanner.Scan(src); err != nil {
			return err
		}
		*o = Some(target)
		return nil
	}

	// Try default assignment for common types (fallback)
	if val, ok := src.(driver.Valuer); ok {
		v, err := val.Value()
		if err != nil {
			return err
		}
		if conv, ok := v.(T); ok {
			*o = Some(conv)
			return nil
		}
	}

	return fmt.Errorf("option.Scan: unsupported type for Option[%T]: %T", o.value, src)
}

// String helpers
type String = Option[string]
type StringPtr = Option[*string]
