package evaluator

import (
	"pir-interpreter/lexer"
	"pir-interpreter/object"
	"pir-interpreter/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"999", 999},
		{"0", 0},
		{"-10", -10},
		{"-0", 0},
		{"5 + 5", 10},
		{"1 + 2 - 4", -1},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}
func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	ns := object.NewNamespace()
	return Eval(program, ns)
}
func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Int)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"ay", true},
		{"nay", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 = 1", true},
		{"1 != 1", false},
		{"1 = 2", false},
		{"1 != 2", true},
		{"ay = ay", true},
		{"nay = nay", true},
		{"ay = nay", false},
		{"ay != nay", true},
		{"nay != ay", true},
		{"(1 < 2) = ay", true},
		{"(1 < 2) = nay", false},
		{"(1 > 2) = ay", false},
		{"(1 > 2) = nay", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Bool)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func TestAAAOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!ay", false},
		{"!nay", true},
		{"!!ay", true},
		{"!!nay", false},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}
func TestGivesStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"gives 1.", 1},
		{"gives 1. 2.", 1},
		{"gives 2 * 2. 5.", 4},
		{"1. gives 2*4. 2.", 8},
		{
			`
        if 10 > 1:
            if 1 > 2:
                gives 2.
            lsif 1 < 3:
                gives 3..

            gives 1..
        `,
			3,
		},
		{
			`
        if 10 > 1:
            if 1 > 2:
                gives 2.
            lsif 1 > 3:
                gives 3..

            gives 1..
        `,
			1,
		},
		{
			`
        if 10 > 1:
            if 1 > 2:
                gives 2.
            lsif 1 > 3:
                gives 3.
            ls:
                gives 10..
            gives 1..
        `,
			10,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + ay.",
			"type mismatch: INT + BOOL",
		},
		{
			"5 + ay. 5.",
			"type mismatch: INT + BOOL",
		},
		{
			"-ay",
			"unknown operator: -BOOL",
		},
		{
			"!3",
			"Unsupported operation !INT",
		},
		{
			"ay + nay.",
			"unknown operator: BOOL + BOOL",
		},
		{
			"5. ay + nay. 5",
			"unknown operator: BOOL + BOOL",
		},
		{
			"a",
			"Identifier not found: a",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"yar a be 5. a.", 5},
		{"yar a be 5 * 5. a.", 25},
		{"yar b be 5. yar imma be b. imma.", 5},
		{"yar a be 5. yar b be a. yar c be a + b + 5. c.", 15},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}
