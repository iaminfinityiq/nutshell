package parser

import "nutshell/runtime"

const (
	BlockStmt = iota
	VariableDeclarationStmt
	AssignmentExpr
	BracketExpr
	BinaryExpr
	UnaryExpr
	CallExpr
	MemberExpr
	IntExpr
	DoubleExpr
	StringExpr
	IdentifierExpr
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
