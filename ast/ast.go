package ast

import (
	"bytes"
	"pir-interpreter/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func NewEmptyProgram() *Program {
	return &Program{Statements: []Statement{}}
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type YarStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ys *YarStatement) statementNode()       {}
func (ys *YarStatement) TokenLiteral() string { return ys.Token.Literal }
func (ys *YarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ys.TokenLiteral() + " ")
	out.WriteString(ys.Name.String())
	out.WriteString(" be ")
	if ys.Value != nil {
		out.WriteString(ys.Value.String())
	}
	out.WriteString(".")
	return out.String()
}

type GivesStatement struct {
	Token token.Token
	Value Expression
}

func (gs *GivesStatement) statementNode()       {}
func (gs *GivesStatement) TokenLiteral() string { return gs.Token.Literal }
func (gs *GivesStatement) String() string {
	var out bytes.Buffer
	out.WriteString(gs.TokenLiteral() + " ")
	if gs.Value != nil {
		out.WriteString(gs.Value.String())
	}
	out.WriteString(".")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type IfStatement struct {
	Token        token.Token
	Conditionals []Conditional
	Alternate    *BlockStatement
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	var out bytes.Buffer
	for _, c := range is.Conditionals {
		out.WriteString(c.String())
	}
	out.WriteString("ls: ")
	if is.Alternate != nil {
		out.WriteString(is.Alternate.String())
	}
	return out.String()
}

type Conditional struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
}

func (c *Conditional) TokenLiteral() string { return c.Token.Literal }
func (c *Conditional) String() string {
	var out bytes.Buffer
	out.WriteString(c.Token.Literal)
	out.WriteString(" (")
	out.WriteString(c.Condition.String())
	out.WriteString("): ")
	out.WriteString(c.Consequence.String())
	return out.String()

}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type BlockStatement struct {
	Token      token.Token // should be : since it starts a block
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type FunctionLiteral struct {
	Token  token.Token
	Params []*Identifier
	Body   *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Params {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
