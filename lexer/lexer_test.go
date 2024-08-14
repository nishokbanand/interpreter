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
		tok := l.NextToken()
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
		tok := l.NextToken()
		if test.Type != tok.Type {
			t.Errorf("expected tokenType %v, received tokenType %v", test.Type, tok.Type)
		}
		if test.Literal != tok.Literal {
			t.Errorf("expected Literal %v, received Literal %v", test.Literal, tok.Literal)
		}
	}
}

func TestLexerAdvanced(t *testing.T) {
	input := `
	!-/*5;
	5 < 10 > 5
	if (5 < 10){
		return true;
	}else {
		return false;
	}
	10 == 10;
	10 != 9;
	`
	tests := []struct {
		Type    token.TokenType
		Literal string
	}{
		{token.NOT, "!"},
		{token.MINUS, "-"},
		{token.DIVIDE, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LESSTHAN, "<"},
		{token.INT, "10"},
		{token.GREATERTHAN, ">"},
		{token.INT, "5"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LESSTHAN, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACES, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACES, "}"},
		{token.ELSE, "else"},
		{token.LBRACES, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACES, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
	}
	l := New(input)
	for _, test := range tests {
		tok := l.NextToken()
		if test.Type != tok.Type {
			t.Errorf("expected tokenType %v, received tokenType %v", test.Type, tok.Type)
		}
		if test.Literal != tok.Literal {
			t.Errorf("expected Literal %v, received Literal %v", test.Literal, tok.Literal)
		}
	}
}
