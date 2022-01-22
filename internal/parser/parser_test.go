package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"monkey/internal/ast"
	"monkey/internal/lexer"
	"testing"
)

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%v", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	letIden, ok := letStmt.Name.(*ast.Identifier)
	if !ok {
		t.Fatalf("letIden not *ast.Identifier. got=%T", letIden)
	}

	if letIden.Value != name {
		t.Errorf("letIden.Value not '%s'. got=%s", name, letIden.Value)
		return false
	}

	if letIden.TokenLiteral() != name {
		t.Errorf("letIden.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func testReturnStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "return" {
		t.Errorf("s.TokenLiteral no 'return'. got=%v", s.TokenLiteral())
		return false
	}

	returnStmt, is := s.(*ast.ReturnStatement)
	if !is {
		t.Errorf("s not *ast.ReturnStatement. got=%T", s)
		return false
	}

	if returnStmt.ReturnValue.TokenLiteral() != name {
		t.Errorf("returnStmt.ReturnValue.TokenLiteral() not '%s'. got=%s", name, returnStmt.ReturnValue.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	if len(p.Errors()) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(p.Errors()))
	for _, err := range p.Errors() {
		t.Errorf("parser error: %q", err)
	}
	t.FailNow()
}

func Test_ParseProgram(t *testing.T) {
	t.Run("code = `let t = 5;`", func(t *testing.T) {
		code := `let t = 5;`
		l := lexer.New(code)
		p := New(l)
		p.ParseProgram()
	})
	t.Run("code = `let t = 5 + 5 + 8 + 1;`", func(t *testing.T) {
		code := `let t = 5 + 5 + 8 + 1;`
		l := lexer.New(code)
		p := New(l)
		ast := p.ParseProgram()
		checkParserErrors(t, p)
		assert.True(t, len(ast.Statements) == 1)
	})

	t.Run("Parse Let Statement", func(t *testing.T) {
		input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
		l := lexer.New(input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(program.Statements) != 3 {
			t.Fatalf("program.Statements does not contain 3 statements. got=%v", len(program.Statements))
		}

		tests := []struct {
			expectedIdentifier string
		}{
			{"x"},
			{"y"},
			{"foobar"},
		}

		for i, tt := range tests {
			stmt := program.Statements[i]
			if !testLetStatement(t, stmt, tt.expectedIdentifier) {
				return
			}
		}
	})

	t.Run("Parse Return Statement", func(t *testing.T) {
		code := `
return 123;
return 1;
return false;
`
		l := lexer.New(code)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 3 {
			t.Fatalf("program.Statement does not contain 3 statements. got=%d", len(program.Statements))
		}

		for _, stmt := range program.Statements {
			returnStmt, ok := stmt.(*ast.ReturnStatement)
			if !ok {
				t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			}

			if returnStmt.TokenLiteral() != "return" {
				t.Errorf("returnStmt.TokenLiteral() not 'return'. got=%v", returnStmt.TokenLiteral())
			}
		}
	})
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
		t.Errorf("Ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program statements != 1. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. get=%T", program)
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.Identifier. got=%T", ident)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("type of exp not handled. got=%T", exp)

	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not *ast.InfixExpressiong. get=%T(%s)", exp, exp)
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

func TestIntegerLiteralExpressions(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program statements != 1. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program)
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %s. got=%d", "foobar", literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("ident.Value not %s. got=%d", "foobar", literal.Value)
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	literal, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if literal.Value != value {
		t.Errorf("literal.Value not %s. got=%d", "foobar", literal.Value)
		return false
	}
	if literal.TokenLiteral() != fmt.Sprintf("%v", value) {
		t.Errorf("ident.Value not %s. got=%d", "foobar", literal.Value)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, bl ast.Expression, value bool) bool {
	literal, ok := bl.(*ast.Boolean)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", bl)
		return false
	}
	if literal.Value != value {
		t.Errorf("literal.Value not %s. got=%v", "foobar", literal.Value)
		return false
	}

	if literal.TokenLiteral() != fmt.Sprintf("%v", value) {
		t.Errorf("ident.Value not %s. got=%v", "foobar", literal.Value)
		return false
	}

	return true

}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue interface{}
	}{
		{"!5;", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("program statements != 1. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not *ast.PrefixExpression. get=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator is not %s. got=%s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input          string
		leftValueType  int
		leftValue      interface{}
		operator       string
		rightValueType int
		rightValue     interface{}
	}{
		{"5 + 5;", 0, int64(5), "+", 0, int64(5)},
		{"5 - 5;", 0, int64(5), "-", 0, int64(5)},
		{"5 * 5;", 0, int64(5), "*", 0, int64(5)},
		{"5 / 5;", 0, int64(5), "/", 0, int64(5)},
		{"5 > 5;", 0, int64(5), ">", 0, int64(5)},
		{"5 < 5;", 0, int64(5), "<", 0, int64(5)},
		{"5 == 5;", 0, int64(5), "==", 0, int64(5)},
		{"5 != 5;", 0, int64(5), "!=", 0, int64(5)},
		{`"5" != "5";`, 1, "5", "!=", 1, "5"},
		{`"5" == "5";`, 1, "5", "!=", 1, "5"},
		{`"5" * 5;`, 1, "5", "!=", 0, int64(5)},
		{`"5" + "5";`, 1, "5", "!=", 0, "5"},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("program statements != 1. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt is not *ast.InfixExpression. get=%T", stmt.Expression)
		}

		if tt.leftValueType == 0 {
			if !testInfixExpression(t, exp, tt.leftValue.(int64), tt.operator, tt.rightValue.(int64)) {
				return
			}
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program statements != 1. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program)
	}

	literal, ok := stmt.Expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
	}
	if !literal.Value {
		t.Errorf("literal.Value not %s. got=%v", "true", literal.Value)
	}
	if literal.TokenLiteral() != "true" {
		t.Errorf("ident.Value not %s. got=%v", "true", literal.Value)
	}
}
func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		Input  string
		Output string
	}{
		{"!(true == true)", "(!(true == true))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"5 * 1", "(5 * 1)"},
		{"5 * 1 + 1", "((5 * 1) + 1)"},
		{"5 * 1 + 1 * 4", "((5 * 1) + (1 * 4))"},
		{"true", "true"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == 2 + 1", "((3 < 5) == (2 + 1))"},
		{"a + add(a * c) + d", "((a + add((a * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.Input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		fmt.Println(program.Statements[0].String())

		if tt.Output != program.Statements[0].String() {
			t.Errorf("did not parse as expected %s. got=%s", tt.Output, program.Statements[0].String())
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program statements != 1. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program)
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("exp not *ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}
func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program statements != 1. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program)
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("exp not *ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative == nil {
		t.Fatal("exp.Alternative was nil")
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("alternative is not 1 statements. got=%d", len(exp.Consequence.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !testIdentifier(t, alternative.Expression, "y") {
		t.Fatalf("exp.Alternative.Statements[0] is not *ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program statements != 1. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program)
	}

	lit, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("exp not *ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(lit.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2. got=%d", len(lit.Parameters))
	}

	if !testLiteralExpression(t, lit.Parameters[0], "x") {
		return
	}
	if !testLiteralExpression(t, lit.Parameters[1], "y") {
		return
	}

	if len(lit.Body.Statements) != 1 {
		t.Fatalf("function body statements wrong. want 1. got=%d", len(lit.Body.Statements))
	}

	bodyStmt := lit.Body.Statements[0].(*ast.ExpressionStatement)
	if !testInfixExpression(t, bodyStmt.Expression, "x", "+", "y") {
		t.Fatalf("lit.Body.Statements[0] not *ast.FunctionLiteral. got=%T", lit.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements not length 1. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		literal, ok := stmt.Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf("stmt not *ast.FunctionLiteral. got=%T", stmt)
		}

		if len(literal.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d got %d", len(tt.expectedParams), len(literal.Parameters))
		}

		for i, ident := range literal.Parameters {
			testLiteralExpression(t, ident, tt.expectedParams[i])
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected program.Statements length 1. got=%d", len(program.Statements))
	}

	exp, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected program.Statements[0] type *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	call, ok := exp.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expected exp.Expression type *ast.CallExpression. got=%T", exp.Expression)
	}

	if !testIdentifier(t, call.Function, "add") {
		return
	}

	if len(call.Arguments) != 3 {
		t.Fatalf("expected call.Arguements len 3. got=%d", len(call.Arguments))
	}

	testLiteralExpression(t, call.Arguments[0], 1)
	testInfixExpression(t, call.Arguments[1], 2, "*", 3)
	testInfixExpression(t, call.Arguments[2], 4, "+", 5)
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected program.Statements length 1. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}
func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return y;", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected program.Statements length 1. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]

		val := stmt.(*ast.ReturnStatement).ReturnValue
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"hello world!"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expected program.Statements[0] type *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world!" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestParsingArrayLiteral(t *testing.T) {
	input := `["1", 2, fn(x) {x * x}]`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected program.Statements length 1. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("expected statement to be of type *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	arr, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Errorf("expected expression to be of type *ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(arr.Elements) != 3 {
		t.Errorf("expected array literal to have 3 arguments. got=%d", len(arr.Elements))
	}

	if arr.Elements[0].String() != "1" {
		t.Errorf("expected arr[0] to equal `1`. got=%s", arr.Elements[0].String())
	}
	testIntegerLiteral(t, arr.Elements[1], 2)
}
