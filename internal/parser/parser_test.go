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

	//	t.Run("multiple error let statements", func(t *testing.T) {
	//		input := `
	//let x 5;
	//let = 10;
	//let 838383;
	//`
	//		l := lexer.New(input)
	//		p := New(l)
	//
	//		program := p.ParseProgram()
	//		checkParserErrors(t, p)
	//
	//		if program == nil {
	//			t.Fatalf("ParseProgram() returned nil")
	//		}
	//
	//		if len(program.Statements) != 3 {
	//			t.Fatalf("program.Statements does not contain 3 statements. got=%v", len(program.Statements))
	//		}
	//
	//		tests := []struct {
	//			expectedIdentifier string
	//		}{
	//			{"x"},
	//			{"y"},
	//			{"foobar"},
	//		}
	//
	//		for i, tt := range tests {
	//			stmt := program.Statements[i]
	//			if !testLetStatement(t, stmt, tt.expectedIdentifier) {
	//				return
	//			}
	//		}
	//	})

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
func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15", "-", 15},
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

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5+5;", 5, "+", 5},
		{"5+5;", 5, "-", 5},
		{"5+5;", 5, "*", 5},
		{"5+5;", 5, "/", 5},
		{"5+5;", 5, ">", 5},
		{"5+5;", 5, "<", 5},
		{"5+5;", 5, "==", 5},
		{"5+5;", 5, "!=", 5},
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

		if !testIntegerLiteral(t, exp.Right, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator is not %s. got=%s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}
