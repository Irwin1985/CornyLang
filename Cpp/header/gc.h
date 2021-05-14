//
// Created by irwin on 14/05/2021.
//

#ifndef CPP_GC_H
#define CPP_GC_H
#include "object.h"
#include "environment.h"

namespace corny {
    class Environment; // forward reference to avoid the circular dependency

    class GarbageCollector {
    public:
        GarbageCollector() {
            this->head = new StringObj("head"); // the main object
        }
        ~GarbageCollector() {
            delete head;
        }
        // Add an object
        void add(Object* obj) {
            obj->next = head->next;
            head->next = obj;
        }
        // Mark an object
        void mark(Object* obj) {
            if (obj->mark == true) return;
            obj->mark = true;
            // An Array must mark all its elements
            if (obj->type == OBJ_ARRAY) {
                for (auto element : ((ArrayObj*)obj)->elements) {
                    mark(element);
                }
            }
            // A Hash table must mark all its elements
            if (obj->type == OBJ_HASH) {
                for (std::pair<std::string, Object*> pair : ((HashObj*)obj)->elements) {
                    Object* element = pair.second;
                    mark(element);
                }
            }
        }
        // overload the mark method to allow Environment
        void mark(Environment* env) {
            // Environments must mark all objects contained in its symbol table.
            for (std::pair<std::string, Object*> p : env->symbolTable) {
                Object* obj = p.second;
                mark(obj);
            }
            // and dont forget its outer environment
            if (env->outer != nullptr) {
                mark(env->outer);
            }
        }
        // sweep method that find all unreferenced objects and delete them.
        void sweep() {
            Object* node = head;
            // search for unmarked objects
            while (node->next != nullptr) {
                if (node->next->mark == false) {
                    Object* temp = node->next;
                    node->next = temp->next;
                    delete temp;
                } else {
                    // this object was reached so unmark it (for the next GC)
                    // and move on to the next.
                    //node->mark = false;
                    node = node->next;
                }
            }
            // now unmark all active objects
            node = head;
            while (node->next != nullptr) {
                node->next->mark = false;
                node = node->next;
            }
        }

        Object* head;
    };
}

#endif //CPP_GC_H
