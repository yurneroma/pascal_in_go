package interpreter

import (
	"log"
	"pascal_in_go/lexer"
	"pascal_in_go/token"
	"strconv"
)

//Interpreter represents the interpreter struct
type Interpreter struct {
	Lexer    lexer.Lexer `json:"lexer"`
	CurToken token.Token `json:"curToken"`
}

// NewInterpreter  init the Interpretetr
func NewInterpreter(lexer lexer.Lexer) *Interpreter {
	tok := lexer.NextToken()
	return &Interpreter{Lexer: lexer, CurToken: tok}
}

// eat function compare the current token type with the passed token
// type and if they match then "eat" the current token
// and assign the next token to the  interpreter's current_token,
// otherwise raise an exception.
func (interpreter *Interpreter) eat(tokenType token.Type) {
	if interpreter.CurToken.Type == tokenType {
		interpreter.CurToken = interpreter.Lexer.NextToken()
	} else {
		log.Fatal("type not match, cur and input  is ", interpreter.CurToken, ":", tokenType)
	}
}

func (interpreter *Interpreter) factor() (res int) {
	tok := interpreter.CurToken
	if tok.Type == token.INTEGER {
		interpreter.eat(token.INTEGER)
		res, _ = strconv.Atoi(tok.Literal)
		return
	}

	if tok.Type == token.LPAREN {
		interpreter.eat(token.LPAREN)
		res = interpreter.Expr()
		interpreter.eat(token.RPAREN)
		return
	}

	return
}

func (interpreter *Interpreter) term() (res int) {
	// context free grammar
	// term : factor ((MUL|DIV)factor)*
	res = interpreter.factor()
	for interpreter.CurToken.Type == token.DIV || interpreter.CurToken.Type == token.MUL {
		tok := interpreter.CurToken
		if tok.Type == token.MUL {
			interpreter.eat(token.MUL)
			res *= interpreter.factor()
		}
		if tok.Type == token.DIV {
			interpreter.eat(token.DIV)
			res /= interpreter.factor()
		}
	}
	return res
}

// Expr implements the arithmetic  expression
func (interpreter *Interpreter) Expr() (res int) {
	// context free grammar
	// calc > 1 + 9 * 2 - 6 / 3
	// expr :  term ((PLUS | MINUS) term )*
	// term :  factor ((MUL | DIV) factor )*
	// factor : INTEGER | Lparenthesized  expr  Rparenthesized

	res = interpreter.term()
	for interpreter.CurToken.Type == token.PLUS || interpreter.CurToken.Type == token.MINUS {
		tok := interpreter.CurToken
		if tok.Type == token.PLUS {
			interpreter.eat(token.PLUS)
			res += interpreter.term()
		}
		if tok.Type == token.MINUS {
			interpreter.eat(token.MINUS)
			res -= interpreter.term()
		}
	}
	return res
}
