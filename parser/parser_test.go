package parser

// Helper for other tests
import (
	"fmt"
	"khanhanh_lang/ast"
	"testing"
)

// TEST LITERALS
// ---------------------------------------------------------------------------------
func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not  *ast.Identifier. got=%T", exp)
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true

}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}
	return true
}

func testBooleanliteral(t *testing.T, il ast.Expression, value bool) bool {
	bool, ok := il.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("il not *ast.BoolLiteral. got=%T", il)
		return false
	}
	if bool.Value != value {
		t.Errorf("integ.Value not %v. got=%v", value, bool.Value)
	}
	if bool.TokenLiteral() != fmt.Sprintf("%v", value) {
		t.Errorf("integ.TokenLiteral not %v. got=%s", value, bool.TokenLiteral())
		return false
	}
	return true

}

//---------------------------------------------------------------------------------

// TEST EXPRESSION
// ---------------------------------------------------------------------------------
func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {

	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))

	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanliteral(t, exp, v)
	}

	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpresion. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testPrefixExpression(
	t *testing.T,
	exp ast.Expression,
	operator string,
	right interface{},
) bool {
	opExp, ok := exp.(*ast.PrefixExpression)
	if !ok {
		t.Errorf("exp is not ast.PrefixExpresion. got=%T(%s)", exp, exp)
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testExpression(t *testing.T, statement ast.Statement) *ast.ExpressionStatement {
	stmt, ok := statement.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"Program Statements[0] is not ast.ExpressionStatement , got=%T",
			statement,
		)
		return nil
	}
	return stmt

}

// --------------------------------------------------------------------------------
// Single test program length and  error check
func testProgramLength(t *testing.T, program *ast.Program) {
	if len(program.Statements) != 1 {
		t.Fatalf("Program has not enough statements. got=%d", program.Statements[0])
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser had %d error", len(errors))
	for _, err := range errors {
		t.Errorf("parser error : %q", err)
	}
	t.FailNow()
}
