package parser

import (
	"khanhanh_lang/ast"
	"khanhanh_lang/lexer"
	"testing"
)

// TEST LITERALS
// ---------------------------------------------------------------------------------------------------
func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	lex := lexer.New(input)
	par := New(lex)

	program := par.ParseProgram()

	checkParserErrors(t, par)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	stmt := testExpression(t, program.Statements[0])
	testIdentifier(t, stmt.Expression, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5"
	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	testProgramLength(t, program)
	stmt := testExpression(t, program.Statements[0])

	testIntegerLiteral(t, stmt.Expression, 5)
}

func TestBooleanLiteralExpression(t *testing.T) {
	input := "true"
	lex := lexer.New(input)
	par := New(lex)
	program := par.ParseProgram()
	checkParserErrors(t, par)
	testProgramLength(t, program)
	stmt := testExpression(t, program.Statements[0])
	testBooleanliteral(t, stmt.Expression, true)

}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func(x,y) { x +  y; }`
	lex := lexer.New(input)
	par := New(lex)
	program := par.ParseProgram()
	checkParserErrors(t, par)
	testProgramLength(t, program)
	stmt := testExpression(t, program.Statements[0])

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "func() {};", expectedParams: []string{}},
		{input: "func(x) {};", expectedParams: []string{"x"}},
		{input: "func(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, test := range tests {
		lex := lexer.New(test.input)
		par := New(lex)
		program := par.ParseProgram()
		testProgramLength(t, program)
		checkParserErrors(t, par)
		stmt := testExpression(t, program.Statements[0])
		function := stmt.Expression.(*ast.FunctionLiteral)
		if len(function.Parameters) != len(test.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(test.expectedParams), len(function.Parameters))
		}
		for i, ident := range test.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}

	}
}

//-------------------------------------------------------------------------------------------------------

// TEST PREFIX
// ------------------------------------------------------------------------------------------------------
func TestParsingPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		val      any
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		testProgramLength(t, program)
		stmt := testExpression(t, program.Statements[0])
		testPrefixExpression(t, stmt.Expression, test.operator, test.val)
	}
}

// -------------------------------------------------------------------------------------------------------------

// TEST INFIX
// -------------------------------------------------------------------------------------------------------------
func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true != false", true, "!=", false},
	}
	for _, test := range tests {
		lex := lexer.New(test.input)
		par := New(lex)
		program := par.ParseProgram()
		testProgramLength(t, program)
		checkParserErrors(t, par)
		stmt := testExpression(t, program.Statements[0])
		testInfixExpression(t, stmt.Expression, test.leftValue, test.operator, test.rightValue)

	}
}

// TEST PRECEDENCE
func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b", "((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},

		{
			"a * b * c",
			"((a * b) * c)",
		},

		{
			"a * b / c",
			"((a * b) / c)",
		},

		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},

		{
			"3 + 4; - 5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			" 5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 +5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},

		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
	}

	for _, test := range tests {
		lex := lexer.New(test.input)
		par := New(lex)
		program := par.ParseProgram()
		checkParserErrors(t, par)
		actual := program.String()
		if actual != test.expected {
			t.Errorf("expected = %q , got=%q", test.expected, actual)
		}
	}
}

//----------------------------------------------------------------------------------------------------------------

// TEST CONDITIONALS
// ----------------------------------------------------------------------------------------------------------------
func TestIfExpression(t *testing.T) {
	input := `if (x < y){x}`
	lex := lexer.New(input)
	par := New(lex)
	program := par.ParseProgram()
	testProgramLength(t, program)
	checkParserErrors(t, par)
	stmt := testExpression(t, program.Statements[0])

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression.got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("Consequence len is not 1. got=%d\n", len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not expression statement. got=%T", exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if exp.Alternative != nil {
		t.Errorf("exp.Alternative was not nil. got=%+v", exp.Alternative)
	}

}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	lex := lexer.New(input)
	par := New(lex)
	program := par.ParseProgram()
	testProgramLength(t, program)
	stmt := testExpression(t, program.Statements[0])

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression.got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("Consequence len is not 1. got=%d\n", len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not expression statement. got=%T", exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("Alternative len is not 1. got=%d\n", len(exp.Consequence.Statements))
	}
	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not expression statement. got=%T", exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1 , 2 * 3 , 4 + 5)"
	lex := lexer.New(input)
	par := New(lex)
	program := par.ParseProgram()
	testProgramLength(t, program)
	stmt := testExpression(t, program.Statements[0])
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}
	if !testIdentifier(t, exp.Function, "add") {
		return
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}
	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}
