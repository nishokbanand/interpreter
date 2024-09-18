package ast

import (
	"bytes"
	"strings"

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
	out.WriteString("{")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString("}")
	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token //fn
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fn *FunctionLiteral) TokenLiteral() string { return fn.Token.Literal }
func (fn *FunctionLiteral) expressionNode()      {}
func (fn *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, param := range fn.Parameters {
		params = append(params, param.String())
	}
	out.WriteString(fn.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString(fn.Body.String())
	return out.String()
}

// add (5,6) or hello(5,func(5,4))
type CallExpression struct {
	Token     token.Token
	Function  ExpressionNode
	Arguments []ExpressionNode
}

func (c *CallExpression) TokenLiteral() string { return c.Token.Literal }

func (c *CallExpression) expressionNode() {}

func (c *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, arg := range c.Arguments {
		args = append(args, arg.String())
	}
	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ","))
	out.WriteString(")")
	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

func (b *Boolean) expressionNode() {}

func (b *Boolean) String() string {
	return b.Token.Literal
}

type String struct {
	Token token.Token
	Value string
}

func (s *String) TokenLiteral() string { return s.Token.Literal }
func (s *String) expressionNode()      {}

func (s *String) String() string { return s.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token
	Elements []ExpressionNode
}

func (arr *ArrayLiteral) TokenLiteral() string { return arr.Token.Literal }
func (arr *ArrayLiteral) expressionNode()      {}
func (arr *ArrayLiteral) String() string {
	var out bytes.Buffer
	var values []string
	for _, ele := range arr.Elements {
		values = append(values, ele.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(values, ","))
	out.WriteString("]")
	return out.String()
}

type IndexExpression struct {
	Token token.Token //[ token
	Left  ExpressionNode
	Index ExpressionNode
}

func (i *IndexExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IndexExpression) expressionNode()      {}
func (i *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString("[")
	out.WriteString(i.Index.String())
	out.WriteString("]")
	out.WriteString(")")
	return out.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[ExpressionNode]ExpressionNode
}

func (h *HashLiteral) TokenLiteral() string { return h.Token.Literal }
func (h *HashLiteral) expressionNode()      {}
func (h *HashLiteral) String() string {
	var out bytes.Buffer
	var elements []string
	for key, value := range h.Pairs {
		elements = append(elements, key.String()+":"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(elements, ","))
	out.WriteString("}")
	return out.String()
}
