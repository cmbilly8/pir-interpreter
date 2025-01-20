package parser

// Using Pratt's Top Down Parsing

import (
	"fmt"
	"pir-interpreter/ast"
	"pir-interpreter/lexer"
	"pir-interpreter/token"
	"strconv"
)

type (
	prefixParseFunc func() ast.Expression
	infixParseFunc  func(ast.Expression) ast.Expression
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.curToken = p.l.NextToken()
	p.peekToken = p.l.NextToken()
	return p
}

func (p *Parser) resolvePrefixParseFunc(tok token.TokenType) prefixParseFunc {
	switch tok {
	case token.IDENT:
		return p.parseIdentifier
	case token.INT:
		return p.parseIntegerLiteral
	case token.AAAA:
		return p.parsePrefixExpression
	case token.MINUS:
		return p.parsePrefixExpression
	case token.TRUE:
		return p.parseBoolean
	case token.FALSE:
		return p.parseBoolean
	case token.LPAREN:
		return p.parseGroupedExpression
	case token.LBRACKET:
		return p.parseArrayLiteral
	case token.LBRACE:
		return p.parseHashMapLiteral
	case token.F:
		return p.parseFunctionLiteral
	case token.STRING:
		return p.parseStringLiteral
	default:
		return nil
	}
}

func (p *Parser) resolveInfixParseFunc(tok token.TokenType) infixParseFunc {
	switch tok {
	case token.PLUS:
		return p.parseInfixExpression
	case token.MINUS:
		return p.parseInfixExpression
	case token.FSLASH:
		return p.parseInfixExpression
	case token.STAR:
		return p.parseInfixExpression
	case token.EQUAL:
		return p.parseInfixExpression
	case token.NOTEQUAL:
		return p.parseInfixExpression
	case token.AND:
		return p.parseInfixExpression
	case token.OR:
		return p.parseInfixExpression
	case token.LESS:
		return p.parseInfixExpression
	case token.GREATER:
		return p.parseInfixExpression
	case token.LPAREN:
		return p.parseCallExpression
	case token.LBRACKET:
		return p.parseIndexExpression
	default:
		return nil
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Next token expected: %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) advanceTokens() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) expectPeekToken(t token.TokenType) bool {
	if p.peekToken.Is(t) {
		p.advanceTokens()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) ParseProgram() *ast.Program {
	programNode := ast.NewEmptyProgram()

	var currentStatement ast.Statement
	for p.curToken.IsNot(token.EOF) {
		currentStatement = p.parseStatement()
		if currentStatement != nil {
			programNode.Statements = append(programNode.Statements, currentStatement)
		}
		p.advanceTokens()

	}
	return programNode
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.YAR:
		return p.parseYarStatement()
	case token.GIVES:
		return p.parseGivesStatement()
	case token.IF:
		return p.parseIfStatement()
	// If we start with an expression
	default:
		startToken := p.curToken
		expr := p.parseExpression(token.PREC_LOWEST)
		if indexAssign, ok := expr.(*ast.IndexExpression); ok && p.peekToken.Is(token.BE) {
			return p.parseIndexAssignment(startToken, indexAssign)
		}

		if p.peekToken.Is(token.PERIOD) {
			p.advanceTokens()
		}
		return &ast.ExpressionStatement{Expression: expr}
	}
}

func (p *Parser) parseIndexAssignment(startToken token.Token, indexAssign *ast.IndexExpression) *ast.IndexAssignment {
	p.advanceTokens()
	p.advanceTokens()
	value := p.parseExpression(token.PREC_LOWEST)
	if p.peekToken.Is(token.PERIOD) {
		p.advanceTokens()
	}
	return &ast.IndexAssignment{
		Token: startToken,
		Left:  indexAssign.Left,
		Index: indexAssign.Index,
		Value: value,
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionCollection(token.RBRACKET)
	return array
}

func (p *Parser) parseHashMapLiteral() ast.Expression {
	hml := &ast.HashMapLiteral{Token: p.curToken}
	hml.MP = make(map[ast.Expression]ast.Expression)
	for p.peekToken.IsNot(token.RBRACE) {
		p.advanceTokens()
		key := p.parseExpression(token.PREC_LOWEST)
		if !p.expectPeekToken(token.COLOGNE) {
			return nil
		}
		p.advanceTokens()
		value := p.parseExpression(token.PREC_LOWEST)
		hml.MP[key] = value
		if !p.peekToken.Is(token.RBRACE) && !p.expectPeekToken(token.COMMA) {
			return nil
		}
	}
	if !p.expectPeekToken(token.RBRACE) {
		return nil
	}
	return hml
}

func (p *Parser) parseExpressionCollection(endToken token.TokenType) []ast.Expression {
	expressions := []ast.Expression{}
	if p.peekToken.Is(endToken) {
		p.advanceTokens()
		return expressions
	}
	p.advanceTokens()
	expressions = append(expressions, p.parseExpression(token.PREC_LOWEST))
	for p.peekToken.Is(token.COMMA) {
		p.advanceTokens()
		p.advanceTokens()
		expressions = append(expressions, p.parseExpression(token.PREC_LOWEST))
	}
	if !p.expectPeekToken(endToken) {
		return nil
	}
	return expressions
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curToken.Is(token.TRUE)}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fLiteral := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeekToken(token.LPAREN) {
		return nil
	}

	fLiteral.Params = p.parseFunctionParams()

	if !p.expectPeekToken(token.COLOGNE) {
		return nil
	}
	p.advanceTokens()

	fLiteral.Body = p.parseBlockStatement()
	return fLiteral
}

