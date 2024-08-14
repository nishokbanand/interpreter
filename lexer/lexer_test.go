package lexer

import (
	"testing"

	"github.com/nishokbanand/interpreter/token"
)

func TestLexerSimple(t *testing.T) {
	input := "-={}();"
	tests := []struct {
		Type    token.TokenType
		Literal string
	}{
		{Type: token.MINUS, Literal: "-"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.LBRACES, Literal: "{"},
		{Type: token.RBRACES, Literal: "}"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.EOF, Literal: ""},
	}
	l := New(input)
	for _, test := range tests {
		tok := l.nextToken()
		if test.Type != tok.Type {
			t.Errorf("expected tokenType %v, received tokenType %v", test.Type, tok.Type)
		}
		if test.Literal != tok.Literal {
			t.Errorf("expected Literal %v, received Literal %v", test.Literal, tok.Literal)
		}
	}
}

func TestLexerIntermediate(t *testing.T) {
	input := `let five = 5;
	let ten = 10;

	let add = fn(x,y){
		x+y;
	};
	let result = add(five,ten);
	`
	tests := []struct {
		Type    token.TokenType
		Literal string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACES, "{"},
		{token.IDENT, "x"},
		{token.SUM, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACES, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := New(input)
	for _, test := range tests {
		tok := l.nextToken()
		if test.Type != tok.Type {
			t.Errorf("expected tokenType %v, received tokenType %v", test.Type, tok.Type)
		}
		if test.Literal != tok.Literal {
			t.Errorf("expected Literal %v, received Literal %v", test.Literal, tok.Literal)
		}
	}
}
