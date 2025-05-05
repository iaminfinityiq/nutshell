package main

import (
	"fmt"
	"nutshell/backend/interpreter"
	"nutshell/backend/objects"
	"nutshell/frontend/lexer"
	"nutshell/frontend/parser"
	"nutshell/runtime"
	//"github.com/davecgh/go-spew/spew"
)

func main() {
	var file_extension string = "nut"
	var test string = `
double a = 1
a = 1.5
`

	var l *lexer.Lexer = lexer.InitLexer("nutshell", file_extension, test)
	var rt runtime.RuntimeResult[*[]*lexer.Token] = l.Tokenize()
	if rt.Error != nil {
		rt.Error.DisplayError()
		return
	}

	var rt2 runtime.RuntimeResult[*parser.Block]
	switch file_extension {
	case "nut":
		var nut_parser *parser.NutParser = parser.InitNutParser(rt.Result)
		rt2 = nut_parser.ParseBlock()
		if rt2.Error != nil {
			rt2.Error.DisplayError()
			return
		}
	case "nutsh":
		var nutsh_parser *parser.NutshParser = parser.InitNutshParser(rt.Result)
		rt2 = nutsh_parser.ParseBlock()
		if rt2.Error != nil {
			rt2.Error.DisplayError()
			return
		}
	}

	var heap *objects.Heap = &objects.Heap{
		Heap: make(map[int]*objects.Object),
		Last: 0,
	}

	var scope *objects.Scope = &objects.Scope{
		Parent:      nil,
		Heap:        heap,
		Scope:       make(map[string]int),
		ConstantMap: make(map[string]bool),
		DataTypeMap: make(map[string]string),
	}

	scope.Declare("type", objects.MakeType(heap, scope, []string{"type"}), true)
	scope.Declare("any", objects.MakeType(heap, scope, []string{"any"}), true)
	scope.Declare("builtin_function", objects.MakeType(heap, scope, []string{"builtin_function"}), true)
	scope.Declare("int", objects.MakeType(heap, scope, []string{"int"}), true)
	scope.Declare("double", objects.MakeType(heap, scope, []string{"double"}), true)
	scope.Declare("string", objects.MakeType(heap, scope, []string{"string"}), true)

	var rt3 runtime.RuntimeResult[*objects.Object] = interpreter.EvaluateBlock(heap, scope, rt2.Result)
	if rt3.Error != nil {
		rt3.Error.DisplayError()
		return
	}

	repr, _ := rt3.Result.Access("repr")
	fmt.Println(repr.Value)
}
