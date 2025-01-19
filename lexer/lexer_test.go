package lexer

import (
	"pir-interpreter/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `yar five be 5.
	yar ten be 10.
	yar add be f(x, y):
		gives x + y..
	yar result be add(five, ten).
	!-/*5.
	5 < 10 > 5.
	if 5 < 10:
		gives ay.
	ls:
		gives nay.
    "yes".
    'yes no'.
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.YAR, "yar"},
		{token.IDENT, "five"},
		{token.BE, "be"},
		{token.INT, "5"},
		{token.PERIOD, "."},
		{token.YAR, "yar"},
		{token.IDENT, "ten"},
		{token.BE, "be"},
		{token.INT, "10"},
		{token.PERIOD, "."},
		{token.YAR, "yar"},
		{token.IDENT, "add"},
		{token.BE, "be"},
		{token.F, "f"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.COLOGNE, ":"},
		{token.GIVES, "gives"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.PERIOD, "."},
		{token.PERIOD, "."},
		{token.YAR, "yar"},
		{token.IDENT, "result"},
		{token.BE, "be"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.PERIOD, "."},
		{token.AAAA, "!"},
		{token.MINUS, "-"},
		{token.FSLASH, "/"},
		{token.STAR, "*"},
		{token.INT, "5"},
		{token.PERIOD, "."},
		{token.INT, "5"},
		{token.LESS, "<"},
		{token.INT, "10"},
		{token.GREATER, ">"},
		{token.INT, "5"},
		{token.PERIOD, "."},
		{token.IF, "if"},
		{token.INT, "5"},
		{token.LESS, "<"},
		{token.INT, "10"},
		{token.COLOGNE, ":"},
		{token.GIVES, "gives"},
		{token.TRUE, "ay"},
		{token.PERIOD, "."},
		{token.LS, "ls"},
		{token.COLOGNE, ":"},
		{token.GIVES, "gives"},
		{token.FALSE, "nay"},
		{token.PERIOD, "."},
		{token.STRING, "yes"},
		{token.PERIOD, "."},
		{token.STRING, "yes no"},
		{token.PERIOD, "."},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected: %q, got: %q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected: %q, got: %q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
