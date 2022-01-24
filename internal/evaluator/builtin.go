package evaluator

import (
	"fmt"
	"monkey/internal/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` is not supported. got %s", args[0].Type())
			}
		},
	},
	"printf": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 {
				return newError("wrong number of arguments. got=%d", len(args))
			}

			argsInterface := make([]interface{}, 0, len(args))
			for i, arg := range args {
				if i > 0 {
					argsInterface = append(argsInterface, arg.Inspect())
				}
			}

			fmt.Printf(args[0].Inspect(), argsInterface...)
			return NULL
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 {
				return newError("wrong number of arguments. got=%d", len(args))
			}

			argsInterface := make([]interface{}, 0, len(args))
			for _, arg := range args {
				argsInterface = append(argsInterface, arg.Inspect())
			}
			fmt.Println(argsInterface...)

			return NULL
		},
	},
}
