//
// Created by irwin on 13/05/2021.
//
#include "../header/evaluator.h"

namespace corny {
    // check whether passed object is null or ErrorObj type.
    bool Evaluator::isError(ObjBase *obj) {
        return obj != nullptr && obj->type == OBJ_ERROR;
    }
    // recursively evaluates the current node.
    ObjBase* Evaluator::eval(Node *node, Environment *env) {
        NodeType type = node->type;
        switch(type) {
            case NT_PROGRAM:
                return evalProgram((ProgramNode*)node, env);
            case NT_BLOCK:
                return evalBlock((BlockNode*)node, env);
            case NT_BOOLEAN:
                return (((BooleanNode*)node)->value) ? TRUE : FALSE;
            case NT_NUMBER:
                return new NumberObj(((NumberNode*)node)->value);
            case NT_NULL:
                return NIL;
            case NT_STRING:
                return new StringObj(((StringNode*)node)->value);
            case NT_LET:
                return evalLet((LetNode*)node, env);
            case NT_RETURN:
                return evalReturn((ReturnNode*)node, env);
            case NT_IDENT:
                return evalIdentifier((IdentNode*)node, env);
            case NT_CALL:
                return evalCallExpr((CallExprNode*)node, env);
            case NT_FUNCTION:
                return evalFunctionLiteral((FunctionNode*)node, env);
            default:
                return new ErrorObj("Unknown Node type.");
        }
    }
    // evalProgram
    ObjBase* Evaluator::evalProgram(ProgramNode *programNode, Environment *env) {
        ObjBase* resultObj;
        for (auto statement : programNode->statements) {
            resultObj = eval(statement, env);
            if (isError(resultObj)) {
                return resultObj;
            }
            else if (resultObj->type == OBJ_RETURN) {
                return ((ReturnObj*)resultObj)->value;
            }
        }

        return resultObj;
    }
    // evalBlock
    ObjBase* Evaluator::evalBlock(BlockNode *blockNode, Environment *env) {
        ObjBase* resultObj;
        for (auto statement : blockNode->statements) {
            resultObj = eval(statement, env);

            if (isError(resultObj) || resultObj->type == OBJ_RETURN)
                return resultObj;
        }
        return resultObj;
    }
    // evalLet
    ObjBase* Evaluator::evalLet(LetNode *letNode, Environment *env) {
        // evaluate the value property
        ObjBase* valueObj = eval(letNode->value, env);
        if (isError(valueObj)) return valueObj;
        // register the symbol
        env->set(letNode->ident->value.literal, valueObj);

        return valueObj;
    }
    // evalReturn
    ObjBase* Evaluator::evalReturn(ReturnNode *returnNode, Environment *env) {
        ObjBase* valueObj = eval(returnNode->value, env);
        if (isError(valueObj)) return valueObj;

        return new ReturnObj(valueObj);
    }
    // evalIdentifier
    ObjBase* Evaluator::evalIdentifier(IdentNode *identNode, Environment *env) {
        std::string identifier = identNode->value.literal;
        ObjBase *valueObj = env->get(identifier);
        if (valueObj == nullptr) {
            return new ErrorObj("variable not defined: " + identifier);
        }
        return valueObj;
    }
    // evalCallExpression
    ObjBase* Evaluator::evalCallExpr(CallExprNode *callExprNode, Environment *env) {
        // this node can contain function call, array call or hash table call.
        // 1. get the object of the callee
        ObjBase* calleeObj = eval(callExprNode->callee, env);
        if (isError(calleeObj)) return calleeObj;
        // 2. evaluate the arguments
        std::vector<ObjBase*> arguments;
        if (callExprNode->arguments.size() > 0) {
            ObjBase* resultObj;
            for (auto argument : callExprNode->arguments) {
                resultObj = eval(argument, env);
                if (isError(resultObj)) return resultObj;
                arguments.emplace_back(resultObj);
            }
        }
        // 3. check the callee type
        ObjType type = calleeObj->type;
        switch (type) {
            case OBJ_FUNCTION:
                return evalFunction((FunctionObj*)calleeObj, arguments);
            case OBJ_ARRAY:
                return evalArrayAccess((ArrayObj*)calleeObj, arguments);
            case OBJ_HASH:
                return evalHashAccess((HashObj*)calleeObj, arguments);
            case OBJ_STRING:
                return evalStringAccess((StringObj*)calleeObj, arguments);
            default:
                return new ErrorObj("Invalid callable object.");
        }

    }
    // evalFunctionLiteral
    ObjBase* Evaluator::evalFunctionLiteral(FunctionNode *functionNode, Environment *env) {
        FunctionObj *functionObj = new FunctionObj();
        functionObj->env = env;
        functionObj->parameters = functionNode->parameters;
        functionObj->body = functionNode->body;

        return functionObj;
    }
    // evalFunction
    ObjBase* Evaluator::evalFunction(FunctionObj *functionObj, std::vector<ObjBase *> arguments) {
        // 1. create new environment for the function
        Environment* newEnv = new Environment();
        newEnv->outer = functionObj->env; // enclose environment

        // 2. check for function arity.
        int numArgs = arguments.size();
        int numParams = functionObj->parameters.size();
        if (numArgs != numParams) return new ErrorObj("Unexpected arguments, got: " + std::to_string(numArgs) + " want: " + std::to_string(numParams));

        // 3. fill the new environment with arguments
        for (int i = 0; i < numParams; i++) {
            newEnv->set(functionObj->parameters.at(i)->value.literal, arguments.at(i));
        }
        // 4. execute the function with new environment
        ObjBase *resultObj = eval(functionObj->body, newEnv);
        if (isError(resultObj)) return resultObj;
        // 5. check for return
        if (resultObj->type == OBJ_RETURN) return ((ReturnObj*)resultObj)->value;

        return resultObj;
    }
    // evalArrayAccess
    ObjBase* Evaluator::evalArrayAccess(ArrayObj *arrayObj, std::vector<ObjBase *> arguments) {
        ObjBase* indexObj = arguments.at(0);
        if (indexObj->type != OBJ_NUMBER) return new ErrorObj("Invalid subscript reference");
        // check for out of bounds
        int index = ((NumberObj*)indexObj)->value;
        if (index <= 0 || index > arrayObj->elements.size()) {
            return new ErrorObj("Index out of bounds");
        }
        ObjBase* valueObj = arrayObj->elements[index];

        return valueObj;
    }
    // evalHashAccess
    ObjBase* Evaluator::evalHashAccess(HashObj *hashObj, std::vector<ObjBase *> arguments) {
        ObjBase* indexObj = arguments.at(0);
        if (indexObj->type != OBJ_STRING) return new ErrorObj("Invalid subscript reference");

        return nullptr;
    }
    // stringAccess
    ObjBase* Evaluator::evalStringAccess(StringObj *stringObj, std::vector<ObjBase *> arguments) {
        ObjBase* indexObj = arguments.at(0);
        if (indexObj->type != OBJ_NUMBER) return new ErrorObj("Invalid subscript reference");
        int index = ((NumberObj*)indexObj)->value;
        // check for out of bounds
        if (index <= 0 || index > stringObj->value.length()) {
            return new ErrorObj("Index out of bounds");
        }

        return new StringObj(std::string(1, stringObj->value[index]));
    }
}

