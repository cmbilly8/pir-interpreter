package main

import (
	"flag"
	"fmt"
	"os"
	"pir-interpreter/evaluator"
	"pir-interpreter/lexer"
	"pir-interpreter/object"
	"pir-interpreter/parser"
	"pir-interpreter/repl"
)

func main() {
	startREPL := flag.Bool("r", false, "Start pir repl")
	flag.Parse()

	if *startREPL {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Error: Please provide a .pir file as an argument.")
		os.Exit(1)
	}
	fileName := args[0]
	code, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error: Problem reading file")
	}

	ns := object.NewNamespace()
	l := lexer.New(string(code))
	p := parser.New(l)
	program := p.ParseProgram()
	evaluated := evaluator.Eval(program, ns)
	print(evaluated.AsString())
}
