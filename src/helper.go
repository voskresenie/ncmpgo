package main

func ptr[T any](t T) *T {
	return &t
}

func eqSlice[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func eqSet[T comparable](a, b []T) bool {
	sa := map[T]struct{}{}
	for _, e := range a {
		sa[e] = struct{}{}
	}
	sb := map[T]struct{}{}
	for _, e := range b {
		sb[e] = struct{}{}
	}
	if len(sa) != len(sb) {
		return false
	}
	for e := range sa {
		if _, ok := sb[e]; !ok {
			return false
		}
	}
	return true
}
