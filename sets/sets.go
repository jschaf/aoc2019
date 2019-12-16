package sets

type set map[int]bool

type IntSet struct {
	ints set
}

func NewIntSet() IntSet {
	return IntSet{make(set)}
}

func (s *IntSet) Has(x int) bool {
	_, ok := s.ints[x]
	return ok
}

func (s *IntSet) Add(x int) {
	s.ints[x] = true
}

func (s *IntSet) AddAll(xs []int) {
	for _, x := range xs {
		s.Add(x)
	}
}

func (s *IntSet) Remove(x int) {
	delete(s.ints, x)
}

func (s *IntSet) Clear() {
	s.ints = make(set)
}

func (s *IntSet) Size() int {
	return len(s.ints)
}

func (s *IntSet) Union(s2 *IntSet) *IntSet {
	res := NewIntSet()
	for x := range s.ints {
		res.Add(x)
	}
	for x := range s2.ints {
		res.Add(x)
	}
	return &res
}

// Add
// Remove
