package token

type TokenType string //

type Token struct {
	Type    TokenType
	Literal string
}

const (
	IDENT = "IDENT"
	INT   = "INT"
	//keywords
	LET      = "LET"
	FUNCTION = "FUNCTION"
	//SYMBOLS
	LPAREN    = "("
	RPAREN    = ")"
	SEMICOLON = ";"
	LBRACES   = "{"
	RBRACES   = "}"
	COMMA     = ","
	//OPERATORS
	SUM      = "+"
	MINUS    = "-"
	ASSIGN   = "="
	MULTIPLY = "*"
	DIVIDE   = "/"
	//MICEL
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"
)

var Keywords = map[string]TokenType{ // maps cannot be created as const
	"fn":  FUNCTION,
	"let": LET,
}
