package evaluate

import (
	"github.com/nishokbanand/interpreter/ast"
	"github.com/nishokbanand/interpreter/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.NULL{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.Boolean:
		return nativeBooltoBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}
	}
	return nil
}

func nativeBooltoBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalStatements(stmts []ast.StatmentNode) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)
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
		return NULL
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
		return NULL
	}
	value := right.(*object.Integer).Value
	return &object.Integer{
		Value: -value,
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evaluateIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBooltoBooleanObject(left == right)
	case operator == "!=":
		return nativeBooltoBooleanObject(left != right)
	default:
		return NULL
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
	return NULL
}
