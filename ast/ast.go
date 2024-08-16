package ast

import "github.com/nishokbanand/interpreter/token"

type Node interface {
	TokenLiteral() string //gets the token literal
}

type StatmentNode interface {
	Node
	statementNode()
}
type ExpressionNode interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []StatmentNode
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) != 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

type LetStatement struct {
	Token token.Token //this will have the LET token
	Name  *Identifier
	Value *ExpressionNode
}

func (let *LetStatement) TokenLiteral() string {
	return let.Token.Literal
}

func (let *LetStatement) statementNode() {}

type Identifier struct {
	Token token.Token //this will have the IDENT toke
	Value string
}

// we say the identifier is a type of expression to make it easy as identifiers could have value producing expressions.

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) expressionNode() {}

type ReturnStatement struct {
	Token token.Token
	Value ExpressionNode
}

func (r *ReturnStatement) TokenLiteral() string { return r.Token.Literal }
func (r *ReturnStatement) statementNode()       {}
