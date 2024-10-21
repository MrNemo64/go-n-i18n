package cli

func copySlice[T any](arr []T, newElement ...T) []T {
	copied := make([]T, len(arr))
	copy(copied, arr)
	return append(copied, newElement...)
}

func Map[T any, R any](v []T, mapper func(*T) R) []R {
	r := make([]R, len(v))
	for i := 0; i < len(v); i++ {
		r[i] = mapper(&v[i])
	}
	return r
}

type Set[T comparable] struct {
	values []T
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{values: make([]T, 0)}
}

func (s *Set[T]) Add(value T) {
	if !s.Contains(value) {
		s.values = append(s.values, value)
	}
}

func (s *Set[T]) Contains(value T) bool {
	for _, v := range s.values {
		if v == value {
			return true
		}
	}
	return false
}

func (s *Set[T]) Get() []T {
	return copySlice(s.values)
}

func (s *Set[T]) Size() int {
	return len(s.values)
}

func (s *Set[T]) AddAll(other *Set[T]) {
	for _, v := range other.values {
		s.Add(v)
	}
}
