#include <iostream>
#include <ctime>
#include <string>

#include "header/lexer.h"
#include "header/parser.h"
#include "header/evaluator.h"
#include "header/environment.h"
#include <vector>

int main() {
    const std::string PROGRAM = "CornyLang";
    const std::string VERSION = "1.0.1";
    const std::string WELCOME = "Please feel free to type some valid commands or expressions!";
    const std::string PROMPT = ">> ";
    const std::string LOGO = R"V0G0N(
        ,\
        \\\,_
         \` ,\
    __,.-" =__)
  ."        )
,_/   ,    \/\_
\_|    )_-\ \_-`
`-----` `--`)V0G0N";
    const std::string OUTPUT = "";
    time_t TIME;
    std::time(&TIME);
    std::cout << PROGRAM << " v" << VERSION;
    std::cout << std::ctime(&TIME) << LOGO << std::endl;
    std::cout << WELCOME << std::endl << "Type: 'help' for more info.\n";
    std::string input;

    // create global variables
    corny::Lexer lexer;
    corny::Parser parser;
    corny::Evaluator evaluator;
    corny::Environment *globalEnv = new corny::Environment();

    // start the REPL
    while (true) {
        std::cout << PROMPT;
        std::getline(std::cin, input);
        if (input.empty())
            break;
        lexer.start(input);
        parser.start(lexer);

        // lexical analysis
        //corny::Lexer lexer(input);
        // syntax analysis
        //corny::Parser parser(lexer);
        // parse the syntax and generate AST node.
        corny::Node* program = parser.parseProgram();
        // evaluator
        //corny::Evaluator evaluator;
        // evaluate the node
        corny::Object* evaluated = evaluator.eval(program, globalEnv);
        // inspect the object and print out the result
        if (evaluated != nullptr) {
            std::cout << evaluated->Inspect() << std::endl;
        }
    }
    return 0;
}
