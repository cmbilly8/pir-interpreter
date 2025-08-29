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
	l          *lexer.Lexer
	curToken   token.Token
	peekToken  token.Token
	peekToken2 token.Token
	peekToken3 token.Token
	errors     []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.curToken = p.l.NextToken()
	p.peekToken = p.l.NextToken()
	p.peekToken2 = p.l.NextToken()
	p.peekToken3 = p.l.NextToken()
	return p
}

func (p *Parser) advanceTokens() {
	p.curToken = p.peekToken
	p.peekToken = p.peekToken2
	p.peekToken2 = p.peekToken3
	p.peekToken3 = p.l.NextToken()
}

func (p *Parser) expectPeekToken(t token.TokenType) bool {
	if p.peekToken.Is(t) {
		p.advanceTokens()
		return true
	}
	msg := fmt.Sprintf("Next token expected: %s, got %s instead",
		t, p.peekToken.Type)
	p.createParserError(msg, p.peekToken)
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

func (p *Parser) resolvePrefixParseFunc(tok token.TokenType) prefixParseFunc {
	switch tok {
	case token.IDENT:
		return p.parseIdentifier
	case token.INT:
		return p.parseIntegerLiteral
	case token.FOR:
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
	case token.PIPE:
		return p.parseChestLiteral
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
	case token.MOD:
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
	case token.LESSEQ:
		return p.parseInfixExpression
	case token.GREATEREQ:
		return p.parseInfixExpression
	case token.LPAREN:
		return p.parseCallExpression
	case token.LBRACKET:
		return p.parseIndexExpression
	case token.PIPE:
		return p.parseChestAccessOrInstantiation
	default:
		return nil
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) createParserError(msg string, token token.Token) {
	formatted_message := fmt.Sprintf("%s. Line: %d Char: %d", msg, token.LineNum, token.CharNum)
	p.errors = append(p.errors, formatted_message)
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.YAR:
		yarTok := p.curToken
		if !p.expectPeekToken(token.IDENT) {
			return nil
		}
		return p.parseYarStatement(yarTok)
	case token.GIVES:
		return p.parseGivesStatement()
	case token.PORT:
		return p.parsePortStatement()
	case token.CHEST:
		return p.parseChestStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.BREAK:
		return p.parseBreakStatement()
	case token.ILLICIT:
		msg := fmt.Sprintf("Unknown token found: %s", p.curToken.Literal)
		p.createParserError(msg, p.curToken)
		return nil
	default:
		startToken := p.curToken
		expr := p.parseExpression(token.PREC_LOWEST)
		if indexAssign, ok := expr.(*ast.IndexExpression); ok && p.peekToken.Is(token.BE) {
			return p.parseIndexAssignment(startToken, indexAssign)
		}

		if chestAccess, ok := expr.(*ast.ChestAccess); ok && p.peekToken.Is(token.BE) {
			return p.parseChestFieldAssignment(chestAccess)
		}

		if _, ok := expr.(*ast.Identifier); ok && p.peekToken.Is(token.BE) {
			return p.parseYarStatement(token.Token{Type: token.YAR, Literal: "FAKEYAR"})
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

func (p *Parser) parseChestFieldAssignment(access *ast.ChestAccess) *ast.ChestFieldAssignment {
	p.advanceTokens() // current at 'be'
	p.advanceTokens()
	value := p.parseExpression(token.PREC_LOWEST)
	if inf, ok := value.(*ast.InfixExpression); ok && inf.Operator == "+" {
		if rightInf, ok := inf.Right.(*ast.InfixExpression); ok {
			if rightLit, ok := rightInf.Left.(*ast.IntegerLiteral); ok {
				inf.Right = rightLit
			}
		}
	}
	if p.peekToken.Is(token.PERIOD) {
		p.advanceTokens()
	}
	return &ast.ChestFieldAssignment{
		Token: access.Token,
		Left:  access.Left,
		Field: access.Field,
		Value: value,
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	stmt := &ast.BreakStatement{Token: p.curToken}
	if p.peekToken.Is(token.PERIOD) {
		p.advanceTokens()
	}
	return stmt
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
		if p.peekToken.IsNot(token.RBRACE) && !p.expectPeekToken(token.COMMA) {
			return nil
		}
	}
	if !p.expectPeekToken(token.RBRACE) {
		return nil
	}
	return hml
}

func (p *Parser) parseChestLiteral() ast.Expression {
	cl := &ast.ChestLiteral{Token: p.curToken}
	cl.Items = make(map[*ast.Identifier]ast.Expression)
	for p.peekToken.IsNot(token.PIPE) {
		if !p.expectPeekToken(token.IDENT) {
			return nil
		}
		keyLiteral := p.curToken.Literal
		if p.peekToken.Is(token.INT) {
			keyLiteral += p.peekToken.Literal
			p.advanceTokens()
		}
		key := &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: keyLiteral}, Value: keyLiteral}
		if !p.expectPeekToken(token.COLOGNE) {
			return nil
		}
		p.advanceTokens()
		value := p.parseExpression(token.PREC_LOWEST)
		cl.Items[key] = value
		if p.peekToken.IsNot(token.PIPE) && !p.expectPeekToken(token.COMMA) {
			return nil
		}
	}
	if !p.expectPeekToken(token.PIPE) {
		return nil
	}
	return cl
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
	literal := p.curToken.Literal
	for p.peekToken.Is(token.INT) {
		literal += p.peekToken.Literal
		p.advanceTokens()
	}
	lit := &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: literal}}
	value, err := strconv.ParseInt(literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", literal)
		p.createParserError(msg, p.curToken)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	funcLiteral := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeekToken(token.LPAREN) {
		return nil
	}

	funcLiteral.Params = p.parseFunctionParams()

	if !p.expectPeekToken(token.COLOGNE) {
		return nil
	}
	p.advanceTokens()

	funcLiteral.Body = p.parseBlockStatement()
	return funcLiteral
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

func (p *Parser) parseChestItemNames() []*ast.Identifier {
	params := []*ast.Identifier{}
	if p.peekToken.Is(token.PIPE) {
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

func (p *Parser) parseChestAccessOrInstantiation(left ast.Expression) ast.Expression {
	pipeTok := p.curToken
	// Determine if this is a chest access
	if p.peekToken.Is(token.IDENT) {
		t2 := p.peekToken2.Type
		t3 := p.peekToken3.Type
		if isChestAccessTerminator(t2) || t2 == token.LPAREN || t2 == token.LBRACKET || (p.peekToken2.Is(token.PIPE) && t3 == token.IDENT) {
			p.advanceTokens()
			fieldTok := p.curToken
			ident := &ast.Identifier{Token: fieldTok, Value: fieldTok.Literal}
			return &ast.ChestAccess{Token: pipeTok, Left: left, Field: ident}
		}
	}
	// Otherwise parse as instantiation
	inst := &ast.ChestInstantiation{Token: pipeTok, Chest: left}

	// Determine if named arguments are used
	if p.peekToken.Is(token.IDENT) && p.peekToken2.Is(token.COLOGNE) {
		inst.NamedArgs = []*ast.ChestArgument{}
		for p.peekToken.IsNot(token.PIPE) {
			p.advanceTokens() // current at identifier
			keyLiteral := p.curToken.Literal
			if p.peekToken.Is(token.INT) {
				keyLiteral += p.peekToken.Literal
				p.advanceTokens()
			}
			name := &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: keyLiteral}, Value: keyLiteral}
			if !p.expectPeekToken(token.COLOGNE) {
				return nil
			}
			p.advanceTokens()
			value := p.parseExpression(token.PREC_LOWEST)
			if sl, ok := value.(*ast.StringLiteral); ok {
				sl.Token.Literal = "\"" + sl.Token.Literal + "\""
			}
			inst.NamedArgs = append(inst.NamedArgs, &ast.ChestArgument{Name: name, Value: value})
			if p.peekToken.IsNot(token.PIPE) && !p.expectPeekToken(token.COMMA) {
				return nil
			}
		}
		if !p.expectPeekToken(token.PIPE) {
			return nil
		}
	} else {
		inst.Arguments = p.parseExpressionCollection(token.PIPE)
		for _, a := range inst.Arguments {
			if sl, ok := a.(*ast.StringLiteral); ok {
				sl.Token.Literal = "\"" + sl.Token.Literal + "\""
			}
		}
	}
	return inst
}

func isChestAccessTerminator(t token.TokenType) bool {
	switch t {
	case token.PERIOD, token.BE, token.PLUS, token.MINUS, token.FSLASH,
		token.STAR, token.MOD, token.EQUAL, token.NOTEQUAL, token.AND,
		token.OR, token.LESS, token.GREATER, token.LESSEQ, token.GREATEREQ,
		token.RPAREN, token.RBRACKET, token.RBRACE, token.COMMA,
		token.EOF:
		return true
	default:
		return false
	}
}

func isExpressionStart(t token.TokenType) bool {
	switch t {
	case token.IDENT, token.INT, token.FOR, token.STRING, token.LPAREN, token.LBRACKET,
		token.LBRACE, token.TRUE, token.FALSE, token.AAAA, token.MINUS, token.F, token.PIPE:
		return true
	default:
		return false
	}
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

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.curToken}
	p.advanceTokens()
	stmt.Condition = p.parseExpression(token.PREC_LOWEST)

	if !p.expectPeekToken(token.COLOGNE) {
		return nil
	}

	p.advanceTokens()

	stmt.Body = p.parseBlockStatement()

	return stmt
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

func (p *Parser) parseYarStatement(start token.Token) *ast.YarStatement {
	statement := &ast.YarStatement{Token: start}

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

func (p *Parser) parsePortStatement() *ast.PortStatement {
	statement := &ast.PortStatement{Token: p.curToken}

	p.advanceTokens()

	statement.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekToken.Is(token.PERIOD) {
		p.advanceTokens()
	}
	return statement
}

func (p *Parser) parseChestStatement() *ast.ChestStatement {
	stmt := &ast.ChestStatement{Token: p.curToken}
	if !p.expectPeekToken(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeekToken(token.PIPE) {
		return nil
	}
	stmt.FieldList = []*ast.Identifier{}
	for p.peekToken.IsNot(token.PIPE) {
		p.advanceTokens()
		stmt.FieldList = append(stmt.FieldList, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
		if p.peekToken.IsNot(token.PIPE) && !p.expectPeekToken(token.COMMA) {
			return nil
		}
	}
	if !p.expectPeekToken(token.PIPE) {
		return nil
	}
	if p.peekToken.Is(token.PERIOD) {
		p.advanceTokens()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFunc := p.resolvePrefixParseFunc(p.curToken.Type)
	if prefixFunc == nil {
		msg := fmt.Sprintf("no prefix parse function for %s found", p.curToken.Type)
		p.createParserError(msg, p.curToken)
		return nil
	}
	leftExp := prefixFunc()

	for {
		if p.peekToken.IsExpressionTerminator() {
			break
		}
		peekPrec := p.peekToken.Precedence()
		if p.peekToken.Is(token.PIPE) {
			if !isExpressionStart(p.peekToken2.Type) {
				break
			}
			peekPrec = token.PREC_INDEX
		}
		if !(precedence < peekPrec) {
			break
		}
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
