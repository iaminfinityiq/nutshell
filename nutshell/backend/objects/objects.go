package objects

import (
	"math"
	"nutshell/runtime"
)

type ArgumentTuple struct {
	PositionStart *runtime.Position
	PositionEnd   *runtime.Position
	Argument      *Object
}

type Object struct {
	DataType   string
	Bases      map[string]bool
	Value      any
	Properties *Scope
	Heap       *Heap
	HeapIndex  int
	Flag       bool
}

func (o *Object) MatchesDataType(data_type string) bool {
	if o.DataType == data_type {
		return true
	}

	return o.Bases[data_type]
}

func (o *Object) Assign(attribute_name string, value *Object) {
	o.Properties.Assign(attribute_name, value)
}

func (o *Object) Access(attribute_name string) (*Object, bool) {
	a, b := o.Properties.ObjectAccess(attribute_name)
	return a, b
}

func (o *Object) Unmark() {
	o.Flag = true
}

func (o *Object) Mark() {
	o.Flag = false
	o.Properties.Mark()
}

type Heap struct {
	Heap map[int]*Object
	Last int
}

func (h *Heap) Add(value *Object) int {
	h.Heap[h.Last] = value
	var returned int = h.Last
	for h.Heap[h.Last] != nil {
		h.Last++
	}

	return returned
}

func (h *Heap) Unmark() {
	for _, v := range h.Heap {
		v.Unmark()
	}
}

func (h *Heap) Sweep() {
	for k, v := range h.Heap {
		if v.Flag {
			h.Last = int(math.Min(float64(h.Last), float64(k)))
			delete(h.Heap, k)
		}
	}
}
