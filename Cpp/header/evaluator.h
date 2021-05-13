//
// Created by irwin on 13/05/2021.
//

#ifndef CPP_EVALUATOR_H
#define CPP_EVALUATOR_H

#include "object.h"
#include "ast.h"
#include "environment.h"

namespace corny {
    class Evaluator {
    public:
        Evaluator() {}
        ~Evaluator() {}
        BooleanObj *TRUE = new BooleanObj(true);
        BooleanObj *FALSE = new BooleanObj(false);
        NullObj *NIL = new NullObj();

        static bool isError(Object* obj);
        Object* eval(Node* node, Environment* env);
        Object* evalProgram(ProgramNode* programNode, Environment* env);
        Object* evalBlock(BlockNode* blockNode, Environment* env);
        Object* evalLet(LetNode* letNode, Environment* env);
        Object* evalReturn(ReturnNode* returnNode, Environment* env);
        Object* evalIdentifier(IdentNode* identNode, Environment* env);
        Object* evalCallExpr(CallExprNode* callExprNode, Environment* env);
        Object* evalFunctionLiteral(FunctionNode* functionNode, Environment* env);
        Object* evalFunction(FunctionObj* functionObj, std::vector<Object*> arguments);
        Object* evalArrayAccess(ArrayObj* arrayObj, std::vector<Object*> arguments);
        Object* evalHashAccess(HashObj* hashObj, std::vector<Object*> arguments);
        Object* evalStringAccess(StringObj* stringObj, std::vector<Object*> arguments);
        Object* evalArrayLiteral(ArrayNode* arrayNode, Environment* env);
        Object* evalHashLiteral(HashNode* hashNode, Environment* env);
        Object* evalUnaryExpression(UnaryNode* unaryNode, Environment* env);
        Object* evalBinaryExpression(BinOpNode* binOpNode, Environment* env);
        Object* evalLogicalExpression(BinOpNode* binOpNode, Environment* env);
        static Object* evalBinaryString(Object* leftObj, TokenType type, Object* rightObj);
        static Object* evalBinaryInteger(Object* leftObj, TokenType type, Object* rightObj);
        static Object* evalBinaryBoolean(Object* leftObj, TokenType type, Object* rightObj);
        Object* evalIfExpression(IfNode* ifNode, Environment* env);
    };
}
#endif //CPP_EVALUATOR_H
