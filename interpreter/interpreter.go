package interpreter

import (
	"fmt"
	"log"
	"pascal_in_go/lexer"
	"pascal_in_go/token"
	"strconv"
)

const INF = 0x3fffffff

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

type BinNode struct {
	left  Expr
	right Expr
	tok   token.Token
}

type NumNode struct {
	tok   token.Token
	value string
}

type Expr interface {
	ToStr() string
}

type Unary struct {
	Op   string
	expr Expr
}

func (binNode BinNode) ToStr() string {
	left := fmt.Sprint(binNode.left)
	right := fmt.Sprint(binNode.right)
	op := fmt.Sprint(binNode.tok)
	return fmt.Sprint(left, right, op)
}

func (numNode NumNode) ToStr() string {
	return fmt.Sprint(numNode.tok)
}

func (unary Unary) ToStr() string {
	return fmt.Sprint(unary)
}
func (interpreter *Interpreter) factor() Expr {
	tok := interpreter.CurToken
	if tok.Type == token.INTEGER {
		interpreter.eat(token.INTEGER)
		res := NumNode{
			tok:   tok,
			value: tok.Literal}
		return res
	}

	if tok.Type == token.LPAREN {
		interpreter.eat(token.LPAREN)
		res := interpreter.AstBuild()
		interpreter.eat(token.RPAREN)
		return res
	}

	if tok.Type == token.MINUS {
		interpreter.eat(token.MINUS)
		expr := interpreter.factor()
		res := Unary{
			Op:   token.MINUS,
			expr: expr}

		return res
	}

	if tok.Type == token.PLUS {
		interpreter.eat(token.PLUS)
		expr := interpreter.factor()
		res := Unary{
			Op:   token.PLUS,
			expr: expr}

		return res
	}
	return nil
}

func (interpreter *Interpreter) term() Expr {
	// context free grammar
	// calc > 1 + 9 * 2 - 6 / 3
	// term : factor ((MUL|DIV)factor)*
	left := interpreter.factor()
	for interpreter.CurToken.Type == token.DIV || interpreter.CurToken.Type == token.MUL {
		tok := interpreter.CurToken
		if tok.Type == token.MUL {
			interpreter.eat(token.MUL)
			rnode := interpreter.factor()
			left = BinNode{left: left, right: rnode, tok: tok}
		}

		if tok.Type == token.DIV {
			interpreter.eat(token.DIV)
			rnode := interpreter.factor()
			left = BinNode{left: left, right: rnode, tok: tok}
		}
	}
	return left
}

//AstBuild implements the ast tree
func (interpreter *Interpreter) AstBuild() Expr {
	/* context free grammar
	program : Compound_statement DOT

	compound_statement :  START   statement_list  END

	statement_list : statement | statement SEMI  statement_list

	statement :  compound_statement | assignment  | empty

	assignment :  variable  ASSIGN expr

	expr : term ((PLUS |  MINUS) term )*

	term : factor ((MUL | DIV) factor )*

	factor :  (PLUS | MINUS) factor  | INTEGER | Lparenthesized expr Rparenthesized | variable

	variable :  ID
	*/

	left := interpreter.term()
	for interpreter.CurToken.Type == token.PLUS || interpreter.CurToken.Type == token.MINUS {
		tok := interpreter.CurToken
		if tok.Type == token.PLUS {
			interpreter.eat(token.PLUS)
			rnode := interpreter.term()
			left = BinNode{left: left, right: rnode, tok: tok}
		}
		if tok.Type == token.MINUS {
			interpreter.eat(token.MINUS)
			rnode := interpreter.term()
			left = BinNode{left: left, right: rnode, tok: tok}
		}
	}

	return left
}

func (interpreter *Interpreter) Expr() float64 {
	ast := interpreter.AstBuild()
	ret := postOrder(ast)
	return ret
}

func postOrder(ast Expr) float64 {
	if ast == nil {
		return 0
	}

	switch t := ast.(type) {
	case BinNode:
		if t.tok.Type == token.PLUS {
			return postOrder(t.left) + postOrder(t.right)
		}

		if t.tok.Type == token.MINUS {
			return postOrder(t.left) - postOrder(t.right)
		}

		if t.tok.Type == token.MUL {
			return postOrder(t.left) * postOrder(t.right)
		}

		if t.tok.Type == token.DIV {
			temp := postOrder(t.right)
			if temp != 0 {
				return postOrder(t.left) / temp
			}
			return INF
		}

	case Unary:
		if t.Op == token.PLUS {
			return postOrder(t.expr)
		}
		if t.Op == token.MINUS {
			return -postOrder(t.expr)
		}

	case NumNode:
		num, _ := strconv.ParseFloat(t.tok.Literal, 64)
		return num

	default:
		fmt.Println("no match")
	}
	return 0
}
