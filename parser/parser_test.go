package parser

import (
	"fmt"
	"testing"

	"github.com/nishokbanand/interpreter/ast"
	"github.com/nishokbanand/interpreter/lexer"
)

func TestLetStatments(t *testing.T) {
	input := `let x = 5;
	let y = 20;
	let z = 123123;
	`
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"z"},
	}
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)
	if program == nil {
		t.Fatalf("parseProgram returned nil")
	}
	fmt.Printf("%#v\n", program)
	if len(program.Statements) != 3 {
		t.Fatalf("expected number of statements :3 , got %d", len(program.Statements))
	}
	for i, test := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, test.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.StatmentNode, expectedIdentifier string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("TokenLiteral expected : let, got %s", stmt.TokenLiteral())
		return false
	}
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("Expected : LetStatement, got %T", stmt)
		return false
	}
	if letStmt.Name.Value != expectedIdentifier {
		t.Errorf("expectedIdentifier : %s , got %s", expectedIdentifier, letStmt.Name.Value)
		return false
	}
	return true
}

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has errors : %d\n", len(errors))
	for i, msg := range errors {
		t.Errorf("Error %d : %s", i, msg)
	}
	t.FailNow()
}
func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10000;
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)
	if program == nil {
		t.Fatalf("parseProgram returned nil")
	}
	fmt.Printf("%#v\n", program)
	if len(program.Statements) != 2 {
		t.Fatalf("expected number of statements :2 , got %d", len(program.Statements))
	}
	for _, stmt := range program.Statements {
		if returnStmt, ok := stmt.(*ast.ReturnStatement); !ok {
			t.Fatalf("expected ReturnStatement, got :%T", returnStmt)
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		}, {
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5*5",
			"(3 + 4)((-5) * 5)",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
