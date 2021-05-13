//
// Created by irwin on 13/05/2021.
//

#ifndef CPP_ENVIRONMENT_H
#define CPP_ENVIRONMENT_H
#include <map>
#include "object.h"

namespace corny {
    class ObjBase;
    class Environment {
    public:
        Environment() {};
        Environment(Environment* outer) {
            this->outer = outer;
        }
        Environment* outer = nullptr;
        std::map<std::string, ObjBase*> symbolTable;
        // register an object in symbol table.
        void set(std::string key, ObjBase* value) {
            symbolTable[key] = value;
        }
        // get an object from symbol table.
        ObjBase* get(std::string key) {
            if (symbolTable.find(key) == symbolTable.end()) {
                if (outer != nullptr) {
                    return outer->get(key);
                } else {
                    return nullptr;
                }
            }
            return symbolTable[key];
        }
    };
}

#endif //CPP_ENVIRONMENT_H
