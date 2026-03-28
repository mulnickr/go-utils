// Package optional provides an Optional[R] type containing `value` and `err` and some utility functions.
package optional

type Optional[R any] struct {
	value *R
	err   error
}

func Ok[R any](value R) Optional[R] {
	return Optional[R]{
		value: &value,
	}
}

func New[R any](value R, err error) Optional[R] {
	return Optional[R]{
		value: &value,
		err:   err,
	}
}

func (o Optional[R]) OrElseGet(other R) R {
	if o.IsErr() {
		return other
	}

	return *o.value
}

func (o Optional[R]) OrElsePanic() R {
	if o.IsErr() {
		panic("value not found")
	}

	return *o.value
}

func (o Optional[R]) Value() R {
	return *o.value
}

func (o Optional[R]) Err() error {
	return o.err
}

func (o Optional[R]) IsErr() bool {
	return o.err != nil
}

func (o Optional[R]) IsPresent() bool {
	return !o.IsErr()
}
