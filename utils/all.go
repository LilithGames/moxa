package utils

func All[T any](items []T, fn func(item T) bool) bool {
	for i := range items {
		if !fn(items[i]) {
			return false
		}
	}
	return true
}

func Any[T any](items []T, fn func(item T) bool) bool {
	for i := range items {
		if fn(items[i]) {
			return true
		}
	}
	return false
}

func MapAll[K comparable, V any](dict map[K]V, fn func(k K, v V) bool) bool {
	for k, v := range dict {
		if !fn(k, v) {
			return false
		}
	}
	return true
}

func MapAny[K comparable, V any](dict map[K]V, fn func(k K, v V) bool) bool {
	for k, v := range dict {
		if fn(k, v) {
			return true
		}
	}
	return false
}
