package interpreter

import (
	"fmt"
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

type Node struct {
	left  *Node
	right *Node
	value string
}

func (interpreter *Interpreter) factor() *Node {
	tok := interpreter.CurToken
	if tok.Type == token.INTEGER {
		interpreter.eat(token.INTEGER)
		res := &Node{
			value: tok.Literal}
		return res
	}

	if tok.Type == token.LPAREN {
		interpreter.eat(token.LPAREN)
		res := interpreter.AstBuild()
		interpreter.eat(token.RPAREN)
		return res
	}

	return nil
}

func (interpreter *Interpreter) term() *Node {
	// context free grammar
	// term : factor ((MUL|DIV)factor)*
	left := interpreter.factor()
	for interpreter.CurToken.Type == token.DIV || interpreter.CurToken.Type == token.MUL {
		tok := interpreter.CurToken
		if tok.Type == token.MUL {
			interpreter.eat(token.MUL)
			rnode := interpreter.factor()
			left = &Node{left: left, right: rnode, value: token.MUL}
		}

		if tok.Type == token.DIV {
			interpreter.eat(token.DIV)
			rnode := interpreter.factor()
			left = &Node{left: left, right: rnode, value: token.DIV}
		}
	}
	return left
}

//AstBuild implements the ast tree
func (interpreter *Interpreter) AstBuild() *Node {
	// context free grammar
	// calc > 1 + 9 * 2 - 6 / 3
	// expr :  term ((PLUS | MINUS) term )*
	// term :  factor ((MUL | DIV) factor )*
	// factor : INTEGER | Lparenthesized  expr  Rparenthesized

	left := interpreter.term()
	for interpreter.CurToken.Type == token.PLUS || interpreter.CurToken.Type == token.MINUS {
		tok := interpreter.CurToken
		if tok.Type == token.PLUS {
			interpreter.eat(token.PLUS)
			rnode := interpreter.term()
			left = &Node{left: left, right: rnode, value: token.PLUS}
		}
		if tok.Type == token.MINUS {
			interpreter.eat(token.MINUS)
			rnode := interpreter.term()
			left = &Node{left: left, right: rnode, value: token.MINUS}
		}
	}

	fmt.Printf("tree is %+v\n", left)
	return left
}

func (interpreter *Interpreter) Expr() float64 {
	ast := interpreter.AstBuild()
	ret := postOrder(ast)
	return ret
}

func postOrder(ast *Node) float64 {
	if ast == nil {
		return 0
	}

	t := ast.value
	num, err := strconv.ParseFloat(t, 64)
	if err == nil {
		return num
	}
	if t == token.PLUS {
		return postOrder(ast.left) + postOrder(ast.right)
	}
	if t == token.MINUS {
		return postOrder(ast.left) - postOrder(ast.right)
	}

	if t == token.MUL {
		return postOrder(ast.left) * postOrder(ast.right)
	}

	if t == token.DIV {

		if postOrder(ast.right) != 0 {
			return postOrder(ast.left) / postOrder(ast.right)
		}

	}

	return 0
}
