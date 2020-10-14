package interpreter

import (
	"fmt"
	"log"
	"pascal_in_go/lexer"
	"pascal_in_go/token"
	"strconv"
)

type Interpreter struct {
	Lexer    lexer.Lexer `json:"lexer"`
	CurToken token.Token `json:"curToken"`
}

func NewInterpreter(lexer lexer.Lexer) *Interpreter {
	tok := lexer.NextToken()
	fmt.Println("initial token", tok)
	return &Interpreter{Lexer: lexer, CurToken: tok}
}

func (interpreter *Interpreter) Eat(tokenType token.Type) {
	if interpreter.CurToken.Type == tokenType {
		interpreter.CurToken = interpreter.Lexer.NextToken()
	} else {
		log.Fatal("type not match, cur and input  is ", interpreter.CurToken, ":", tokenType)
	}
}

func (interpreter *Interpreter) factor() int {
	tok := interpreter.CurToken
	interpreter.Eat(token.INTEGER)
	res, _ := strconv.Atoi(tok.Literal)
	return res
}
func (interpreter *Interpreter) Expr() (res int) {
	res = interpreter.factor()
	for interpreter.CurToken.Type == token.DIV || interpreter.CurToken.Type == token.MUL {
		tok := interpreter.CurToken
		if tok.Type == token.MUL {
			interpreter.Eat(token.MUL)
			res *= interpreter.factor()
		}
		if tok.Type == token.DIV {
			interpreter.Eat(token.DIV)
			res /= interpreter.factor()
		}
	}
	return res
}
