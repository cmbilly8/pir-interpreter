package evaluator

import (
	"math/rand"
	"pir-interpreter/object"
	"pir-interpreter/writer"
	"slices"
)

func resolveBuiltin(id string) *object.Builtin {
	builtin := &object.Builtin{}

	switch id {
	case "len":
		builtin.Fn = len_f
	case "peek":
		builtin.Fn = peek
	case "pop":
		builtin.Fn = pop
	case "push":
		builtin.Fn = push
	case "insert":
		builtin.Fn = insert
	case "isMTValue":
		builtin.Fn = isMT
	case "ahoy":
		builtin.Fn = ahoy
	case "empty":
		builtin.Fn = empty
	case "maybe":
		builtin.Fn = maybe
	default:
		return nil
	}

	return builtin
}

/*
func name(args ...object.Object) object.Object {
	if len(args) != 0 {
		return newEvaluationError("wrong number of args. got=%d, expected=0",
			len(args))
	}

	switch arg := args[0].(type) {
	case *object.SupportedArgType:
		return &object.ObjectType{with logic}
	default:
		return newEvaluationError("argument to `name` not supported, got %s",
			args[0].Type())
	}
}
*/
//arrays
func len_f(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newEvaluationError("wrong number of args. got=%d, expected=1",
			len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return nativeIntToIntObj(int64(len(arg.Value)))
	case *object.Array:
		return nativeIntToIntObj(int64(len(arg.Elements)))
	default:
		return newEvaluationError("argument to `len` not supported, got %s",
			args[0].Type())
	}
}

func empty(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newEvaluationError("wrong number of args. got=%d, expected=1",
			len(args))
	}

	switch arg := args[0].(type) {
	case *object.HashMap:
		arg.MP = make(map[object.HashKey]object.KVP)
		return arg
	case *object.Array:
		arg.Elements = make([]object.Object, 0)
		return arg
	default:
		return newEvaluationError("argument to `empty` not supported, got %s",
			args[0].Type())
	}
}

func peek(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newEvaluationError("wrong number of arguments. got=%d, expected=1",
			len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newEvaluationError("argument to `peek` must be ARRAY, got %s",
			args[0].Type())
	}
	arr := args[0].(*object.Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[len(arr.Elements)-1]
	}
	return MT
}

func pop(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newEvaluationError("wrong number of arguments. got=%d, expected=1",
			len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newEvaluationError("argument to `pop` must be ARRAY, got %s",
			args[0].Type())
	}
	arr := args[0].(*object.Array)
	if len(arr.Elements) > 0 {
		last := arr.Elements[len(arr.Elements)-1]
		arr.Elements = arr.Elements[:len(arr.Elements)-1]
		return last
	}
	return MT
}

func push(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newEvaluationError("wrong number of arguments. got=%d, expected=2",
			len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newEvaluationError("first argument to `push` must be ARRAY, got %s",
			args[0].Type())
	}
	arr := args[0].(*object.Array)
	obj := args[1]
	arr.Elements = append(arr.Elements, obj)
	return obj
}

func insert(args ...object.Object) object.Object {
	if len(args) != 3 {
		return newEvaluationError("wrong number of arguments. got=%d, expected=3",
			len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newEvaluationError("first argument to `insert` must be ARRAY, got %s",
			args[0].Type())
	}
	if args[1].Type() != object.INT_OBJ {
		return newEvaluationError("second argument to `insert` must be INT, got %s",
			args[0].Type())
	}
	arr := args[0].(*object.Array)
	i := args[1].(*object.Int).Value
	if i < 0 || i > int64(len(arr.Elements))-1 {
		return newEvaluationError("index out of bound. index=%d, len=%d", i, len(arr.Elements))
	}
	obj := args[2]
	arr.Elements = slices.Insert(arr.Elements, int(i), obj)
	return obj
}

func isMT(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newEvaluationError("wrong number of arguments. got=%d, expected=1",
			len(args))
	}
	return nativeBoolToBoolObj(args[0] == MT)
}

func maybe(args ...object.Object) object.Object {
	if len(args) > 0 {
		return newEvaluationError("wrong number of arguments. got=%d, expected=0",
			len(args))
	}
	return nativeBoolToBoolObj(rand.Intn(2) == 0)
}

/*
func ahoy(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.AsString())
	}
	return MT
}
*/

func ahoy(args ...object.Object) object.Object {
	var str string
	for _, arg := range args {
		str = arg.AsString() + "\n"
		writer.WriteOutput(str)
	}
	return MT
}
