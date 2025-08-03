package option

type Option[T any] struct {
	value  T
	isSome bool
}

func Some[T any](v T) Option[T] {
	return Option[T]{value: v, isSome: true}
}

func FromPtr[T any](v *T) Option[T] {
	if v == nil {
		return Option[T]{isSome: false}
	} else {
		return Option[T]{value: *v, isSome: true}
	}
}

func None[T any]() Option[T] {
	return Option[T]{isSome: false}
}

func (o Option[_]) IsSome() bool {
	return o.isSome
}

func (o Option[T]) Value() T {
	return o.value
}

type String = Option[string]
type StringPtr = Option[*string]
