package utils

import "golang.org/x/exp/constraints"

func GetValueOrDefault[T constraints.Integer | constraints.Float | string](value T, def T) T {
	var zeroValue T
	if value == zeroValue {
		return def
	}
	return value
}
