package token

// String is used as TokenType to easier for debug and distinguish
type TokenType string

const (
	EOF = "EOF" // End of file

	ILLEGAL = "ILLEGAL" // Signifies unknown token

	// Identifier  and  Literal
	IDENT  = "INDENT" // Identifier like foo , bar
	INT    = "INT"    // Interger literal like 1 , 2 , 3
	STRING = "STRING"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	BANG     = "!"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keyword
	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	EQ       = "=="
	NOT_EQ   = "!="
	LT_EQ    = "<="
	GT_EQ    = ">="
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// Look up the keyword table , if Identifier indeed a keyword , then return the keyword constant
// other while IDENT constant
func LookUpKeyword(ident string) TokenType {
	if token, ok := keywords[ident]; ok {
		return token
	}
	return IDENT
}
