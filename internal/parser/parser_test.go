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

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
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
		ast := p.ParseProgram()
		fmt.Printf("%+v %+v", ast, ast.Statements[0])
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
