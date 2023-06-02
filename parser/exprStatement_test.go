package parser

import (
	"khanhanh_lang/lexer"
	"testing"
)

// Test Literals
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

//-------------------------------------------------------------------------------------------------------

// TEST PREFIX
// -----------------------------------------------------------------------
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

// ----------------------------------------------------------------------------------
// TEST INFIX
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
