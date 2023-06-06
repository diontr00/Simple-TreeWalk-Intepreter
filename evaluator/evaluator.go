package evaluator

import (
	"bytes"
	"fmt"
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

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func Eval(node ast.Node, tracker *object.Tracker) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, tracker)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, tracker)
	case *ast.PrefixExpression:
		right := Eval(node.Right, tracker)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, tracker)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, tracker)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBool(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.BlockStatement:
		return evalBlockStatement(node, tracker)
	case *ast.IfExpression:
		return evalIfExpression(node, tracker)
	case *ast.ReturnStatement:
		// evaluate expression associated with the ast return statement  , and then wrap the result insid  the object ReturnValue to keep track
		val := Eval(node.ReturnValue, tracker)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, tracker)
		if isError(val) {
			return val
		}
		tracker.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, tracker)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: tracker, Body: body}
	case *ast.CallExpression:
		function := Eval(node.Function, tracker)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, tracker)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	}
	return nil

}

// traver the AST , and parse each statement
func evalProgram(stmts []ast.Statement, tracker *object.Tracker) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement, tracker)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}

	}
	return result
}

// Dealing with nested block effect
func evalBlockStatement(block *ast.BlockStatement, tracker *object.Tracker) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, tracker)
		if result != nil {
			result_type := result.Type()
			if result_type == object.RETURN_VALUE_OBJ || result_type == object.ERROR_OBJ {
				return result
			}
			// here only return object.returnValue and not  object.returnValue.Value (Not upwrapped)
			// so the nested block has chance to evaluated
		}

	}
	return result
}

// Evaluate scope and calling  associated function
// -------------------------------------------------------------------------------
func evalExpressions(exps []ast.Expression, env *object.Tracker) []object.Object {
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
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}
	extendEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Tracker {
	env := object.NewEnclosedTracker(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
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
		return newError("[Error]: Unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefix(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("[Error]: Unknown operator -%s", right.Type())
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

func evalIfExpression(ie *ast.IfExpression, tracker *object.Tracker) object.Object {
	condition := Eval(ie.Condition, tracker)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, tracker)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, tracker)
	} else {
		return NIL
	}
}
func isTruthy(obj object.Object) bool {
	switch obj {
	case NIL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalIdentifier(node *ast.Identifier, tracker *object.Tracker) object.Object {
	val, ok := tracker.Get(node.Value)
	if !ok {
		return newError("[Error]: Identifier not found: " + node.Value)
	}
	return val
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
		return newError("[Error]: Mismatch %s %s %s", left.Type(), operator, right.Type())
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
		return newError("[Error]: Unknown %s %s %s", left.Type(), operator, right.Type())
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
		return newError("[Error]: Unknown operator %s", operator)
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
		return newError("[Error]: Unknwon operator %s %s %s", left.Type(), operator, right.Type())
	}
}
