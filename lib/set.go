package lib

type Set[T comparable] struct {
	Map  map[T]bool
	Size int
}

func NewSet[T comparable]() *Set[T] {
	s := new(Set[T])
	s.Map = map[T]bool{}
	return s
}

func (s *Set[T]) Add(item T) {
	_, exists := s.Map[item]
	if !exists {
		s.Size++
	}
	s.Map[item] = true
}

func (s *Set[T]) Delete(item T) {
	_, exists := s.Map[item]
	if !exists {
		return
	}
	s.Size--
	delete(s.Map, item)
}

func (s *Set[T]) Has(item T) bool {
	v, exists := s.Map[item]
	return exists && v == true
}

func (s *Set[T]) Iterator() <-chan T {
	c := make(chan T)
	go func() {
		for k, _ := range s.Map {
			c <- k
		}
		close(c)
	}()
	return c
}

func (s *Set[T]) Items() []T {
	var items []T
	for k, _ := range s.Map {
		items = append(items, k)
	}
	return items
}
