//
// Created by irwin on 11/05/2021.
//

#ifndef CPP_UTIL_H
#define CPP_UTIL_H
#include "token.h"

namespace corny {
    std::string getTokenStr(TokenType type) {
        switch (type) {
            case END:
                return "END";
            case IDENT:
                return "IDENT";
            case NUMBER:
                return "NUMBER";
            case STRING:
                return "STRING";
            case BOOLEAN:
                return "BOOLEAN";
            case NIL:
                return "NIL";
            case ASSIGN:
                return "ASSIGN";
            case LESS:
                return "LESS";
            case LESS_EQ:
                return "LESS_EQ";
            case GREATER:
                return "GREATER";
            case GREATER_EQ:
                return "GREATER_EQ";
            case EQUAL:
                return "EQUAL";
            case NOT_EQ:
                return "NOT_EQ";
            case PLUS:
                return "PLUS";
            case MINUS:
                return "MINUS";
            case MUL:
                return "MUL";
            case DIV:
                return "DIV";
            case POW:
                return "POW";
            case AND:
                return "AND";
            case OR:
                return "OR";
            case NOT:
                return "NOT";
            case DOT:
                return "DOT";
            case QUESTION:
                return "QUESTION";
            case COMMA:
                return "COMMA";
            case COLON:
                return "COLON";
            case SEMICOLON:
                return "SEMICOLON";
            case LPAREN:
                return "LPAREN";
            case RPAREN:
                return "RPAREN";
            case LBRACE:
                return "LBRACE";
            case RBRACE:
                return "RBRACE";
            case LBRACKET:
                return "LBRACKET";
            case RBRACKET:
                return "RBRACKET";
            case FUNCTION:
                return "FUNCTION";
            case LET:
                return "LET";
            default:
                return "ILLEGAL";
        }
    }
}

#endif //CPP_UTIL_H
