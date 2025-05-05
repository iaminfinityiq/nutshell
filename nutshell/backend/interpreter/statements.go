package interpreter

import (
	"fmt"
	"nutshell/backend/objects"
	"nutshell/frontend/parser"
	"nutshell/runtime"
)

func EvaluateBlock(heap *objects.Heap, scope *objects.Scope, ast_node *parser.Block) runtime.RuntimeResult[*objects.Object] {
	var last_evaluated *objects.Object = nil
	for _, statement := range ast_node.Body {
		var rt runtime.RuntimeResult[*objects.Object] = Evaluate(heap, scope, statement)
		if rt.Error != nil {
			return runtime.RuntimeResult[*objects.Object]{
				Result: nil,
				Error:  rt.Error,
			}
		}

		last_evaluated = rt.Result
		objects.Mark(scope)
		objects.Sweep(heap)
	}

	return runtime.RuntimeResult[*objects.Object]{
		Result: last_evaluated,
		Error:  nil,
	}
}

func EvaluateVariableDeclaration(heap *objects.Heap, scope *objects.Scope, ast_node *parser.VariableDeclaration) runtime.RuntimeResult[*objects.Object] {
	var node parser.Statement = (*ast_node.DataType).(parser.Statement)
	var rt runtime.RuntimeResult[*objects.Object] = Evaluate(heap, scope, &node)
	if rt.Error != nil {
		return runtime.RuntimeResult[*objects.Object]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	if !rt.Result.MatchesDataType("type") {
		var err runtime.Error = runtime.TypeError(node.StartPosition(), node.EndPosition(), fmt.Sprintf("Expected the type of a variable in a variable declaration to be a type, not %s", rt.Result.DataType))
		return runtime.RuntimeResult[*objects.Object]{
			Result: nil,
			Error:  &err,
		}
	}

	var data_type string = (rt.Result.Value).([]string)[0]
	node = (*ast_node.Value).(parser.Statement)
	rt = Evaluate(heap, scope, &node)
	if rt.Error != nil {
		return runtime.RuntimeResult[*objects.Object]{
			Result: nil,
			Error:  rt.Error,
		}
	}

	if !rt.Result.MatchesDataType(data_type) {
		var err runtime.Error = runtime.TypeError(node.StartPosition(), node.EndPosition(), fmt.Sprintf("Expected the type of the value to be a %s, not %s", data_type, rt.Result.DataType))
		return runtime.RuntimeResult[*objects.Object]{
			Result: nil,
			Error:  &err,
		}
	}

	var result bool = scope.Declare(ast_node.VariableName, rt.Result, ast_node.IsConstant)
	if !result {
		var err runtime.Error = runtime.VariableError(ast_node.StartPosition(), ast_node.EndPosition(), fmt.Sprintf("Variable with name %s already exists!", ast_node.VariableName))
		return runtime.RuntimeResult[*objects.Object]{
			Result: nil,
			Error:  &err,
		}
	}

	scope.DataTypeMap[ast_node.VariableName] = data_type
	return runtime.RuntimeResult[*objects.Object]{
		Result: nil,
		Error:  nil,
	}
}
