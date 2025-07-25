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

type PortStatement struct {
	Token token.Token
	Name  *Identifier
}

func (ps *PortStatement) statementNode()       {}
func (ps *PortStatement) TokenLiteral() string { return ps.Token.Literal }
func (ps *PortStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ps.TokenLiteral() + " ")
	if ps.Name != nil {
		out.WriteString(ps.Name.String())
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
	if is.Alternate != nil {
		out.WriteString("ls: ")
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

type BreakStatement struct {
	Token token.Token
}

func (b *BreakStatement) statementNode()       {}
func (b *BreakStatement) TokenLiteral() string { return b.Token.Literal }
func (b *BreakStatement) String() string       { return b.Token.Literal }

type BlockStatement struct {
	Token      token.Token // should be : since it starts a block
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
		out.WriteString(".")
	}
	out.WriteString(")")
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

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

type IndexAssignment struct {
	Token token.Token
	Left  Expression
	Index Expression
	Value Expression
}

func (ie *IndexAssignment) statementNode()       {}
func (ia *IndexAssignment) TokenLiteral() string { return ia.Token.Literal }
func (ia *IndexAssignment) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ia.Left.String())
	out.WriteString("[")
	out.WriteString(ia.Index.String())
	out.WriteString("] be ")
	out.WriteString(ia.Value.String())
	out.WriteString(")")
	return out.String()
}

type HashMapLiteral struct {
	Token token.Token
	MP    map[Expression]Expression
}

func (hml *HashMapLiteral) expressionNode()      {}
func (hml *HashMapLiteral) TokenLiteral() string { return hml.Token.Literal }
func (hml *HashMapLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range hml.MP {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type ForStatement struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString("4 ")
	out.WriteString(fs.Condition.String())
	out.WriteString(": ")
	out.Write([]byte(fs.Body.String()))
	return out.String()
}
