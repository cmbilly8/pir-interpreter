//go:build js && wasm
// +build js,wasm

package main

import (
	"pir-interpreter/evaluator"
	"pir-interpreter/lexer"
	"pir-interpreter/object"
	"pir-interpreter/parser"
	"pir-interpreter/writer"
	"syscall/js"
)

func evaluate(code string) {
	ns := object.NewNamespace()
	l := lexer.New(string(code))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		errors := p.Errors()
		for _, msg := range errors {
			writer.WriteOutput("\t" + msg + "\n")
		}
	}
	evaluated := evaluator.Eval(program, ns)
	if evaluated == evaluator.MT {
		return
	}
	writer.WriteOutput(evaluated.AsString())
}

func processInput(_ js.Value, args []js.Value) interface{} {
	writer.ClearOutput()

	if len(args) < 1 {
		writer.WriteOutput("Error: No input provided\n")
		return nil
	}

	input := args[0].String()
	evaluate(input)

	return nil
}

func getOutput(_ js.Value, _ []js.Value) interface{} {
	return js.ValueOf(writer.GetOutput())
}

func main() {
	js.Global().Set("processInput", js.FuncOf(processInput))
	js.Global().Set("getOutput", js.FuncOf(getOutput))

	select {}
}
