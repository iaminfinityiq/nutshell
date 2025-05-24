package main

import (
	"bufio"
	"fmt"
	"nutshell/backend/interpreter"
	"nutshell/backend/objects"
	"nutshell/frontend/lexer"
	"nutshell/frontend/parser"
	"nutshell/runtime"
	"os"
	//"github.com/davecgh/go-spew/spew"
)

func main() {
	file, err := os.Open("code_syntaxes/hello_world/hello_world.nut")
	var file_extension string = "nut"
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()
	var scanner bufio.Scanner = *bufio.NewScanner(file)
	var data string = ""
	for scanner.Scan() {
		data += scanner.Text()
		data += "\n"
	}

	data = data[0 : len(data)-1]

	var l *lexer.Lexer = lexer.InitLexer("hello_world", file_extension, data)
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
		Last: 1,
	}

	var scope *objects.Scope = &objects.Scope{
		Parent:      nil,
		Heap:        heap,
		Scope:       make(map[string]int),
		ConstantMap: make(map[string]bool),
		DataTypeMap: make(map[string]string),
	}

	// data types
	scope.Declare("type", objects.MakeType(heap, scope, []string{"type"}), true)
	scope.Declare("any", objects.MakeType(heap, scope, []string{"any"}), true)
	scope.Declare("builtin_function", objects.MakeType(heap, scope, []string{"builtin_function"}), true)
	scope.Declare("int", objects.MakeType(heap, scope, []string{"int"}), true)
	scope.Declare("double", objects.MakeType(heap, scope, []string{"double"}), true)
	scope.Declare("string", objects.MakeType(heap, scope, []string{"string"}), true)
	scope.Declare("void", objects.MakeType(heap, scope, []string{"void"}), true)

	// constant values
	scope.Declare("null", objects.MakeNull(heap, scope), true)

	// built-in functions
	scope.Declare("print", objects.MakeBuiltInFunction(heap, scope, &objects.BuiltInFunctionPair{
		Name: "print",
		Function: func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*objects.ArgumentTuple) runtime.RuntimeResult[*objects.Object] {
			var argument_length int = len(*arguments)
			for i, argument := range *arguments {
				repr_function, _ := argument.Argument.Access("repr")

				var repr *objects.Object
				if repr_function.DataType == "builtin_function" {
					var arguments []*objects.ArgumentTuple = []*objects.ArgumentTuple{argument}
					var rt runtime.RuntimeResult[*objects.Object] = repr_function.Value.(*objects.BuiltInFunctionPair).Function(position_start, position_end, &arguments)
					if rt.Error != nil {
						return runtime.RuntimeResult[*objects.Object]{
							Result: nil,
							Error:  rt.Error,
						}
					}

					repr = rt.Result
				}

				fmt.Print(repr.Value)

				if i < argument_length-1 {
					fmt.Print(" ")
				}
			}

			null, _ := scope.Access("null")
			return runtime.RuntimeResult[*objects.Object]{
				Result: null,
				Error:  nil,
			}
		},
	}), true)

	scope.Declare("println", objects.MakeBuiltInFunction(heap, scope, &objects.BuiltInFunctionPair{
		Name: "println",
		Function: func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*objects.ArgumentTuple) runtime.RuntimeResult[*objects.Object] {
			var argument_length int = len(*arguments)
			for i, argument := range *arguments {
				repr_function, _ := argument.Argument.Access("repr")

				var repr *objects.Object
				if repr_function.DataType == "builtin_function" {
					var arguments []*objects.ArgumentTuple = []*objects.ArgumentTuple{argument}
					var rt runtime.RuntimeResult[*objects.Object] = repr_function.Value.(*objects.BuiltInFunctionPair).Function(position_start, position_end, &arguments)
					if rt.Error != nil {
						return runtime.RuntimeResult[*objects.Object]{
							Result: nil,
							Error:  rt.Error,
						}
					}

					repr = rt.Result
				}

				fmt.Print(repr.Value)

				if i < argument_length-1 {
					fmt.Print(" ")
				}
			}

			null, _ := scope.Access("null")

			fmt.Print("\n")
			return runtime.RuntimeResult[*objects.Object]{
				Result: null,
				Error:  nil,
			}
		},
	}), true)

	scope.Declare("input", objects.MakeBuiltInFunction(heap, scope, &objects.BuiltInFunctionPair{
		Name: "input",
		Function: func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*objects.ArgumentTuple) runtime.RuntimeResult[*objects.Object] {
			var argument_length int = len(*arguments)
			for i, argument := range *arguments {
				repr_function, _ := argument.Argument.Access("repr")

				var repr *objects.Object
				if repr_function.DataType == "builtin_function" {
					var args []*objects.ArgumentTuple = []*objects.ArgumentTuple{argument}
					var rt runtime.RuntimeResult[*objects.Object] = repr_function.Value.(*objects.BuiltInFunctionPair).Function(position_start, position_end, &args)
					if rt.Error != nil {
						return runtime.RuntimeResult[*objects.Object]{
							Result: nil,
							Error:  rt.Error,
						}
					}

					repr = rt.Result
				}

				fmt.Print(repr.Value)

				if i < argument_length-1 {
					fmt.Print(" ")
				}
			}

			var input string
			fmt.Scanln(&input)

			return runtime.RuntimeResult[*objects.Object]{
				Result: objects.MakeString(heap, scope, input),
				Error:  nil,
			}
		},
	}), true)

	var zero_values *objects.ZeroValues = &objects.ZeroValues{
		Hashmap: &objects.Scope{
			Parent:      nil,
			Heap:        heap,
			Scope:       make(map[string]int),
			ConstantMap: make(map[string]bool),
			DataTypeMap: make(map[string]string),
		},
	}

	zero_values.DeclareZeroValue("int", objects.MakeInt(heap, scope, 0))
	zero_values.DeclareZeroValue("double", objects.MakeDouble(heap, scope, 0))
	zero_values.DeclareZeroValue("string", objects.MakeString(heap, scope, ""))
	zero_values.DeclareZeroValue("void", objects.MakeNull(heap, scope))

	any_type, _ := scope.Access("any")
	zero_values.DeclareZeroValue("type", any_type)
	zero_values.DeclareZeroValue("any", objects.MakeNull(heap, scope))
	zero_values.DeclareZeroValue("builtin_function", objects.MakeBuiltInFunction(heap, scope, &objects.BuiltInFunctionPair{
		Name: "null_function",
		Function: func(position_start *runtime.Position, position_end *runtime.Position, arguments *[]*objects.ArgumentTuple) runtime.RuntimeResult[*objects.Object] {
			return runtime.RuntimeResult[*objects.Object]{
				Result: objects.MakeNull(heap, scope),
				Error:  nil,
			}
		},
	}))

	var rt3 runtime.RuntimeResult[*objects.Object] = interpreter.EvaluateBlock(heap, scope, zero_values, rt2.Result)
	if rt3.Error != nil {
		rt3.Error.DisplayError()
		return
	}
}
