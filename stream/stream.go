// Package stream mimics Java's Stream API
// TODO(RM): Do I really care to do this? YES
package stream

import (
	"errors"
	"slices"
)

func ForEach[S any](s []S, fn func(i int, s S)) {
	for i, v := range s {
		fn(i, v)
	}
}

func Map[S, R any](s []S, fn func(s S) R) []R {
	var result []R
	for _, obj := range s {
		result = append(result, fn(obj))
	}

	return result
}

func Max[S any](s []S, fn func(a, b S) int) S {
	copy := slices.Clone(s)
	Sorted(copy, func(a, b S) int {
		return fn(a, b)
	})

	return copy[len(copy)-1]
}

func Min[S any](s []S, fn func(a, b S) int) S {
	copy := slices.Clone(s)
	Sorted(copy, func(a, b S) int {
		return fn(a, b)
	})

	return copy[0]
}

func Sorted[S any](s []S, fn func(a, b S) int) {
	slices.SortFunc(s, fn)
}

func Some[S any](s []S, fn func(s S) bool) bool {
	return slices.ContainsFunc(s, fn)
}

func Every[S any](s []S, fn func(s S) bool) bool {
	for _, obj := range s {
		if !fn(obj) {
			return false
		}
	}

	return true
}

func Never[S any](s []S, fn func(s S) bool) bool {
	inv := func(s S) bool {
		return fn(s)
	}

	return !slices.ContainsFunc(s, inv)
}

// Find returns the all values in a slice that satisfies the function passed.
func Find[S any](s []S, fn func(s S) bool) ([]S, error) {
	var found []S

	for _, v := range s {
		if fn(v) {
			found = append(found, v)
		}
	}

	if len(found) < 1 {
		return found, errors.New("no values found")
	}

	return found, nil
}

// FindOnly returns the only element that satisfies fn.
// Returns an error if there are no values or multiple values.
func FindOnly[S any](s []S, fn func(s S) bool) (S, error) {
	result, err := Find(s, fn)
	if err != nil {
		return result[0], err
	}

	if len(result) > 1 {
		return result[0], errors.New("multiple values found")
	}

	return result[0], nil
}

// FindFirst returns the first element that satisfies fn.
// Returns an error if no values are found
func FindFirst[S any](s []S, fn func(s S) bool) (S, error) {
	result, err := Find(s, fn)
	if err != nil {
		return result[0], err
	}

	return result[0], err
}

// FindLast returns the last element that satifies fn.
// Returns an error if no values are found
func FindLast[S any](s []S, fn func(s S) bool) (S, error) {
	result, err := Find(s, fn)
	if err != nil {
		return result[0], err
	}

	return result[len(result)-1], nil
}
