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

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.Boolean:
		return nativeBooltoBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalStatements(node.Statements)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ReturnStatement:
		val := Eval(node.Value)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
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

func evalProgram(stmts []ast.StatmentNode) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)
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
func evalStatements(stmts []ast.StatmentNode) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)
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

func evalIfExpression(node *ast.IfExpression) object.Object {
	condition := Eval(node.Condition)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(node.Consequence)
	} else if node.Alternative != nil {
		return Eval(node.Alternative)
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
