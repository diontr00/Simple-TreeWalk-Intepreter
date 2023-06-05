package evaluator

import (
	"bytes"
	"khanhanh_lang/ast"
	"khanhanh_lang/object"
	"strconv"
)

// Since we dont want to create distinct , true/false object evertytime (there is only two possible value )
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NIL   = &object.Nil{}
)

func nativeBool(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBool(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	}

	return nil
}

// traver the AST , and parse each statement
func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement)
	}
	return result
}

// --------------------------------------------------------------------
// PREFIX EXPRESSION
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangPrefix(right)
	case "-":
		return evalMinusPrefix(right)
	default:
		return NIL
	}
}

func evalMinusPrefix(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NIL
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalBangPrefix(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NIL:
		return TRUE
	default:
		return FALSE
	}
}

// ---------------------------------------------------------------------
// INFIX EXPRESSION
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntOperator(operator, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBoolOperator(operator, left, right)
	case left.Type() == object.STRING_OBJ:
		return evalStringOperator(operator, left, right)

	default:
		return NIL
	}
}

// boolean operator
func evalBoolOperator(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value

	switch operator {
	case "==":
		return nativeBool(leftValue == rightValue)
	case "!=":
		return nativeBool(leftValue != rightValue)
	default:
		return NIL
	}
}

// string operator
func evalStringOperator(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.String).Value
	var rightValue string
	switch right := right.(type) {
	case *object.String:
		rightValue = right.Value
	case *object.Boolean:
		return NIL
	case *object.Integer:
		rightValue = strconv.Itoa(int(right.Value))
	default:
		return NIL
	}

	switch operator {
	case "+":
		return &object.String{Value: leftValue + rightValue}
	case "*":
		n, err := strconv.Atoi(rightValue)
		if err != nil {
			return NIL
		}
		var out bytes.Buffer
		for i := 0; i < n; i++ {
			out.WriteString(leftValue)
		}
		return &object.String{Value: out.String()}
	case "==":
		return nativeBool(leftValue == rightValue)
	case "!=":
		return nativeBool(leftValue != rightValue)
	default:
		return NIL
	}

}

// Integer operator
func evalIntOperator(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "==":
		return nativeBool(leftValue == rightValue)
	case "!=":
		return nativeBool(leftValue != rightValue)
	case "<":
		return nativeBool(leftValue < rightValue)
	case ">":
		return nativeBool(leftValue > rightValue)
	case "<=":
		return nativeBool(leftValue <= rightValue)
	case ">=":
		return nativeBool(leftValue >= rightValue)
	default:
		return NIL
	}
}
