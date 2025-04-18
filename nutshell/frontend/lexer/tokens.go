package lexer

import "nutshell/runtime"

const (

	// Whitespaces
	EOF = iota
	Semicolon

	// Operators
	Plus
	Minus
	Multiply
	Divide
	Power
	Modulo

	// ()/[]/{}
	LeftParenthese
	RightParenthese

	// Value types
	Int
	Double
	Identifier
)

type Token struct {
	StartPosition *runtime.Position
	EndPosition   *runtime.Position
	TokenType     int
	Value         string
}
