package evaluator

import (
	"dito/src/object"
	"fmt"
	"math"
)

// Builtins : map of builtin functions
var Builtins = map[string]*object.Builtin{
	"len":   &object.Builtin{Fn: ditoLen},
	"type":  &object.Builtin{Fn: ditoType},
	"print": &object.Builtin{Fn: ditoPrint},
	"sqrt":  &object.Builtin{Fn: ditoSqrt},
	"iota":  &object.Builtin{Fn: ditoIota},
	"int":   &object.Builtin{Fn: ditoInt},
	"log":   &object.Builtin{Fn: ditoLog},
	"log2":  &object.Builtin{Fn: ditoLog2},
	"log10": &object.Builtin{Fn: ditoLog10},
	"cos":   &object.Builtin{Fn: ditoCos},
	"sin":   &object.Builtin{Fn: ditoSin},
	"Tan":   &object.Builtin{Fn: ditoTan},
	// "abs": &object.Builtin{Fn: validDitoAbs},
}

// func EvalBuiltinFn(fn *object.Builtin, args ...Object) Object {
// 	if argLen := len(args); argLen > fn.ArgsMax || argLen < fn.ArgsMin {
// 		return newError("Wrong   number of arguments. got=%d, want=1", len(args))
// 	}

// 	return fn.Fn(args...)
// }

func ditoTan(args ...object.Object) object.Object {
	switch arg := args[0].(type) {
	case *object.Integer:
		return object.NewDitoFloat(math.Tan(float64(arg.Value)))
	case *object.Float:
		return object.NewDitoFloat(math.Tan(arg.Value))
	default:
		return newError("Argument to `Tan` not supported, got %s", args[0].Type())
	}
}

func ditoSin(args ...object.Object) object.Object {
	switch arg := args[0].(type) {
	case *object.Integer:
		return object.NewDitoFloat(math.Sin(float64(arg.Value)))
	case *object.Float:
		return object.NewDitoFloat(math.Sin(arg.Value))
	default:
		return newError("Argument to `Sin` not supported, got %s", args[0].Type())
	}
}

func ditoCos(args ...object.Object) object.Object {
	switch arg := args[0].(type) {
	case *object.Integer:
		return object.NewDitoFloat(math.Cos(float64(arg.Value)))
	case *object.Float:
		return object.NewDitoFloat(math.Cos(arg.Value))
	default:
		return newError("Argument to `Cos` not supported, got %s", args[0].Type())
	}
}

func ditoLog(args ...object.Object) object.Object {
	switch arg := args[0].(type) {
	case *object.Integer:
		return object.NewDitoFloat(math.Log(float64(arg.Value)))
	case *object.Float:
		return object.NewDitoFloat(math.Log(arg.Value))
	default:
		return newError("Argument to `log` not supported, got %s", args[0].Type())
	}
}

func ditoLog10(args ...object.Object) object.Object {
	switch arg := args[0].(type) {
	case *object.Integer:
		return object.NewDitoFloat(math.Log10(float64(arg.Value)))
	case *object.Float:
		return object.NewDitoFloat(math.Log10(arg.Value))
	default:
		return newError("Argument to `log10` not supported, got %s", args[0].Type())
	}
}

func ditoLog2(args ...object.Object) object.Object {
	switch arg := args[0].(type) {
	case *object.Integer:
		return object.NewDitoFloat(math.Log2(float64(arg.Value)))
	case *object.Float:
		return object.NewDitoFloat(math.Log2(arg.Value))
	default:
		return newError("Argument to `log2` not supported, got %s", args[0].Type())
	}
}

func ditoInt(args ...object.Object) object.Object {
	switch arg := args[0].(type) {
	case *object.Integer:
		return arg
	case *object.Float:
		return object.NewDitoInteger(int64(arg.Value))
	default:
		return newError("Argument to `Int` not supported, got %s", args[0].Type())
	}
}

func ditoIota(args ...object.Object) object.Object {
	var arg1, arg2, arg3 int64
	switch len(args) {
	case 1:
		if args[0].Type() != object.IntergerObj {
			return newError("Argument to `iota` not supported, got=%s. want=Int", args[0].Type())
		}
		arg1 = 0
		arg2 = args[0].(*object.Integer).Value
		arg3 = 1
	case 2:
		if args[0].Type() != object.IntergerObj {
			return newError("Argument to `iota` not supported, got=%s. want=Int", args[0].Type())
		}
		if args[1].Type() != object.IntergerObj {
			return newError("Argument to `iota` not supported, got=%s. want=Int", args[1].Type())
		}
		arg1 = args[0].(*object.Integer).Value
		arg2 = args[1].(*object.Integer).Value
		arg3 = 1
	case 3:
		if args[0].Type() != object.IntergerObj {
			return newError("Argument to `iota` not supported, got=%s. want=Int", args[0].Type())
		}
		if args[1].Type() != object.IntergerObj {
			return newError("Argument to `iota` not supported, got=%s. want=Int", args[1].Type())
		}
		if args[2].Type() != object.IntergerObj {
			return newError("Argument to `iota` not supported, got=%s. want=Int", args[2].Type())
		}
		arg1 = args[0].(*object.Integer).Value
		arg2 = args[1].(*object.Integer).Value
		arg3 = args[2].(*object.Integer).Value
	default:
		return newError("`iota`: wrong number of args. want=(1-3) got=%d", len(args))
	}
	var result []object.Object
	for i := arg1; i < arg2; i += arg3 {
		result = append(result, object.NewDitoInteger(i))
	}
	return &object.Array{Len: (arg2 - arg1) / arg3, Elements: result}
}

func ditoType(args ...object.Object) object.Object {
	return &object.DitoString{Value: args[0].Type()}
}

func ditoPrint(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Printf("%s ", arg.Inspect())
	}
	fmt.Println()
	return nil
}

func ditoSqrt(args ...object.Object) object.Object {
	switch arg := args[0].(type) {
	case *object.Integer:
		return object.NewDitoFloat(math.Sqrt(float64(arg.Value)))
	case *object.Float:
		return object.NewDitoFloat(math.Sqrt(arg.Value))
	default:
		return newError("Argument to `sqrt` not supported, got %s", args[0].Type())
	}
}

func ditoLen(args ...object.Object) object.Object {
	switch arg := args[0].(type) {
	case *object.DitoString:
		return object.NewDitoInteger(int64(len(arg.Value)))
	case *object.Array:
		return object.NewDitoInteger(arg.Len)
	default:
		return newError("Argument to `Len` not supported, got %s", args[0].Type())
	}
}