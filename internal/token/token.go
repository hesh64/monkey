package token

// todo for now the token types we want to define are for the code sample
//  down here
/*`
let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
}

let result = add(five, ten);
`*/

type (
	TokenType = string
	Token     struct {
		Type    TokenType
		Literal string
	}
)

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers
	IDENT  = "IDENT" // token type for all the user defined identifiers
	INT    = "INT"   // integer data type
	STRING = "STRING"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	PERIOD    = "."
	COMMA     = ","
	COLON     = ":"
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var (
	keywords = map[string]TokenType{
		"let":    LET,
		"fn":     FUNCTION,
		"true":   TRUE,
		"false":  FALSE,
		"if":     IF,
		"else":   ELSE,
		"return": RETURN,
	}
)

// LookupIdent checks the keywords table to see whether the given identifier is in fact a keyword. If it is, it returns
// the keyword's TokenType constant. If it isn't, we just get back token.IDENT, which is the TokenType for all
// the user defined identifiers.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}

var (
	operators = map[string]bool{
		ASSIGN:   true,
		PLUS:     true,
		MINUS:    true,
		BANG:     true,
		ASTERISK: true,
		SLASH:    true,
		LT:       true,
		GT:       true,
		EQ:       true,
		NOT_EQ:   true,
	}
)

func IsOperator(op string) bool {
	_, has := operators[op]
	return has
}
