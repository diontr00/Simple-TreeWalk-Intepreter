package ast

import (
	"bytes"
	"khanhanh_lang/token"
)

// TREE
// ---------------------------------------------------------------------------------
// Every node in AST has to implement this
type Node interface {
	TokenLiteral() string //Should return the literal value of the token
	String() string
}

type Statement interface {
	Node
	statementNode() // empty but need  impl to satisfy as the statement
}

type Expression interface {
	Node
	expressionNode() // Empty but need impl to satisfy as expression
}

type Program struct {
	Statements []Statement // Program is simply just a series of statements
}

// Return the whole program back as string to debug and test
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Satisfy both the Node and Expression  interface
type Identifier struct {
	Token token.Token // the IDENT token
	Value string
}

func (i *Identifier) String() string       { return i.Value }
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (l *Identifier) expressionNode()      {}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

//----------------------------------------------------------------------------------

// LITERALS
// ---------------------------------------------------------------------------------
// represent integer
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// represent boolean
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

//-----------------------------------------------------------------------------------

// STATEMENTS
// ----------------------------------------------------------------------------------
type LetStatement struct {
	Token token.Token // The LET token
	Name  *Identifier // To keep the nodeType small , we use the Identifier for both binding variable value and identifier in the right part , but to clarify that the identifier used in this case doesn't produce a value  but  still satisfy the Expression interface
	Value Expression  // point to the expression
}

func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	Token       token.Token // the RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// only the wrapper around the expression ,since  it is totally fine to write either
// let x = 1; // let stament  OR  x + 2; // expression by itself
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

//------------------------------------------------------------------------------------

// EXPRESSION
// -----------------------------------------------------------------------------------
// represent prefix expression
type PrefixExpression struct {
	Token    token.Token // prefix token such as !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// represent infix expression
type InfixExpression struct {
	Token    token.Token // The operator token +
	Left     Expression
	Right    Expression
	Operator string
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" ")
	out.WriteString(ie.Operator)
	out.WriteString(" ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}
