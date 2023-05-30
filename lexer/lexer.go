package lexer

import "khanhanh_lang/token"

type Lexer struct {
	input        string // input literal to lex into token
	position     int    // position in the input (point to  ch)
	readPosition int    // read position in the input (after ch )  , always point to the next position where we're going to read from
	ch           rune   // current character support UTF8
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// read the next character and advance our position in the input string
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // NUL
	} else {
		l.ch = rune(l.input[l.readPosition])
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// use to tokenize the input string current char and call readChar internally to advance the readPosition to the next char
func (l *Lexer) NextToken() token.Token {
	var resultToken token.Token
	l.skipWhiteSpace()
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			resultToken = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			resultToken = newToken(token.ASSIGN, l.ch)
		}
	case '-':
		resultToken = newToken(token.MINUS, l.ch)
	case '+':
		resultToken = newToken(token.PLUS, l.ch)
	case '*':
		resultToken = newToken(token.ASTERISK, l.ch)
	case '/':
		resultToken = newToken(token.SLASH, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			resultToken = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			resultToken = newToken(token.BANG, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			resultToken = token.Token{Type: token.LT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			resultToken = newToken(token.LT, l.ch)
		}
	case '>':
		resultToken = newToken(token.GT, l.ch)
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			resultToken = token.Token{Type: token.GT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			resultToken = newToken(token.GT, l.ch)
		}
	case '(':
		resultToken = newToken(token.LPAREN, l.ch)
	case ')':
		resultToken = newToken(token.RPAREN, l.ch)
	case '{':
		resultToken = newToken(token.LBRACE, l.ch)
	case '}':
		resultToken = newToken(token.RBRACE, l.ch)
	case ';':
		resultToken = newToken(token.SEMICOLON, l.ch)
	case ',':
		resultToken = newToken(token.COMMA, l.ch)
	case 0:
		resultToken.Literal = ""
		resultToken.Type = token.EOF
	case '"':
		resultToken.Type = token.STRING
		resultToken.Literal = l.readString()

	default:
		if isLetter(l.ch) {
			resultToken.Literal = l.readIdentifier()
			// either identifier or the keyword
			resultToken.Type = token.LookUpKeyword(resultToken.Literal)
			return resultToken
		} else if isDigit(l.ch) {
			resultToken.Literal = l.readNumber()
			resultToken.Type = token.INT
			return resultToken
		} else {
			resultToken = newToken(token.ILLEGAL, l.ch)
		}

	}
	l.readChar()
	return resultToken

}

// Check if the character is digit
func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'

}

// Check whether the character is in (a->z , A -> Z or _)
func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// Continute to read the identifier until we get the non literal  position
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// Read Number
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]

}

// Read String literal
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]

}

// ignore the white space in the input
func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// Peek ahead not moving the readPosition or position
func (l *Lexer) peekChar() rune {
	if l.position > len(l.input) {
		return 0
	} else {
		return rune(l.input[l.readPosition])
	}

}

func newToken(t token.TokenType, l rune) token.Token {
	return token.Token{Type: t, Literal: string(l)}
}
