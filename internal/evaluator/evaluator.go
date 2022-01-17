package evaluator

import (
	"fmt"
	"monkey/internal/ast"
	"monkey/internal/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		o := &object.Integer{Value: node.Value}
		return o
	case *ast.Boolean:
		if node.Value {
			return TRUE
		}
		return FALSE
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
		return evalBlockStatement(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}
	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("Unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default: // we handled the only two possible falsy values (False & Null)
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	// Why allocate a new object if you can just update?
	right.(*object.Integer).Value = right.(*object.Integer).Value * -1
	return right
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	switch operator {
	case "+":
		return &object.Integer{Value: left.(*object.Integer).Value + right.(*object.Integer).Value}
	case "-":
		return &object.Integer{Value: left.(*object.Integer).Value - right.(*object.Integer).Value}
	case "*":
		return &object.Integer{Value: left.(*object.Integer).Value * right.(*object.Integer).Value}
	case "/":
		// todo handle error?
		return &object.Integer{Value: left.(*object.Integer).Value / right.(*object.Integer).Value}
	case "==":
		return nativeBoolToBooleanObject(left.(*object.Integer).Value == right.(*object.Integer).Value)
	case "!=":
		return nativeBoolToBooleanObject(left.(*object.Integer).Value != right.(*object.Integer).Value)
	case "<":
		return nativeBoolToBooleanObject(left.(*object.Integer).Value < right.(*object.Integer).Value)
	case ">":
		return nativeBoolToBooleanObject(left.(*object.Integer).Value > right.(*object.Integer).Value)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func nativeBoolToBooleanObject(b bool) object.Object {
	if b {
		return TRUE
	}

	return FALSE
}

func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
	switch operator {
	case "==":
		return nativeBoolToBooleanObject(left.(*object.Boolean).Value == right.(*object.Boolean).Value)
	case "!=":
		return nativeBoolToBooleanObject(left.(*object.Boolean).Value != right.(*object.Boolean).Value)
	case "<":
		leftVal := 0
		if left.(*object.Boolean).Value {
			leftVal = 1
		}
		rightVal := 0
		if right.(*object.Boolean).Value {
			rightVal = 1
		}
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		leftVal := 0
		if left.(*object.Boolean).Value {
			leftVal = 1
		}
		rightVal := 0
		if right.(*object.Boolean).Value {
			rightVal = 1
		}
		return nativeBoolToBooleanObject(leftVal > rightVal)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBoolToInt(boolean object.Object) object.Object {
	if boolean.(*object.Boolean).Value {
		return &object.Integer{Value: 1}
	}

	return &object.Integer{Value: 0}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	if left.Type() != right.Type() {
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return evalIntegerInfixExpression(operator, left, right)
	}

	if left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ {
		return evalBooleanInfixExpression(operator, left, right)
	}

	//if left.Type() == object.BOOLEAN_OBJ && right.Type() == object.INTEGER_OBJ {
	//	return evalIntegerInfixExpression(operator, evalBoolToInt(left), right)
	//}

	//if left.Type() == object.INTEGER_OBJ && right.Type() == object.BOOLEAN_OBJ {
	//	return evalIntegerInfixExpression(operator, left, evalBoolToInt(right))
	//}

	return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		e := Eval(ie.Consequence)
		return e
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil && (result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ) {
			return result
		}
	}

	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}
