package functools

func Map[T, U any](ts []T, fn func(T) U) []U {
	us := make([]U, len(ts))

	for i, t := range ts {
		us[i] = fn(t)
	}

	return us
}
