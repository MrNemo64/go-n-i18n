package cli

func copySlice[T any](arr []T, newElement ...T) []T {
	copied := make([]T, len(arr))
	copy(copied, arr)
	return append(copied, newElement...)
}
