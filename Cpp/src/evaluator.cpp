//
// Created by irwin on 13/05/2021.
//
#include "../header/evaluator.h"

namespace corny {
    // check whether passed object is null or ErrorObj type.
    bool Evaluator::isError(Object *obj) {
        return obj != nullptr && obj->type == OBJ_ERROR;
    }
    // recursively evaluates the current node.
    Object* Evaluator::eval(Node *node, Environment *env) {
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
            case NT_UNARY:
                return evalUnaryExpression((UnaryNode*)node, env);
            case NT_BINARY:
                return evalBinaryExpression((BinOpNode*)node, env);
            case NT_STRING:
                return new StringObj(((StringNode*)node)->value);
            case NT_LET:
                return evalLet((LetNode*)node, env);
            case NT_RETURN:
                return evalReturn((ReturnNode*)node, env);
            case NT_IDENT:
                return evalIdentifier((IdentNode*)node, env);
            case NT_ARRAY:
                return evalArrayLiteral((ArrayNode*)node, env);
            case NT_CALL:
                return evalCallExpr((CallExprNode*)node, env);
            case NT_FUNCTION:
                return evalFunctionLiteral((FunctionNode*)node, env);
            case NT_IF:
                return evalIfExpression((IfNode*)node, env);
            default:
                return new ErrorObj("Unknown Node type.");
        }
    }
    // evalProgram
    Object* Evaluator::evalProgram(ProgramNode *programNode, Environment *env) {
        Object* resultObj;
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
    Object* Evaluator::evalBlock(BlockNode *blockNode, Environment *env) {
        Object* resultObj;
        for (auto statement : blockNode->statements) {
            resultObj = eval(statement, env);

            if (isError(resultObj) || resultObj->type == OBJ_RETURN)
                return resultObj;
        }
        return resultObj;
    }
    // evalLet
    Object* Evaluator::evalLet(LetNode *letNode, Environment *env) {
        // evaluate the value property
        Object* valueObj = eval(letNode->value, env);
        if (isError(valueObj)) return valueObj;
        // register the symbol
        env->set(letNode->ident->value.literal, valueObj);

        return valueObj;
    }
    // evalReturn
    Object* Evaluator::evalReturn(ReturnNode *returnNode, Environment *env) {
        Object* valueObj = eval(returnNode->value, env);
        if (isError(valueObj)) return valueObj;

        return new ReturnObj(valueObj);
    }
    // evalIfExpression
    Object* Evaluator::evalIfExpression(IfNode *ifNode, Environment *env) {
        Object* conditionObj = eval(ifNode->condition, env);
        if (isError(conditionObj)) return conditionObj;
        if (conditionObj->type != OBJ_BOOLEAN) return new ErrorObj("Invalid data type for if condition");
        // evaluate if or else based on condition value.
        if (((BooleanObj*)conditionObj)->value == true) {
            return eval(ifNode->consequence, env);
        } else {
            if (ifNode->alternative != nullptr) {
                return eval(ifNode->alternative, env);
            } else {
                return NIL;
            }
        }
    }
    // evalUnaryExpression
    Object* Evaluator::evalUnaryExpression(UnaryNode *unaryNode, Environment *env) {
        Object* rightObj = eval(unaryNode->left, env);
        if (isError(rightObj)) return rightObj;
        // check the token operator
        switch (unaryNode->opToken.type) {
            case TT_MINUS:
                if (rightObj->type != OBJ_NUMBER) return new ErrorObj("Invalid data type");
                return new NumberObj(((NumberObj*)rightObj)->value * -1);
            case TT_NOT:
                if (rightObj->type != OBJ_BOOLEAN) return new ErrorObj("Invalid data type");
                return (((BooleanObj*)rightObj)->value == true) ? FALSE : TRUE;
            default:
                return new ErrorObj("Invalid operator: " + unaryNode->opToken.literal);
        }
    }
    // evalBinaryExpression
    Object* Evaluator::evalBinaryExpression(BinOpNode *binOpNode, Environment *env) {
        if (binOpNode->opToken.type == TT_AND || binOpNode->opToken.type == TT_OR) {
            return evalLogicalExpression(binOpNode, env);
        }
        // we need to know the both operands types
        Object* leftObj = eval(binOpNode->left, env);
        if (isError(leftObj)) return leftObj;
        Object* rightObj = eval(binOpNode->right, env);
        if (isError(rightObj)) return rightObj;
        // now based on their types we perform the correct operations
        if (leftObj->type == OBJ_STRING && rightObj->type == OBJ_STRING) {
            return evalBinaryString(leftObj, binOpNode->opToken.type, rightObj);
        }
        if (leftObj->type == OBJ_NUMBER && rightObj->type == OBJ_NUMBER) {
            return evalBinaryInteger(leftObj, binOpNode->opToken.type, rightObj);
        }
        if (leftObj->type == OBJ_BOOLEAN && rightObj->type == OBJ_NUMBER) {
            return evalBinaryBoolean(leftObj, binOpNode->opToken.type, rightObj);
        }
        return new ErrorObj("Operator not supported: " + binOpNode->opToken.literal);
    }
    // evalLogicalExpression
    Object* Evaluator::evalLogicalExpression(BinOpNode *binOpNode, Environment *env) {
        // evaluate the left hand operator
        Object* leftObj = eval(binOpNode->left, env);
        if (isError(leftObj)) return leftObj;
        if (leftObj->type != OBJ_BOOLEAN) return new ErrorObj("Invalid left hand type operand");
        if (binOpNode->opToken.type == TT_AND) {
            // if leftObj is false then there's nothing else to do.
            if (((BooleanObj*)leftObj)->value == true) return FALSE;
            // otherwise we need to evaluate the right hand operator
            Object* rightObj = eval(binOpNode->right, env);
            if (isError(rightObj)) return rightObj;
            if (rightObj->type != OBJ_BOOLEAN) return new ErrorObj("Invalid right hand type operand");

            return (((BooleanObj*)leftObj)->value == true) ? TRUE : FALSE;
        }
        else if (binOpNode->opToken.type == TT_OR){
            // if leftObj is true then there's nothing else to do.
            if (((BooleanObj*)leftObj)->value == true) return TRUE;
            // otherwise we need to evaluate the right hand operator
            Object* rightObj = eval(binOpNode->right, env);
            if (isError(rightObj)) return rightObj;
            if (rightObj->type != OBJ_BOOLEAN) return new ErrorObj("Invalid right hand type operand");

            return (((BooleanObj*)leftObj)->value == true) ? TRUE : FALSE;
        }
        return new ErrorObj("Operator not supported: " + binOpNode->opToken.literal);
    }
    // evalBinaryString
    Object* Evaluator::evalBinaryString(Object* leftObj, TokenType type, Object* rightObj) {
        switch (type) {
            case TT_PLUS:
                return new StringObj(((StringObj*)leftObj)->value + ((StringObj*)rightObj)->value);
            default:
                return new ErrorObj("Invalid operator");
        }
    }
    // evalBinaryInteger
    Object* Evaluator::evalBinaryInteger(Object* leftObj, TokenType type, Object* rightObj) {
        switch (type) {
            case TT_PLUS:
                return new NumberObj(((NumberObj*)leftObj)->value + ((NumberObj*)rightObj)->value);
            case TT_MINUS:
                return new NumberObj(((NumberObj*)leftObj)->value - ((NumberObj*)rightObj)->value);
            case TT_MUL:
                return new NumberObj(((NumberObj*)leftObj)->value * ((NumberObj*)rightObj)->value);
            case TT_DIV:
            {
                double rightVal = ((NumberObj*)rightObj)->value;
                if (rightVal <= 0) return new ErrorObj("Division by zero.");
                return new NumberObj(((NumberObj*)leftObj)->value / rightVal);
            }
            case TT_LESS:
                return new BooleanObj(((NumberObj*)leftObj)->value < ((NumberObj*)rightObj)->value);
            case TT_GREATER:
                return new BooleanObj(((NumberObj*)leftObj)->value > ((NumberObj*)rightObj)->value);
            case TT_LESS_EQ:
                return new BooleanObj(((NumberObj*)leftObj)->value <= ((NumberObj*)rightObj)->value);
            case TT_GREATER_EQ:
                return new BooleanObj(((NumberObj*)leftObj)->value >= ((NumberObj*)rightObj)->value);
            case TT_EQUAL:
                return new BooleanObj(((NumberObj*)leftObj)->value == ((NumberObj*)rightObj)->value);
            case TT_NOT_EQ:
                return new BooleanObj(((NumberObj*)leftObj)->value != ((NumberObj*)rightObj)->value);
            default:
                return new ErrorObj("Invalid operator.");
        }
    }
    // evalBinaryBoolean
    Object* Evaluator::evalBinaryBoolean(Object* leftObj, TokenType type, Object* rightObj) {
        // this is the trick: transform boolean to int and call binary integer operations.
        NumberObj *leftNumber = new NumberObj((((BooleanObj*)leftObj)->value == true) ? 1 : 0);
        NumberObj *rightNumber = new NumberObj((((BooleanObj*)rightObj)->value == true) ? 1 : 0);
        return evalBinaryInteger(leftNumber, type, rightNumber);
    }
    // evalIdentifier
    Object* Evaluator::evalIdentifier(IdentNode *identNode, Environment *env) {
        std::string identifier = identNode->value.literal;
        Object *valueObj = env->get(identifier);
        if (valueObj == nullptr) {
            return new ErrorObj("variable not defined: " + identifier);
        }
        return valueObj;
    }
    // evalCallExpression
    Object* Evaluator::evalCallExpr(CallExprNode *callExprNode, Environment *env) {
        // this node can contain function call, array call or hash table call.
        // 1. get the object of the callee
        Object* calleeObj = eval(callExprNode->callee, env);
        if (isError(calleeObj)) return calleeObj;
        // 2. evaluate the arguments
        std::vector<Object*> arguments;
        if (callExprNode->arguments.size() > 0) {
            Object* resultObj;
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
    Object* Evaluator::evalFunctionLiteral(FunctionNode *functionNode, Environment *env) {
        FunctionObj *functionObj = new FunctionObj();
        functionObj->env = env;
        functionObj->parameters = functionNode->parameters;
        functionObj->body = functionNode->body;

        return functionObj;
    }
    // evalArrayLiteral
    Object* Evaluator::evalArrayLiteral(ArrayNode *arrayNode, Environment *env) {
        ArrayObj *arrayObj = new ArrayObj();
        // evaluate the elements of the array
        if (arrayNode->elements.size() > 0) {
            Object* resultObj;
            for (auto element : arrayNode->elements) {
                resultObj = eval(element, env);
                if (isError(resultObj)) return resultObj;
                // push the element into the array.
                arrayObj->elements.emplace_back(resultObj);
            }
        }
        return arrayObj;
    }
    // evalHashLiteral
    Object* Evaluator::evalHashLiteral(HashNode *hashNode, Environment *env) {
        HashObj* hashObj = new HashObj();
        int index = 0;
        Object* keyObj, *valueObj;
        // loop through keys and their values
        for (auto keyNode : hashNode->keys) {
            keyObj = eval(keyNode, env);
            if (isError(keyObj)) return keyObj;
            // validate the OBJ_STRING data type
            if (keyObj->type != OBJ_STRING) return new ErrorObj("Invalid data type for key");
            // evaluate the value
            valueObj = eval(hashNode->values.at(index), env);
            if (isError(valueObj)) return valueObj;
            // save the key-value in data type
            hashObj->elements[((StringObj*)keyObj)->value] = valueObj;
        }
        return hashObj;
    }
    // evalFunction
    Object* Evaluator::evalFunction(FunctionObj *functionObj, std::vector<Object *> arguments) {
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
        Object *resultObj = eval(functionObj->body, newEnv);
        if (isError(resultObj)) return resultObj;
        // 5. check for return
        if (resultObj->type == OBJ_RETURN) return ((ReturnObj*)resultObj)->value;

        return resultObj;
    }
    // evalArrayAccess
    Object* Evaluator::evalArrayAccess(ArrayObj *arrayObj, std::vector<Object *> arguments) {
        Object* indexObj = arguments.at(0);
        if (indexObj->type != OBJ_NUMBER) return new ErrorObj("Invalid subscript reference");
        // check for out of bounds
        int index = ((NumberObj*)indexObj)->value;
        if (index <= 0 || index >= arrayObj->elements.size()) {
            return new ErrorObj("Index out of bounds");
        }
        Object* valueObj = arrayObj->elements[index];

        return valueObj;
    }
    // evalHashAccess
    Object* Evaluator::evalHashAccess(HashObj *hashObj, std::vector<Object *> arguments) {
        Object* indexObj = arguments.at(0);
        if (indexObj->type != OBJ_STRING) return new ErrorObj("Invalid subscript reference");

        return nullptr;
    }
    // stringAccess
    Object* Evaluator::evalStringAccess(StringObj *stringObj, std::vector<Object *> arguments) {
        Object* indexObj = arguments.at(0);
        if (indexObj->type != OBJ_NUMBER) return new ErrorObj("Invalid subscript reference");
        int index = ((NumberObj*)indexObj)->value;
        // check for out of bounds
        if (index <= 0 || index > stringObj->value.length()) {
            return new ErrorObj("Index out of bounds");
        }

        return new StringObj(std::string(1, stringObj->value[index]));
    }
}