package lexer

import (
	"pascal_in_go/token"
	"unicode"
)

type Lexer struct {
	Text    string `json:"text"`
	Pos     int    `json:"pos"`
	CurChar byte   `json:"curChar"`
}

func NewLexer(text string) Lexer {
	return Lexer{Text: text, Pos: 0, CurChar: 0}
}

func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token
	for lexer.CurChar != 0 {
		if unicode.IsSpace(rune(lexer.CurChar)) {
			lexer.skipWhiteSpace()
			continue
		}

		if unicode.IsDigit(rune(lexer.CurChar)) {
			val := lexer.integer()
			tok.Type = token.INTEGER
			tok.Literal = val
			return tok
		}
		if lexer.CurChar == '+' {
			tok.Type = token.PLUS
			tok.Literal = "+"
			lexer.advance()
			return tok
		}

		if lexer.CurChar == '-' {
			tok.Type = token.MINUS
			tok.Literal = "-"
			lexer.advance()
			return tok
		}

		if lexer.CurChar == 0 {
			tok.Type = token.EOF
			tok.Literal = ""
			lexer.advance()
			return tok
		}
	}
	return newToken(token.ILLEGAL, lexer.CurChar)
}

func (lexer *Lexer) integer() string {
	result := ""
	for lexer.CurChar != 0 && unicode.IsDigit(rune(lexer.CurChar)) {
		result += string(lexer.CurChar)
		lexer.advance()
	}

	return result
}

func (lexer *Lexer) skipWhiteSpace() {
	for lexer.CurChar == ' ' || lexer.CurChar == '\t' || lexer.CurChar == '\r' || lexer.CurChar == '\n' {
		lexer.advance()
	}
}

func (lexer *Lexer) advance() {
	lexer.Pos++
	if lexer.Pos > len(lexer.Text)-1 {
		lexer.CurChar = 0
	} else {
		lexer.CurChar = lexer.Text[lexer.Pos]
	}

}

func newToken(tokenType token.Type, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
