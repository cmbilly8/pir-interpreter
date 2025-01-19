package main

import (
	"flag"
	"fmt"
	"os"
	"pir-interpreter/repl"
	"pir-interpreter/visualizer"
)

func main() {
	visualize := flag.Bool("ast", false, "Generate and open the AST visualization PNG")
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

	if *visualize {
		visualizer.Visualize(string(code))
	}
}
