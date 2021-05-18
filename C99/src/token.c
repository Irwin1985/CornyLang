//
// Created by irwin on 18/05/2021.
//
#include <string.h>
#include "../header/token.h"

// newToken
Token newToken(TokenType type, const char* literal) {
    Token token = {.type = type, .literal = literal};
    return token;
}
// isKeyword
TokenType isKeyword(const char *literal) {
    if (strcmp(literal, "let") == 0) {
        return TT_LET;
    }
    else if (strcmp(literal, "fn") == 0) {
        return TT_FUNCTION;
    }
    else if (strcmp(literal, "return") == 0) {
        return TT_RETURN;
    }
    else if (strcmp(literal, "and") == 0) {
        return TT_AND;
    }
    else if (strcmp(literal, "or") == 0) {
        return TT_OR;
    }
    else if (strcmp(literal, "true") == 0) {
        return TT_TRUE;
    }
    else if (strcmp(literal, "false") == 0) {
        return TT_FALSE;
    }
    else if (strcmp(literal, "null") == 0) {
        return TT_NULL;
    }
    return TT_IDENT;
}

