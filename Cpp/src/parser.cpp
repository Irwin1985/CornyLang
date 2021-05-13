//
// Created by irwin on 12/05/2021.
//
#include "../header/parser.h"

namespace corny {
    void Parser::advance(TokenType type) {
        if (curToken.type == type){
            nextToken();
            return;
        }
        std::cout << "Lexer Error: Unexpected token. Got: " << curToken.type << ", want: " << type << std::endl;
    }
    void Parser::nextToken() {
        curToken = peekToken;
        peekToken = lexer.nextToken();
    }
    // program ::= ( statement )*
    Node* Parser::parseProgram() {
        ProgramNode* program = new ProgramNode();

        while (curToken.type != TT_EOF) {
            program->statements.emplace_back(parseStatement());
        }

        return program;
    }
    // parseBlock ::= '{' (statement)* '}'
    Node* Parser::parseBlock() {
        BlockNode* blockNode = new BlockNode();

        advance(TT_LBRACE);
        while (curToken.type != TT_RBRACE) {
            blockNode->statements.emplace_back(parseStatement());
        }
        advance(TT_RBRACE);

        return blockNode;
    }
    // parseStatement ::= parseLetStatement | parseReturnStatement | parseExpression
    Node* Parser::parseStatement() {
        Node* statement;
        if (curToken.type == TT_LET) {
            statement = parseLetStatement();
        }
        else if (curToken.type == TT_RETURN) {
            statement = parseReturnStatement();
        }
        else {
            statement = parseExpression();
        }
        // optionally skip the semicolon
        if (curToken.type == TT_SEMICOLON) {
            advance(TT_SEMICOLON);
        }
        return statement;
    }
    // parseLetStatement ::= 'let' IDENT '=' parseExpression
    Node* Parser::parseLetStatement() {
        LetNode* letNode = new LetNode();

        advance(TT_LET);

        letNode->ident = (IdentNode*)parseIdentifier();
        advance(TT_ASSIGN);
        letNode->value = parseExpression();

        return letNode;
    }
    // parseReturnStatement ::= 'return' parseExpression
    Node* Parser::parseReturnStatement() {
        ReturnNode* returnNode = new ReturnNode();

        advance(TT_RETURN);
        returnNode->value = parseExpression();

        return returnNode;
    }
    // parseExpression ::= parseAssignment
    Node* Parser::parseExpression() {
        return parseAssignment();
    }
    // parseAssignment ::= parseLogicOr
    Node* Parser::parseAssignment() {
        return parseLogicOr();
    }
    // parseLogicOr ::= parseLogicAnd ('and' parseLogicAnd)*
    Node* Parser::parseLogicOr() {
        Node* node = parseLogicAnd();
        while (curToken.type == TT_OR) {
            Token token = curToken;
            advance(TT_OR);
            node = new BinOpNode(node, token, parseLogicAnd());
        }
        return node;
    }
    // parseLogicAnd ::= parseEquality ('and' parseEquality)*
    Node* Parser::parseLogicAnd() {
        Node* node = parseEquality();
        while (curToken.type == TT_AND) {
            Token token = curToken;
            advance(TT_AND);
            node = new BinOpNode(node, token, parseEquality());
        }
        return node;
    }
    // parseEquality ::= parseComparison (('==', '!=') parseComparison)*
    Node* Parser::parseEquality() {
        Node* node = parseComparison();
        while (curToken.type == TT_EQUAL || curToken.type == TT_NOT_EQ) {
            Token token = curToken;
            advance(token.type);
            node = new BinOpNode(node, token, parseComparison());
        }
        return node;
    }
    // parseComparison ::= parseTerm (('<'|'<='|'>'|'>=') parseTerm)*
    Node* Parser::parseComparison() {
        Node* node = parseTerm();
        while (curToken.type == TT_LESS || curToken.type == TT_LESS_EQ ||
        curToken.type == TT_GREATER || curToken.type == TT_GREATER_EQ)
        {
            Token token = curToken;
            advance(token.type);
            node = new BinOpNode(node, token, parseTerm());
        }
        return node;
    }
    // parseTerm ::= parseFactor (('+'|'-') parseFactor)*
    Node* Parser::parseTerm() {
        Node* node = parseFactor();
        while (curToken.type == TT_PLUS || curToken.type == TT_MINUS) {
            Token token = curToken;
            advance(token.type);
            node = new BinOpNode(node, token, parseFactor());
        }
        return node;
    }
    // parseFactor ::= parseUnary (('*'|'/' parseUnary)*
    Node* Parser::parseFactor() {
        Node* node = parseUnary();
        while (curToken.type == TT_MUL || curToken.type == TT_DIV) {
            Token token = curToken;
            advance(token.type);
            node = new BinOpNode(node, token, parseUnary());
        }
        return node;
    }
    // parseUnary ::= ('-'|'!') parseUnary | parseCall
    Node* Parser::parseUnary() {
        while (curToken.type == TT_MINUS || curToken.type == TT_NOT) {
            Token token = curToken;
            advance(token.type);
            return new UnaryNode(token, parseUnary()); // call itself
        }
        return parseCall();
    }
    // parseCall ::= parsePrimary | parseCallExpression
    Node* Parser::parseCall() {
        Node* node = parsePrimary();
        while (true) {
            if (curToken.type == TT_LPAREN || curToken.type == TT_LBRACKET) {
                node = parseCallExpr(node, curToken);
            }
            else if (curToken.type == TT_QUESTION) {
                node = parseTernaryExpr(node);
            }
            else {
                break;
            }
        }
        return node;
    }
    // parsePrimary ::= NUMBER | TRUE | FALSE | STRING | NULL | IDENT | FUNCTION | IF-ELSE
    Node* Parser::parsePrimary() {
        Token token = curToken;
        switch (token.type) {
            case TT_NUMBER:
                advance(TT_NUMBER);
                return new NumberNode(std::stod(token.literal));
            case TT_TRUE:
                advance(TT_TRUE);
                return new BooleanNode(true);
            case TT_FALSE:
                advance(TT_FALSE);
                return new BooleanNode(false);
            case TT_NULL:
                advance(TT_NULL);
                return new NullNode();
            case TT_IDENT:
                return parseIdentifier();
            case TT_LPAREN: {
                advance(TT_LPAREN);
                Node* resultExpr = parseExpression();
                advance(TT_RPAREN);
                return resultExpr;
            }
            case TT_STRING:
                advance(TT_STRING);
                return new StringNode(token.literal);
            case TT_FUNCTION:
                return parseFunctionLiteral();
            case TT_LBRACKET:
                return parseArrayLiteral();
            case TT_LBRACE:
                return parseHashLiteral();
            case TT_IF:
                return parseIfExpr();
            default:
                std::cout << "Unknown expression token: " << token.literal << std::endl;
                std::exit(1);
        }
    }
    // parseCallExpr ::= FUNCTIONCALL | ARRAYCALL | HASHCALL
    Node* Parser::parseCallExpr(Node *callee, Token token) {
        CallExprNode *callExprNode = new CallExprNode();
        callExprNode->callee = callee;

        if (token.type == TT_LBRACKET) { // array or hash call
            advance(TT_LBRACKET);
            if (curToken.type == TT_RBRACKET) {
                std::cout << "Invalid subscript reference for array or hash types\n";
                std::exit(1);
            }
            callExprNode->arguments.emplace_back(parseExpression());
            advance(TT_RBRACKET);
        }
        else if (token.type == TT_LPAREN) { // function call
            advance(TT_LPAREN);
            if (curToken.type != TT_RPAREN) {
                callExprNode->arguments.emplace_back(parseExpression());
                while (curToken.type == TT_COMMA) {
                    advance(TT_COMMA);
                    callExprNode->arguments.emplace_back(parseExpression());
                }
            }
            advance(TT_RPAREN);
        }
        return callExprNode;
    }
    // parseTernaryExpr
    Node* Parser::parseTernaryExpr(Node *condition) {
        IfNode *ifNode = new IfNode();

        ifNode->condition = parseExpression();

        advance(TT_QUESTION); // advance '?'
        ifNode->consequence = parseBlock();

        advance(TT_COLON);
        ifNode->alternative = parseBlock();

        return ifNode;
    }
    // parseFunctionLiteral
    Node* Parser::parseFunctionLiteral() {
        FunctionNode *functionNode = new FunctionNode();
        advance(TT_FUNCTION);
        advance(TT_LPAREN);
        if (curToken.type != TT_RPAREN) {
            functionNode->parameters.emplace_back((IdentNode*)parseIdentifier());
            while (curToken.type == TT_COMMA) {
                advance(TT_COMMA);
                functionNode->parameters.emplace_back((IdentNode*)parseIdentifier());
            }
        }
        advance(TT_RPAREN);
        functionNode->body = (BlockNode*)parseBlock();

        return functionNode;
    }
    // parseArrayLiteral
    Node* Parser::parseArrayLiteral() {
        ArrayNode* arrayNode = new ArrayNode();
        advance(TT_LBRACKET);
        if (curToken.type != TT_RBRACKET) {
            arrayNode->elements.emplace_back(parseExpression());
            while (curToken.type == TT_COMMA) {
                advance(TT_COMMA);
                arrayNode->elements.emplace_back(parseExpression());
            }
        }
        advance(TT_RBRACKET);

        return arrayNode;
    }
    // parseHashLiteral
    Node* Parser::parseHashLiteral() {
        HashNode* hashNode = new HashNode();
        advance(TT_LBRACE);
        if (curToken.type != TT_RBRACE) {
            // parse key-value pair expressions
            hashNode->keys.emplace_back((IdentNode*)parseExpression());
            advance(TT_COLON); // advance ':'
            hashNode->values.emplace_back(parseExpression());

            while (curToken.type == TT_COMMA) {
                advance(TT_COMMA);
                // parse key-value pair expressions (refactor here in the future)
                hashNode->keys.emplace_back((IdentNode*)parseExpression());
                advance(TT_COLON); // advance ':'
                hashNode->values.emplace_back(parseExpression());
            }
        }
        advance(TT_RBRACE);
        return hashNode;
    }
    // parseIfExpression
    Node* Parser::parseIfExpr() {
        IfNode* ifNode = new IfNode();

        advance(TT_IF);

        advance(TT_LPAREN);
        ifNode->condition = parseExpression();
        advance(TT_RPAREN);

        ifNode->consequence = parseBlock();

        if (curToken.type = TT_ELSE) {
            advance(TT_ELSE);
            ifNode->alternative = parseBlock();
        }

        return ifNode;
    }
    // parserIdentifier
    Node* Parser::parseIdentifier() {
        IdentNode* identNode = new IdentNode();
        identNode->value = curToken;

        advance(TT_IDENT); // skip the IDENT token.

        // new feature will be placed here...

        return identNode;
    }
}

