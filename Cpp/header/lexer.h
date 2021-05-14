//
// Created by irwin on 11/05/2021.
//

#ifndef CPP_LEXER_H
#define CPP_LEXER_H

#include <string>
#include "token.h"

namespace corny {
    class Lexer {
    public:
        Lexer() {};
        ~Lexer() {};

        std::string input;
        int pos = 0;
        char current_char = '\0';

        // methods prototype
        void start(std::string input);
        static bool isDigit(char chr);
        static bool isAlpha(char chr);
        static bool isLetter(char chr);
        static bool isString(char chr);
        std::string getNumber();
        void checkEOF();
        void advance();
        char peek();
        void skipComments();
        void skipWhitespace();
        Token parseInteger();
        Token parseString(char delimiter);
        Token parseIdent();
        Token nextToken();
    };
}

#endif //CPP_LEXER_H
