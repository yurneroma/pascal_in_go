package interpreter

import (
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
	return &Interpreter{Lexer: lexer, CurToken: tok}
}

func (interpreter *Interpreter) Eat(tokenType token.Type) {
	if interpreter.CurToken.Type == tokenType {
		interpreter.CurToken = interpreter.Lexer.NextToken()
	} else {
		log.Fatal("type not match, cur and input type is ", interpreter.CurToken.Type, ":", tokenType)
	}
}

func (interpreter *Interpreter) term() int {
	tok := interpreter.CurToken
	interpreter.Eat(token.INTEGER)
	res, _ := strconv.Atoi(tok.Literal)
	return res
}

func (interpreter *Interpreter) Expr() (res int) {
	interpreter.CurToken = interpreter.Lexer.NextToken()
	res = interpreter.term()
	for interpreter.CurToken.Type == token.MINUS || interpreter.CurToken.Type == token.PLUS {
		tok := interpreter.CurToken
		if tok.Type == token.PLUS {
			interpreter.Eat(token.PLUS)
			res += interpreter.term()
		}
		if tok.Type == token.MINUS {
			interpreter.Eat(token.MINUS)
			res -= interpreter.term()
		}
	}
	return res
}
