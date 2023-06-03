package ast

import (
	"bytes"
	"khanhanh_lang/token"
	"strings"
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

// Represent string
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(`"`)
	out.WriteString(sl.Value)
	out.WriteString(`"`)
	return out.String()
}

// represent boolean
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

// Function is the first class citizen , so it is an expression too
type FunctionLiteral struct {
	Token      token.Token //the fn token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, para := range fl.Parameters {
		params = append(params, para.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

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

// Represent block of statements, just like program contains statements, block also contain statements but in smaller scale
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, statement := range bs.Statements {
		out.WriteString(statement.String())
	}
	return out.String()
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

// Represent the if statement , that can contain an optional else
type IfExpression struct {
	Token       token.Token // IF token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ife *IfExpression) expressionNode()      {}
func (ife *IfExpression) TokenLiteral() string { return ife.Token.Literal }
func (ife *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ife.Condition.String())
	out.WriteString(" ")
	out.WriteString(ife.Consequence.String())

	if ife.Alternative != nil {
		out.WriteString("else")
		out.WriteString(ife.Alternative.String())
	}
	return out.String()
}

// Call expression consist  of an expression that result in a function when evaluated and a list of expression that are the arguments to this function call
// as add(2 , 3) is valid
// add(2 + 3 + 3  * 2) is also valid
// multiply(2  , add( 2 , 3)) // is also valid
type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
