package evaluator

import (
	"monkey/internal/lexer"
	"monkey/internal/object"
	"monkey/internal/parser"
	"reflect"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5", 10},
		{"5 + 5 * 2", 15},
		{"5 + 5 * 2 / 2", 10},
		{"5 + (5 * 2) / 2", 10},
		{"(5 + (5 * 2)) / 2", 7},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}
func TestEvalBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnv()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong Integer. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong Boolean. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 == 2", false},
		{"1 == 1", true},
		{"1 != 2", true},
		{"1 != 1", false},
		{"(1 != 1) == false", true},
		{"(1 != 1) == !false", false},
		{"true == true", true},
		{"true == false", false},
		{"true != false", true},
		{"true < false", false},
		{"true > false", true},
		//{"true + true == 2", true},
		//////{"true + false == 1", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. get=%T (%+v)", obj, obj)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`if (1 == 1) {
			if (2 == 2) {
				return 10;
			}
		
			return 0;
		}`, 10},
		{`
let x = fn(y) { fn(z) { return z + y} };
x(5)(5);
`, 10},
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
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`if (10 > 1) {
                     if (10 > 1) { true + false; }

		             return 1;
		          }`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, go=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b; c;", 10},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFuncObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not FUnction. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameter count. parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
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
		{"let identity = fn(x) { x; }; identity(5)", 5},
		{"let identity = fn(x) { return x; }; identity(5)", 5},
		{"let double = fn(x) { x * 2; }; double(5)", 10},
		{"let add = fn(x, y) { x + y; }; add(5, add(1, 2))", 8},
		{"fn(x) { x; }(5);", 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = fn(x) {
  fn(y) { x + y};
};

let addTwo = newAdder(2);
addTwo(2);`
	testIntegerObject(t, testEval(input), 4)
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong String. got=%s, want=%s", result.Value, expected)
		return false
	}

	return true
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"a"`, "a"},
		{`"b"`, "b"},
		{`"cc"`, "cc"},
		{`"hello world"`, "hello world"},
		{`"a" + "b"`, "ab"},
		{`"a" + ""`, "a"},
		{`"" + "a"`, "a"},
		{`"a" * 2`, "aa"},
		{`"ab" * 2`, "abab"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}

	testsBool := []struct {
		input    string
		expected bool
	}{
		{`"ab" * 2 == "abab"`, true},
		{`"ab" * 2 != "abab"`, false},
		{`"ab" * 2 != "aaa"`, true},
		{`"ab" * 2 == "aaa"`, false},
	}

	for _, tt := range testsBool {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` is not supported. got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. get=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != tt.expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
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
		t.Errorf("array has wrong number of elements. got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestHashLiterals(t *testing.T) {
	input := "{1:1}"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("object is not *object.Hash. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Pairs) != 1 {
		t.Errorf("hash has wrong number of elements. got=%d", len(result.Pairs))
	}

	for _, v := range result.Pairs {
		testIntegerObject(t, v.Key, 1)
		testIntegerObject(t, v.Value, 1)
	}
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[0,1,2][0]", 0},
		{"[0,1,2][1]", 1},
		{"[0,1,2][2]", 2},
		{`["a","b","c"][0]`, "a"},
		{`["a","b","c"][1]`, "b"},
		{`["a","b","c"][2]`, "c"},
		{`let myArray = [1, "b", fn() { "1312" }]; myArray[0]`, 1},
		{`let myArray = [1, "b", fn() { "1312" }]; myArray[1]`, "b"},
		{`let myArray = [1, "b", (fn() { "1312" })]; myArray[2]()`, "1312"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		if reflect.TypeOf(tt.expected).Kind() == reflect.Int {
			testIntegerObject(t, evaluated, int64(tt.expected.(int)))
		} else {
			testStringObject(t, evaluated, tt.expected.(string))
		}
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"{1: 1}[1]", 1},
		{"[0,1,2][1]", 1},
		{"[0,1,2][2]", 2},
		{`["a","b","c"][0]`, "a"},
		{`["a","b","c"][1]`, "b"},
		{`["a","b","c"][2]`, "c"},
		{`let myArray = [1, "b", fn() { "1312" }]; myArray[0]`, 1},
		{`let myArray = [1, "b", fn() { "1312" }]; myArray[1]`, "b"},
		{`let myArray = [1, "b", (fn() { "1312" })]; myArray[2]()`, "1312"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		if reflect.TypeOf(tt.expected).Kind() == reflect.Int {
			testIntegerObject(t, evaluated, int64(tt.expected.(int)))
		} else {
			testStringObject(t, evaluated, tt.expected.(string))
		}
	}
}
