package evaluator

import (
	"fmt"
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
		{"(1 > 2) and nay", false},
		{"ay and nay", false},
		{"nay and nay", false},
		{"ay and ay", true},
		{"ay or nay", true},
		{"nay or ay", true},
		{"nay or nay", false},
	}
	for _, tt := range tests {
		fmt.Println(tt.input)
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
		{
			`"ay" - "matey"`,
			"unknown operator: STRING - STRING",
		},
		{
			"[1, 2, 3][4]",
			"index out of bounds. len=3, index=4",
		},
		{
			"[1, 2, 3][-1]",
			"index out of bounds. len=3, index=-1",
		},
		{
			`{"a": "b"}[f(x): x..]`,
			"Object not hashable. Type=FUNCTION",
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

func TestFunctionObject(t *testing.T) {
	input := "f(x):  x + 2.."
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Params) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Params)
	}
	if fn.Params[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Params[0])
	}
	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"yar identity be f(x): x.. identity(5)", 5},
		{"yar identity be f(x): gives x.. identity(5).", 5},
		{"yar double be f(x): x * 2.. double(5).", 10},
		{"yar add be f(x, y): x + y.. add(5, 5).", 10},
		{"yar add be f(x, y): x + y.. add(5 + 5, add(5, 5)).", 20},
		{"f(x): x..(5)", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
    yar foo be f(x):
        f(y): 
            x + y.
            .
        .
    .
    yar bar be foo(2).
    bar(2).`
	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"ay matey?"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "ay matey?" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"ay" + " " + "matey"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "ay matey" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringCompare(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"'hello' = 'hello'", true},
		{"'hello' = \"hello\"", true},
		{"'hello' = 'hlo'", false},
		{"'hello' != 'hello'", false},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBuiltins(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INT"},
		{`len("one", "two")`, "wrong number of args. got=2, expected=1"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}
	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"yar i be 0. [1][i].",
			1,
		},
		{
			"[1, 2, 3][1 + 1].",
			3,
		},
		{
			"yar arrray be [1, 2, 3]. arrray[2].",
			3,
		},
		{
			"yar a be [1, 2, 3]. a[0] + a[1] + a[2].",
			6,
		},
		{
			"yar a be [1, 2, 3]. yar i be a[0]. a[i].",
			2,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testMTObject(t, evaluated)
		}
	}
}

func testMTObject(t *testing.T, obj object.Object) bool {
	if obj != MT {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestHashLiterals(t *testing.T) {
	input := `yar two be "two".
        {
        "one": 10 - 9,
        two: 1 + 1,
        "thr" + "ee": 6 / 2,
        4: 4,
        ay: 5,
        nay: 6
        }`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.HashMap)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).Hash():   1,
		(&object.String{Value: "two"}).Hash():   2,
		(&object.String{Value: "three"}).Hash(): 3,
		(&object.Int{Value: 4}).Hash():          4,
		AY.Hash():                               5,
		NAY.Hash():                              6,
	}
	if len(result.MP) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.MP))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := result.MP[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"hi": 5}["hi"]`,
			5,
		},
		{
			`{"hi": 5}["ho"]`,
			nil,
		},
		{
			`yar key be "yo". {"yo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{ay: 5}[ay]`,
			5,
		},
		{
			`{ay: 5}[ay]`,
			5,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testMTObject(t, evaluated)
		}
	}
}

func TestIndexAssign(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`yar x be {"hi": 5}. x['hi'] be 6. x['hi'].`,
			6,
		},
		{
			`yar x be {"hi": 5}. x['ho'] be 6. x['ho'].`,
			6,
		},
		{
			`yar x be {"hi": 5}. x['he'] be 6. x['hi']`,
			5,
		},
		{
			`yar x be [1,2,3]. x[0] be 6. x[0].`,
			6,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testMTObject(t, evaluated)
		}
	}
}
