package interpreter

import (
	"nutshell/backend/objects"
	"nutshell/frontend/parser"
	"nutshell/runtime"
)

func EvaluateBlock(heap *objects.Heap, scope *objects.Scope, ast_node *parser.Block) runtime.RuntimeResult[*objects.Object] {
	var last_evaluated *objects.Object = nil
	for _, statement := range *ast_node.Body {
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
