package parser

import (
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser had %d error", len(errors))
	for _, err := range errors {
		t.Errorf("parser error : %q", err)
	}
	t.FailNow()
}
