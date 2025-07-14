package evaluator

import (
	"fmt"
	"pir-interpreter/ast"
	"pir-interpreter/object"
	"strconv"
)

var (
	AY    = &object.Bool{Value: true}
	NAY   = &object.Bool{Value: false}
	MT    = &object.MT{}
	BREAK = &object.Break{}
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
	case *ast.PortStatement:
		return evalPortStatementNode()
	case *ast.YarStatement:
		return evalYarStatementNode(node, ns)
	case *ast.ForStatement:
		return evalForStatementNode(node, ns)
	case *ast.FunctionLiteral:
		return evalFuncLiteral(node, ns)
	case *ast.CallExpression:
		return evalFuncCallNode(node, ns)
	case *ast.IndexExpression:
		return evalIndexExpressionNode(node, ns)
	case *ast.IndexAssignment:
		return evalIndexAssignmentNode(node, ns)
	case *ast.ArrayLiteral:
		return evalArrayLiteralNode(node, ns)
	case *ast.StringLiteral:
		return nativeStringToStringObj(node.Value)
	case *ast.HashMapLiteral:
		return evalHashMapLiteralNode(node, ns)
	case *ast.BreakStatement:
		return BREAK
	}
	return MT
}

func evalIndexAssignmentNode(node *ast.IndexAssignment, ns *object.Namespace) object.Object {
	left := Eval(node.Left, ns)
	if object.IsError(left) {
		return left
	}
	index := Eval(node.Index, ns)
	if object.IsError(index) {
		return index
	}
	value := Eval(node.Value, ns)
	if object.IsError(value) {
		return value
	}

	return evalIndexAssignment(left, index, value)

}

func evalIndexAssignment(left, index, val object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INT_OBJ:
		return evalArrayIndexAssignment(left, index, val)
	case left.Type() == object.HASHMAP_OBJ:
		return evalHashMapIndexAssignment(left, index, val)
	default:
		return newEvaluationError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexAssignment(left, index, val object.Object) object.Object {
	arr := left.(*object.Array)
	i := index.(*object.Int).Value
	mx := int64(len(arr.Elements) - 1)
	if i < 0 || i > mx {
		return newEvaluationError("index out of bounds. len=%d, index=%d", len(arr.Elements), i)
	}
	arr.Elements[i] = val
	return MT
}

func evalHashMapIndexAssignment(left, key, val object.Object) object.Object {
	if hashMap, ok := left.(*object.HashMap); ok {
		preHashKey, ok := key.(object.Hashable)
		if !ok {
			return newEvaluationError("Object not hashable: type=%T", key)
		}
		hashKey := preHashKey.Hash()
		hashMap.MP[hashKey] = object.KVP{Key: key, Value: val}
	}
	return MT
}

func evalHashMapLiteralNode(node *ast.HashMapLiteral, ns *object.Namespace) object.Object {
	hm := make(map[object.HashKey]object.KVP)
	for keyNode, valueNode := range node.MP {
		key := Eval(keyNode, ns)
		if object.IsError(key) {
			return key
		}
		preHashKey, ok := key.(object.Hashable)
		if !ok {
			return newEvaluationError("Object not hashable. Type=%s", key.Type())
		}

		value := Eval(valueNode, ns)

		if object.IsError(value) {
			return value
		}
		hashKey := preHashKey.Hash()
		hm[hashKey] = object.KVP{Key: key, Value: value}
	}
	return &object.HashMap{MP: hm}
}

func evalIndexExpressionNode(node *ast.IndexExpression, ns *object.Namespace) object.Object {
	left := Eval(node.Left, ns)
	if object.IsError(left) {
		return left
	}
	index := Eval(node.Index, ns)
	if object.IsError(index) {
		return index
	}
	return evalIndexExpression(left, index)
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INT_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASHMAP_OBJ:
		return evalHashMapIndexExpression(left, index)
	default:
		return newEvaluationError("index operator not supported: %s", left.Type())
	}
}

func evalHashMapIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.HashMap)
	key, ok := index.(object.Hashable)
	if !ok {
		return newEvaluationError("Object not hashable. Type=%s", index.Type())
	}
	kvp, ok := hashObject.MP[key.Hash()]
	if !ok {
		return MT
	}
	return kvp.Value

}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arr := array.(*object.Array)
	i := index.(*object.Int).Value
	mx := int64(len(arr.Elements) - 1)
	if i < 0 || i > mx {
		return newEvaluationError("index out of bounds. len=%d, index=%d", len(arr.Elements), i)
	}
	return arr.Elements[i]
}

func evalArrayLiteralNode(node *ast.ArrayLiteral, ns *object.Namespace) object.Object {
	elements := evalExpressions(node.Elements, ns)
	if len(elements) == 1 && object.IsError(elements[0]) {
		return elements[0]
	}
	return &object.Array{Elements: elements}
}

func nativeStringToStringObj(str string) object.Object {
	return &object.String{Value: str}
}

