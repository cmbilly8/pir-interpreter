package token

const (
	ILLICIT = "ILLICIT"
	EOF     = "EOF"
	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1343456
	// Operators
	BE       = "BE"
	PLUS     = "+"
	MINUS    = "-"
	STAR     = "*"
	FSLASH   = "/"
	AAAA     = "!"
	LESS     = "<"
	GREATER  = ">"
	EQUAL    = "="
	NOTEQUAL = "!="
	// Delimiters
	SQUOTE    = "'"
	COMMA     = ","
	PERIOD    = "."
	SEMICOLON = ";"
	COLOGNE   = ":"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	// Keywords
	F       = "F"
	YAR     = "YAR"
	GIVES   = "GIVES"
	IF      = "IF"
	LSIF    = "LSIF"
	LS      = "LS"
	CHANTEY = "CHANTEY"
	AVAST   = "AVAST"
	OR      = "OR"
	TRUE    = "TRUE"
	FALSE   = "FALSE"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	LineNum int
}

func (tok *Token) Is(t TokenType) bool {
	return tok.Type == t
}

func (tok *Token) IsNot(t TokenType) bool {
	return tok.Type != t
}

func (tok *Token) IsBlockTerminator() bool {
	switch tok.Type {
	case LS, LSIF, PERIOD, EOF:
		return true
	default:
		return false
	}
}

const (
	_ int = iota
	LOWEST
	EQUALS      // =
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

func (tok *Token) Precedence() int {
	switch tok.Type {
	case EQUAL, NOTEQUAL:
		return EQUALS
	case LESS, GREATER:
		return LESSGREATER
	case PLUS, MINUS:
		return SUM
	case FSLASH, STAR:
		return PRODUCT
	case LPAREN:
		return CALL
	default:
		return LOWEST
	}
}

func LookupIdent(ident string) TokenType {
	switch ident {
	case "f":
		return F
	case "yar":
		return YAR
	case "gives":
		return GIVES
	case "be":
		return BE
	case "if":
		return IF
	case "lsif":
		return LSIF
	case "ls":
		return LS
	case "chantey":
		return CHANTEY
	case "avast":
		return AVAST
	case "or":
		return OR
	case "ay":
		return TRUE
	case "nay":
		return FALSE
	default:
		return IDENT
	}
}
