//go:build !js && !wasm
// +build !js,!wasm

package main

import (
	"flag"
	"fmt"
	"io"
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
		fmt.Printf("Error while reading input file: %s", err)
		os.Exit(1)
	}

	l := lexer.New(string(code))
	p := parser.New(l)
	programTreeRoot := p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Println("Errors while parsing program:")
		printParserErrors(os.Stdout, p.Errors())
		os.Exit(1)
	}

	ns := object.NewNamespace()
	evaluated := evaluator.Eval(programTreeRoot, ns)
	if evaluated.Type() != object.MT_OBJ {
		fmt.Println(evaluated.AsString())
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
