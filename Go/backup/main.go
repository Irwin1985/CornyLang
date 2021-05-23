package main

import (
	"CornyLang/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)

	//input := `let frutas = []; frutas.push("Peras", "Manzanas", "Piñas");`
	//repl.TestEvaluator(input)
	//repl.TestingParser(input)
	//repl.ScanTokens(input)
}
