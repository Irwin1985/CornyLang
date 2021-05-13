//
// Created by irwin on 13/05/2021.
//

#ifndef CPP_OBJECT_H
#define CPP_OBJECT_H
#include <iostream>
#include <vector>
#include "environment.h"
#include "ast.h"

namespace corny {
    enum ObjType {
        OBJ_ERROR,
        OBJ_NUMBER,
        OBJ_BOOLEAN,
        OBJ_STRING,
        OBJ_NULL,
        OBJ_FUNCTION,
        OBJ_ARRAY,
        OBJ_HASH,
        OBJ_RETURN,
    };
    // ObjBase class where all system objects inherit from.
    class ObjBase {
    public:
        ObjBase() {}
        ~ObjBase() {}
        ObjType type;
        virtual std::string Inspect() = 0;
    };
    // ErrorObj
    class ErrorObj : public ObjBase {
    public:
        ErrorObj(std::string message) {
            this->message = message;
            this->type = OBJ_ERROR;
        }
        std::string message;

        std::string Inspect() {
            return message;
        }
    };
    // ReturnObj
    class ReturnObj : public ObjBase {
    public:
        ReturnObj() {
            this->type = OBJ_RETURN;
        }
        ReturnObj(ObjBase* value) {
            this->value = value;
            this->type = OBJ_RETURN;
        }

        ObjBase* value;
        // inspect
        std::string Inspect() {
            return value->Inspect();
        }
    };
    // FunctionObj
    class FunctionObj : public ObjBase {
    public:
        FunctionObj() {
            this->type = OBJ_FUNCTION;
        }
        std::vector<IdentNode*> parameters;
        BlockNode* body;
        Environment *env;

        std::string Inspect() {
            return "function: ok";
        }
    };
    // ArrayObj
    class ArrayObj : public ObjBase {
    public:
        ArrayObj() {
            this->type = OBJ_ARRAY;
        }
        std::vector<ObjBase*> elements;
        std::string Inspect() {
            return "array";
        }
    };
    // HashObj
    class HashObj : public ObjBase {
    public:
        HashObj() {
            this->type = OBJ_HASH;
        }
        std::vector<ObjBase*> keys;
        std::vector<ObjBase*> values;

        std::string Inspect() {
            return "hash";
        }
    };
    // NumberObj
    class NumberObj : public ObjBase {
    public:
        NumberObj() {
            this->type = OBJ_NUMBER;
        }
        NumberObj(double value) {
            this->value = value;
            this->type = OBJ_NUMBER;
        }
        double value = 0.0;
        std::string Inspect() {
            return std::to_string(value);
        }
    };
    // BooleanObj
    class BooleanObj : public ObjBase {
    public:
        BooleanObj() {
            this->type = OBJ_BOOLEAN;
        }
        BooleanObj(bool value) {
            this->value = value;
            this->type = OBJ_BOOLEAN;
        }
        bool value;
        std::string Inspect() {
            return (value == true) ? "true" : "false";
        }
    };
    // NullObj
    class NullObj : public ObjBase {
    public:
        NullObj() {
            this->type = OBJ_NULL;
        }
        std::string Inspect() {
            return "null";
        }
    };
    // StringObj
    class StringObj : public ObjBase {
    public:
        StringObj() {
            this->type = OBJ_STRING;
        }
        StringObj(std::string value) {
            this->value = value;
            this->type = OBJ_STRING;
        }
        std::string value;
        std::string Inspect() {
            return '\"'+ value + '\"';
        }
    };
}

#endif //CPP_OBJECT_H
