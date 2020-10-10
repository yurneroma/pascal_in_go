package token

const (
	INTEGER = "INTEGER"
	PLUS    = "PLUS"
	MINUS   = "MINUS"
	DIV     = "DIV"
	MUL     = "MUL"
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"
)

//Type represents the type of a token
type Type string

// Token represents the atom token, with a type and literal value
type Token struct {
	Literal string `json:"literal"`
	Type    Type   `json:"type"`
}
