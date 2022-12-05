package utils

func ToSlice[T any](items ...T) []T {
	result := make([]T, len(items))
	for i := range items {
		result[i] = items[i]
	}
	return result
}