func (p *Parser) parseFunctionParams() []*ast.Identifier {
	params := []*ast.Identifier{}
	if p.peekToken.Is(token.RPAREN) {
		p.advanceTokens()
		return params
	}
	p.advanceTokens()

	params = append(params, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	for p.peekToken.Is(token.COMMA) {
		p.advanceTokens()
		p.advanceTokens()
		params = append(params, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	if !p.expectPeekToken(token.RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionCollection(token.RPAREN)
	return exp
}

func (p *Parser) parseIndexExpression(collection ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: collection}
	p.advanceTokens()
	exp.Index = p.parseExpression(token.PREC_LOWEST)
	if !p.expectPeekToken(token.RBRACKET) {
		return nil
	}
	return exp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	bs := &ast.BlockStatement{Token: p.curToken, Statements: make([]ast.Statement, 0)}
	// LS, LSIF, PERIOD, EOF
	for !p.curToken.IsBlockTerminator() {
		bs.Statements = append(bs.Statements, p.parseStatement())
		p.advanceTokens()
	}

	if p.curToken.Is(token.GIVES) {
		bs.Statements = append(bs.Statements, p.parseGivesStatement())

	}

	return bs

}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	statement := &ast.IfStatement{Token: p.curToken}

	//Handle ifs and else ifs
	for p.curToken.IsNot(token.PERIOD) && p.curToken.IsNot(token.LS) {
		if p.curToken.IsNot(token.IF) && p.curToken.IsNot(token.LSIF) {
			return nil
		}
		currConditional := ast.Conditional{Token: p.curToken}
		p.advanceTokens()
		currConditional.Condition = p.parseExpression(token.PREC_LOWEST)
		if !p.expectPeekToken(token.COLOGNE) {
			return nil
		}
		p.advanceTokens()
		currConditional.Consequence = p.parseBlockStatement()
		statement.Conditionals = append(statement.Conditionals, currConditional)
	}

	if p.curToken.Is(token.LS) {
		if !p.expectPeekToken(token.COLOGNE) {
			return nil
		}
		p.advanceTokens()
		statement.Alternate = p.parseBlockStatement()
	} else {
		statement.Alternate = nil
	}

	return statement
}

func (p *Parser) parseYarStatement() *ast.YarStatement {
	statement := &ast.YarStatement{Token: p.curToken}

	if !p.expectPeekToken(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeekToken(token.BE) {
		return nil
	}

	p.advanceTokens()

	statement.Value = p.parseExpression(token.PREC_LOWEST)

	if p.peekToken.Is(token.PERIOD) {
		p.advanceTokens()
	}

	return statement
}

func (p *Parser) parseGivesStatement() *ast.GivesStatement {
	statement := &ast.GivesStatement{Token: p.curToken}

	p.advanceTokens()

	if p.curToken.IsNot(token.PERIOD) {
		statement.Value = p.parseExpression(token.PREC_LOWEST)
	} else {
		statement.Value = nil
	}

	p.advanceTokens()

	return statement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(token.PREC_LOWEST)
	if p.peekToken.Is(token.PERIOD) {
		p.advanceTokens()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFunc := p.resolvePrefixParseFunc(p.curToken.Type)
	if prefixFunc == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefixFunc()

	for precedence < p.peekToken.Precedence() && !p.peekToken.IsExpressionTerminator() {
		infixFunc := p.resolveInfixParseFunc(p.peekToken.Type)
		if infixFunc == nil {
			return leftExp
		}
		p.advanceTokens()
		leftExp = infixFunc(leftExp)
	}
	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.advanceTokens()
	expression.Right = p.parseExpression(token.PREC_PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curToken.Precedence()
	p.advanceTokens()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advanceTokens()
	exp := p.parseExpression(token.PREC_LOWEST)
	if !p.expectPeekToken(token.RPAREN) {
		return nil
	}

	return exp
}
