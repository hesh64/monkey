package parser

import (
	"fmt"
	"monkey/internal/ast"
	"monkey/internal/lexer"
	"monkey/internal/token"
)

type (
	Parser struct {
		l         *lexer.Lexer
		curToken  *token.Token
		peekToken *token.Token
		errors    []string
	}
)

func (p *Parser) nextToken() {
	//if p.peekToken.Type != token.EOF {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
	//}
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: make([]ast.Statement, 0),
	}

	for !p.curTokenIs(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

// todo this isn't a correct implementation
//func (p *Parser) parseExpression() ast.Expression {
//	if p.curToken.Type == token.INT {
//		if token.IsOperator(p.peekToken.Type) && p.peekToken.Type != token.ASSIGN {
//			return p.parseOperatorExpression()
//		} else if p.peekToken.Type == token.SEMICOLON {
//			return p.parseIntegerLiteral()
//		}
//	} else if p.curToken.Type == token.LPAREN {
//		p.parseGroupedExpression()
//	}
//
//	return nil
//}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

//func (p *Parser) parseOperatorExpression() *ast.OperatorExpression {
//	opExp := &ast.OperatorExpression{}
//	opExp.Left = p.parseIntegerLiteral()
//	p.nextToken()
//	opExp.Token = p.curToken
//	p.nextToken()
//	opExp.Right = p.parseExpression()
//
//	return opExp
//}

func (p *Parser) parseGroupedExpression() ast.Expression { return nil }

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = p.parseIdentifier()
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	//
	p.nextToken() // move from = to expression

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
func (p *Parser) parseIfStatement() ast.Statement { return nil }

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	if p.peekToken.Type != t {
		p.errors = append(p.errors, fmt.Sprintf("expected next token to be %s, got %s instead",
			t, p.peekToken.Type))
	}
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	return p
}
