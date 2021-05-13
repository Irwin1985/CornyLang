//
// Created by irwin on 12/05/2021.
//

#ifndef CPP_AST_H
#define CPP_AST_H
#include <iostream>
#include <vector>
#include "token.h"

namespace corny {
    // NodeType
    enum NodeType {
        NT_PROGRAM,
        NT_BLOCK,
        NT_NUMBER,
        NT_LET,
        NT_RETURN,
        NT_FUNCTION,
        NT_ARRAY,
        NT_HASH,
        NT_CALL,
        NT_IDENT,
        NT_BINARY,
        NT_UNARY,
        NT_STRING,
        NT_BOOLEAN,
        NT_NULL,
        NT_IF,
    };
    // Node base class: all nodes will inherit from it.
    class Node {
    public:
        Node() {};
        ~Node() {};
        NodeType type;
        virtual std::string toString() = 0;
    };
    // ProgramNode
    class ProgramNode : public Node {
    public:
        ProgramNode() {
            this->type = NT_PROGRAM;
        }
        ~ProgramNode() {
            if (statements.size() > 0) {
                for (auto statement : statements) {
                    delete statement;
                }
            }
        }
        std::vector<Node*> statements;

        std::string toString() {
            std::string result;
            if (statements.size() > 0) {
                for (Node* statement : statements) {
                    result += statement->toString() + '\n';
                }
            }
            return result;
        }
    };
    // BlockNode
    class BlockNode : public Node {
    public:
        BlockNode() {
            this->type = NT_BLOCK;
        }
        std::vector<Node*> statements;
        std::string toString() {
            std::string result = "\n{\n";
            if (statements.size() > 0) {
                for (Node* statement : statements) {
                    result += "   " + statement->toString() + ";\n";
                }
            }
            result += "}";
            return result;
        }
    };
    // ReturnNode
    class ReturnNode : public Node {
    public:
        ReturnNode() {
            this->type = NT_RETURN;
        }
        ReturnNode(Node* value) {
            this->value = value;
            this->type = NT_RETURN;
        }
        Node* value;
        std::string toString() {
            return "return " + value->toString();
        }
    };
    // BinaryNode
    class BinOpNode : public Node {
    public:
        BinOpNode() {
            this->type = NT_BINARY;
        }
        BinOpNode(Node* left, Token opToken, Node* right) {
            this->left = left;
            this->opToken = opToken;
            this->right = right;
            this->type = NT_BINARY;
        }
        Node* left;
        Token opToken;
        Node* right;
        std::string toString() {
            return left->toString() + " " + opToken.literal + " " + right->toString();
        }
    };
    // UnaryNode
    class UnaryNode : public Node {
    public:
        UnaryNode() {
            this->type = NT_UNARY;
        }
        UnaryNode(Token opToken, Node* left) {
            this->left = left;
            this->opToken = opToken;
            this->type = NT_UNARY;
        }
        Node* left;
        Token opToken;
        std::string toString() {
            return left->toString() + " " + opToken.literal;
        }
    };
    // CallExprNode
    class CallExprNode : public Node {
    public:
        CallExprNode() {
            this->type = NT_CALL;
        }
        Node* callee;
        std::vector<Node*> arguments;

        std::string toString() {
            std::string result = callee->toString();
            if (arguments.size() > 0) {
                std::string openDelim, closeDelim;
                if (callee->type == NT_ARRAY || callee->type == NT_HASH) {
                    openDelim = "[";
                    closeDelim = "]";
                } else {
                    openDelim = "(";
                    closeDelim = ")";
                }
                result += openDelim;
                int counter = 0;
                for (Node* argument : arguments) {
                    counter += 1;
                    if (counter > 1) result += ',';
                    result += argument->toString();
                }
                result += closeDelim;
            }
            return result;
        }
    };
    // IfNode
    class IfNode : public Node {
    public:
        IfNode() {
            this->type = NT_IF;
        }
        Node* condition;
        Node* consequence;
        Node* alternative;
        std::string toString() {
            std::string result = "if(";
            result += condition->toString() + ")";
            // evaluate consequence
            if (consequence != nullptr) {
                result += consequence->toString();
            }
            // evaluate alternative
            if (alternative != nullptr) {
                result += "else" + alternative->toString();
            }

            return result;
        }
    };
    // NumberNode
    class NumberNode : public Node {
    public:
        NumberNode() {
            this->type = NT_NUMBER;
        };
        NumberNode(double value) {
            this->value = value;
            this->type = NT_NUMBER;
        }
        double value = 0.00;
        std::string toString() {
            return std::to_string(value);
        }
    };
    // StringNode
    class StringNode : public Node {
    public:
        StringNode() {
            this->type = NT_STRING;
        };
        StringNode(std::string value) {
            this->value = value;
            this->type = NT_STRING;
        }
        std::string value;
        std::string toString() {
            return '\"' + value + '\"';
        }
    };
    // BooleanNode
    class BooleanNode : public Node {
    public:
        BooleanNode() {
            this->type = NT_BOOLEAN;
        };
        BooleanNode(bool value) {
            this->value = value;
            this->type = NT_BOOLEAN;
        }
        bool value;
        std::string toString() {
            return (value == true) ? "true" : "false";
        }
    };
    // IdentNode
    class IdentNode : public Node {
    public:
        IdentNode() {
            this->type = NT_IDENT;
        }
        IdentNode(Token value) {
            this->value = value;
            this->type = NT_IDENT;
        }
        Token value;

        std::string toString() {
            return value.literal;
        }
    };
    // NullNode
    class NullNode : public Node {
    public:
        NullNode() {
            this->type = NT_NULL;
        };
        std::string toString() {
            return "null";
        }
    };
    // LetNode
    class LetNode : public Node {
    public:
        LetNode() {
            this->type = NT_LET;
        }
        IdentNode* ident;
        Node* value;
        std::string toString() {
            return "let " + ident->toString() + " = " + value->toString();
        }
    };
    // FunctionNode
    class FunctionNode : public Node {
    public:
        FunctionNode() {
            this->type = NT_FUNCTION;
        }
        std::vector<IdentNode*> parameters;
        BlockNode* body;

        std::string toString() {
            std::string result = "fn(";
            if (parameters.size() > 0) {
                int counter = 0;
                for (auto parameter : parameters) {
                    counter += 1;
                    if (counter > 1) result += ",";
                    result += parameter->toString();
                }
            }
            result += ")" + body->toString();

            return result;
        }
    };
    // ArrayNode
    class ArrayNode : public Node {
    public:
        ArrayNode() {
            this->type = NT_ARRAY;
        }
        std::vector<Node*> elements;

        std::string toString() {
            std::string result = "[";
            if (elements.size() > 0) {
                int counter = 0;
                for (Node* element : elements) {
                    counter += 1;
                    if (counter > 1) result += ",";
                    result += element->toString();
                }
            }
            result += "]";
            return result;
        }
    };
    // HashNode
    class HashNode : public Node {
    public:
        HashNode() {
            this->type = NT_HASH;
        }
        std::vector<IdentNode*> keys;
        std::vector<Node*> values;

        std::string toString() {
            std::string result = "{";
            if (keys.size() > 0) {
                for (int i = 0; i < keys.size(); i++) {
                    if (i >= 1) result += ",";
                    result += keys.at(i)->toString() + ":" + values.at(i)->toString();
                }
            }
            result += "}";

            return result;
        }
    };
}

#endif //CPP_AST_H
