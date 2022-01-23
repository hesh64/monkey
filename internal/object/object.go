package object

import (
	"bytes"
	"fmt"
	"monkey/internal/ast"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	STRING_OBJ       = "STRING"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

type (
	Object interface {
		Type() ObjectType
		Inspect() string
	}

	Hashable interface {
		Object
		HashKey() HashKey
	}

	Integer struct {
		Value int64
	}

	String struct {
		Value string
	}

	Boolean struct {
		Value bool
	}
)

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) HashKey() HashKey {
	return HashKey(fmt.Sprintf("%s_%d", i.Type(), i.Value))
}

func (s *String) Inspect() string {
	return s.Value
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (s *String) HashKey() HashKey {
	return HashKey(fmt.Sprintf("%s_%s", s.Type(), s.Value))
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

type Null struct{}

func (*Null) Inspect() string {
	return "null"
}

func (*Null) Type() ObjectType {
	return NULL_OBJ
}

type (
	ReturnValue struct {
		Value Object
	}
)

func (r *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type BuiltinFunction func(arg ...Object) Object
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elts := make([]string, 0, len(a.Elements))
	for _, obj := range a.Elements {
		elts = append(elts, obj.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elts, ", "))
	out.WriteString("]")

	return out.String()
}

type HashKey string

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	elts := make([]string, 0, len(h.Pairs))
	for _, v := range h.Pairs {
		elts = append(elts, fmt.Sprintf("%s: %s", v.Key.Inspect(), v.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(elts, ", "))
	out.WriteString("}")

	return out.String()
}
