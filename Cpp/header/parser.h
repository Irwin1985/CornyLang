//
// Created by irwin on 12/05/2021.
//

#ifndef CPP_PARSER_H
#define CPP_PARSER_H
#include "lexer.h"
#include "ast.h"

namespace corny {
    class Parser {
    public:
        Parser() {}
        Lexer lexer;
        Token curToken;
        Token peekToken;
        // methods
        void start(Lexer lexer);
        void advance(TokenType type);
        void nextToken();
        Node* parseProgram();
        Node* parseBlock();
        Node* parseStatement();
        Node* parseLetStatement();
        Node* parseReturnStatement();
        Node* parseExpression();
        Node* parseAssignment();
        Node* parseLogicOr();
        Node* parseLogicAnd();
        Node* parseEquality();
        Node* parseComparison();
        Node* parseTerm();
        Node* parseFactor();
        Node* parseUnary();
        Node* parseCall();
        Node* parsePrimary();
        Node* parseFunctionLiteral();
        Node* parseArrayLiteral();
        Node* parseHashLiteral();
        Node* parseIfExpr();
        Node* parseTernaryExpr(Node* condition);
        Node* parseCallExpr(Node* callee, Token token);
        Node* parseIdentifier();
    };
}

#endif //CPP_PARSER_H
