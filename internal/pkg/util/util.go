package util

func ToPtr[T comparable](value T) *T {
	return &value
}
