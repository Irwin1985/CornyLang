package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

// keywords
var keywords = map[string]TokenType{
	"fn":     TT_FUNCTION,
	"let":    TT_LET,
	"true":   TT_TRUE,
	"false":  TT_FALSE,
	"and":    TT_AND,
	"null":   TT_NULL,
	"return": TT_RETURN,
	"if":     TT_IF,
	"else":   TT_ELSE,
}

func IsKeyword(key string) TokenType {
	if tok, ok := keywords[key]; ok {
		return tok
	} else {
		return TT_IDENT
	}
}

const (
	TT_EOF = "EOF"

	// Identifier + literal
	TT_IDENT  = "IDENT"
	TT_NUMBER = "NUMBER"
	TT_STRING = "STRING"

	// Operators
	TT_ASSIGN = "="

	// Relational operators
	TT_LESS       = "<"
	TT_LESS_EQ    = "<="
	TT_GREATER    = ">"
	TT_GREATER_EQ = ">="
	TT_EQUAL      = "=="
	TT_NOT_EQ     = "!="

	// Arithmetic operators
	TT_PLUS  = "+"
	TT_MINUS = "-"
	TT_MUL   = "*"
	TT_DIV   = "/"
	TT_POW   = "^"

	// Logical operators
	TT_AND = "and"
	TT_OR  = "or"
	TT_NOT = "!"

	// Delimiters
	TT_COMMA     = ","
	TT_COLON     = ":"
	TT_SEMICOLON = ";"
	TT_DOT       = "."
	TT_QUESTION  = "?"

	// Special characters
	TT_LPAREN   = "("
	TT_RPAREN   = ")"
	TT_LBRACE   = "{"
	TT_RBRACE   = "}"
	TT_LBRACKET = "["
	TT_RBRACKET = "]"

	// Keywords
	TT_FUNCTION = "FUNCTION"
	TT_LET      = "LET"
	TT_RETURN   = "RETURN"
	TT_IF       = "IF"
	TT_ELSE     = "ELSE"
	TT_TRUE     = "TRUE"
	TT_FALSE    = "FALSE"
	TT_NULL     = "NULL"
)