func evalFuncCallNode(node *ast.CallExpression, ns *object.Namespace) object.Object {
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
	switch f := f.(type) {
	case *object.Function:
		localNS := newFunctionNamespace(f, args)
		result := Eval(f.Body, localNS)
		return extractGivesValue(result)
	case *object.Builtin:
		return f.Fn(args...)
	default:
		return newEvaluationError("Not a function: %s", f.Type())
	}
}

func newFunctionNamespace(f *object.Function, args []object.Object) *object.Namespace {
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
	for _, exp := range exps {
		evaluated := Eval(exp, ns)
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
	if val, ok := ns.Get(node.Value); ok {
		return val
	}
	if builtin := resolveBuiltin(node.Value); builtin != nil {
		return builtin
	}

	return newEvaluationError("Identifier not found: %s", node.Value)
}

func evalYarStatementNode(node *ast.YarStatement, ns *object.Namespace) object.Object {
	val := Eval(node.Value, ns)
	if object.IsError(val) {
		return val
	}
	ns.Set(node.Name.Value, val)
	return &object.MT{}
}

func evalGivesStatementNode(node *ast.GivesStatement, ns *object.Namespace) object.Object {
	value := Eval(node.Value, ns)
	if object.IsError(value) {
		return value
	}
	return &object.GivesValue{Value: value}
}

func evalPortStatementNode() object.Object {
	return newEvaluationError("Port statement not implemented yet.")
}

func evalForStatementNode(node *ast.ForStatement, ns *object.Namespace) object.Object {
	condition := Eval(node.Condition, ns)
	if object.IsError(condition) {
		return condition
	}

	if condition.Type() != object.BOOL_OBJ {
		return newEvaluationError("4 statement condition is not boolen. Got type=%s", condition.Type())
	}

	for condition == AY {
		if node.Body != nil {
			result := Eval(node.Body, ns)
			if object.IsError(result) || result.Type() == object.GIVES_VALUE_OBJ {
				return result
			}

			if result == BREAK {
				return MT
			}
		}

		condition = Eval(node.Condition, ns)
		if object.IsError(condition) {
			return condition
		}
	}
	return MT
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
			if rt == object.GIVES_VALUE_OBJ || rt == object.ERROR_OBJ || rt == object.BREAK_OBJ {
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
	case left.Type() == object.STRING_OBJ && right.Type() == object.INT_OBJ:
		return evalStringInfixExpression(left, node.Operator, castIntToString(right))
	case left.Type() == object.INT_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(castIntToString(left), node.Operator, right)
	case left.Type() != right.Type():
		return newEvaluationError("type mismatch: %s %s %s",
			left.Type(), node.Operator, right.Type())
	default:
		return newEvaluationError("No infix expression for: %s %s %s",
			left.Type(), node.Operator, right.Type())
	}

}

func castIntToString(obj object.Object) object.Object {
	intObj, ok := obj.(*object.Int)
	if ok {
		return &object.String{Value: strconv.FormatInt(intObj.Value, 10)}
	}
	return newEvaluationError("Error while casting INT to STRING. Expected INT, got %s", obj.Type())

}

func evalBoolInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	// Don't deref value since the object would have to be equal anyway
	switch operator {
	case "=":
		return nativeBoolToBoolObj(left == right)
	case "<>":
		return nativeBoolToBoolObj(left != right)
	case "and":
		return nativeBoolToBoolObj(left == AY && right == AY)
	case "or":
		return nativeBoolToBoolObj(left == AY || right == AY)
	default:
		return newEvaluationError("unknown operator: %s %s %s",
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
	case "<>":
		return nativeBoolToBoolObj(leftVal != rightVal)
	default:
		return newEvaluationError("unknown operator: %s %s %s",
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
	case "mod":
		return nativeIntToIntObj(leftVal % rightVal)
	case "/":
		return nativeIntToIntObj(leftVal / rightVal)
	case "=":
		return nativeBoolToBoolObj(leftVal == rightVal)
	case "<>":
		return nativeBoolToBoolObj(leftVal != rightVal)
	case ">":
		return nativeBoolToBoolObj(leftVal > rightVal)
	case "<":
		return nativeBoolToBoolObj(leftVal < rightVal)
	case "<=":
		return nativeBoolToBoolObj(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBoolObj(leftVal >= rightVal)
	default:
		return newEvaluationError("unknown operator: %s %s %s",
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
		return newEvaluationError("unknown operator: %s%s", node.Operator, operand.Type())
	}
}

func evalMathmaticalNegateExpression(operand object.Object) object.Object {
	if operand.Type() != object.INT_OBJ {
		return newEvaluationError("unknown operator: -%s", operand.Type())
	}

	new_val := operand.(*object.Int).Value * -1
	return &object.Int{Value: new_val}
}

func evalLogicalNegateExpression(operand object.Object) object.Object {
	if operand.Type() != object.BOOL_OBJ {
		return newEvaluationError("Unsupported operation !%s", operand.Type())
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

func newEvaluationError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
