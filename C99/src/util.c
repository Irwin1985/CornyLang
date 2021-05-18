//
// Created by irwin on 18/05/2021.
//
#include "../header/util.h"

const char *getTokenType(TokenType type) {
    switch (type) {
        case TT_EOF:
            return "EOF";
        case TT_STRING:
            return "STRING";
        case TT_NUMBER:
            return "NUMBER";
        case TT_DIV:
            return "DIV";
        case TT_MUL:
            return "MUL";
        case TT_MINUS:
            return "MINUS";
        case TT_PLUS:
            return "PLUS";
        case TT_IDENT:
            return "IDENT";
        case TT_FUNCTION:
            return "FUNCTION";
        case TT_LET:
            return "LET";
        case TT_RETURN:
            return "RETURN";
        case TT_COMMA:
            return "COMMA";
        case TT_DOT:
            return "COMMA";
        case TT_COLON:
            return "COLON";
        case TT_SEMICOLON:
            return "SEMICOLON";
        case TT_LPAREN:
            return "LPAREN";
        case TT_RPAREN:
            return "RPAREN";
        case TT_LBRACE:
            return "LBRACE";
        case TT_RBRACE:
            return "RBRACE";
        case TT_LBRACKET:
            return "LBRACKET";
        case TT_RBRACKET:
            return "RBRACKET";
        case TT_AND:
            return "AND";
        case TT_OR:
            return "OR";
        case TT_NOT:
            return "NOT";
        case TT_LESS:
            return "LESS";
        case TT_GREATER:
            return "GREATER";
        case TT_LESS_EQ:
            return "LESS_EQ";
        case TT_GREATER_EQ:
            return "GREATER_EQ";
        case TT_ASSIGN:
            return "ASSIGN";
        case TT_EQUAL:
            return "EQUAL";
        case TT_NOT_EQ:
            return "NOT_EQ";
        case TT_TRUE:
            return "TRUE";
        case TT_FALSE:
            return "FALSE";
        case TT_NULL:
            return "NULL";
        default:
            return "UNKNOWN";
    }
}

