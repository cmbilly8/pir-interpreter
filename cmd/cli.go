//go:build !js && !wasm
// +build !js,!wasm

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
	"pir-interpreter/writer"
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
		fmt.Printf("Error while reading input file: %s", err)
		os.Exit(1)
	}

	l := lexer.New(string(code))
	p := parser.New(l)
	programTreeRoot := p.ParseProgram()

	if len(p.Errors()) != 0 {
		errors := p.Errors()
		for _, msg := range errors {
			writer.WriteOutput("\t" + msg + "\n")
		}
	}

	ns := object.NewNamespace()
	evaluated := evaluator.Eval(programTreeRoot, ns)
	if evaluated.Type() != object.MT_OBJ {
		writer.WriteOutput(evaluated.AsString())
	}
	fmt.Print(writer.GetOutput())
}
