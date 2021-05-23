package evaluator

import (
	"CornyLang/object"
	"fmt"
)

var builtins = map[string]*object.BuiltinObj{
	"size": {
		Func: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.NewError(fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)))
			}
			switch arg := args[0].(type) {
			case *object.StringObj:
				return &object.NumberObj{Value: float64(len(arg.Value) - 1)}
			case *object.ArrayObj:
				return &object.NumberObj{Value: float64(len(arg.Elements) - 1)}
			case *object.HashObj:
				return &object.NumberObj{Value: float64(len(arg.Elements) - 1)}
			default:
				return object.NewError(fmt.Sprintf("wrong argument to `size` not supported, got=%s", args[0].Type()))
			}
		},
		ExtraArgs: false,
	},
	"type": {
		Func: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.NewError(fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)))
			}
			switch args[0].Type() {
			case object.STRING_OBJ:
				return &object.StringObj{Value: "C"}
			case object.NUMBER_OBJ:
				return &object.StringObj{Value: "N"}
			case object.NULL_OBJ:
				return &object.StringObj{Value: "X"}
			case object.FUNCTION_OBJ:
				return &object.StringObj{Value: "F"}
			case object.ARRAY_OBJ:
				return &object.StringObj{Value: "A"}
			case object.HASH_OBJ:
				return &object.StringObj{Value: "H"}
			default:
				return &object.StringObj{Value: "U"}
			}
		},
		ExtraArgs: false,
	},
	"push": {
		Func: func(args ...object.Object) object.Object {
			if arrayObj, ok := args[0].(*object.ArrayObj); ok {
				// create a new array
				var newArrayObj = &object.ArrayObj{
					Elements: []object.Object{},
				}
				// copy old elements if exists
				if len(arrayObj.Elements) > 0 {
					newArrayObj.Elements = append(newArrayObj.Elements, arrayObj.Elements...)
				}

				// append new elements if there are
				if len(args) > 1 {
					for i := 1; i < len(args); i++ {
						newArrayObj.Elements = append(newArrayObj.Elements, args[i])
					}
					return newArrayObj
				}
				newArrayObj.Elements = append(newArrayObj.Elements, &object.NullObj{})

				return newArrayObj
			}
			return object.NewError(fmt.Sprintf("wrong argument to `push` not supported, got=%s", args[0].Type()))
		},
		ExtraArgs: true,
	},
}
