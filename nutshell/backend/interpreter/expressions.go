package interpreter

import (
	"fmt"
	"nutshell/backend/objects"
	"nutshell/frontend/lexer"
	"nutshell/frontend/parser"
	"nutshell/runtime"
)

func EvaluateBracketExpression(heap *objects.Heap, scope *objects.Scope, ast_node *parser.BracketExpression) runtime.RuntimeResult[*objects.Object] {
	var expression parser.Statement = (*ast_node.Value).(parser.Statement)
	var rt runtime.RuntimeResult[*objects.Object] = Evaluate(heap, scope, &expression)
	if rt.Error != nil {
		return runtime.RuntimeResult[*objects.Object]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	return runtime.RuntimeResult[*objects.Object]{
		Result: rt.Result,
		Error:  nil,
	}
}

func EvaluateAssignmentExpression(heap *objects.Heap, scope *objects.Scope, ast_node *parser.AssignmentExpression) runtime.RuntimeResult[*objects.Object] {
	if (*ast_node.Callee).Kind() == parser.IdentifierExpr {
		var left_node parser.Identifier = (*ast_node.Callee).(parser.Identifier)
		if scope.ConstantMap[left_node.VariableName] {
			var err runtime.Error = runtime.VariableError(ast_node.StartPosition(), ast_node.EndPosition(), fmt.Sprintf("Cannot reassign variable %s because it is a constant!", left_node.VariableName))
			return runtime.RuntimeResult[*objects.Object]{
				Result: nil,
				Error:  &err,
			}
		}

		var right_node parser.Statement = (*ast_node.Value).(parser.Statement)
		var rt runtime.RuntimeResult[*objects.Object] = Evaluate(heap, scope, &right_node)
		if rt.Error != nil {
			return runtime.RuntimeResult[*objects.Object]{
				Result: nil,
				Error:  rt.Error,
			}
		}

		data_type, ok := scope.DataTypeMap[left_node.VariableName]
		if ok {
			if !rt.Result.MatchesDataType(data_type) {
				var err runtime.Error = runtime.TypeError(right_node.StartPosition(), right_node.EndPosition(), fmt.Sprintf("Expected the type of the value to be a %s, not %s", data_type, rt.Result.DataType))
				return runtime.RuntimeResult[*objects.Object]{
					Result: nil,
					Error:  &err,
				}
			}
		}

		scope.Assign(left_node.VariableName, rt.Result)
		return runtime.RuntimeResult[*objects.Object]{
			Result: rt.Result,
			Error:  nil,
		}
	}

	var err runtime.Error = runtime.SyntaxError(ast_node.StartPosition(), ast_node.EndPosition(), "Invalid syntax!")
	return runtime.RuntimeResult[*objects.Object]{
		Result: nil,
		Error:  &err,
	}
}

func EvaluateIdentifier(heap *objects.Heap, scope *objects.Scope, ast_node *parser.Identifier) runtime.RuntimeResult[*objects.Object] {
	value, ok := scope.Access(ast_node.VariableName)
	if !ok {
		var err runtime.Error = runtime.VariableError(ast_node.StartPosition(), ast_node.EndPosition(), fmt.Sprintf("Cannot access variable %s because it does not exist", ast_node.VariableName))
		return runtime.RuntimeResult[*objects.Object]{
			Result: nil,
			Error:  &err,
		}
	}

	return runtime.RuntimeResult[*objects.Object]{
		Result: value,
		Error:  nil,
	}
}

func EvaluateInt(heap *objects.Heap, scope *objects.Scope, ast_node *parser.Int) runtime.RuntimeResult[*objects.Object] {
	var returned *objects.Object = objects.MakeInt(heap, scope, ast_node.Value)
	return runtime.RuntimeResult[*objects.Object]{
		Result: returned,
		Error:  nil,
	}
}

func EvaluateDouble(heap *objects.Heap, scope *objects.Scope, ast_node *parser.Double) runtime.RuntimeResult[*objects.Object] {
	var returned *objects.Object = objects.MakeDouble(heap, scope, ast_node.Value)
	return runtime.RuntimeResult[*objects.Object]{
		Result: returned,
		Error:  nil,
	}
}

func EvaluateBinaryExpression(heap *objects.Heap, scope *objects.Scope, ast_node *parser.BinaryExpression) runtime.RuntimeResult[*objects.Object] {
	var node parser.Statement = (*ast_node.Left).(parser.Statement)
	var rt runtime.RuntimeResult[*objects.Object] = Evaluate(heap, scope, &node)
	if rt.Error != nil {
		return runtime.RuntimeResult[*objects.Object]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	var left *objects.Object = rt.Result

	node = (*ast_node.Right).(parser.Statement)
	rt = Evaluate(heap, scope, &node)
	if rt.Error != nil {
		return runtime.RuntimeResult[*objects.Object]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	var right *objects.Object = rt.Result

	switch ast_node.Operator {
	case lexer.Plus:
		add_attribute, ok := left.Access("add")
		if !ok {
			var err runtime.Error = runtime.TypeError(ast_node.StartPosition(), ast_node.EndPosition(), fmt.Sprintf("Cannot perform operation '+' on %s and %s", left.DataType, right.DataType))
			return runtime.RuntimeResult[*objects.Object]{
				Result: nil,
				Error:  &err,
			}
		}

		if add_attribute.DataType == "builtin_function" {
			var add_function func(*runtime.Position, *runtime.Position, *[]*objects.ArgumentTuple) runtime.RuntimeResult[*objects.Object] = add_attribute.Value.(func(*runtime.Position, *runtime.Position, *[]*objects.ArgumentTuple) runtime.RuntimeResult[*objects.Object])

			var arguments []*objects.ArgumentTuple = []*objects.ArgumentTuple{
				&objects.ArgumentTuple{
					PositionStart: (*ast_node.Left).StartPosition(),
					PositionEnd:   (*ast_node.Left).EndPosition(),
					Argument:      left,
				},
				&objects.ArgumentTuple{
					PositionStart: (*ast_node.Right).StartPosition(),
					PositionEnd:   (*ast_node.Right).EndPosition(),
					Argument:      right,
				},
			}

			rt = add_function(ast_node.StartPosition(), ast_node.EndPosition(), &arguments)
			if rt.Error != nil {
				return runtime.RuntimeResult[*objects.Object]{
					Result: nil,
					Error:  rt.Error,
				}
			}

			return runtime.RuntimeResult[*objects.Object]{
				Result: rt.Result,
				Error:  nil,
			}
		}
	case lexer.Minus:
		subtract_attribute, ok := left.Access("subtract")
		if !ok {
			var err runtime.Error = runtime.TypeError(ast_node.StartPosition(), ast_node.EndPosition(), fmt.Sprintf("Cannot perform operation '-' on %s and %s", left.DataType, right.DataType))
			return runtime.RuntimeResult[*objects.Object]{
				Result: nil,
				Error:  &err,
			}
		}

		if subtract_attribute.DataType == "builtin_function" {
			var subtract_function func(*runtime.Position, *runtime.Position, *[]*objects.ArgumentTuple) runtime.RuntimeResult[*objects.Object] = subtract_attribute.Value.(func(*runtime.Position, *runtime.Position, *[]*objects.ArgumentTuple) runtime.RuntimeResult[*objects.Object])

			var arguments []*objects.ArgumentTuple = []*objects.ArgumentTuple{
				&objects.ArgumentTuple{
					PositionStart: (*ast_node.Left).StartPosition(),
					PositionEnd:   (*ast_node.Left).EndPosition(),
					Argument:      left,
				},
				&objects.ArgumentTuple{
					PositionStart: (*ast_node.Right).StartPosition(),
					PositionEnd:   (*ast_node.Right).EndPosition(),
					Argument:      right,
				},
			}

			rt = subtract_function(ast_node.StartPosition(), ast_node.EndPosition(), &arguments)
			if rt.Error != nil {
				return runtime.RuntimeResult[*objects.Object]{
					Result: nil,
					Error:  rt.Error,
				}
			}

			return runtime.RuntimeResult[*objects.Object]{
				Result: rt.Result,
				Error:  nil,
			}
		}
	case lexer.Multiply:
		multiply_attribute, ok := left.Access("multiply")
		if !ok {
			var err runtime.Error = runtime.TypeError(ast_node.StartPosition(), ast_node.EndPosition(), fmt.Sprintf("Cannot perform operation '*' on %s and %s", left.DataType, right.DataType))
			return runtime.RuntimeResult[*objects.Object]{
				Result: nil,
				Error:  &err,
			}
		}

		if multiply_attribute.DataType == "builtin_function" {
			var multiply_fraction func(*runtime.Position, *runtime.Position, *[]*objects.ArgumentTuple) runtime.RuntimeResult[*objects.Object] = multiply_attribute.Value.(func(*runtime.Position, *runtime.Position, *[]*objects.ArgumentTuple) runtime.RuntimeResult[*objects.Object])

			var arguments []*objects.ArgumentTuple = []*objects.ArgumentTuple{
				&objects.ArgumentTuple{
					PositionStart: (*ast_node.Left).StartPosition(),
					PositionEnd:   (*ast_node.Left).EndPosition(),
					Argument:      left,
				},
				&objects.ArgumentTuple{
					PositionStart: (*ast_node.Right).StartPosition(),
					PositionEnd:   (*ast_node.Right).EndPosition(),
					Argument:      right,
				},
			}

			rt = multiply_fraction(ast_node.StartPosition(), ast_node.EndPosition(), &arguments)
			if rt.Error != nil {
				return runtime.RuntimeResult[*objects.Object]{
					Result: nil,
					Error:  rt.Error,
				}
			}

			return runtime.RuntimeResult[*objects.Object]{
				Result: rt.Result,
				Error:  nil,
			}
		}
	case lexer.Divide:
		divide_attribute, ok := left.Access("divide")
		if !ok {
			var err runtime.Error = runtime.TypeError(ast_node.StartPosition(), ast_node.EndPosition(), fmt.Sprintf("Cannot perform operation '/' on %s and %s", left.DataType, right.DataType))
			return runtime.RuntimeResult[*objects.Object]{
				Result: nil,
				Error:  &err,
			}
		}

		if divide_attribute.DataType == "builtin_function" {
			var divide_function func(*runtime.Position, *runtime.Position, *[]*objects.ArgumentTuple) runtime.RuntimeResult[*objects.Object] = divide_attribute.Value.(func(*runtime.Position, *runtime.Position, *[]*objects.ArgumentTuple) runtime.RuntimeResult[*objects.Object])

			var arguments []*objects.ArgumentTuple = []*objects.ArgumentTuple{
				&objects.ArgumentTuple{
					PositionStart: (*ast_node.Left).StartPosition(),
					PositionEnd:   (*ast_node.Left).EndPosition(),
					Argument:      left,
				},
				&objects.ArgumentTuple{
					PositionStart: (*ast_node.Right).StartPosition(),
					PositionEnd:   (*ast_node.Right).EndPosition(),
					Argument:      right,
				},
			}

			rt = divide_function(ast_node.StartPosition(), ast_node.EndPosition(), &arguments)
			if rt.Error != nil {
				return runtime.RuntimeResult[*objects.Object]{
					Result: nil,
					Error:  rt.Error,
				}
			}

			return runtime.RuntimeResult[*objects.Object]{
				Result: rt.Result,
				Error:  nil,
			}
		}
	case lexer.Modulo:
		divide_attribute, ok := left.Access("modulo")
		if !ok {
			var err runtime.Error = runtime.TypeError(ast_node.StartPosition(), ast_node.EndPosition(), "Cannot perform operation '%' on "+left.DataType+" and "+right.DataType)
			return runtime.RuntimeResult[*objects.Object]{
				Result: nil,
				Error:  &err,
			}
		}

		if divide_attribute.DataType == "builtin_function" {
			var divide_function func(*runtime.Position, *runtime.Position, *[]*objects.ArgumentTuple) runtime.RuntimeResult[*objects.Object] = divide_attribute.Value.(func(*runtime.Position, *runtime.Position, *[]*objects.ArgumentTuple) runtime.RuntimeResult[*objects.Object])

			var arguments []*objects.ArgumentTuple = []*objects.ArgumentTuple{
				&objects.ArgumentTuple{
					PositionStart: (*ast_node.Left).StartPosition(),
					PositionEnd:   (*ast_node.Left).EndPosition(),
					Argument:      left,
				},
				&objects.ArgumentTuple{
					PositionStart: (*ast_node.Right).StartPosition(),
					PositionEnd:   (*ast_node.Right).EndPosition(),
					Argument:      right,
				},
			}

			rt = divide_function(ast_node.StartPosition(), ast_node.EndPosition(), &arguments)
			if rt.Error != nil {
				return runtime.RuntimeResult[*objects.Object]{
					Result: nil,
					Error:  rt.Error,
				}
			}

			return runtime.RuntimeResult[*objects.Object]{
				Result: rt.Result,
				Error:  nil,
			}
		}
	}

	return runtime.RuntimeResult[*objects.Object]{
		Result: nil,
		Error:  nil,
	}
}
