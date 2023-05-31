package parser

import (
	"fmt"
	"khanhanh_lang/ast"
	"khanhanh_lang/lexer"
	"khanhanh_lang/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token // same as position in the lexer, but instead of point to current ch , it point to current token
	peekToken token.Token // same as the readPosition in the lexer , but  instead of point to next ch, it point to the next token (both cur and Peek are needed for decision making)
	errors    []string
}

// Advancce the current Token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Create new Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	// Read two token so the current and peak token are both set
	p.nextToken()
	p.nextToken()
	return p
}

// Construct the root node
// Continute to parse until the end of the input , where parseStatement does the its job
// then push the parsed statement into program tree  which represent by the slice of Statements
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

// Check Current Token Type
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// Check Peek Token Type
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// move token if peekToken is expected
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}

}
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token : %s , but get %s", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
