package token

type TokenType string //

type Token struct {
	Type    TokenType
	Literal string
}

const (
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"
	//keywords
	LET      = "LET"
	FUNCTION = "FUNCTION"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	//SYMBOLS
	LPAREN    = "("
	RPAREN    = ")"
	SEMICOLON = ";"
	LBRACES   = "{"
	RBRACES   = "}"
	COMMA     = ","
	//OPERATORS
	SUM         = "+"
	MINUS       = "-"
	ASSIGN      = "="
	EQ          = "=="
	NOT         = "!"
	NOT_EQ      = "!="
	ASTERISK    = "*"
	DIVIDE      = "/"
	LESSTHAN    = "<"
	GREATERTHAN = ">"
	//MICEL
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"
)

var Keywords = map[string]TokenType{ // maps cannot be created as const
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}
