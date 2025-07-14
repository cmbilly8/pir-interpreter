// lexer/lexer.go
package lexer

import "pir-interpreter/token"

type Lexer struct {
	input         string
	position      int
	readPosition  int
	ch            byte
	curLine       int
	curCharOfLine int
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.curCharOfLine++
}

func (l *Lexer) peekNext() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}

}

func (l *Lexer) NextToken() token.Token {
	var currentToken token.Token

	l.ignoreWhitespace()
	l.ignoreComment()

	switch l.ch {
	case '+':
		currentToken = l.newToken(token.PLUS, "+")
	case '-':
		currentToken = l.newToken(token.MINUS, "-")
	case '*':
		currentToken = l.newToken(token.STAR, "*")
	case '/':
		currentToken = l.newToken(token.FSLASH, "/")
	case '!':
		currentToken = l.newToken(token.AAAA, "!")
	case '<':
		if l.peekNext() == '>' {
			l.readChar()
			currentToken = l.newToken(token.NOTEQUAL, "<>")
		} else if l.peekNext() == '=' {
			l.readChar()
			currentToken = l.newToken(token.LESSEQ, "<=")
		} else {
			currentToken = l.newToken(token.LESS, "<")
		}
	case '>':
		if l.peekNext() == '=' {
			l.readChar()
			currentToken = l.newToken(token.GREATEREQ, ">=")
		} else {
			currentToken = l.newToken(token.GREATER, ">")
		}
	case '=':
		currentToken = l.newToken(token.EQUAL, "=")
	case ',':
		currentToken = l.newToken(token.COMMA, ",")
	case '.':
		currentToken = l.newToken(token.PERIOD, ".")
	case ';':
		currentToken = l.newToken(token.SEMICOLON, ";")
	case ':':
		currentToken = l.newToken(token.COLOGNE, ":")
	case '(':
		currentToken = l.newToken(token.LPAREN, "(")
	case ')':
		currentToken = l.newToken(token.RPAREN, ")")
	case '{':
		currentToken = l.newToken(token.LBRACE, "{")
	case '}':
		currentToken = l.newToken(token.RBRACE, "}")
	case '[':
		currentToken = l.newToken(token.LBRACKET, "[")
	case ']':
		currentToken = l.newToken(token.RBRACKET, "]")
	case '4':
		currentToken = l.newToken(token.FOR, "4")
	case '\'', '"':
		str := l.readString()
		currentToken = l.newToken(token.STRING, str)
	case 0:
		currentToken = l.newToken(token.EOF, "")
	default:
		if isCharLetter(l.ch) {
			literal := l.readIdentifier()
			tokType := token.LookupIdent(literal)
			return l.newToken(tokType, literal)
		} else if isCharNumber(l.ch) {
			literal := l.readNumber()
			return l.newToken(token.INT, literal)
		} else {
			currentToken = l.newToken(token.ILLICIT, string(l.ch))
		}
	}
	l.readChar()
	return currentToken
}

func (l *Lexer) readString() string {
	// Could be either ' or "
	endQuote := l.ch
	l.readChar()
	start := l.position
	for l.ch != byte(endQuote) && l.ch != 0 {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) ignoreComment() {
	if l.ch == '$' {
		for l.ch != '\n' {
			l.readChar()
		}
		l.ignoreWhitespace()
	}
}

func (l *Lexer) ignoreWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.curLine += 1
			l.curCharOfLine = 0
		}
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	firstIndex := l.position
	for isCharLetter(l.ch) {
		l.readChar()
	}
	return l.input[firstIndex:l.position]

}

func isCharLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	firstIndex := l.position
	for isCharNumber(l.ch) {
		l.readChar()
	}
	return l.input[firstIndex:l.position]
}

func isCharNumber(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) newToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: literal,
		LineNum: l.curLine,
		CharNum: l.curCharOfLine,
	}
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}
