package interpreter

import (
	"fmt"
	"pascal_in_go/ast"
	"pascal_in_go/parser"
	"pascal_in_go/token"
	"strconv"
)

//Interpreter represents the interpreter struct
type Interpreter struct {
	Parser *parser.Parser
}

func NewInterpreter(parser *parser.Parser) *Interpreter {
	return &Interpreter{Parser: parser}

}
func (inp *Interpreter) Expr() float64 {
	astTree := inp.Parser.AstBuild()
	fmt.Println("asttree: ", astTree.ToStr())
	ret := postOrder(astTree)
	return ret
}

func postOrder(astTree ast.Expr) float64 {
	if astTree == nil {
		return 0
	}

	switch t := astTree.(type) {
	case ast.BinNode:
		if t.Tok.Type == token.PLUS {
			return postOrder(t.Left) + postOrder(t.Right)
		}

		if t.Tok.Type == token.MINUS {
			return postOrder(t.Left) - postOrder(t.Right)
		}

		if t.Tok.Type == token.MUL {
			return postOrder(t.Left) * postOrder(t.Right)
		}

		if t.Tok.Type == token.DIV {
			temp := postOrder(t.Right)
			if temp != 0 {
				return postOrder(t.Left) / temp
			}
			return parser.INF
		}

	case ast.Unary:
		if t.Op == token.PLUS {
			return postOrder(t.Expr)
		}
		if t.Op == token.MINUS {
			return -postOrder(t.Expr)
		}

	case ast.NumNode:
		num, _ := strconv.ParseFloat(t.Tok.Literal, 64)
		return num

	default:
		fmt.Println("no match")
	}
	return 0
}
