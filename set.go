package godig

type set[T comparable] map[T]struct{}

func (s set[T]) Add(item T) {
	s[item] = struct{}{}
}

func (s set[T]) Remove(item T) {
	delete(s, item)
}

func (s set[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}
