package pointer

func Pointer[T any](v T) *T {
	return &v
}

func Value[T any](v *T) T {
	var t T
	if v != nil {
		t = *v
	}

	return t
}
