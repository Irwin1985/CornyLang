//
// Created by irwin on 11/05/2021.
//

#include "../header/token.h"

namespace corny {
    std::map<std::string, TokenType> keywords ({
        {"fn",  TT_FUNCTION},
        {"let", TT_LET},
        {"true", TT_TRUE},
        {"false", TT_FALSE},
        {"return", TT_RETURN},
        {"if", TT_IF},
        {"else", TT_ELSE},
    });
    // isKeyword: check if 'ident' is a keyword or a ident token.
    TokenType isKeyword(std::string ident) {
        if (keywords.find(ident) == keywords.end()) {
            return TT_IDENT;
        }
        return keywords[ident];
    }
}
