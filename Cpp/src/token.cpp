//
// Created by irwin on 11/05/2021.
//

#include "../header/token.h"

namespace corny {
    std::map<std::string, TokenType> keywords ({
        {"fn", FUNCTION},
        {"let", LET},
    });
    // isKeyword: check if 'ident' is a keyword or a ident token.
    TokenType isKeyword(std::string ident) {
        if (keywords.find(ident) == keywords.end()) {
            return IDENT;
        }
        return keywords[ident];
    }
}
