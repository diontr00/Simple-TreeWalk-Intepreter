package lexer

import (
	"khanhanh_lang/token"
	"testing"
)

func TestBasicToken(t *testing.T) {

	input := `=+-*/(){},;<>!9"🥳"if else return true false == != <= >=`
	tests := []struct {
		expectedTokenType token.TokenType
		expectedLiteral   string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.BANG, "!"},
		{token.INT, "9"},
		{token.STRING, "🥳"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.EQ, "=="},
		{token.NOT_EQ, "!="},
		{token.LT_EQ, "<="},
		{token.GT_EQ, ">="},
	}

	l := New(input)

	for _, test := range tests {

		test_token := l.NextToken()
		if test_token.Type != test.expectedTokenType {
			t.Fatalf(
				"wrong token type , expect : %q , got %q",
				test.expectedTokenType,
				test_token.Type,
			)
		}
		if test_token.Literal != test.expectedLiteral {
			t.Fatalf(
				"wrong literal , expected : %q , got %q",
				test.expectedLiteral,
				test_token.Literal,
			)
		}

	}

}

func TestExtendNextToken(t *testing.T) {
	input := `let five  =5;
  let ten = 10;
  let add = func(x, y) {
  x + y;
  };
  let result = add(five, ten);
  let unicode = "🥲 🥳";
  ! x > 5 ; 
  if ( 2 < 3) { 
      return true;
  } else { 
    return false;
  } 
  1 == 1; 
  2 != 1;  
  1 >= 1; 
  2 <= 2; 
  `

	lex := New(input)

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "func"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "unicode"},
		{token.ASSIGN, "="},
		{token.STRING, "🥲 🥳"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.IDENT, "x"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "2"},
		{token.LT, "<"},
		{token.INT, "3"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "1"},
		{token.EQ, "=="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.INT, "2"},
		{token.NOT_EQ, "!="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.INT, "1"},
		{token.GT_EQ, ">="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.INT, "2"},
		{token.LT_EQ, "<="},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	for _, test := range tests {
		test_token := lex.NextToken()
		if test_token.Type != test.expectedType {
			t.Fatalf(
				"[extend-test] -- Wrong type :  expect %q , got  %q ",
				test.expectedType,
				test_token.Type,
			)
		}
		if test_token.Literal != test.expectedLiteral {
			t.Fatalf(
				"[extend-test] -- Wrong Literal : expect %q , got %q",
				test.expectedLiteral,
				test_token.Literal,
			)
		}

	}

}
