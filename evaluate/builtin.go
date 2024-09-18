package evaluate

import (
	"fmt"

	"github.com/nishokbanand/interpreter/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong Number of args, want 1, got %d", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to 'len' not supported, got %s", arg.Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong Number of args, want 1, got %d", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return arg.Elements[0]
			default:
				return newError("argument to 'len' not supported, got %s", arg.Type())
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong Number of args, want 1, got %d", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return arg.Elements[len(arg.Elements)-1]
			default:
				return newError("argument to 'last' not supported, got %s", arg.Type())
			}
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong Number of args, want 1, got %d", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				newArr := make([]object.Object, length-1, length-1)
				copy(newArr, arg.Elements[1:length])
				return &object.Array{Elements: newArr}
			default:
				return newError("argument to 'rest' not supported, got %s", arg.Type())
			}
		},
	},

	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("Wrong Number of args, want 2, got %d", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				newArr := make([]object.Object, length+1, length+1)
				copy(newArr, arg.Elements[:])
				newArr[length] = args[1]
				return &object.Array{Elements: newArr}
			default:
				return newError("argument to 'rest' not supported, got %s", arg.Type())
			}
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
