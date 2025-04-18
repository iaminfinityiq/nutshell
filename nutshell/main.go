package main

import (
	"nutshell/frontend/lexer"
	"nutshell/frontend/parser"
	"nutshell/runtime"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	var file_extension string = "nut"
	var test string = `(1 + 2) * 3`

	var l *lexer.Lexer = lexer.InitLexer("nutshell", file_extension, test)
	var rt runtime.RuntimeResult[*[]*lexer.Token] = l.Tokenize()
	if rt.Error != nil {
		rt.Error.DisplayError()
		return
	}

	switch file_extension {
	case "nut":
		var nut_parser *parser.NutParser = parser.InitNutParser(rt.Result)
		var rt runtime.RuntimeResult[*parser.Block] = nut_parser.ParseBlock()
		if rt.Error != nil {
			rt.Error.DisplayError()
			return
		}

		spew.Dump(rt.Result)
	case "nutsh":
		var nutsh_parser *parser.NutshParser = parser.InitNutshParser(rt.Result)
		var rt runtime.RuntimeResult[*parser.Block] = nutsh_parser.ParseBlock()
		if rt.Error != nil {
			rt.Error.DisplayError()
			return
		}
	}
}
