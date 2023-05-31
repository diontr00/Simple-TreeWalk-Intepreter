package parser

import (
	"khanhanh_lang/ast"
	"khanhanh_lang/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `
  let  x = 5; 
  let  y = 10; 
  let  foobar = 101234; 
  `
	lex := lexer.New(input)
	par := New(lex)
	program := par.ParseProgram()
	checkParserErrors(t, par)
	if program == nil {
		t.Fatalf("ParseProgarm return nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf(
			"program.Statements does not  contain 3 statements. got=%d",
			len(program.Statements),
		)
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"}, {"y"}, {"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}

}

func testLetStatement(t *testing.T, s ast.Statement, expectedName string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not  'let'. got=%q", s.TokenLiteral())
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != expectedName {
		t.Errorf("letStmt.Name.Value not '%s'. got='%s'", expectedName, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != expectedName {
		t.Errorf("s.Name not '%s'. got = %s", expectedName, letStmt.Name)
		return false
	}
	return true
}
