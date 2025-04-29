package interpreter

import (
	"nutshell/backend/objects"
	"nutshell/frontend/parser"
	"nutshell/runtime"
)

func Evaluate(heap *objects.Heap, scope *objects.Scope, ast_node *parser.Statement) runtime.RuntimeResult[*objects.Object] {
	var rt runtime.RuntimeResult[*objects.Object]
	switch (*ast_node).Kind() {
	case parser.BlockStmt:
		var node parser.Block = (*ast_node).(parser.Block)
		rt = EvaluateBlock(heap, scope, &node)
	case parser.BracketExpr:
		var node parser.BracketExpression = (*ast_node).(parser.BracketExpression)
		rt = EvaluateBracketExpression(heap, scope, &node)
	case parser.IntExpr:
		var node parser.Int = (*ast_node).(parser.Int)
		rt = EvaluateInt(heap, scope, &node)
	case parser.DoubleExpr:
		var node parser.Double = (*ast_node).(parser.Double)
		rt = EvaluateDouble(heap, scope, &node)
	case parser.BinaryExpr:
		var node parser.BinaryExpression = (*ast_node).(parser.BinaryExpression)
		rt = EvaluateBinaryExpression(heap, scope, &node)
	}

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
