package objects

import (
	"fmt"
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
	Parent     int
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
	if o.DataType == "builtin_function" {
		switch attribute_name {
		case "call":
			var new *Object = MakeBuiltInFunction(o.Heap, o.Properties, o.Value.(*BuiltInFunctionPair))
			new.Parent = o.HeapIndex
			return new, true
		case "name":
			var name *Object = MakeString(o.Heap, o.Properties, o.Value.(*BuiltInFunctionPair).Name)
			name.Parent = o.HeapIndex
			return name, true
		case "repr":
			var new *Object = MakeBuiltInFunction(o.Heap, o.Properties, &BuiltInFunctionPair{
				Name: "repr",
				Function: func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
					if len(*arguments) != 1 {
						var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 1 argument in 'repr' function, got %d/1", len(*arguments)))
						return runtime.RuntimeResult[*Object]{
							Result: nil,
							Error:  &err,
						}
					}

					return runtime.RuntimeResult[*Object]{
						Result: MakeString(o.Heap, o.Properties, fmt.Sprintf("<builtin_function %s>", o.Value.(*BuiltInFunctionPair).Name)),
						Error:  nil,
					}
				},
			})
			new.Parent = o.HeapIndex
			return new, true
		default:
			return nil, false
		}
	}

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
