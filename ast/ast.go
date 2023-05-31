package ast

import "khanhanh_lang/token"

// Every node in AST has to implement this
type Node interface {
	TokenLiteral() string //Should return the literal value of the token
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement // Program is simply just a series of statements
}

// Satisfy both the Node and Expression  interface
type Identifier struct {
	Token token.Token // the IDENT token
	Value string
}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (l *Identifier) expressionNode()      {}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token // The LET token
	Name  *Identifier // To keep the nodeType small , we use the Identifier for both binding variable value and identifier in the right part , but to clarify that the identifier used in this case doesn't produce a value  but  still satisfy the Expression interface
	Value Expression  // point to the expression
}

func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) statementNode()       {}

type ReturnStatement struct {
	Token       token.Token // the RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
