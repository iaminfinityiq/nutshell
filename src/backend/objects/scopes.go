package objects

type Scope struct {
	Parent      *Scope
	Heap        *Heap
	Scope       map[string]int
	ConstantMap map[string]bool
	DataTypeMap map[string]string
}

func (s *Scope) Declare(variable_name string, value *Object, constant bool) bool {
	_, ok := s.Scope[variable_name]
	if ok {
		return false
	}

	s.Assign(variable_name, value)
	s.ConstantMap[variable_name] = constant
	return true
}

func (s *Scope) Assign(variable_name string, value *Object) {
	if value == nil {
		s.Scope[variable_name] = 0
	} else {
		s.Scope[variable_name] = value.HeapIndex
	}
}

func (s *Scope) Access(variable_name string) (*Object, bool) {
	heap_index, ok := s.Scope[variable_name]
	if !ok {
		if s.Parent == nil {
			return nil, false
		}

		a, b := s.Parent.Access(variable_name)
		return a, b
	}

	returned, ok := s.Heap.Heap[heap_index]
	return returned, ok
}

func (s *Scope) ObjectAccess(variable_name string) (*Object, bool) {
	heap_index, ok := s.Scope[variable_name]
	if !ok {
		return nil, false
	}

	returned, ok := s.Heap.Heap[heap_index]
	return returned, ok
}

func (s *Scope) Mark() {
	for _, v := range s.Scope {
		if v == 0 {
			continue
		}

		if s.Heap.Heap[v].Flag {
			s.Heap.Heap[v].Mark()
		}
	}
}
