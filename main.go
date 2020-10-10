package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"pascal_in_go/token"
	"strconv"
	"unicode"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("read error")
	}

	interpreter := NewInterpreter(text)
	result := interpreter.Expr()
	fmt.Println(result)

}

func NewInterpreter(text string) *Interpreter {
	return &Interpreter{Text: text, Pos: 0, CurChar: text[0]}
}

func newToken(tokenType token.Type, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (interpreter *Interpreter) NextToken() token.Token {
	var tok token.Token
	for interpreter.CurChar != 0 {
		if unicode.IsSpace(rune(interpreter.CurChar)) {
			interpreter.SkipWhiteSpace()
			continue
		}

		if unicode.IsDigit(rune(interpreter.CurChar)) {
			val := interpreter.integer()
			tok.Type = token.INTEGER
			tok.Literal = val
			return tok
		}
		if interpreter.CurChar == '+' {
			tok.Type = token.PLUS
			tok.Literal = "+"
			interpreter.advance()
			return tok
		}

		if interpreter.CurChar == '-' {
			tok.Type = token.MINUS
			tok.Literal = "-"
			interpreter.advance()
			return tok
		}

		if interpreter.CurChar == 0 {
			tok.Type = token.EOF
			tok.Literal = ""
			interpreter.advance()
			return tok
		}
	}
	return newToken(token.ILLEGAL, interpreter.CurChar)
}

func (interpreter *Interpreter) integer() string {
	result := ""
	for interpreter.CurChar != 0 && unicode.IsDigit(rune(interpreter.CurChar)) {
		result += string(interpreter.CurChar)
		interpreter.advance()
	}

	return result
}

func (interpreter *Interpreter) SkipWhiteSpace() {
	for interpreter.CurChar == ' ' || interpreter.CurChar == '\t' || interpreter.CurChar == '\r' || interpreter.CurChar == '\n' {
		interpreter.advance()
	}
}

func (interpreter *Interpreter) advance() {
	interpreter.Pos++
	if interpreter.Pos > len(interpreter.Text)-1 {
		interpreter.CurChar = 0
	} else {
		interpreter.CurChar = interpreter.Text[interpreter.Pos]
	}

}

func (interpreter *Interpreter) Eat(tokenType token.Type) {
	if interpreter.CurToken.Type == tokenType {
		interpreter.CurToken = interpreter.NextToken()
	} else {
		log.Fatal("type not match, cur and input type is ", interpreter.CurToken.Type, ":", tokenType)
	}
}

func (interpreter *Interpreter) Expr() (res int) {
	left := interpreter.NextToken()
	interpreter.CurToken = left
	interpreter.Eat(token.INTEGER)
	op := interpreter.CurToken
	if op.Type == token.PLUS {
		interpreter.Eat(token.PLUS)

	} else {
		interpreter.Eat(token.MINUS)

	}
	right := interpreter.CurToken
	interpreter.Eat(token.INTEGER)

	intl, _ := strconv.Atoi(left.Literal)
	intr, _ := strconv.Atoi(right.Literal)
	if op.Type == token.PLUS {
		res = intl + intr
	} else {
		res = intl - intr
	}
	return
}

type Interpreter struct {
	Text     string      `json:"text"`
	Pos      int         `json:"pos"`
	CurToken token.Token `json:"curToken"`
	CurChar  byte        `json:"curChar"`
}
