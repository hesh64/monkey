package ast

import (
	"bytes"
	"monkey/internal/token"
)

type (
	Node interface {
		TokenLiteral() string
		String() string
	}

	Statement interface {
		Node
		statementNode()
	}

	Expression interface {
		Node
		expressionNode()
	}

	// Program is the root node of any program we will every parse. Statements are what a program composed of.
	Program struct {
		Statements []Statement
	}
)

func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	}

	return p.Statements[0].TokenLiteral()
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type (
	// Statement implementers

	// LetStatement is a let declaration ast node
	LetStatement struct {
		Token *token.Token // the token to which this statement points to
		Name  Expression   // name of the variable
		Value Expression
	}

	// ReturnStatement is a return statement ast node
	ReturnStatement struct {
		Token       *token.Token // the token to which this statement points to
		ReturnValue Expression
	}
	// ExpressionStatement is any type of expression
	// ex:
	// foobar;
	// foobar + 1;
	// 1 + 1;
	// 1 + add(1, foobar);
	ExpressionStatement struct {
		Token      *token.Token // the first token of the expression
		Expression Expression
	}

	// Expression implementer

	// Identifier holds the value of a user or language defined identifier name
	Identifier struct {
		Token *token.Token
		Value string
	}

	IntegerLiteral struct {
		Token *token.Token
		Value int64
	}

	// PrefixExpression like "-1", "!found"
	PrefixExpression struct {
		Token    *token.Token
		Operator string
		Right    Expression // this will point to the 1 or found
	}

	InfixExpression struct {
		Token    *token.Token
		Operator string
		Left     Expression
		Right    Expression
	}
)

func (l *LetStatement) statementNode()       {}
func (l *LetStatement) TokenLiteral() string { return l.Token.Literal }
func (l *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String())
	out.WriteString(" = ")

	if l.Value != nil {
		out.WriteString(l.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

func (r *ReturnStatement) statementNode()       {}
func (r *ReturnStatement) TokenLiteral() string { return r.Token.Literal }
func (r *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(r.TokenLiteral() + " ")

	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

func (r *ExpressionStatement) statementNode()       {}
func (r *ExpressionStatement) TokenLiteral() string { return r.Token.Literal }
func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}

	return ""
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return i.Token.Literal }

func (i *PrefixExpression) expressionNode()      {}
func (i *PrefixExpression) TokenLiteral() string { return i.Token.Literal }
func (i *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Operator)
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}
func (i *InfixExpression) expressionNode()      {}
func (i *InfixExpression) TokenLiteral() string { return i.Token.Literal }
func (i *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}
