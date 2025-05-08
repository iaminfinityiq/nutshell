package parser

import (
	"nutshell/frontend/lexer"
	"nutshell/runtime"
)

type AssignmentExpression struct {
	Callee        *Expression
	Value         *Expression
	PositionStart *runtime.Position
	PositionEnd   *runtime.Position
}

func (a AssignmentExpression) StartPosition() *runtime.Position {
	return a.PositionStart
}

func (a AssignmentExpression) EndPosition() *runtime.Position {
	return a.PositionEnd
}

func (a AssignmentExpression) Kind() int {
	return AssignmentExpr
}

func (a AssignmentExpression) ExpressionConfirm() {

}

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

type CallExpression struct {
	Callee               *Expression
	Arguments            []*Expression
	LeftParentheseToken  *lexer.Token
	RightParentheseToken *lexer.Token
}

func (c CallExpression) StartPosition() *runtime.Position {
	var returned runtime.Position = (*c.Callee).StartPosition().Copy()
	return &returned
}

func (c CallExpression) EndPosition() *runtime.Position {
	var returned runtime.Position = c.RightParentheseToken.EndPosition.Copy()
	return &returned
}

func (c CallExpression) Kind() int {
	return CallExpr
}

func (c CallExpression) ExpressionConfirm() {

}

type Int struct {
	Value    int64
	IntToken *lexer.Token
}

func (i Int) StartPosition() *runtime.Position {
	var returned runtime.Position = i.IntToken.StartPosition.Copy()
	return &returned
}

func (i Int) EndPosition() *runtime.Position {
	var returned runtime.Position = i.IntToken.EndPosition.Copy()
	return &returned
}

func (i Int) Kind() int {
	return IntExpr
}

func (i Int) ExpressionConfirm() {

}

type Double struct {
	Value       float64
	DoubleToken *lexer.Token
}

func (d Double) StartPosition() *runtime.Position {
	var returned runtime.Position = d.DoubleToken.StartPosition.Copy()
	return &returned
}

func (d Double) EndPosition() *runtime.Position {
	var returned runtime.Position = d.DoubleToken.EndPosition.Copy()
	return &returned
}

func (d Double) Kind() int {
	return DoubleExpr
}

func (d Double) ExpressionConfirm() {

}

type String struct {
	Value       string
	StringToken *lexer.Token
}

func (s String) StartPosition() *runtime.Position {
	var returned runtime.Position = s.StringToken.StartPosition.Copy()
	return &returned
}

func (s String) EndPosition() *runtime.Position {
	var returned runtime.Position = s.StringToken.EndPosition.Copy()
	return &returned
}

func (s String) Kind() int {
	return StringExpr
}

func (s String) ExpressionConfirm() {

}

type Identifier struct {
	VariableName    string
	IdentifierToken *lexer.Token
}

func (i Identifier) StartPosition() *runtime.Position {
	var returned runtime.Position = i.IdentifierToken.StartPosition.Copy()
	return &returned
}

func (i Identifier) EndPosition() *runtime.Position {
	var returned runtime.Position = i.IdentifierToken.EndPosition.Copy()
	return &returned
}

func (i Identifier) Kind() int {
	return IdentifierExpr
}

func (i Identifier) ExpressionConfirm() {

}
