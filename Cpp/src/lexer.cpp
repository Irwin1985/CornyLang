//
// Created by irwin on 11/05/2021.
//
#include <iostream>
#include "../header/lexer.h"

namespace corny {
    // checkEOF
    void Lexer::checkEOF() {
        if (current_char == NONE) {
            std::cout << "Unexpected End of File.\n";
            std::exit(1);
        }
    }
    // isDigit
    bool Lexer::isDigit(char chr) {
        return '0' <= chr && chr <= '9';
    }
    // isAlpha
    bool Lexer::isAlpha(char chr) {
        return 'a' <= chr && chr <= 'z' || 'A' <= chr && chr <= 'Z';
    }
    // isLetter
    bool Lexer::isLetter(char chr) {
        return isAlpha(chr) || isDigit(chr) || chr == '_';
    }
    // isString
    bool Lexer::isString(char chr) {
        return chr == '"' || chr == '\'';
    }
    // advance
    void Lexer::advance() {
        pos += 1;
        if (pos >= input.length()) current_char = NONE;
        current_char = input[pos];
    }
    // peek
    char Lexer::peek() {
        int peekPos = pos + 1;
        if (peekPos >= input.length()) return NONE;
        return input[peekPos];
    }
    // skipComments
    void Lexer::skipComments() {
        while (current_char != NONE && current_char != '\n') {
            advance();
        }
        //checkEOF();
        advance(); // the new line
    }
    // skipWhitespace
    void Lexer::skipWhitespace() {
        while (current_char != NONE && isspace(current_char)) {
            advance();
        }
        //checkEOF();
    }
    // getNumber
    std::string Lexer::getNumber() {
        std::string numStr;
        while (current_char != NONE && isDigit(current_char)) {
            numStr += current_char;
            advance();
        }
        //checkEOF();
        return numStr;
    }
    // parseInteger
    Token Lexer::parseInteger() {
        std::string lexeme = getNumber();
        if (current_char == '.' && isDigit(peek())) {
            lexeme += '.';
            advance(); // skip the '.'
            lexeme += getNumber();
        }
        return Token(NUMBER, lexeme);
    }
    // parseString
    Token Lexer::parseString(char delimiter) {
        advance(); // skip the delimiter
        std::string lexeme;
        while (current_char != NONE && current_char != delimiter) {
            lexeme += current_char;
            advance();
        }
        //checkEOF();
        advance(); // skip the closing delimiter

        return Token(STRING, lexeme);
    }
    // parseIdentifier
    Token Lexer::parseIdent() {
        std::string lexeme;
        while (current_char != NONE && isLetter(current_char)) {
            lexeme += current_char;
            advance();
        }
        return Token(corny::isKeyword(lexeme), lexeme);
    }
    // nextToken
    Token Lexer::nextToken() {
        while (current_char != NONE) {
            if (isspace(current_char)) {
                skipWhitespace();
                continue;
            }
            if (current_char == '/' && peek() == '/') {
                skipComments();
                continue;
            }
            if (isDigit(current_char)) return parseInteger();
            if (isAlpha(current_char)) return parseIdent();
            if (isString(current_char)) return parseString(current_char);
            // Arithmetic operator
            if (current_char == '+') {
                advance();
                return Token(PLUS, '+');
            }
            if (current_char == '-') {
                advance();
                return Token(MINUS, '-');
            }
            if (current_char == '*') {
                advance();
                return Token(MUL, '*');
            }
            if (current_char == '/') {
                advance();
                return Token(DIV, '/');
            }
            if (current_char == '^') {
                advance();
                return Token(POW, '^');
            }
            // Relational operators
            if (current_char == '<') {
                advance();
                if (current_char == '=') {
                    advance();
                    return Token(LESS_EQ, "<=");
                }
                return Token(LESS, '<');
            }
            if (current_char == '>') {
                advance();
                if (current_char == '=') {
                    advance();
                    return Token(GREATER_EQ, ">=");
                }
                return Token(GREATER, '>');
            }
            if (current_char == '!') {
                advance();
                if (current_char == '=') {
                    advance();
                    return Token(NOT_EQ, "!=");
                }
                return Token(NOT, '!');
            }
            if (current_char == '=') {
                advance();
                if (current_char == '=') {
                    advance();
                    return Token(EQUAL, "==");
                }
                return Token(ASSIGN, '=');
            }
            // Special characters
            if (current_char == ',') {
                advance();
                return Token(COMMA, ',');
            }
            if (current_char == ':') {
                advance();
                return Token(COLON, ':');
            }
            if (current_char == ';') {
                advance();
                return Token(SEMICOLON, ';');
            }
            if (current_char == '{') {
                advance();
                return Token(LBRACE, '{');
            }
            if (current_char == '}') {
                advance();
                return Token(RBRACE, '}');
            }
            if (current_char == '[') {
                advance();
                return Token(LBRACKET, '[');
            }
            if (current_char == ']') {
                advance();
                return Token(RBRACKET, ']');
            }
            if (current_char == '(') {
                advance();
                return Token(LPAREN, '(');
            }
            if (current_char == ')') {
                advance();
                return Token(RPAREN, ')');
            }

            std::cout << "unknown character: " << current_char << std::endl;
            std::exit(1);
        }
        return Token(END, "");
    }
}

