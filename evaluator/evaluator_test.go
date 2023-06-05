package evaluator

import (
	"fmt"
	"khanhanh_lang/lexer"
	"khanhanh_lang/object"
	"khanhanh_lang/parser"
	"testing"
)

// Literals
func TestLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		// Literal
		{"5", int64(5)},
		{"10", int64(10)},
		{`"hello"`, "hello"},
		{`"world5"`, "world5"},
		{`"🤣"`, "🤣"},
		{"1023012312", int64(1023012312)},
		{"true", true},
		// Prefix
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{`!"hello"`, false},
		{"-5", int64(-5)},
		{"--10", int64(10)},
		{"-true", nil},
		// Infix
		{"5 + 5 + 5 + 5 -10", int64(10)},
		{"2 * 2 * 2 * 2 * 2", int64(32)},
		{"-50 + 100 - 50", int64(0)},
		{"20 + 2 *-10", int64(0)},
		{"5 * 2 + 10", int64(20)},
		{"5 + 2 * 10", int64(25)},
		{"20 + 2 * -10", int64(0)},
		{"50 / 2 * 2 + 10", int64(60)},
		{"2 * (5 + 10)", int64(30)},
		{"3 * 3 * 3 + 10", int64(37)},
		{"3 * (3 * 3) + 10", int64(37)},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", int64(50)},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 == 2", false},
		{"1 <= 2", true},
		{"1 < 1", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{`"hello" + "world"`, "helloworld"},
		{`"hello" == "hello"`, true},
		{`"hello" != "hello"`, false},
		{`"hello" * 2`, "hellohello"},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testTypeObject(t, evaluated, test.expected)
	}
}

func testEval(input string) object.Object {
	lex := lexer.New(input)
	par := parser.New(lex)
	program := par.ParseProgram()
	return Eval(program)
}

func testTypeObject(t *testing.T, obj object.Object, expect any) bool {
	if obj == nil {
		t.Errorf("Obj is nil , expected value(%v)", expect)
		return false
	}
	switch typi := expect.(type) {
	case int64:
		result, ok := obj.(*object.Integer)
		if !ok {
			t.Errorf("object is not an integer. got=%T (%+v)", obj, obj)
		}
		if result.Value != expect {
			t.Errorf("object has wrong value. got=%d , want=%d", result.Value, expect)
			return false
		}

	case bool:
		result, ok := obj.(*object.Boolean)
		if !ok {
			t.Errorf("object is not a boolean. got=%T (%+v)", obj, obj)
			return false
		}
		if result.Value != expect {
			t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expect)
			return false
		}
	case string:
		result, ok := obj.(*object.String)
		if !ok {
			t.Errorf("object is not a string. got=%T (%+v)", obj, obj)
			return false
		}
		if result.Value != expect {
			t.Errorf("object has wrong value. got=%s , want=%s", result.Value, expect)
		}
		return false

	case nil:
		result, ok := obj.(*object.Nil)
		if !ok {
			t.Errorf("Expected nil , got=%T", obj)
			return false
		}
		if result.Inspect() != "nil" {
			t.Errorf("expect nil as the value  but got : %s", result.Inspect())
		}
		return false
	default:
		fmt.Printf("Not yet implemented for %T", typi)
		return false

	}
	return false
}
