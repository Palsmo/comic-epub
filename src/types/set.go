package types

// hook
type Set[T comparable] struct {
	set[T]
}

// set data type
type set[T comparable] struct {
	m map[T]struct{}
}

// set factory
func NewSet[T comparable](values ...T) Set[T] {
	s := set[T]{}
	s.m = make(map[T]struct{})
	s.Add(values...)
	return Set[T]{s}
}

// add to set
func (s set[T]) Add(values ...T) {
	for _, v := range values {
		s.m[v] = struct{}{}
	}
}

// remove from set
func (s set[T]) Remove(values ...T) {
	for _, v := range values {
		delete(s.m, v)
	}
}

// contains in set
func (s set[T]) Contains(value T) bool {
	_, ok := s.m[value]
	return ok
}
