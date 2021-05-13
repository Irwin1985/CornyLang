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

        static bool isError(ObjBase* obj);
        ObjBase* eval(Node* node, Environment* env);
        ObjBase* evalProgram(ProgramNode* programNode, Environment* env);
        ObjBase* evalBlock(BlockNode* blockNode, Environment* env);
        ObjBase* evalLet(LetNode* letNode, Environment* env);
        ObjBase* evalReturn(ReturnNode* returnNode, Environment* env);
        ObjBase* evalIdentifier(IdentNode* identNode, Environment* env);
        ObjBase* evalCallExpr(CallExprNode* callExprNode, Environment* env);
        ObjBase* evalFunctionLiteral(FunctionNode* functionNode, Environment* env);
        ObjBase* evalFunction(FunctionObj* functionObj, std::vector<ObjBase*> arguments);
        ObjBase* evalArrayAccess(ArrayObj* arrayObj, std::vector<ObjBase*> arguments);
        ObjBase* evalHashAccess(HashObj* hashObj, std::vector<ObjBase*> arguments);
        ObjBase* evalStringAccess(StringObj* stringObj, std::vector<ObjBase*> arguments);
    };
}
#endif //CPP_EVALUATOR_H
