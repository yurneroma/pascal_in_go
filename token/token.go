package token

const (
	INTEGER = "INTEGER"
	PLUS    = "PLUS"
	MINUS   = "MINUS"
	DIV     = "DIV"
	MUL     = "MUL"
	EOF     = "EOF"
	LPAREN  = "LPAREN"
	RPAREN  = "RPAREN"
	ILLEGAL = "ILLEGAL"
	END     = "END"
	DOT     = "DOT"
	ASSIGN  = "ASSGIGN"
	SEMI    = "SEMI"
	ID      = "ID"
	BEGIN   = "BEGIN"
)

//Type represents the type of a token
type Type string

// Token represents the atom token, with a type and literal value
type Token struct {
	Type    Type   `json:"type"`
	Literal string `json:"literal"`
}
