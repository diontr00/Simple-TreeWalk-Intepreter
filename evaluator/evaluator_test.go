package evaluator

import (
	"fmt"
	"khanhanh_lang/lexer"
	"khanhanh_lang/object"
	"khanhanh_lang/parser"
	"testing"
)

type ErrorMesssage string
type FunctionObject struct {
	params []string
	body   string
}

// Literals
func TestEvaluator(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		// Literal
		//-----------------------------------------------------------------------------------------------
		{"5", int64(5)},
		{"10", int64(10)},
		{`"hello"`, "hello"},
		{`"world5"`, "world5"},
		{`"🤣"`, "🤣"},
		{"1023012312", int64(1023012312)},
		{"true", true},
		// Prefix
		//-----------------------------------------------------------------------------------------------
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{`!"hello"`, false},
		{"-5", int64(-5)},
		{"--10", int64(10)},
		// Infix
		//-----------------------------------------------------------------------------------------------
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
		// Condition
		//-----------------------------------------------------------------------------------------------
		{"if(true) {10}", int64(10)},
		{"if(false){10}", nil},
		{"if (1) {10}", int64(10)},
		{"if (1 < 2) { 10 }", int64(10)},
		{"if (1 > 2) {10} else {20}", int64(20)},
		{"if (2 > 1) {10} else {20}", int64(10)},
		// Return
		//----------------------------------------------------------------------------------------------
		{"return 10", int64(10)},
		{"return  10; 9;", int64(10)},
		{"return 2 * 5; 9;", int64(10)},
		{"9; return 2 * 5 ; 9", int64(10)},
		{"if(10 < 11) { if (9 > 2){return 9} return 10}", int64(9)},
		//Error
		//-----------------------------------------------------------------------------------------------
		{"5 + true;", ErrorMesssage("[Error]: Mismatch INTEGER + BOOLEAN")},
		{"5 + true; 5;", ErrorMesssage("[Error]: Mismatch INTEGER + BOOLEAN")},
		{"-true", ErrorMesssage("[Error]: Unknown operator -BOOLEAN")},
		{"true + false", ErrorMesssage("[Error]: Unknown operator: BOOLEAN + BOOLEAN")},
		{"5; true + false; 5", ErrorMesssage("[Error]: Unknown operator: BOOLEAN + BOOLEAN")},
		{
			"if (10 > 1) {true + false}",
			ErrorMesssage("[Error]: Unknown operator : BOOLEAN + BOOLEAN"),
		},
		{"foobar", ErrorMesssage("[Warn]: Identifier not found")},
		//Binding
		//-------------------------------------------------------------------------------------------------
		{"let a = 5; a;", int64(5)},
		{"let a = 5 * 5; a;", int64(25)},
		{"let a = 5; let b = a; b;", int64(5)},
		{"let a = 5; let b = a; let c = a + b + 5; c;", int64(15)},
		{"foobar", ErrorMesssage("[Error]: Identifier not found: foobar")},
		// Function
		//-------------------------------------------------------------------------------------------------
		{"func(x){x + 2}", FunctionObject{params: []string{"x"}, body: "(x + 2)"}},
		{"let identity = func(x) {x ;}; identity(5);", int64(5)},
		{"let identity = func(x) { return x; }; identity(5);", int64(5)},
		{"let double = func(x) { x * 2; }; double(5);", int64(10)},
		{"let add = func(x, y) { x + y; }; add(5, 5);", int64(10)},
		{"let add = func(x, y) { x + y; }; add(5 + 5, add(5, 5));", int64(20)},
		{"func(x) { x; }(5)", int64(5)},
		{"let Add=func(x){func(y) { x + y };}; let addTwo= Add(2);addTwo(2)", int64(4)},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testTypeObject(t, evaluated, test.expected)
	}
}

func testEval(input string) object.Object {
	lex := lexer.New(input)
	par := parser.New(lex)
	tracker := object.NewTracker()
	program := par.ParseProgram()
	return Eval(program, tracker)
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
	case ErrorMesssage:
		_, ok := obj.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", obj, obj)
			t.Errorf("Expect : %s", expect)
		}
	case FunctionObject:
		fn, ok := obj.(*object.Function)
		if !ok {
			t.Fatalf("object is not a function. got=%T (%+v)", obj, obj)
		}
		expect_fn, ok := expect.(FunctionObject)
		if !ok {
			t.Fatalf("Cannot compare expect of type =%T and functionOBject", expect_fn)
		}
		if len(fn.Parameters) != len(expect_fn.params) {
			t.Fatalf("Wrong expected parameters length , got=%d , expect=%d", len(fn.Parameters), len(expect_fn.params))
		}
		if expect_fn.body != fn.Body.String() {
			t.Fatalf("body is not %q. got=%q", expect_fn.body, fn.Body.String())
		}

	default:
		fmt.Printf("Not yet implemented for %T \n", typi)
		return false

	}
	return false
}
