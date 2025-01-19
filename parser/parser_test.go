package parser

import (
	"fmt"
	"pir-interpreter/ast"
	"pir-interpreter/lexer"
	"testing"
)

func printErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestYarStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"yar x be 5.", "x", 5},
		{"yar y be ay.", "y", true},
		{"yar foobar be y.", "foobar", "y"},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		printErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}
		stmt := program.Statements[0]
		if !testYarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
		val := stmt.(*ast.YarStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestGivesStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue any
	}{
		{"gives 5.", 5},
		{"gives ay.", true},
		{"gives foobar.", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		printErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		givesStatement, ok := stmt.(*ast.GivesStatement)
		if !ok {
			t.Fatalf("stmt not *ast.GivesStatement. got=%T", stmt)
		}
		if givesStatement.TokenLiteral() != "gives" {
			t.Fatalf("givesStatement.TokenLiteral not 'gives', got %q",
				givesStatement.TokenLiteral())
		}
		if !testLiteralExpression(t, givesStatement.Value, tt.expectedValue) {
			return
		}
	}
}

func testYarStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "yar" {
		t.Errorf("s.TokenLiteral not 'yar'. got=%q", s.TokenLiteral())
		return false
	}
	yarStatement, ok := s.(*ast.YarStatement)
	if !ok {
		t.Errorf("s not *ast.YarStatement. got=%T", s)
		return false
	}
	if yarStatement.Name.Value != name {
		t.Errorf("yarStatement.Name.Value not '%s'. got=%s", name, yarStatement.Name.Value)
		return false
	}
	if yarStatement.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, yarStatement.Name)
		return false
	}
	return true
}

func TestIdentifierExpression(t *testing.T) {
	input := "variableName."
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	printErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "variableName" {
		t.Errorf("ident.Value not %s. got=%s", "variableName", ident.Value)
	}
	if ident.TokenLiteral() != "variableName" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "variableName",
			ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5."
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	printErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
			literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    any
	}{
		{"!5.", "!", 5},
		{"-15.", "-", 15},
		{"!ay.", "!", true},
		{"!nay.", "!", false},
	}
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		printErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}
	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5.", 5, "+", 5},
		{"5 - 5.", 5, "-", 5},
		{"5 * 5.", 5, "*", 5},
		{"5 / 5.", 5, "/", 5},
		{"5 > 5.", 5, ">", 5},
		{"5 < 5.", 5, "<", 5},
		{"5 = 5.", 5, "=", 5},
		{"5 != 5.", 5, "!=", 5},
		{"ay = ay", true, "=", true},
		{"ay != nay", true, "!=", false},
		{"nay = nay", false, "=", false},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		printErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}
		if !testLiteralExpression(t, exp.Left, tt.leftValue) {
			return
		}
		if !testLiteralExpression(t, exp.Right, tt.rightValue) {
			return
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - g",
			"(((a + (b * c)) + (d / e)) - g)",
		},
		{
			"3 + 4. -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 = 3 < 4",
			"((5 > 4) = (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 = 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) = ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 = 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) = ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 = false",
			"((3 > 5) = false)",
		},
		{
			"3 < 5 = true",
			"((3 < 5) = true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(ay = ay)",
			"(!(ay = ay))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / a + g)",
			"add((((a + b) + ((c * d) / a)) + g))",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		printErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected any,
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolean(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any,
	operator string, right any) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func testBoolean(t *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}
	if boolean.Value != value {
		t.Errorf("boolean.Value not %T. got=%T", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != getBoolLiteral(value) {
		t.Errorf("boolean.TokenLiteral not %t. got=%s", value,
			boolean.TokenLiteral())
		return false
	}
	return true
}

func getBoolLiteral(b bool) string {
	if b {
		return "ay"
	}
	return "nay"
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"ay.", true},
		{"nay.", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		printErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}
		if boolean.Value != tt.expectedBoolean {
			t.Errorf("boolean.Value not %t. got=%t", tt.expectedBoolean,
				boolean.Value)
		}
	}
}

func TestFullConditionalStatement(t *testing.T) {
	input := `
    if x < y:
        x.
    lsif x > y:
        y.
    ls:
        10.
    .
    `
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	printErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ConditionalStatement. got=%T",
			program.Statements[0])
	}

	if len(stmt.Conditionals) != 2 {
		t.Errorf("Expected length of stmt.Conditionals to be 2. got=%d\n",
			len(stmt.Conditionals))
	}

	firstConditional := stmt.Conditionals[0]
	if !testInfixExpression(t, firstConditional.Condition, "x", "<", "y") {
		return
	}

	if len(firstConditional.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(firstConditional.Consequence.Statements))
	}

	firstConsequence, ok := firstConditional.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			firstConditional.Consequence.Statements[0])
	}

	if !testIdentifier(t, firstConsequence.Expression, "x") {
		return
	}

	secondConditional := stmt.Conditionals[1]
	if !testInfixExpression(t, secondConditional.Condition, "x", ">", "y") {
		return
	}

	if len(secondConditional.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(secondConditional.Consequence.Statements))
	}

	secondConsequence, ok := secondConditional.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			secondConditional.Consequence.Statements[0])
	}

	if !testIdentifier(t, secondConsequence.Expression, "y") {
		return
	}

	if len(stmt.Alternate.Statements) != 1 {
		t.Errorf("stmt.Alternate is not 1 statements. got=%d\n",
			len(stmt.Alternate.Statements))
	}

	alternative, ok := stmt.Alternate.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statments[0] is not ast.ExpressionStatement. got %T",
			stmt.Alternate.Statements[0])
	}
	if !testIntegerLiteral(t, alternative.Expression, 10) {
		return
	}
}

func TestIfStatement(t *testing.T) {
	input := `
    if x < y:
        x.
    .
    `
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	printErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ConditionalStatement. got=%T",
			program.Statements[0])
	}

	if len(stmt.Conditionals) != 1 {
		t.Errorf("Expected length of stmt.Conditionals to be 1. got=%d\n",
			len(stmt.Conditionals))
	}

	firstConditional := stmt.Conditionals[0]
	if !testInfixExpression(t, firstConditional.Condition, "x", "<", "y") {
		return
	}

	if len(firstConditional.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(firstConditional.Consequence.Statements))
	}

	firstConsequence, ok := firstConditional.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			firstConditional.Consequence.Statements[0])
	}

	if !testIdentifier(t, firstConsequence.Expression, "x") {
		return
	}

	if stmt.Alternate != nil {
		t.Errorf("stmt.Alternate is not nil")
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `f(x, y): x + y..`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	printErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}
	if len(function.Params) != 2 {
		t.Fatalf("function literap.nextToken()l parameters wrong. want 2, got=%d\n",
			len(function.Params))
	}
	testLiteralExpression(t, function.Params[0], "x")
	testLiteralExpression(t, function.Params[1], "y")
	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "f(): gives..", expectedParams: []string{}},
		{input: "f(x): gives..", expectedParams: []string{"x"}},
		{input: "f(x, y, z): gives..", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		printErrors(t, p)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)
		if len(function.Params) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Params))
		}
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Params[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)."
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	printErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}
	if !testIdentifier(t, exp.Function, "add") {
		return
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}
	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"ay mate.".`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	printErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != "ay mate." {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}
