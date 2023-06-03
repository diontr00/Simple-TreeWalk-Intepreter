package parser

import (
	"khanhanh_lang/ast"
	"khanhanh_lang/lexer"
	"testing"
)

func TestReturnStatement(t *testing.T) {
	input := `return 5;
  return 10;
  return 993322;
  `

	lex := lexer.New(input)
	par := New(lex)
	program := par.ParseProgram()
	if len(program.Statements) != 3 {
		t.Fatalf("program statement does not contain 3 statement. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return' , got %q", returnStmt.TokenLiteral())
		}

	}
}
