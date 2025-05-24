package objects

import (
	"fmt"
	"nutshell/runtime"
	"strconv"
)

type BuiltInFunctionPair struct {
	Name     string
	Function func(*runtime.Position, *runtime.Position, *[]*ArgumentTuple) runtime.RuntimeResult[*Object]
}

func MakeBuiltInFunction(heap *Heap, scope *Scope, value *BuiltInFunctionPair) *Object {
	var returned *Object = &Object{
		DataType: "builtin_function",
		Bases:    make(map[string]bool),
		Value:    value,
		Properties: &Scope{
			Parent:      scope,
			Heap:        heap,
			Scope:       make(map[string]int),
			ConstantMap: make(map[string]bool),
			DataTypeMap: make(map[string]string),
		},
		Heap: heap,
		Flag: true,
	}

	returned.Bases["any"] = true

	returned.Properties.ConstantMap["call"] = true
	returned.Properties.ConstantMap["repr"] = true

	var heap_index int = heap.Add(returned)
	returned.HeapIndex = heap_index

	returned.Assign("call", nil)
	returned.Assign("repr", nil)

	return returned
}

func MakeType(heap *Heap, scope *Scope, value []string) *Object {
	var returned *Object = &Object{
		DataType: "type",
		Bases:    make(map[string]bool),
		Value:    value,
		Properties: &Scope{
			Parent:      scope,
			Heap:        heap,
			Scope:       make(map[string]int),
			ConstantMap: make(map[string]bool),
			DataTypeMap: make(map[string]string),
		},
	}

	returned.Bases["any"] = true
	returned.Properties.ConstantMap["repr"] = true

	var heap_index int = heap.Add(returned)
	returned.HeapIndex = heap_index

	returned.Assign("repr", MakeBuiltInFunction(heap, scope, &BuiltInFunctionPair{
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
				Result: MakeString(heap, scope, fmt.Sprintf("<type %s>", ((*arguments)[0].Argument.Value.([]string))[0])),
				Error:  nil,
			}
		},
	}))

	heap.Heap[returned.Properties.Scope["repr"]].Parent = returned.HeapIndex

	return returned
}

func MakeNull(heap *Heap, scope *Scope) *Object {
	var returned *Object = &Object{
		DataType: "void",
		Bases:    make(map[string]bool),
		Value:    nil,
		Properties: &Scope{
			Parent:      scope,
			Heap:        heap,
			Scope:       make(map[string]int),
			ConstantMap: make(map[string]bool),
			DataTypeMap: make(map[string]string),
		},

		Heap: heap,
		Flag: true,
	}

	returned.Bases["any"] = true

	returned.Properties.ConstantMap["repr"] = true

	var heap_index int = heap.Add(returned)
	returned.HeapIndex = heap_index

	returned.Assign("repr", MakeBuiltInFunction(heap, scope, &BuiltInFunctionPair{
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
				Result: MakeString(heap, scope, "null"),
				Error:  nil,
			}
		},
	}))

	heap.Heap[returned.Properties.Scope["repr"]].Parent = returned.HeapIndex

	return returned
}

func MakeString(heap *Heap, scope *Scope, value string) *Object {
	var returned *Object = &Object{
		DataType: "string",
		Bases:    make(map[string]bool),
		Value:    value,
		Properties: &Scope{
			Parent:      scope,
			Heap:        heap,
			Scope:       make(map[string]int),
			ConstantMap: make(map[string]bool),
			DataTypeMap: make(map[string]string),
		},
		Heap: heap,
		Flag: true,
	}

	returned.Bases["any"] = true

	returned.Properties.ConstantMap["repr"] = true

	var heap_index int = heap.Add(returned)
	returned.HeapIndex = heap_index

	returned.Properties.ConstantMap["add"] = true

	returned.Assign("add", MakeBuiltInFunction(heap, scope, &BuiltInFunctionPair{
		Name: "add",
		Function: func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
			if len(*arguments) != 2 {
				var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 2 argument in 'add' function, got %d/2", len(*arguments)))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var left string = (*arguments)[0].Argument.Value.(string)
			var right_type string = (*arguments)[1].Argument.DataType
			if right_type != "string" {
				var err runtime.Error = runtime.TypeError(position_start, position_end, fmt.Sprintf("Cannot perform operation '+' on string and %s", right_type))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var right string = (*arguments)[1].Argument.Value.(string)
			return runtime.RuntimeResult[*Object]{
				Result: MakeString(heap, scope, left+right),
				Error:  nil,
			}
		},
	}))

	returned.Assign("repr", MakeBuiltInFunction(heap, scope, &BuiltInFunctionPair{
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
				Result: (*arguments)[0].Argument,
				Error:  nil,
			}
		},
	}))

	heap.Heap[returned.Properties.Scope["add"]].Parent = returned.HeapIndex
	heap.Heap[returned.Properties.Scope["repr"]].Parent = returned.HeapIndex

	return returned
}

