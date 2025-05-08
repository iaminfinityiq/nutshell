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
	String
	Identifier

	// Keywords
	Let
	Var
	Const

	// Other symbols
	Equals
	Comma
)

type Token struct {
	StartPosition *runtime.Position
	EndPosition   *runtime.Position
	TokenType     int
	Value         string
}

func CreateToken(start_position *runtime.Position, end_position *runtime.Position, token_type int, value string) *Token {
	return &Token{
		StartPosition: start_position,
		EndPosition: end_position,
		TokenType: token_type,
		Value: value,
	}
}