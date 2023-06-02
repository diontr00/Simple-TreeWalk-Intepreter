package parser

import (
	"fmt"
	"khanhanh_lang/ast"
	"khanhanh_lang/token"
	"strconv"
)

// PRECEDENCE
// --------------------------------------------------------------------------
// our order of precedence for operation
const (
	_ int = iota
	LOWEST
	EQUAL       // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

// map infix operation precedence
var precedences = map[token.TokenType]int{
	token.EQ:       EQUAL,
	token.NOT_EQ:   EQUAL,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

// Check the current token precedence
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Check next  token precedence
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

//---------------------------------------------------------------------

// PARSING LITERAL
// ---------------------------------------------------------------------
// The parsing function that register in prefixParseFns for token IDENT
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// The parsing function that register in prefixParseFns for token INT
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

//----------------------------------------------------------------------------

// PARSING PREFIX
// ----------------------------------------------------------------------------
func (p *Parser) noPrefixParsfnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found ", t)
	p.errors = append(p.errors, msg)
}

// Create prefixExpression node by setting current token and calling nextToken to get the operand
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

// PARSING INFIX
// ------------------------------------------------------------------------
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression

}

//--------------------------------------------------------------------------

// MAIN
// -------------------------------------------------------------------------
// default expression parsing route to here
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt

}

// Prat Parsing to parse the right precedence
// What we are trying to do here is  nested ast.InfixExpression , so that the left most node always executed first
// For example 1 + 2 + 3 -> ((1 + 2) + 3)
// so the root node is ast.infixParseFns , which have two child , one is ast.IntegerLiteral which is three and the other is another ast.InfixExpression , the process continute
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParsfnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)

	}

	return leftExp
}
