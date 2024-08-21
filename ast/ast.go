package ast

import (
	"bytes"

	"github.com/nishokbanand/interpreter/token"
)

type Node interface {
	TokenLiteral() string //gets the token literal
	String() string
}

type StatmentNode interface {
	Node
	statementNode()
	String() string
}
type ExpressionNode interface {
	Node
	expressionNode()
	String() string
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

func (p *Program) String() string {
	out := &bytes.Buffer{}
	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
	}
	return out.String()
}

type LetStatement struct {
	Token token.Token //this will have the LET token
	Name  *Identifier
	Value ExpressionNode
}

func (let *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(let.TokenLiteral() + " ")
	out.WriteString(let.Name.String())
	out.WriteString(" = ")
	if let.Value != nil {
		out.WriteString(let.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

func (let *LetStatement) TokenLiteral() string {
	return let.Token.Literal
}

func (let *LetStatement) statementNode() {}

type Identifier struct {
	Token token.Token //this will have the IDENT toke
	Value string
}

func (i *Identifier) String() string {
	var out bytes.Buffer
	out.WriteString(i.Value)
	return out.String()
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

func (r *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(r.TokenLiteral())
	if r.Value != nil {
		out.WriteString(r.Value.String())
	}
	return out.String()
}

func (r *ReturnStatement) TokenLiteral() string { return r.Token.Literal }
func (r *ReturnStatement) statementNode()       {}

type ExpressionStatement struct {
	Token      token.Token
	Expression ExpressionNode
}

func (e *ExpressionStatement) statementNode()       {}
func (e *ExpressionStatement) TokenLiteral() string { return e.Token.Literal }
func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct { //two prefix Exp are there ! and -
	Token    token.Token
	Operator string
	Right    ExpressionNode
}

func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Operator string
	Left     ExpressionNode
	Right    ExpressionNode
}

func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

type IfExpression struct {
	Token       token.Token //IF
	Condition   ExpressionNode
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ife *IfExpression) TokenLiteral() string { return ife.Token.Literal }
func (ife *IfExpression) expressionNode()      {}
func (ife *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ife.Token.Literal)
	out.WriteString("(")
	out.WriteString(ife.Condition.String())
	out.WriteString(")")
	out.WriteString(ife.Consequence.String())
	if ife.Alternative != nil {
		out.WriteString("else")
		out.WriteString(ife.Alternative.String())
	}
	return out.String()
}

type BlockStatement struct {
	Token      token.Token // { token
	Statements []StatmentNode
}

func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) expressionNode()      {}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
