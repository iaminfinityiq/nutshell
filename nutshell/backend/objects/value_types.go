package objects

import (
	"fmt"
	"nutshell/runtime"
)

func MakeBuiltInFunction(heap *Heap, value func(*runtime.Position, *runtime.Position, *[]*ArgumentTuple) runtime.RuntimeResult[*Object]) *Object {
	var returned *Object = &Object{
		DataType: "builtin_function",
		Value:    value,
		Properties: &Scope{
			Heap:  heap,
			Scope: make(map[string]int),
		},
		Heap: heap,
		Flag: true,
	}

	var heap_index int = heap.Add(returned)
	returned.HeapIndex = heap_index

	returned.Assign("call", returned)

	return returned
}

func MakeInt(heap *Heap, value int64) *Object {
	var returned *Object = &Object{
		DataType: "int",
		Value:    value,
		Properties: &Scope{
			Heap:  heap,
			Scope: make(map[string]int),
		},
		Heap: heap,
		Flag: true,
	}

	var heap_index int = heap.Add(returned)
	returned.HeapIndex = heap_index

	returned.Assign("add", MakeBuiltInFunction(heap, func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
		if len(*arguments) != 2 {
			var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 2 argument in 'add' function, got %d/1", len(*arguments)))
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
				Result: MakeInt(heap, left+right),
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
				Result: MakeDouble(heap, float64(left)+right),
				Error:  nil,
			}
		default:
			var err runtime.Error = runtime.TypeError((*arguments)[0].PositionStart, (*arguments)[1].PositionEnd, fmt.Sprintf("Cannot perform operation '+' on int and %s", right_type))
			return runtime.RuntimeResult[*Object]{
				Result: nil,
				Error:  &err,
			}
		}
	}))

	returned.Assign("subtract", MakeBuiltInFunction(heap, func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
		if len(*arguments) != 2 {
			var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 2 argument in 'subtract' function, got %d/1", len(*arguments)))
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
				Result: MakeInt(heap, left-right),
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
				Result: MakeDouble(heap, float64(left)-right),
				Error:  nil,
			}
		default:
			var err runtime.Error = runtime.TypeError((*arguments)[0].PositionStart, (*arguments)[1].PositionEnd, fmt.Sprintf("Cannot perform operation '-' on int and %s", right_type))
			return runtime.RuntimeResult[*Object]{
				Result: nil,
				Error:  &err,
			}
		}
	}))

	returned.Assign("multiply", MakeBuiltInFunction(heap, func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
		if len(*arguments) != 2 {
			var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 2 argument in 'multiply' function, got %d/1", len(*arguments)))
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
				Result: MakeInt(heap, left*right),
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
				Result: MakeDouble(heap, float64(left)*right),
				Error:  nil,
			}
		default:
			var err runtime.Error = runtime.TypeError((*arguments)[0].PositionStart, (*arguments)[1].PositionEnd, fmt.Sprintf("Cannot perform operation '*' on int and %s", right_type))
			return runtime.RuntimeResult[*Object]{
				Result: nil,
				Error:  &err,
			}
		}
	}))

	returned.Assign("divide", MakeBuiltInFunction(heap, func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
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
			Result: MakeDouble(heap, left/right),
			Error:  nil,
		}
	}))

	returned.Assign("modulo", MakeBuiltInFunction(heap, func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
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
			var err runtime.Error = runtime.TypeError((*arguments)[0].PositionStart, (*arguments)[1].PositionEnd, "Cannot perform operation '%' on int and " + right_type)
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
			Result: MakeDouble(heap, left - float64(int_part) * right),
			Error:  nil,
		}
	}))

	return returned
}

func MakeDouble(heap *Heap, value float64) *Object {
	var returned *Object = &Object{
		DataType: "double",
		Value:    value,
		Properties: &Scope{
			Heap:  heap,
			Scope: make(map[string]int),
		},
		Heap: heap,
		Flag: true,
	}

	var heap_index int = heap.Add(returned)
	returned.HeapIndex = heap_index

	returned.Assign("add", MakeBuiltInFunction(heap, func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
		if len(*arguments) != 2 {
			var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 2 argument in 'add' function, got %d/1", len(*arguments)))
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
			Result: MakeDouble(heap, left+right),
			Error:  nil,
		}
	}))

	returned.Assign("subtract", MakeBuiltInFunction(heap, func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
		if len(*arguments) != 2 {
			var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 2 argument in 'subtract' function, got %d/1", len(*arguments)))
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
			Result: MakeDouble(heap, left-right),
			Error:  nil,
		}
	}))

	returned.Assign("multiply", MakeBuiltInFunction(heap, func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
		if len(*arguments) != 2 {
			var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 2 argument in 'multiply' function, got %d/1", len(*arguments)))
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
			Result: MakeDouble(heap, left*right),
			Error:  nil,
		}
	}))

	returned.Assign("divide", MakeBuiltInFunction(heap, func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
		if len(*arguments) != 2 {
			var err runtime.Error = runtime.ArgumentError(position_start, position_end, fmt.Sprintf("Expected 2 argument in 'divide' function, got %d/1", len(*arguments)))
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
			Result: MakeDouble(heap, left/right),
			Error:  nil,
		}
	}))

	returned.Assign("modulo", MakeBuiltInFunction(heap, func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*ArgumentTuple) runtime.RuntimeResult[*Object] {
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
			var err runtime.Error = runtime.TypeError((*arguments)[0].PositionStart, (*arguments)[1].PositionEnd, "Cannot perform operation '%' on double and " + right_type)
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
			Result: MakeDouble(heap, left - float64(int_part) * right),
			Error:  nil,
		}
	}))

	return returned
}
