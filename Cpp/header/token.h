//
// Created by irwin on 11/05/2021.
//

#ifndef CPP_TOKEN_H
#define CPP_TOKEN_H
#include <map>

namespace corny {
    const char NONE = '\0';
    enum TokenType {
        TT_EOF,
        // Identifiers + literals
        TT_IDENT, // add, foobar, x, y, ...
        TT_NUMBER,
        TT_STRING,
        TT_TRUE,
        TT_FALSE,
        TT_NULL,

        // Operators
        TT_ASSIGN,

        // Relational operators
        TT_LESS,
        TT_LESS_EQ,
        TT_GREATER,
        TT_GREATER_EQ,
        TT_EQUAL,
        TT_NOT_EQ,

        // Arithmetic operators
        TT_PLUS,
        TT_MINUS,
        TT_MUL,
        TT_DIV,
        TT_POW,

        // Logical Operators
        TT_AND,
        TT_OR,
        TT_NOT,

        // Delimiters
        TT_COMMA,
        TT_COLON,
        TT_SEMICOLON,
        TT_DOT,
        TT_QUESTION,

        // Special characters
        TT_LPAREN,
        TT_RPAREN,
        TT_LBRACE,
        TT_RBRACE,
        TT_LBRACKET,
        TT_RBRACKET,

        // keywords
        TT_FUNCTION,
        TT_LET,
        TT_RETURN,
        TT_IF,
        TT_ELSE,
    };
    // keywords dictionary: we need to be able to identify an identifier from a keyword.
    extern std::map<std::string, TokenType> keywords;

    /**
     * Token class: the building block of every programming language.
     */
    class Token {
    public:
        Token() {};
        Token(TokenType type, std::string literal) {
            this->type = type;
            this->literal = literal;
        }
        Token(TokenType type, char literal) {
            this->type = type;
            this->literal = std::string(1, literal);
        }
        ~Token(){};

        TokenType type;
        std::string literal;
    };

    // isKeyword: check if 'ident' is a keyword or a ident token.
    TokenType isKeyword(std::string ident);
}

#endif //CPP_TOKEN_H
