package parser

import (
	"fmt"
	"monkey/internal/ast"
	"monkey/internal/lexer"
	"monkey/internal/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x or !x
	CALL        // myFunctions(X)
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(expression ast.Expression) ast.Expression

	Parser struct {
		l              *lexer.Lexer
		errors         []string
		curToken       *token.Token
		peekToken      *token.Token
		prefixParseFns map[token.TokenType]prefixParseFn
		infixParseFns  map[token.TokenType]infixParseFn
	}
)

var (
	precedences = map[token.TokenType]int{
		token.EQ:       EQUALS,
		token.NOT_EQ:   EQUALS,
		token.LT:       LESSGREATER,
		token.GT:       LESSGREATER,
		token.PLUS:     SUM,
		token.MINUS:    SUM,
		token.SLASH:    PRODUCT,
		token.ASTERISK: PRODUCT,
	}
)

// We plan on implementing a pratt parser. This requires associating 2 functions to every type of operator

// helper for registering a prefix parser function
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// helper for registering an infix parser function
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// nextToken moves the value inside of peekToken into curToken
// then reads the next token into peekToken
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// curTokenIs returns true if the curToken type is of that token.TokenType passed
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs return true if the peekToken type is of that token.TokenType passed
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek return an error if the peekToken type is of that token.TokenType passed
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Literal]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Literal]; ok {
		return p
	}

	return LOWEST
}

// peekError appends an error ot the parsers error object.
func (p *Parser) peekError(t token.TokenType) {
	if p.peekToken.Type != t {
		p.errors = append(p.errors, fmt.Sprintf("expected next token to be %s, got %s instead",
			t, p.peekToken.Type))
	}
}

// Errors a helper for extracting all the errors accumulated by the parser during parsing.
func (p *Parser) Errors() []string {
	return p.errors
}

// ParseProgram iterate through the lexer to produce an AST representation of the code
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

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	p.errors = append(p.errors,
		fmt.Sprintf("no prefix parser function for %s found", t))
}

// Statement Parsers

// parseLetStatement parses a let statement
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

	p.nextToken() // move from = to expression

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement parses a return statement
func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Expression Statement is a something like "5 + 1;", "1 + 1"
// It's an expression that is placed on its own without a let or return statement (in the case of monkey)
// but in the general case without a statement.
func (p *Parser) parseExpressionStatements() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseStatement decides which statement parser function to call based on the token type.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatements()
	}
}

// End of Statement Parsers

// parseIdentifier parses an identifier keyword. Identifiers map back to a value.
// ex: "myAge * 2" myAge is an identifier.
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

// parseExpression handles the parsing of any expression. Expression is a statement that produces a value.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	inf := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	inf.Right = p.parseExpression(precedence)

	return inf
}

// parseIntegerLiteral parses an integer literal like "1, 2, 4" This is a helper of the parseExpression method
func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.curToken}

	intValue, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("could not parse %q as integer", p.curToken.Literal))
		return nil
	}

	literal.Value = intValue
	return literal
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	exp.Right = p.parseExpression(PREFIX)
	return exp
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         []string{},
		prefixParseFns: map[token.TokenType]prefixParseFn{},
		infixParseFns:  map[token.TokenType]infixParseFn{},
	}

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()

	return p
}
