package parser

import (
	"khanhanh_lang/ast"
	"khanhanh_lang/token"
)

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt

}
