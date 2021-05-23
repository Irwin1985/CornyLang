package repl

import (
	"CornyLang/evaluator"
	"CornyLang/lexer"
	"CornyLang/object"
	"CornyLang/parser"
	"CornyLang/token"
	"bufio"
	"fmt"
	"io"
	"time"
)

func TestingParser(input string) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.Program()
	if program != nil {
		fmt.Printf(" %s\n", program.String())
	}
}

func ScanTokens(input string) {
	l := lexer.New(input)
	tok := l.NextToken()
	for tok.Type != token.TT_EOF {
		fmt.Printf("Type: %s, Literal: %s\n", tok.Type, tok.Literal)
		tok = l.NextToken()
	}
}

func TestEvaluator(input string) {
	var globalEnv = object.NewEnvironment()

	l := lexer.New(input)
	p := parser.New(l)
	program := p.Program()
	evaluated := evaluator.Eval(program, globalEnv)

	if evaluated != nil {
		fmt.Printf("%s\n", evaluated.Inspect())
	}
}

const (
	PROGRAM = "CornyLang"
	VERSION = "1.0.1"
	WELCOME = "Please feel free to type some valid commands or expressions!"
	PROMPT  = ">> "
	AUTOR   = "Irwin Rodríguez <rodriguez.irwin@gmail.com>"
	LOGO    = `
         ,\
         \\\,_
          \' ,\
     __,.-" =__)
   ."        )
,_/   ,    \/\_
\_|    )_-\ \_-'
'-----' '--'
	`
)

func Start(in io.Reader, out io.Writer) {

	scanner := bufio.NewScanner(in)
	// PRINT HEADER
	fmt.Printf("%s v%s %s\n%s\n%s\n", PROGRAM, VERSION, time.Now().String(), LOGO, WELCOME)
	var globalEnv = object.NewEnvironment()
	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		input := scanner.Text()
		if input == "quit" {
			break
		}
		if input == "version" {
			fmt.Print(VERSION)
			continue
		}
		if input == "autor" {
			fmt.Print(AUTOR)
			continue
		}

		l := lexer.New(input)
		p := parser.New(l)
		program := p.Program()
		evaluated := evaluator.Eval(program, globalEnv)

		if evaluated != nil {
			fmt.Printf("%s\n", evaluated.Inspect())
		}
	}
}
