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

func parseProgramFromInput(input string) (*ast.Program, *Parser) {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	return program, p
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
		program, p := parseProgramFromInput(tt.input)
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

func TestIndexAssignStatements(t *testing.T) {
	tests := []struct {
		input                        string
		expectedCollectionIdentifier string
		expectedIndex                any
		expectedValue                any
	}{
		{"x['1'] be 5.", "x", "1", 5},
		{"b[1] be '3'.", "b", 1, "3"},
	}
	for _, tt := range tests {
		program, p := parseProgramFromInput(tt.input)
		printErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.IndexAssignment)
		if !ok {
			t.Fatalf("program.Statements[0] is not an ast.IndexAssignment, got=%T", program.Statements[0])
		}

		id, ok := stmt.Left.(*ast.Identifier)
		if !ok {
			t.Fatalf("stmt.Left is not an ast.Identifier, got=%T", stmt.Left)
		}
		if id.Value != tt.expectedCollectionIdentifier {
			t.Fatalf("IndexAssignment collection identifier not %q, got=%q", tt.expectedCollectionIdentifier, id.Value)
		}

		switch index := stmt.Index.(type) {
		case *ast.IntegerLiteral:
			expectedIndex, ok := tt.expectedIndex.(int)
			if !ok {
				t.Fatalf("Expected index type mismatch: expected int, got %T", tt.expectedIndex)
			}
			if index.Value != int64(expectedIndex) {
				t.Fatalf("Expected index to be %d, got %d", expectedIndex, index.Value)
			}
		case *ast.StringLiteral:
			expectedIndex, ok := tt.expectedIndex.(string)
			if !ok {
				t.Fatalf("Expected index type mismatch: expected string, got %T", tt.expectedIndex)
			}
			if index.Value != expectedIndex {
				t.Fatalf("Expected index to be %q, got %q", expectedIndex, index.Value)
			}
		default:
			t.Fatalf("Unexpected index type, got=%T", stmt.Index)
		}

		switch v := stmt.Value.(type) {
		case *ast.IntegerLiteral:
			expectedValue, ok := tt.expectedValue.(int)
			if !ok {
				t.Fatalf("Expected value type mismatch: expected=%T, got=%T", tt.expectedValue, v.Value)
			}
			if int(v.Value) != expectedValue {
				t.Fatalf("IndexAssignment value not %v, got=%v", expectedValue, v.Value)
			}
		case *ast.StringLiteral:
			expectedValue, ok := tt.expectedValue.(string)
			if !ok {
				t.Fatalf("Expected value type mismatch: expected=%T, got=%T", tt.expectedValue, v.Value)
			}
			if v.Value != expectedValue {
				t.Fatalf("IndexAssignment value not %q, got=%q", expectedValue, v.Value)
			}
		default:
			t.Fatalf("stmt.Value is not a valid literal type, got=%T", stmt.Value)
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
		program, p := parseProgramFromInput(tt.input)
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
	program, p := parseProgramFromInput(input)
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
	program, p := parseProgramFromInput(input)
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
		program, p := parseProgramFromInput(tt.input)
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
		program, p := parseProgramFromInput(tt.input)
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
			"ay and ay or nay",
			"((ay and ay) or nay)",
		},
		{
			"ay or nay and nay",
			"(ay or (nay and nay))",
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
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"foo(a * b[2], b[1], 2 * [1, 2][1])",
			"foo((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}
	for _, tt := range tests {
		program, p := parseProgramFromInput(tt.input)
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
		program, p := parseProgramFromInput(tt.input)
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
	program, p := parseProgramFromInput(input)
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
	program, p := parseProgramFromInput(input)
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
	program, p := parseProgramFromInput(input)
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
		program, p := parseProgramFromInput(tt.input)
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
	program, p := parseProgramFromInput(input)
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
	program, p := parseProgramFromInput(input)
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

func TestParsingArrayLiterals(t *testing.T) {
	input := "[0, 3 + 3, 9*2]."
	program, p := parseProgramFromInput(input)
	printErrors(t, p)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}
	testIntegerLiteral(t, array.Elements[0], 0)
	testInfixExpression(t, array.Elements[1], 3, "+", 3)
	testInfixExpression(t, array.Elements[2], 9, "*", 2)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "arrrray[1 * 2]"
	program, p := parseProgramFromInput(input)
	printErrors(t, p)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}
	if !testIdentifier(t, indexExp.Left, "arrrray") {
		return
	}
	if !testInfixExpression(t, indexExp.Index, 1, "*", 2) {
		return
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	program, p := parseProgramFromInput(input)
	printErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashMapLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashMapLiteral. got=%T", stmt.Expression)
	}
	if len(hash.MP) != 3 {
		t.Errorf("hash.MP has wrong length. got=%d", len(hash.MP))
	}
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for key, value := range hash.MP {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashMapLiteral(t *testing.T) {
	input := "{}"
	program, p := parseProgramFromInput(input)
	printErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashMapLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.MP) != 0 {
		t.Errorf("hash.MP has wrong length. got=%d", len(hash.MP))
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`
	program, p := parseProgramFromInput(input)
	printErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashMapLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashMapLiteral. got=%T", stmt.Expression)
	}
	if len(hash.MP) != 3 {
		t.Errorf("hash.MP has wrong length. got=%d", len(hash.MP))
	}
	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}
	for key, value := range hash.MP {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}
		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}
		testFunc(value)
	}
}

func TestForStatementParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "for x < 10: x..",
			expected: "for (x < 10): (x.)",
		},
		{
			input:    "for ay: x+y. y..",
			expected: "for ay: ((x + y).y.)",
		},
	}

	for _, tt := range tests {
		program, p := parseProgramFromInput(tt.input)
		printErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ForStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ForStatement. got=%T", program.Statements[0])
		}

		actual := stmt.String()
		if actual != tt.expected {
			t.Errorf("ForStatement.String() mismatch. Expected=%q, Got=%q", tt.expected, actual)
		}
	}
}
