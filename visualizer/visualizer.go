package visualizer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"pir-interpreter/ast"
	"pir-interpreter/lexer"
	"pir-interpreter/parser"
)

type NodeIDGenerator struct {
	id int
}

func (gen *NodeIDGenerator) NextID() string {
	gen.id++
	return fmt.Sprintf("node%d", gen.id)
}

func visualizeAST(node ast.Node, gen *NodeIDGenerator, out *os.File, parent string) {
	if node == nil {
		return
	}

	nodeID := gen.NextID()
	label := fmt.Sprintf("%q", node.TokenLiteral())
	fmt.Fprintf(out, "  %s [label=%s];\n", nodeID, label)

	if parent != "" {
		fmt.Fprintf(out, "  %s -> %s;\n", parent, nodeID)
	}

	switch n := node.(type) {
	case *ast.Program:
		for _, stmt := range n.Statements {
			visualizeAST(stmt, gen, out, nodeID)
		}
	case *ast.YarStatement:
		visualizeAST(n.Name, gen, out, nodeID)
		visualizeAST(n.Value, gen, out, nodeID)
	case *ast.GivesStatement:
		visualizeAST(n.Value, gen, out, nodeID)
	case *ast.ExpressionStatement:
		visualizeAST(n.Expression, gen, out, nodeID)
	case *ast.IfStatement:
		for _, cond := range n.Conditionals {
			visualizeAST(&cond, gen, out, nodeID)
		}
		visualizeAST(n.Alternate, gen, out, nodeID)
	case *ast.Conditional:
		visualizeAST(n.Condition, gen, out, nodeID)
		visualizeAST(n.Consequence, gen, out, nodeID)
	case *ast.BlockStatement:
		if n != nil {
			for _, stmt := range n.Statements {
				visualizeAST(stmt, gen, out, nodeID)
			}
		}
	case *ast.Identifier:
		// Leaf node, nothing to recurse
	case *ast.IntegerLiteral:
		// Leaf node, nothing to recurse
	case *ast.Boolean:
		// Leaf node, nothing to recurse
	case *ast.PrefixExpression:
		visualizeAST(n.Right, gen, out, nodeID)
	case *ast.InfixExpression:
		visualizeAST(n.Left, gen, out, nodeID)
		visualizeAST(n.Right, gen, out, nodeID)
	case *ast.FunctionLiteral:
		for _, param := range n.Params {
			visualizeAST(param, gen, out, nodeID)
		}
		visualizeAST(n.Body, gen, out, nodeID)
	case *ast.CallExpression:
		visualizeAST(n.Function, gen, out, nodeID)
		for _, arg := range n.Arguments {
			visualizeAST(arg, gen, out, nodeID)
		}
	}
}

func Visualize(input string) {
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()

	dotFileName := "ast.dot"
	dotFile, err := os.Create(dotFileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer dotFile.Close()

	fmt.Fprintln(dotFile, "digraph AST {")
	fmt.Fprintln(dotFile, "  node [shape=box];")

	gen := &NodeIDGenerator{}
	visualizeAST(program, gen, dotFile, "")

	fmt.Fprintln(dotFile, "}")
	pngFileName := "ast.png"
	cmd := exec.Command("dot", "-Tpng", dotFileName, "-o", pngFileName)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error generating PNG:", err)
		return
	}
	fmt.Println("AST visualization PNG written to", pngFileName)

	var openCmd *exec.Cmd
	switch {
	case filepath.Base(os.Getenv("OS")) == "Windows_NT":
		openCmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", pngFileName)
	case filepath.Base(os.Getenv("OSTYPE")) == "darwin": // macOS
		openCmd = exec.Command("open", pngFileName)
	default: // Assume Linux
		openCmd = exec.Command("xdg-open", pngFileName)
	}
	err = openCmd.Start()
	if err != nil {
		fmt.Println("Error opening PNG:", err)
	}
}