func MakeInt(heap *Heap, scope *Scope, value int64) *Object {
	var returned *Object = &Object{
		DataType: "int",
		Bases:    make(map[string]bool),
		Value:    value,
		Properties: &Scope{
			Parent:      scope,
			Heap:        heap,
			Scope:       make(map[string]int),
			ConstantMap: make(map[string]bool),
			DataTypeMap: make(map[string]string),
		},
		Heap: heap,
		Flag: true,
	}

	returned.Bases["any"] = true

	var heap_index int = heap.Add(returned)
	returned.HeapIndex = heap_index

	returned.Properties.ConstantMap["add"] = true
	returned.Properties.ConstantMap["subtract"] = true
	returned.Properties.ConstantMap["multiply"] = true
	returned.Properties.ConstantMap["divide"] = true
	returned.Properties.ConstantMap["modulo"] = true
	returned.Properties.ConstantMap["repr"] = true

	returned.Assign("add", MakeBuiltInFunction(heap, scope, &BuiltInFunctionPair{
		Name: "add",
		Function: func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
			if len(*arguments) != 2 {
				var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 2 argument in 'add' function, got %d/2", len(*arguments)))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var left int64 = (*arguments)[0].Argument.Value.(int64)
			var right_type string = (*arguments)[1].Argument.DataType

			switch right_type {
			case "int":
				var right int64 = (*arguments)[1].Argument.Value.(int64)
				return runtime.RuntimeResult[*Object]{
					Result: MakeInt(heap, scope, left+right),
					Error:  nil,
				}
			case "double":
				var right float64
				if right_type == "int" {
					right = float64((*arguments)[1].Argument.Value.(int64))
				} else {
					right = (*arguments)[1].Argument.Value.(float64)
				}

				return runtime.RuntimeResult[*Object]{
					Result: MakeDouble(heap, scope, float64(left)+right),
					Error:  nil,
				}
			default:
				var err runtime.Error = runtime.TypeError((*arguments)[0].PositionStart, (*arguments)[1].PositionEnd, fmt.Sprintf("Cannot perform operation '+' on int and %s", right_type))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}
		},
	}))

	returned.Assign("subtract", MakeBuiltInFunction(heap, scope, &BuiltInFunctionPair{
		Name: "subtract",
		Function: func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
			if len(*arguments) != 2 {
				var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 2 argument in 'subtract' function, got %d/2", len(*arguments)))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var left int64 = (*arguments)[0].Argument.Value.(int64)
			var right_type string = (*arguments)[1].Argument.DataType

			switch right_type {
			case "int":
				var right int64 = (*arguments)[1].Argument.Value.(int64)
				return runtime.RuntimeResult[*Object]{
					Result: MakeInt(heap, scope, left-right),
					Error:  nil,
				}
			case "double":
				var right float64
				if right_type == "int" {
					right = float64((*arguments)[1].Argument.Value.(int64))
				} else {
					right = (*arguments)[1].Argument.Value.(float64)
				}
				return runtime.RuntimeResult[*Object]{
					Result: MakeDouble(heap, scope, float64(left)-right),
					Error:  nil,
				}
			default:
				var err runtime.Error = runtime.TypeError((*arguments)[0].PositionStart, (*arguments)[1].PositionEnd, fmt.Sprintf("Cannot perform operation '-' on int and %s", right_type))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}
		},
	}))

	returned.Assign("multiply", MakeBuiltInFunction(heap, scope, &BuiltInFunctionPair{
		Name: "multiply",
		Function: func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
			if len(*arguments) != 2 {
				var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 2 argument in 'multiply' function, got %d/2", len(*arguments)))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var left int64 = (*arguments)[0].Argument.Value.(int64)
			var right_type string = (*arguments)[1].Argument.DataType

			switch right_type {
			case "int":
				var right int64 = (*arguments)[1].Argument.Value.(int64)
				return runtime.RuntimeResult[*Object]{
					Result: MakeInt(heap, scope, left*right),
					Error:  nil,
				}
			case "double":
				var right float64
				if right_type == "int" {
					right = float64((*arguments)[1].Argument.Value.(int64))
				} else {
					right = (*arguments)[1].Argument.Value.(float64)
				}

				return runtime.RuntimeResult[*Object]{
					Result: MakeDouble(heap, scope, float64(left)*right),
					Error:  nil,
				}
			default:
				var err runtime.Error = runtime.TypeError((*arguments)[0].PositionStart, (*arguments)[1].PositionEnd, fmt.Sprintf("Cannot perform operation '*' on int and %s", right_type))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}
		},
	}))

	returned.Assign("divide", MakeBuiltInFunction(heap, scope, &BuiltInFunctionPair{
		Name: "divide",
		Function: func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
			if len(*arguments) != 2 {
				var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 2 arguments in 'divide' function, got %d/2", len(*arguments)))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var left float64 = float64((*arguments)[0].Argument.Value.(int64))
			var right_type string = (*arguments)[1].Argument.DataType
			if right_type != "int" && right_type != "double" {
				var err runtime.Error = runtime.TypeError((*arguments)[0].PositionStart, (*arguments)[1].PositionEnd, fmt.Sprintf("Cannot perform operation '/' on int and %s", right_type))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var right float64
			if right_type == "int" {
				right = float64((*arguments)[1].Argument.Value.(int64))
			} else {
				right = (*arguments)[1].Argument.Value.(float64)
			}

			if right == 0 {
				var err runtime.Error = runtime.MathError((*arguments)[1].PositionStart, (*arguments)[1].PositionEnd, fmt.Sprintf("Cannot divide %v by 0", left))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			return runtime.RuntimeResult[*Object]{
				Result: MakeDouble(heap, scope, left/right),
				Error:  nil,
			}
		},
	}))

	returned.Assign("modulo", MakeBuiltInFunction(heap, scope, &BuiltInFunctionPair{
		Name: "modulo",
		Function: func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
			if len(*arguments) != 2 {
				var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 2 arguments in 'modulo' function, got %d/2", len(*arguments)))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var left float64 = float64((*arguments)[0].Argument.Value.(int64))
			var right_type string = (*arguments)[1].Argument.DataType
			if right_type != "int" && right_type != "double" {
				var err runtime.Error = runtime.TypeError((*arguments)[0].PositionStart, (*arguments)[1].PositionEnd, "Cannot perform operation '%' on int and "+right_type)
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var right float64
			if right_type == "int" {
				right = float64((*arguments)[1].Argument.Value.(int64))
			} else {
				right = (*arguments)[1].Argument.Value.(float64)
			}

			if right == 0 {
				var err runtime.Error = runtime.MathError((*arguments)[1].PositionStart, (*arguments)[1].PositionEnd, fmt.Sprintf("Cannot mod %v by 0", left))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var result float64 = left / right
			var int_part int64 = int64(result)

			return runtime.RuntimeResult[*Object]{
				Result: MakeDouble(heap, scope, left-float64(int_part)*right),
				Error:  nil,
			}
		},
	}))

	returned.Assign("repr", MakeBuiltInFunction(heap, scope, &BuiltInFunctionPair{
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
				Result: MakeString(heap, scope, strconv.FormatInt((*arguments)[0].Argument.Value.(int64), 10)),
				Error:  nil,
			}
		},
	}))

	heap.Heap[returned.Properties.Scope["add"]].Parent = returned.HeapIndex
	heap.Heap[returned.Properties.Scope["subtract"]].Parent = returned.HeapIndex
	heap.Heap[returned.Properties.Scope["multiply"]].Parent = returned.HeapIndex
	heap.Heap[returned.Properties.Scope["divide"]].Parent = returned.HeapIndex
	heap.Heap[returned.Properties.Scope["modulo"]].Parent = returned.HeapIndex
	heap.Heap[returned.Properties.Scope["repr"]].Parent = returned.HeapIndex

	return returned
}

func MakeDouble(heap *Heap, scope *Scope, value float64) *Object {
	var returned *Object = &Object{
		DataType: "double",
		Bases:    make(map[string]bool),
		Value:    value,
		Properties: &Scope{
			Parent:      scope,
			Heap:        heap,
			Scope:       make(map[string]int),
			ConstantMap: make(map[string]bool),
			DataTypeMap: make(map[string]string),
		},
		Heap: heap,
		Flag: true,
	}

	returned.Bases["any"] = true

	var heap_index int = heap.Add(returned)
	returned.HeapIndex = heap_index

	returned.Properties.ConstantMap["add"] = true
	returned.Properties.ConstantMap["subtract"] = true
	returned.Properties.ConstantMap["multiply"] = true
	returned.Properties.ConstantMap["divide"] = true
	returned.Properties.ConstantMap["modulo"] = true
	returned.Properties.ConstantMap["repr"] = true

	returned.Assign("add", MakeBuiltInFunction(heap, scope, &BuiltInFunctionPair{
		Name: "add",
		Function: func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
			if len(*arguments) != 2 {
				var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 2 argument in 'add' function, got %d/2", len(*arguments)))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var left float64 = (*arguments)[0].Argument.Value.(float64)
			var right_type string = (*arguments)[1].Argument.DataType

			if right_type != "int" && right_type != "double" {
				var err runtime.Error = runtime.TypeError((*arguments)[0].PositionStart, (*arguments)[1].PositionEnd, fmt.Sprintf("Cannot perform operation '+' on double and %s", right_type))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var right float64
			if right_type == "int" {
				right = float64((*arguments)[1].Argument.Value.(int64))
			} else {
				right = (*arguments)[1].Argument.Value.(float64)
			}

			return runtime.RuntimeResult[*Object]{
				Result: MakeDouble(heap, scope, left+right),
				Error:  nil,
			}
		},
	}))

	returned.Assign("subtract", MakeBuiltInFunction(heap, scope, &BuiltInFunctionPair{
		Name: "subtract",
		Function: func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
			if len(*arguments) != 2 {
				var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 2 argument in 'subtract' function, got %d/2", len(*arguments)))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var left float64 = (*arguments)[0].Argument.Value.(float64)
			var right_type string = (*arguments)[1].Argument.DataType

			if right_type != "int" && right_type != "double" {
				var err runtime.Error = runtime.TypeError((*arguments)[0].PositionStart, (*arguments)[1].PositionEnd, fmt.Sprintf("Cannot perform operation '-' on double and %s", right_type))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var right float64
			if right_type == "int" {
				right = float64((*arguments)[1].Argument.Value.(int64))
			} else {
				right = (*arguments)[1].Argument.Value.(float64)
			}

			return runtime.RuntimeResult[*Object]{
				Result: MakeDouble(heap, scope, left-right),
				Error:  nil,
			}
		},
	}))

	returned.Assign("multiply", MakeBuiltInFunction(heap, scope, &BuiltInFunctionPair{
		Name: "multiply",
		Function: func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
			if len(*arguments) != 2 {
				var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 2 argument in 'multiply' function, got %d/2", len(*arguments)))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var left float64 = (*arguments)[0].Argument.Value.(float64)
			var right_type string = (*arguments)[1].Argument.DataType

			if right_type != "int" && right_type != "double" {
				var err runtime.Error = runtime.TypeError((*arguments)[0].PositionStart, (*arguments)[1].PositionEnd, fmt.Sprintf("Cannot perform operation '*' on double and %s", right_type))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var right float64
			if right_type == "int" {
				right = float64((*arguments)[1].Argument.Value.(int64))
			} else {
				right = (*arguments)[1].Argument.Value.(float64)
			}

			return runtime.RuntimeResult[*Object]{
				Result: MakeDouble(heap, scope, left*right),
				Error:  nil,
			}
		},
	}))

	returned.Assign("divide", MakeBuiltInFunction(heap, scope, &BuiltInFunctionPair{
		Name: "divide",
		Function: func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
			if len(*arguments) != 2 {
				var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 2 argument in 'divide' function, got %d/2", len(*arguments)))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var left float64 = (*arguments)[0].Argument.Value.(float64)
			var right_type string = (*arguments)[1].Argument.DataType

			if right_type != "int" && right_type != "double" {
				var err runtime.Error = runtime.TypeError((*arguments)[0].PositionStart, (*arguments)[1].PositionEnd, fmt.Sprintf("Cannot perform operation '/' on double and %s", right_type))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var right float64
			if right_type == "int" {
				right = float64((*arguments)[1].Argument.Value.(int64))
			} else {
				right = (*arguments)[1].Argument.Value.(float64)
			}

			if right == 0 {
				var err runtime.Error = runtime.MathError((*arguments)[1].PositionStart, (*arguments)[1].PositionEnd, fmt.Sprintf("Cannot divide %v by 0", left))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			return runtime.RuntimeResult[*Object]{
				Result: MakeDouble(heap, scope, left/right),
				Error:  nil,
			}
		},
	}))

	returned.Assign("modulo", MakeBuiltInFunction(heap, scope, &BuiltInFunctionPair{
		Name: "modulo",
		Function: func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
			if len(*arguments) != 2 {
				var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 2 arguments in 'modulo' function, got %d/2", len(*arguments)))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var left float64 = (*arguments)[0].Argument.Value.(float64)
			var right_type string = (*arguments)[1].Argument.DataType
			if right_type != "int" && right_type != "double" {
				var err runtime.Error = runtime.TypeError((*arguments)[0].PositionStart, (*arguments)[1].PositionEnd, "Cannot perform operation '%' on double and "+right_type)
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var right float64
			if right_type == "int" {
				right = float64((*arguments)[1].Argument.Value.(int64))
			} else {
				right = (*arguments)[1].Argument.Value.(float64)
			}

			if right == 0 {
				var err runtime.Error = runtime.MathError((*arguments)[1].PositionStart, (*arguments)[1].PositionEnd, fmt.Sprintf("Cannot mod %v by 0", right))
				return runtime.RuntimeResult[*Object]{
					Result: nil,
					Error:  &err,
				}
			}

			var result float64 = left / right
			var int_part int64 = int64(result)

			return runtime.RuntimeResult[*Object]{
				Result: MakeDouble(heap, scope, left-float64(int_part)*right),
				Error:  nil,
			}
		},
	}))

	returned.Assign("repr", MakeBuiltInFunction(heap, scope, &BuiltInFunctionPair{
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
				Result: MakeString(heap, scope, strconv.FormatFloat((*arguments)[0].Argument.Value.(float64), 'f', -1, 64)),
				Error:  nil,
			}
		},
	}))

	heap.Heap[returned.Properties.Scope["add"]].Parent = returned.HeapIndex
	heap.Heap[returned.Properties.Scope["subtract"]].Parent = returned.HeapIndex
	heap.Heap[returned.Properties.Scope["multiply"]].Parent = returned.HeapIndex
	heap.Heap[returned.Properties.Scope["divide"]].Parent = returned.HeapIndex
	heap.Heap[returned.Properties.Scope["modulo"]].Parent = returned.HeapIndex
	heap.Heap[returned.Properties.Scope["repr"]].Parent = returned.HeapIndex

	return returned
}
