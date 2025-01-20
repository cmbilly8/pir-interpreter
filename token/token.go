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
	LBRACKET  = "["
	RBRACKET  = "]"
	// Keywords
	F      = "F"
	YAR    = "YAR"
	GIVES  = "GIVES"
	IF     = "IF"
	LSIF   = "LSIF"
	LS     = "LS"
	OR     = "OR"
	AND    = "AND"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	STRING = "STRING"
	FOR    = "FOR"
	BREAK  = "BREAK"
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

func (tok *Token) IsExpressionTerminator() bool {
	switch tok.Type {
	case PERIOD, COLOGNE, BE:
		return true
	default:
		return false
	}
}

const (
	_ int = iota
	PREC_LOWEST
	PREC_LOGIC_OR
	PREC_LOGIC_AND
	PREC_EQUALS      // =
	PREC_LESSGREATER // > or <
	PREC_SUM         // +
	PREC_PRODUCT     // *
	PREC_PREFIX      // -X or !X
	PREC_CALL        // myFunction(X)
	PREC_INDEX       // array[i]
)

func (tok *Token) Precedence() int {
	switch tok.Type {
	case EQUAL, NOTEQUAL:
		return PREC_EQUALS
	case OR:
		return PREC_LOGIC_OR
	case AND:
		return PREC_LOGIC_AND
	case LESS, GREATER:
		return PREC_LESSGREATER
	case PLUS, MINUS:
		return PREC_SUM
	case FSLASH, STAR:
		return PREC_PRODUCT
	case LPAREN:
		return PREC_CALL
	case LBRACKET:
		return PREC_INDEX
	default:
		return PREC_LOWEST
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
	case "or":
		return OR
	case "and":
		return AND
	case "ay":
		return TRUE
	case "nay":
		return FALSE
	case "for":
		return FOR
	case "break":
		return BREAK
	default:
		return IDENT
	}
}
