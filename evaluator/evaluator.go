package evaluator

import (
	"fmt"
	"pir-interpreter/ast"
	"pir-interpreter/object"
)

var (
	AY    = &object.Bool{Value: true}
	NAY   = &object.Bool{Value: false}
	MT    = &object.MT{}
	MAYBE = &object.Maybe{}
)

func Eval(node ast.Node, ns *object.Namespace) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgramNode(node, ns)
	case *ast.BlockStatement:
		return evalBlockStatement(node, ns)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, ns)
	case *ast.IfStatement:
		return evalIfStatementNode(node, ns)
	case *ast.IntegerLiteral:
		return nativeIntToIntObj(node.Value)
	case *ast.Boolean:
		return nativeBoolToBoolObj(node.Value)
	case *ast.Identifier:
		return evalIdentifier(node, ns)
	case *ast.PrefixExpression:
		return evalPrefixExpressionNode(node, ns)
	case *ast.InfixExpression:
		return evalInfixExpressionNode(node, ns)
	case *ast.GivesStatement:
		return evalGivesStatementNode(node, ns)
	case *ast.YarStatement:
		return evalYarStatementNode(node, ns)
	case *ast.FunctionLiteral:
		return evalFuncLiteral(node, ns)
	case *ast.CallExpression:
		return evalFuncCall(node, ns)
	case *ast.StringLiteral:
		return nativeStringToStringObj(node.Value)
	}
	return MT

}

func nativeStringToStringObj(str string) object.Object {
	return &object.String{Value: str}
}

func evalFuncCall(node *ast.CallExpression, ns *object.Namespace) object.Object {
	f := Eval(node.Function, ns)
	if object.IsError(f) {
		return f
	}
	args := evalExpressions(node.Arguments, ns)
	if len(args) == 1 && object.IsError(args[0]) {
		return args[0]
	}
	return callFunc(f, args)
}

func callFunc(f object.Object, args []object.Object) object.Object {
	function, ok := f.(*object.Function)
	if !ok {
		return newError("Not a function: %s", f.Type())
	}
	localNS := newLocalNamespace(function, args)
	result := Eval(function.Body, localNS)
	return extractGivesValue(result)
}

func newLocalNamespace(f *object.Function, args []object.Object) *object.Namespace {
	localNS := object.NewNestedNamespace(f.NS)
	for i, param := range f.Params {
		localNS.Set(param.Value, args[i])
	}
	return localNS
}

func extractGivesValue(obj object.Object) object.Object {
	if givesValue, ok := obj.(*object.GivesValue); ok {
		return givesValue.Value
	}
	return obj

}

func evalExpressions(exps []ast.Expression, ns *object.Namespace) []object.Object {
	var objs []object.Object
	// Can't wait until I forget about this
	for i := len(exps) - 1; i >= 0; i-- {
		evaluated := Eval(exps[i], ns)
		if object.IsError(evaluated) {
			return []object.Object{evaluated}
		}
		objs = append(objs, evaluated)
	}
	return objs
}

func evalFuncLiteral(node *ast.FunctionLiteral, ns *object.Namespace) object.Object {
	params := node.Params
	body := node.Body
	return &object.Function{Params: params, NS: ns, Body: body}
}

func evalIdentifier(node *ast.Identifier, ns *object.Namespace) object.Object {
	val, ok := ns.Get(node.Value)
	if !ok {
		return newError("Identifier not found: %s", node.Value)
	}
	return val
}

func evalYarStatementNode(node *ast.YarStatement, ns *object.Namespace) object.Object {
	val := Eval(node.Value, ns)
	if object.IsError(val) {
		return val
	}
	ns.Set(node.Name.Value, val)
	return &object.MT{}
}

func evalStatements(stmts []ast.Statement, ns *object.Namespace) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement, ns)
		if returnValue, ok := result.(*object.GivesValue); ok {
			return returnValue.Value
		}
	}
	return result
}

func evalGivesStatementNode(node *ast.GivesStatement, ns *object.Namespace) object.Object {
	value := Eval(node.Value, ns)
	if object.IsError(value) {
		return value
	}
	return &object.GivesValue{Value: value}
}

