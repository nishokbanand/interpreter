package evaluate

import (
	"fmt"

	"github.com/nishokbanand/interpreter/ast"
	"github.com/nishokbanand/interpreter/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.NULL{}
)

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.Boolean:
		return nativeBooltoBooleanObject(node.Value)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalStatements(node.Statements, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.ReturnStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.String:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		ele := evalExpressions(node.Elements, env)
		if len(ele) == 1 && isError(ele[0]) {
			return ele[0]
		}
		return &object.Array{Elements: ele}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.HashLiteral:
		return evalHash(node, env)
	}
	return nil
}

func nativeBooltoBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalProgram(stmts []ast.StatmentNode, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt, env)
		if result != nil {
			switch result := result.(type) {
			case *object.ReturnValue:
				return result.Value
			case *object.Error:
				return result
			}
		}
	}
	return result
}
func evalStatements(stmts []ast.StatmentNode, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt, env)
		if result != nil {
			rt := result.Type()
			if rt == object.ERROR_OBJ || rt == object.RETURN_OBJ {
				return result
			}
		}
	}
	return result
}

func evalPrefixExpression(Operator string, right object.Object) object.Object {
	switch Operator {
	case "!":
		return evaluateBangExpression(right)
	case "-":
		return evaluateMinusExpression(right)
	default:
		return newError("Unknown operator %s %s", Operator, right.Type())
	}
}

func evaluateBangExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return NULL
	}
}

func evaluateMinusExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("Unknown operator -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{
		Value: -value,
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newError("Operands are not of the same type : %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evaluateIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evaluateStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBooltoBooleanObject(left == right)
	case operator == "!=":
		return nativeBooltoBooleanObject(left != right)
	default:
		return newError("Unknown Operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evaluateIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	//int producing operators
	case "+":
		return &object.Integer{Value: (leftVal + rightVal)}
	case "-":
		return &object.Integer{Value: (leftVal - rightVal)}
	case "*":
		return &object.Integer{Value: (leftVal * rightVal)}
	case "/":
		return &object.Integer{Value: (leftVal / rightVal)}
		//boolean producing operators
	case "<":
		return nativeBooltoBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBooltoBooleanObject(leftVal > rightVal)
	case "!=":
		return nativeBooltoBooleanObject(leftVal != rightVal)
	case "==":
		return nativeBooltoBooleanObject(leftVal == rightVal)
	}
	return newError("Unknown Operator %s %s %s", left.Type(), operator, right.Type())
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(condition object.Object) bool {
	switch condition {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		return false
	default:
		return true
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found %s", node.Value)
}

func evalExpressions(exps []ast.ExpressionNode, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendedFuncEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrappedValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function %s :", fn.Type())
	}
}

func extendedFuncEnv(f *object.Function, args []object.Object) *object.Environment {
	extendedEnv := object.NewEnclosedEnvironment(f.Env)
	for idx, param := range f.Parameters {
		extendedEnv.Set(param.Value, args[idx])
	}
	return extendedEnv
}

func unwrappedValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evaluateStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	switch operator {
	//int producing operators
	case "+":
		return &object.String{Value: (leftVal + rightVal)}
	default:
		return newError("Unknown Operator %s %s %s", left, operator, right)
	}
}
func evalIndexExpression(left object.Object, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index Operator not supported %s ", left)
	}
}

func evalArrayIndexExpression(arr object.Object, index object.Object) object.Object {
	arrayObject := arr.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 || idx > max {
		return newError("Out of bound Error :%d greater than %d", idx, max)
	}
	return arrayObject.Elements[idx]
}

func evalHash(hash *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for keyNode, valueNode := range hash.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("Cannot use as HashKey %s", key.Type())
		}
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}
		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{
			Key:   key,
			Value: value,
		}
	}
	return &object.Hash{
		Pairs: pairs,
	}
}

func evalHashIndexExpression(hash object.Object, index object.Object) object.Object {
	hashObject, ok := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("Index unusable as hashKey %s", hash.Type())
	}
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}
