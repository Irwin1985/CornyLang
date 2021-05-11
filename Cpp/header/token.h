//
// Created by irwin on 11/05/2021.
//

#ifndef CPP_TOKEN_H
#define CPP_TOKEN_H
#include <map>

namespace corny {
    const char NONE = '\0';
    enum TokenType {
        END,
        // Identifiers + literals
        IDENT, // add, foobar, x, y, ...
        NUMBER,
        STRING,
        BOOLEAN,
        NIL,

        // Operators
        ASSIGN,

        // Relational operators
        LESS,
        LESS_EQ,
        GREATER,
        GREATER_EQ,
        EQUAL,
        NOT_EQ,

        // Arithmetic operators
        PLUS,
        MINUS,
        MUL,
        DIV,
        POW,

        // Logical Operators
        AND,
        OR,
        NOT,

        // Delimiters
        COMMA,
        COLON,
        SEMICOLON,
        DOT,
        QUESTION,

        // Special characters
        LPAREN,
        RPAREN,
        LBRACE,
        RBRACE,
        LBRACKET,
        RBRACKET,

        // keywords
        FUNCTION,
        LET,
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
