//
// Created by irwin on 18/05/2021.
//

#ifndef C99_LEXER_H
#define C99_LEXER_H
#include "token.h"

typedef struct struct_lexer Lexer;
struct struct_lexer {
    int pos;
    const char *input;
    char current_char;
};

int isDigit(char chr);
int isAlpha(char chr);
int isIdent(char chr);
int isSpace(char chr);
void advance(Lexer* lexer);
void skipWhitespace(Lexer *lexer);
Token getNumber(Lexer *lexer);
Token getString(Lexer *lexer, char strDelim);
Token getIdentifier(Lexer *lexer);
Token nextToken(Lexer *lexer);
Lexer *newLexer(const char* input); // it could be like constructor
void deleteLexer(Lexer* lexer); // and this could be like destructor

#endif //C99_LEXER_H
