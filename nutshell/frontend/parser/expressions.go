package parser

import (
	"nutshell/frontend/lexer"
	"nutshell/runtime"
)

type BracketExpression struct {
	Value                *Expression
	LeftParentheseToken  *lexer.Token
	RightParentheseToken *lexer.Token
}

func (b BracketExpression) StartPosition() *runtime.Position {
	var returned runtime.Position = (*b.LeftParentheseToken).StartPosition.Copy()
	return &returned
}

func (b BracketExpression) EndPosition() *runtime.Position {
	var returned runtime.Position = (*b.RightParentheseToken).EndPosition.Copy()
	return &returned
}

func (b BracketExpression) Kind() int {
	return BracketExpr
}

func (b BracketExpression) ExpressionConfirm() {

}

type BinaryExpression struct {
	Left     *Expression
	Operator int
	Right    *Expression
}

func (b BinaryExpression) StartPosition() *runtime.Position {
	var returned runtime.Position = (*b.Left).StartPosition().Copy()
	return &returned
}

func (b BinaryExpression) EndPosition() *runtime.Position {
	var returned runtime.Position = (*b.Right).EndPosition().Copy()
	return &returned
}

func (b BinaryExpression) Kind() int {
	return BinaryExpr
}

func (b BinaryExpression) ExpressionConfirm() {

}

type UnaryExpression struct {
	Sign           int
	StartSignToken *lexer.Token
	Value          *Expression
}

func (u UnaryExpression) StartPosition() *runtime.Position {
	var returned runtime.Position = u.StartSignToken.StartPosition.Copy()
	return &returned
}

func (u UnaryExpression) EndPosition() *runtime.Position {
	var returned runtime.Position = (*u.Value).EndPosition().Copy()
	return &returned
}

func (u UnaryExpression) Kind() int {
	return UnaryExpr
}

func (u UnaryExpression) ExpressionConfirm() {

}

type Int struct {
	Value    int64
	Position *runtime.Position
}

func (i Int) StartPosition() *runtime.Position {
	var returned runtime.Position = i.Position.Copy()
	return &returned
}

func (i Int) EndPosition() *runtime.Position {
	var returned runtime.Position = i.Position.Copy()
	returned.Advance(nil)
	return &returned
}

func (i Int) Kind() int {
	return IntExpr
}

func (i Int) ExpressionConfirm() {

}

type Double struct {
	Value    float64
	Position *runtime.Position
}

func (d Double) StartPosition() *runtime.Position {
	var returned runtime.Position = d.Position.Copy()
	return &returned
}

func (d Double) EndPosition() *runtime.Position {
	var returned runtime.Position = d.Position.Copy()
	returned.Advance(nil)
	return &returned
}

func (d Double) Kind() int {
	return DoubleExpr
}

func (d Double) ExpressionConfirm() {

}
