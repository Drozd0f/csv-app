package comparator

import (
	"errors"
)

var ErrNotInSlice = errors.New("value not in slice")

func IdxSlice[T comparable](val T, sliceVal []T) (int, error) {
	for idx, v := range sliceVal {
		if v == val {
			return idx, nil
		}
	}

	return 0, ErrNotInSlice
}
