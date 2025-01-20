package repl

import (
	"bufio"
	"fmt"
	"io"
	"pir-interpreter/evaluator"
	"pir-interpreter/lexer"
	"pir-interpreter/object"
	"pir-interpreter/parser"
)

const PROMPT = "8^) "

func Start(in io.Reader, out io.Writer) {
	fmt.Printf("Starting the interactive pir interpreter ye dirty seadog...\n")
	scanner := bufio.NewScanner(in)
	ns := object.NewNamespace()
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		if line == "exit" || line == "bye" {
			io.WriteString(out, "Goodbye my friend\n")
			return
		}
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, ns)
		if evaluated != nil && evaluated != evaluator.MT {
			io.WriteString(out, evaluated.AsString())
			io.WriteString(out, "\n")
		}
	}
}
func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
