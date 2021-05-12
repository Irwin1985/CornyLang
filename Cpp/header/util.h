//
// Created by irwin on 11/05/2021.
//

#ifndef CPP_UTIL_H
#define CPP_UTIL_H
#include "token.h"

namespace corny {
    std::string getTokenStr(TokenType type) {
        switch (type) {
            case TT_EOF:
                return "TT_EOF";
            case TT_IDENT:
                return "TT_IDENT";
            case TT_NUMBER:
                return "TT_NUMBER";
            case TT_STRING:
                return "TT_STRING";
            case TT_TRUE:
                return "TT_TRUE";
            case TT_FALSE:
                return "TT_FALSE";
            case TT_NULL:
                return "TT_NULL";
            case TT_ASSIGN:
                return "TT_ASSIGN";
            case TT_LESS:
                return "TT_LESS";
            case TT_LESS_EQ:
                return "TT_LESS_EQ";
            case TT_GREATER:
                return "TT_GREATER";
            case TT_GREATER_EQ:
                return "TT_GREATER_EQ";
            case TT_EQUAL:
                return "TT_EQUAL";
            case TT_NOT_EQ:
                return "TT_NOT_EQ";
            case TT_PLUS:
                return "TT_PLUS";
            case TT_MINUS:
                return "TT_MINUS";
            case TT_MUL:
                return "TT_MUL";
            case TT_DIV:
                return "TT_DIV";
            case TT_POW:
                return "TT_POW";
            case TT_AND:
                return "TT_AND";
            case TT_OR:
                return "TT_OR";
            case TT_NOT:
                return "TT_NOT";
            case TT_DOT:
                return "TT_DOT";
            case TT_QUESTION:
                return "TT_QUESTION";
            case TT_COMMA:
                return "TT_COMMA";
            case TT_COLON:
                return "TT_COLON";
            case TT_SEMICOLON:
                return "TT_SEMICOLON";
            case TT_LPAREN:
                return "TT_LPAREN";
            case TT_RPAREN:
                return "TT_RPAREN";
            case TT_LBRACE:
                return "TT_LBRACE";
            case TT_RBRACE:
                return "TT_RBRACE";
            case TT_LBRACKET:
                return "TT_LBRACKET";
            case TT_RBRACKET:
                return "TT_RBRACKET";
            case TT_FUNCTION:
                return "TT_FUNCTION";
            case TT_LET:
                return "TT_LET";
            case TT_RETURN:
                return "TT_RETURN";
            case TT_IF:
                return "TT_IF";
            case TT_ELSE:
                return "TT_ELSE";
            default:
                return "ILLEGAL";
        }
    }
}

#endif //CPP_UTIL_H
