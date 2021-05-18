//
// Created by irwin on 18/05/2021.
//
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "../header/lexer.h"

// isDigit
int isDigit(char chr) {
    return '0' <= chr && chr <= '9';
}
// isAlpha
int isAlpha(char chr) {
    return 'a' <= chr && chr <= 'z' || 'A' <= chr && chr <= 'Z';
}
// isIdent
int isIdent(char chr) {
    return isAlpha(chr) || isDigit(chr) || chr == '_';
}
// isSpace
int isSpace(char chr) {
    return chr == ' ' || chr == '\t' || chr == '\n';
}
// advance
void advance(Lexer* lexer) {
    lexer->pos += 1;
    if (lexer->pos > strlen(lexer->input)) {
        lexer->current_char = EOF;
    } else {
        lexer->current_char = lexer->input[lexer->pos];
    }
}
// skipWhitespace
void skipWhitespace(Lexer *lexer) {
    while (lexer->current_char != EOF && isSpace(lexer->current_char)) {
        advance(lexer);
    }
}
// getNumber
Token getNumber(Lexer *lexer) {
    int start_position = lexer->pos;
    while (lexer->current_char != EOF && isDigit(lexer->current_char)) {
        advance(lexer);
    }
    // check for floating point
    if (lexer->current_char == '.') {
        advance(lexer); // skip the '.'
        while (lexer->current_char != EOF && isDigit(lexer->current_char)) {
            advance(lexer);
        }
    }

    int nLength = lexer->pos - start_position;
    char *lexeme = malloc(nLength + 1);
    memcpy(lexeme, lexer->input + start_position, nLength);
    if (lexeme == NULL) {
        printf("Unavailable memory for use.\n");
        exit(1);
    }
    lexeme[nLength] = 0;

    return newToken(TT_NUMBER, lexeme);
}
// getString
Token getString(Lexer *lexer, char strDelim) {
    advance(lexer); // skip first delim
    int start_position = lexer->pos;
    while (lexer->current_char != EOF && lexer->current_char != strDelim) {
        advance(lexer);
    }

    int nLength = lexer->pos - start_position;
    char *lexeme = malloc(nLength + 1);
    if (lexeme == NULL) {
        printf("Unavailable memory for use.\n");
        exit(1);
    }
    memcpy(lexeme, lexer->input + start_position, nLength);
    lexeme[nLength] = '\0'; // close the array string.
    advance(lexer); // skip last delim
    return newToken(TT_STRING, lexeme);
}
// getIdentifier
Token getIdentifier(Lexer *lexer) {
    int start_position = lexer->pos;
    while (lexer->current_char != EOF && isIdent(lexer->current_char)) {
        advance(lexer);
    }
    int nLength = lexer->pos - start_position;
    char *lexeme = malloc(nLength + 1);
    if (lexeme == NULL) {
        printf("Unavailable memory for use.\n");
        exit(1);
    }
    // copy memory from input
    memcpy(lexeme, lexer->input + start_position, nLength);
    lexeme[nLength] = '\0'; // EOF character
    TokenType type = isKeyword(lexeme);

    return newToken(type, lexeme);
}
// nextToken
Token nextToken(Lexer *lexer) {
    while (lexer->current_char) {
        if (isSpace(lexer->current_char)) {
            skipWhitespace(lexer);
            continue;
        }
        if (isDigit(lexer->current_char))
            return getNumber(lexer);
        if (lexer->current_char == '\"' || lexer->current_char == '\'')
            return getString(lexer, lexer->current_char);
        if (isIdent(lexer->current_char))
            return getIdentifier(lexer);
        // 1 length characters
        if (lexer->current_char == '+') {
            advance(lexer);
            return newToken(TT_PLUS, "+");
        }
        if (lexer->current_char == '-') {
            advance(lexer);
            return newToken(TT_MINUS, "-");
        }
        if (lexer->current_char == '*') {
            advance(lexer);
            return newToken(TT_MUL, "*");
        }
        if (lexer->current_char == '/') {
            advance(lexer);
            return newToken(TT_DIV, "/");
        }
        if (lexer->current_char == ',') {
            advance(lexer);
            return newToken(TT_COMMA, ",");
        }
        if (lexer->current_char == ';') {
            advance(lexer);
            return newToken(TT_SEMICOLON, ";");
        }
        if (lexer->current_char == ':') {
            advance(lexer);
            return newToken(TT_COLON, ":");
        }
        if (lexer->current_char == '[') {
            advance(lexer);
            return newToken(TT_LBRACKET, "[");
        }
        if (lexer->current_char == ']') {
            advance(lexer);
            return newToken(TT_RBRACKET, "]");
        }
        if (lexer->current_char == '{') {
            advance(lexer);
            return newToken(TT_LBRACE, "{");
        }
        if (lexer->current_char == '}') {
            advance(lexer);
            return newToken(TT_RBRACE, "}");
        }
        if (lexer->current_char == '(') {
            advance(lexer);
            return newToken(TT_LPAREN, "(");
        }
        if (lexer->current_char == ')') {
            advance(lexer);
            return newToken(TT_RPAREN, ")");
        }
        if (lexer->current_char == '.') {
            advance(lexer);
            return newToken(TT_DOT, ".");
        }
        // double digit characters
        if (lexer->current_char == '<') {
            advance(lexer);
            if (lexer->current_char == '=') {
                advance(lexer);
                return newToken(TT_LESS_EQ, "<=");
            }
            return newToken(TT_LESS, "<");
        }
        if (lexer->current_char == '>') {
            advance(lexer);
            if (lexer->current_char == '=') {
                advance(lexer);
                return newToken(TT_GREATER_EQ, ">=");
            }
            return newToken(TT_GREATER, ">");
        }
        if (lexer->current_char == '!') {
            advance(lexer);
            if (lexer->current_char == '=') {
                advance(lexer);
                return newToken(TT_NOT_EQ, "!=");
            }
            return newToken(TT_NOT, "!");
        }
        if (lexer->current_char == '=') {
            advance(lexer);
            if (lexer->current_char == '=') {
                advance(lexer);
                return newToken(TT_EQUAL, "==");
            }
            return newToken(TT_ASSIGN, "=");
        }
        printf("lexer.c: unknown character '%c'", lexer->current_char);
        exit(1);
    }
    return newToken(TT_EOF, "");
}
// newLexer
Lexer *newLexer(const char* input) {
    Lexer *lexer = (Lexer*)malloc(sizeof(Lexer)); // casted to lexer
    if (!lexer) {
        printf("lexer.c: No available memory for use\n");
        exit(0);
    }
    // update the lexer pointer
    lexer->pos = 0;
    lexer->input = input;
    lexer->current_char = lexer->input[lexer->pos];

    return lexer;
}
// deleteLexer
void deleteLexer(Lexer *lexer) {
    free(lexer);
}