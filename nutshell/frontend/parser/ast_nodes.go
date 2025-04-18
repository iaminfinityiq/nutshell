package parser

import "nutshell/runtime"

const (
	BlockStmt = iota
	BracketExpr
	BinaryExpr
	UnaryExpr
	IntExpr
	DoubleExpr
)

type Statement interface {
	Kind() int
	StartPosition() *runtime.Position
	EndPosition() *runtime.Position
}

type Expression interface {
	Statement
	ExpressionConfirm()
}
