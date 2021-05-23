package main

import (
	"CornyLang/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)

	// input := `let trabajador = class
	// {
	// 	calcularHoras = fn(horas, dias, precio){
	// 		let total = horas * dias * precio;
	// 		return total;
	// 	}
	// };
	// trabajador.type();`
	//input := `let persona = class { nombre = "miguel"; }; persona.apellido;`
	//repl.TestEvaluator(input)
	//repl.TestingParser(input)
	//repl.ScanTokens(input)
}
