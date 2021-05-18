//
// Created by irwin on 18/05/2021.
//

#ifndef C99_TOKEN_H
#define C99_TOKEN_H

typedef enum {
    TT_EOF,
    TT_PLUS,
    TT_MINUS,
    TT_MUL,
    TT_DIV,
    TT_IDENT,
    TT_NUMBER,
    TT_STRING,

    // punctuation symbol
    TT_COMMA,
    TT_DOT,
    TT_COLON,
    TT_SEMICOLON,
    TT_LPAREN,
    TT_RPAREN,
    TT_LBRACE,
    TT_RBRACE,
    TT_LBRACKET,
    TT_RBRACKET,

    // relational operator
    TT_NOT,
    TT_LESS,
    TT_GREATER,
    TT_LESS_EQ,
    TT_GREATER_EQ,
    TT_ASSIGN,
    TT_EQUAL,
    TT_NOT_EQ,

    // keywords
    TT_LET,
    TT_FUNCTION,
    TT_RETURN,
    TT_AND,
    TT_OR,
    TT_TRUE,
    TT_FALSE,
    TT_NULL,

} TokenType;

// create the '_token' type first: note the final 'Token'
typedef struct token_struct Token;
// then we can create the struct type using above custom type 'type_token'
struct token_struct {
    TokenType type;
    void *literal; // void pointer (this could point to any type)
};
Token newToken(TokenType type, const char* literal);
TokenType isKeyword(const char* literal);
#endif //C99_TOKEN_H
