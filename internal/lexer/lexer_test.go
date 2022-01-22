package lexer

import (
	"log"
	"monkey/internal/token"
	"testing"
)

const code = `let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
}

let result = add(five, ten);`

type (
	TestCase struct {
		ExpectedType    token.TokenType
		ExpectedLiteral string
	}
)

var (
	testNextTokenMockData = map[string]struct {
		input string
		tests []TestCase
	}{
		"delimiters": {
			input: `=+(){},;`, tests: []TestCase{
				{token.ASSIGN, "="},
				{token.PLUS, "+"},
				{token.LPAREN, "("},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.RBRACE, "}"},
				{token.COMMA, ","},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		"monkey code": {
			input: `let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;
if (5 < 10) {
	return true;
} else {
	return false;
}
10 == 10;
10 != 9;
"foobar"
"foo bar"
[1, 2];
`,
			tests: []TestCase{
				{token.LET, "let"},
				{token.IDENT, "five"},
				{token.ASSIGN, "="},
				{token.INT, "5"},
				{token.SEMICOLON, ";"},
				{token.LET, "let"},
				{token.IDENT, "ten"},
				{token.ASSIGN, "="},
				{token.INT, "10"},
				{token.SEMICOLON, ";"},
				{token.LET, "let"},
				{token.IDENT, "add"},
				{token.ASSIGN, "="},
				{token.FUNCTION, "fn"},
				{token.LPAREN, "("},
				{token.IDENT, "x"},
				{token.COMMA, ","},
				{token.IDENT, "y"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.IDENT, "x"},
				{token.PLUS, "+"},
				{token.IDENT, "y"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.SEMICOLON, ";"},
				{token.LET, "let"},
				{token.IDENT, "result"},
				{token.ASSIGN, "="},
				{token.IDENT, "add"},
				{token.LPAREN, "("},
				{token.IDENT, "five"},
				{token.COMMA, ","},
				{token.IDENT, "ten"},
				{token.RPAREN, ")"},
				{token.SEMICOLON, ";"},
				{token.BANG, "!"},
				{token.MINUS, "-"},
				{token.SLASH, "/"},
				{token.ASTERISK, "*"},
				{token.INT, "5"},
				{token.SEMICOLON, ";"},
				{token.INT, "5"},
				{token.LT, "<"},
				{token.INT, "10"},
				{token.GT, ">"},
				{token.INT, "5"},
				{token.SEMICOLON, ";"},
				{token.IF, "if"},
				{token.LPAREN, "("},
				{token.INT, "5"},
				{token.LT, "<"},
				{token.INT, "10"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.RETURN, "return"},
				{token.TRUE, "true"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.ELSE, "else"},
				{token.LBRACE, "{"},
				{token.RETURN, "return"},
				{token.FALSE, "false"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.INT, "10"},
				{token.EQ, "=="},
				{token.INT, "10"},
				{token.SEMICOLON, ";"},
				{token.INT, "10"},
				{token.NOT_EQ, "!="},
				{token.INT, "9"},
				{token.SEMICOLON, ";"},
				{token.STRING, "foobar"},
				{token.STRING, "foo bar"},
				{token.LBRACKET, "["},
				{token.INT, "1"},
				{token.COMMA, ","},
				{token.INT, "2"},
				{token.RBRACKET, "]"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
	}
)

func TestNextToken(t *testing.T) {
	index := 0
	for title, test := range testNextTokenMockData {
		input := test.input
		tests := test.tests

		t.Run(title, func(t *testing.T) {
			input := input
			tests := tests

			l := New(input)
			ourIndex := index

			for i, tt := range tests {
				tok := l.NextToken()
				if tok == nil {
					log.Fatalf("test[%v:%v] - NextToken returned nil expected=%v, got=%v", ourIndex, i, tt, tok)
				}

				if tok.Type != tt.ExpectedType {
					log.Fatalf("test[%v:%v] - tokenType wrong. expected=%q, got=%q", ourIndex, i, tt.ExpectedType, tok.Type)
				}

				if tok.Literal != tt.ExpectedLiteral {
					log.Fatalf("test[%v:%v] - literal wrong. expected=%q, got=%q", ourIndex, i, tt.ExpectedType, tok.Literal)
				}
			}
		})
		index++
	}
}
