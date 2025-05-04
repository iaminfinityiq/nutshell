package parser

import (
	"nutshell/frontend/lexer"
	"nutshell/runtime"
)

type Block struct {
	Body []*Statement
}

func (b Block) StartPosition() *runtime.Position {
	if len(b.Body) == 0 {
		return nil
	}

	var returned runtime.Position = (*b.Body[0]).StartPosition().Copy()
	return &returned
}

func (b Block) EndPosition() *runtime.Position {
	if len(b.Body) == 0 {
		return nil
	}

	var returned runtime.Position = (*b.Body[len(b.Body)-1]).StartPosition().Copy()
	return &returned
}

func (b Block) Kind() int {
	return BlockStmt
}

func InitBlock() *Block {
	var body []*Statement = []*Statement{}
	var returned Block = Block{
		Body: body,
	}

	return &returned
}

type VariableDeclaration struct {
	LetToken *lexer.Token
	VariableName string
	Value *Expression
	IsConstant bool
}

func (v VariableDeclaration) StartPosition() *runtime.Position {
	return v.LetToken.StartPosition
}

func (v VariableDeclaration) EndPosition() *runtime.Position {
	return (*v.Value).EndPosition()
}

func (v VariableDeclaration) Kind() int {
	return VariableDeclarationStmt
}
