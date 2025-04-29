package objects

type Scope struct {
	Heap  *Heap
	Scope map[string]int
}

func (s *Scope) Assign(variable_name string, value *Object) {
	s.Scope[variable_name] = value.HeapIndex
}

func (s *Scope) Access(variable_name string) (*Object, bool) {
	heap_index, ok := s.Scope[variable_name]
	if !ok {
		return nil, false
	}

	returned, ok := s.Heap.Heap[heap_index]
	return returned, ok
}

func (s *Scope) Mark() {
	for _, v := range s.Scope {
		if s.Heap.Heap[v].Flag {
			s.Heap.Heap[v].Mark()
		}
	}
}
