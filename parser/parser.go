package parser

import (
	"fmt"
	"os"
	"strconv"

	"github.com/nishokbanand/interpreter/ast"
	"github.com/nishokbanand/interpreter/lexer"
	"github.com/nishokbanand/interpreter/token"
)

var precedences = map[token.TokenType]int{
	token.EQ:          EQUALS,
	token.NOT_EQ:      EQUALS,
	token.LESSTHAN:    LESSGREATER,
	token.GREATERTHAN: LESSGREATER,
	token.SUM:         SUM,
	token.MINUS:       SUM,
	token.DIVIDE:      PRODUCT,
	token.ASTERISK:    PRODUCT,
	token.LPAREN:      CALL,
}

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

type (
	PrefixFns func() ast.ExpressionNode
	InfixFns  func(ast.ExpressionNode) ast.ExpressionNode
)

type Parser struct {
	l         *lexer.Lexer
	errros    []string
	currToken token.Token
	peekToken token.Token
	prefixfns map[token.TokenType]PrefixFns
	infixfns  map[token.TokenType]InfixFns
}

func (p *Parser) registerPrefixFns(tokType token.TokenType, preFn PrefixFns) {
	p.prefixfns[tokType] = preFn
}
func (p *Parser) registerInfixFns(tokType token.TokenType, infFn InfixFns) {
	p.infixfns[tokType] = infFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errros: []string{}}
	p.nextToken()
	p.nextToken()
	p.prefixfns = make(map[token.TokenType]PrefixFns)
	p.registerPrefixFns(token.IDENT, p.parseIdentifier)
	p.registerPrefixFns(token.INT, p.parseIntergerExpression)
	p.registerPrefixFns(token.NOT, p.parsePrefixExpression)
	p.registerPrefixFns(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixFns(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefixFns(token.IF, p.parseIfExpression)
	p.registerPrefixFns(token.FUNCTION, p.parseFunctionExpression)
	p.registerPrefixFns(token.TRUE, p.parseBoolean)
	p.registerPrefixFns(token.FALSE, p.parseBoolean)
	//infix
	p.infixfns = make(map[token.TokenType]InfixFns)
	p.registerInfixFns(token.SUM, p.parseInfixExpression)
	p.registerInfixFns(token.MINUS, p.parseInfixExpression)
	p.registerInfixFns(token.EQ, p.parseInfixExpression)
	p.registerInfixFns(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfixFns(token.DIVIDE, p.parseInfixExpression)
	p.registerInfixFns(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixFns(token.LESSTHAN, p.parseInfixExpression)
	p.registerInfixFns(token.GREATERTHAN, p.parseInfixExpression)
	p.registerInfixFns(token.LPAREN, p.parseCallExpression)
	return p
}

func (p *Parser) Errors() []string {
	return p.errros
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token %v , got token %v", t, p.peekToken.Type)
	p.errros = append(p.errros, msg)
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.StatmentNode{}
	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.StatmentNode {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() ast.StatmentNode {
	stmt := &ast.LetStatement{
		Token: p.currToken,
	}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() ast.StatmentNode {
	stmt := &ast.ReturnStatement{Token: p.currToken}
	//skipping expressions
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) expectPeek(expectedToken token.TokenType) bool {
	if p.peekToken.Type != expectedToken {
		p.peekError(expectedToken)
		return false
	}
	p.nextToken()
	return true
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token: p.currToken,
	}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return stmt

}

func (p *Parser) parseExpression(precedent int) ast.ExpressionNode {
	prefix := p.prefixfns[p.currToken.Type]
	if prefix == nil {
		msg := fmt.Sprintf("no prefix func found for %v", p.currToken.Type)
		p.errros = append(p.errros, msg)
		return nil
	}
	leftExp := prefix()
	for p.peekToken.Type != token.SEMICOLON && precedent < p.peekPrecedence() {
		infix := p.infixfns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseIdentifier() ast.ExpressionNode {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseIntergerExpression() ast.ExpressionNode {
	stmt := &ast.IntegerLiteral{
		Token: p.currToken,
	}
	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot convert %s to integer", p.currToken.Literal)
		p.errros = append(p.errros, err.Error())
	}
	stmt.Value = value
	return stmt
}

func (p *Parser) parsePrefixExpression() ast.ExpressionNode {
	stmt := &ast.PrefixExpression{Token: p.currToken, Operator: p.currToken.Literal}
	p.nextToken()
	stmt.Right = p.parseExpression(PREFIX)
	return stmt
}

func (p *Parser) parseInfixExpression(left ast.ExpressionNode) ast.ExpressionNode {
	stmt := &ast.InfixExpression{Token: p.currToken, Operator: p.currToken.Literal, Left: left}
	precedence := p.currPrecendence()
	p.nextToken()
	stmt.Right = p.parseExpression(precedence)
	return stmt
}

func (p *Parser) currPrecendence() int {
	if precedence, ok := precedences[p.currToken.Type]; ok {
		return precedence
	}
	return LOWEST
}
func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) parseGroupedExpression() ast.ExpressionNode {
	p.nextToken()
	stmt := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return stmt
}

func (p *Parser) parseIfExpression() ast.ExpressionNode {
	stmt := &ast.IfExpression{
		Token: p.currToken,
	}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	p.nextToken()
	if !p.expectPeek(token.LBRACES) {
		return nil
	}
	stmt.Consequence = p.parseBlockStatement()
	if p.peekToken.Type == token.ELSE {
		p.nextToken()
		if !p.expectPeek(token.LBRACES) {
			return nil
		}
		stmt.Alternative = p.parseBlockStatement()
	}
	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.currToken,
	}
	block.Statements = []ast.StatmentNode{}
	for p.currToken.Type != token.RPAREN && p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseFunctionExpression() ast.ExpressionNode {
	fnLiteral := &ast.FunctionLiteral{
		Token: p.currToken,
	}
	if p.expectPeek(token.LPAREN) {
		return nil
	}
	fnLiteral.Parameters = p.parseFunctionParameters()
	if p.expectPeek(token.LBRACES) {
		return nil
	}
	fnLiteral.Body = *p.parseBlockStatement()
	return fnLiteral
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	if p.expectPeek(token.RPAREN) {
		//so we skip )
		p.nextToken()
		return identifiers
	}
	identifier := &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
	p.nextToken()
	identifiers = append(identifiers, identifier)
	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		identifiers = append(identifiers, &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal})
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseCallExpression(function ast.ExpressionNode) ast.ExpressionNode {
	exp := &ast.CallExpression{
		Token:    p.currToken,
		Function: function,
	}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.ExpressionNode {
	args := []ast.ExpressionNode{}
	if p.expectPeek(token.RPAREN) {
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args

}

func (p *Parser) parseBoolean() ast.ExpressionNode {
	return &ast.Boolean{
		Token: p.currToken,
		Value: p.currToken.Type == token.TRUE,
	}
}
