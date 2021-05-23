package lexer

import (
	"CornyLang/token"
	"fmt"
	"os"
)

type Lexer struct {
	input        string
	pos          int
	current_char byte
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input, pos: 0}
	lexer.current_char = lexer.input[lexer.pos]
	return lexer
}

func isLetter(chr byte) bool {
	return isAlpha(chr) || isDigit(chr) || chr == '_'
}

func isAlpha(chr byte) bool {
	return 'a' <= chr && chr <= 'z' || 'A' <= chr && chr <= 'Z'
}

func isSpace(chr byte) bool {
	return chr == ' ' || chr == '\t' || chr == '\n'
}

func isDigit(chr byte) bool {
	return '0' <= chr && chr <= '9'
}

func (l *Lexer) advance() {
	l.pos += 1
	if l.pos >= len(l.input) {
		l.current_char = 0
	} else {
		l.current_char = l.input[l.pos]
	}
}

func (l *Lexer) skipWhiteSpace() {
	for isSpace(l.current_char) {
		l.advance()
	}
}

func (l *Lexer) getNumber() token.Token {
	var start_position = l.pos

	for l.current_char != 0 && isDigit(l.current_char) {
		l.advance()
	}
	if l.current_char == '.' {
		l.advance()
		for l.current_char != 0 && isDigit(l.current_char) {
			l.advance()
		}
	}

	return token.Token{Type: token.TT_NUMBER, Literal: l.input[start_position:l.pos]}
}

func (l *Lexer) getString(strDelim byte) token.Token {
	l.advance() // skip the delimiter
	var start_position = l.pos
	for l.current_char != 0 && l.current_char != strDelim {
		l.advance()
	}
	var lexeme = l.input[start_position:l.pos]
	l.advance() // skip ending delimiter

	return token.Token{Type: token.TT_STRING, Literal: lexeme}
}

func (l *Lexer) getIdentifier() token.Token {
	var start_position = l.pos
	for l.current_char != 0 && isLetter(l.current_char) {
		l.advance()
	}
	var lexeme = l.input[start_position:l.pos]
	var tokenType = token.IsKeyword(lexeme)
	return token.Token{Type: tokenType, Literal: lexeme}
}

func (l *Lexer) NextToken() token.Token {
	for l.current_char != 0 {
		if isSpace(l.current_char) {
			l.skipWhiteSpace()
			continue
		}
		if isDigit(l.current_char) {
			return l.getNumber()
		}
		if isLetter(l.current_char) {
			return l.getIdentifier()
		}
		if l.current_char == '"' || l.current_char == '\'' {
			return l.getString(l.current_char)
		}
		if l.current_char == '+' {
			l.advance()
			return newToken(token.TT_PLUS, "+")
		}
		if l.current_char == '-' {
			l.advance()
			return newToken(token.TT_MINUS, "-")
		}
		if l.current_char == '*' {
			l.advance()
			return newToken(token.TT_MUL, "*")
		}
		if l.current_char == '/' {
			l.advance()
			return newToken(token.TT_DIV, "/")
		}
		if l.current_char == ',' {
			l.advance()
			return newToken(token.TT_COMMA, ",")
		}
		if l.current_char == ';' {
			l.advance()
			return newToken(token.TT_SEMICOLON, ";")
		}
		if l.current_char == ':' {
			l.advance()
			return newToken(token.TT_COLON, ":")
		}
		if l.current_char == '[' {
			l.advance()
			return newToken(token.TT_LBRACKET, "[")
		}
		if l.current_char == ']' {
			l.advance()
			return newToken(token.TT_RBRACKET, "]")
		}
		if l.current_char == '{' {
			l.advance()
			return newToken(token.TT_LBRACE, "{")
		}
		if l.current_char == '}' {
			l.advance()
			return newToken(token.TT_RBRACE, "}")
		}
		if l.current_char == '(' {
			l.advance()
			return newToken(token.TT_LPAREN, "(")
		}
		if l.current_char == ')' {
			l.advance()
			return newToken(token.TT_RPAREN, ")")
		}
		if l.current_char == '.' {
			l.advance()
			return newToken(token.TT_DOT, ".")
		}
		if l.current_char == '?' {
			l.advance()
			return newToken(token.TT_QUESTION, "?")
		}
		// double digit characters
		if l.current_char == '<' {
			l.advance()
			if l.current_char == '=' {
				l.advance()
				return newToken(token.TT_LESS_EQ, "<=")
			}
			return newToken(token.TT_LESS, "<")
		}
		if l.current_char == '>' {
			l.advance()
			if l.current_char == '=' {
				l.advance()
				return newToken(token.TT_GREATER_EQ, ">=")
			}
			return newToken(token.TT_GREATER, ">")
		}
		if l.current_char == '!' {
			l.advance()
			if l.current_char == '=' {
				l.advance()
				return newToken(token.TT_NOT_EQ, "!=")
			}
			return newToken(token.TT_NOT, "!")
		}
		if l.current_char == '=' {
			l.advance()
			if l.current_char == '=' {
				l.advance()
				return newToken(token.TT_EQUAL, "==")
			}
			return newToken(token.TT_ASSIGN, "=")
		}
		fmt.Printf("unknown character %c", l.current_char)
		os.Exit(4)
	}
	return token.Token{Type: token.TT_EOF, Literal: ""}
}

func newToken(ttype token.TokenType, lexeme string) token.Token {
	return token.Token{Type: ttype, Literal: lexeme}
}