func evalIfStatementNode(node *ast.IfStatement, ns *object.Namespace) object.Object {
	for _, conditional := range node.Conditionals {
		if cond := Eval(conditional.Condition, ns); cond == AY {
			if object.IsError(cond) {
				return cond
			}

			return Eval(conditional.Consequence, ns)
		}
	}
	if node.Alternate != nil {
		return Eval(node.Alternate, ns)
	}
	return MT
}

func evalBlockStatement(bs *ast.BlockStatement, ns *object.Namespace) object.Object {
	var result object.Object
	for _, statement := range bs.Statements {
		result = Eval(statement, ns)
		if result != nil {
			rt := result.Type()
			if rt == object.GIVES_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func evalInfixExpressionNode(node *ast.InfixExpression, ns *object.Namespace) object.Object {
	left := Eval(node.Left, ns)
	if object.IsError(left) {
		return left
	}

	right := Eval(node.Right, ns)
	if object.IsError(right) {
		return right
	}

	switch {
	case left.Type() == object.INT_OBJ && right.Type() == object.INT_OBJ:
		return evalIntInfixExpression(left, node.Operator, right)
	case left.Type() == object.BOOL_OBJ && right.Type() == object.BOOL_OBJ:
		return evalBoolInfixExpression(left, node.Operator, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(left, node.Operator, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), node.Operator, right.Type())
	default:
		return newError("No infix expression for: %s %s %s",
			left.Type(), node.Operator, right.Type())
	}

}

func evalBoolInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	// Don't deref value since the object would have to be equal anyway
	switch operator {
	case "=":
		return nativeBoolToBoolObj(left == right)
	case "!=":
		return nativeBoolToBoolObj(left != right)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return nativeStringToStringObj(leftVal + rightVal)
	case "=":
		return nativeBoolToBoolObj(leftVal == rightVal)
	case "!=":
		return nativeBoolToBoolObj(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIntInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftVal := left.(*object.Int).Value
	rightVal := right.(*object.Int).Value

	switch operator {
	case "+":
		return nativeIntToIntObj(leftVal + rightVal)
	case "-":
		return nativeIntToIntObj(leftVal - rightVal)
	case "*":
		return nativeIntToIntObj(leftVal * rightVal)
	case "/":
		return nativeIntToIntObj(leftVal / rightVal)
	case "=":
		return nativeBoolToBoolObj(leftVal == rightVal)
	case "!=":
		return nativeBoolToBoolObj(leftVal != rightVal)
	case ">":
		return nativeBoolToBoolObj(leftVal > rightVal)
	case "<":
		return nativeBoolToBoolObj(leftVal < rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalPrefixExpressionNode(node *ast.PrefixExpression, ns *object.Namespace) object.Object {
	operand := Eval(node.Right, ns)
	if object.IsError(operand) {
		return operand
	}

	switch node.Operator {
	case "!":
		return evalLogicalNegateExpression(operand)
	case "-":
		return evalMathmaticalNegateExpression(operand)
	default:
		return newError("unknown operator: %s%s", node.Operator, operand.Type())
	}
}

func evalMathmaticalNegateExpression(operand object.Object) object.Object {
	if operand.Type() != object.INT_OBJ {
		return newError("unknown operator: -%s", operand.Type())
	}

	new_val := operand.(*object.Int).Value * -1
	return &object.Int{Value: new_val}
}

func evalLogicalNegateExpression(operand object.Object) object.Object {
	if operand.Type() != object.BOOL_OBJ {
		return newError("Unsupported operation !%s", operand.Type())
	}

	if operand == AY {
		return NAY
	}
	return AY
}

func nativeIntToIntObj(i int64) object.Object {
	return &object.Int{Value: i}
}

func nativeBoolToBoolObj(b bool) object.Object {
	if b {
		return AY
	}
	return NAY
}

func evalProgramNode(program *ast.Program, ns *object.Namespace) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, ns)
		switch result := result.(type) {
		case *object.GivesValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
