#include <stdio.h>
#include <stdlib.h>
#include "header/lexer.h"
#include "header/util.h"

void repl();

int main() {
    repl();
    return 0;
}

// repl
void repl() {
    printf("CornyLang v1.0.1\n");
    char* input = NULL;
    size_t len = 0;
    while (1) {
        printf("> ");
        if (!getline(&input, &len, stdin)) break;
        // create the lexer
        Lexer *lexer = newLexer(input);
        Token tok = nextToken(lexer);
        while (tok.type != TT_EOF) {
            printf("type: %s, literal: '%s'\n", getTokenType(tok.type), tok.literal);
            tok = nextToken(lexer);
        }
        deleteLexer(lexer);
    }
}