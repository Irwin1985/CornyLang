//
// Created by irwin on 13/05/2021.
//

#ifndef CPP_OBJECT_H
#define CPP_OBJECT_H
#include <iostream>
#include <vector>
#include <map>
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
    // Object class where all system objects inherit from.
    class Object {
    public:
        Object() {
            this->next = nullptr;
            this->mark = false;
        }
        ~Object() {}
        ObjType type;
        Object* next;
        bool mark;
        virtual std::string Inspect() = 0;
    };
    // ErrorObj
    class ErrorObj : public Object {
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
    class ReturnObj : public Object {
    public:
        ReturnObj() {
            this->type = OBJ_RETURN;
        }
        ReturnObj(Object* value) {
            this->value = value;
            this->type = OBJ_RETURN;
        }

        Object* value;
        // inspect
        std::string Inspect() {
            return value->Inspect();
        }
    };
    // FunctionObj
    class FunctionObj : public Object {
    public:
        FunctionObj() {
            this->type = OBJ_FUNCTION;
        }
        ~FunctionObj() {
            for (auto parameter : parameters) {
                delete parameter;
            }
            delete body;
        }
        std::vector<IdentNode*> parameters;
        BlockNode* body;
        Environment *env;

        std::string Inspect() {
            return "function: ok";
        }
    };
    // ArrayObj
    class ArrayObj : public Object {
    public:
        ArrayObj() {
            this->type = OBJ_ARRAY;
        }
        ~ArrayObj() {
            for (auto element : elements) {
                delete element;
            }
        }
        std::vector<Object*> elements;
        std::string Inspect() {
            return "array";
        }
    };
    // HashObj
    class HashObj : public Object {
    public:
        HashObj() {
            this->type = OBJ_HASH;
        }
        ~HashObj() {
            for (auto it = elements.begin(); it != elements.end();) {
                delete it->second;
                it = elements.erase(it);
            }
        }
        std::map<std::string, Object*> elements;

        std::string Inspect() {
            return "hash";
        }
    };
    // NumberObj
    class NumberObj : public Object {
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
    class BooleanObj : public Object {
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
    class NullObj : public Object {
    public:
        NullObj() {
            this->type = OBJ_NULL;
        }
        std::string Inspect() {
            return "null";
        }
    };
    // StringObj
    class StringObj : public Object {
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
