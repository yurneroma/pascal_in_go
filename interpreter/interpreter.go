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
	VarMap map[string]interface{}
}

func NewInterpreter(parser *parser.Parser) *Interpreter {
	return &Interpreter{Parser: parser, VarMap: make(map[string]interface{})}
}
func (inp *Interpreter) Expr() map[string]interface{} {
	astTree := inp.Parser.Program()
	inp.visit(astTree)
	return inp.VarMap
}

func (inp *Interpreter) visit(astTree ast.Expr) float64 {
	if astTree == nil {
		return 0
	}

	switch t := astTree.(type) {
	case ast.BinNode:
		if t.Tok.Type == token.PLUS {
			return inp.visit(t.Left) + inp.visit(t.Right)
		}

		if t.Tok.Type == token.MINUS {
			return inp.visit(t.Left) - inp.visit(t.Right)
		}

		if t.Tok.Type == token.MUL {
			return inp.visit(t.Left) * inp.visit(t.Right)
		}

		if t.Tok.Type == token.DIV {
			temp := inp.visit(t.Right)
			if temp != 0 {
				return inp.visit(t.Left) / temp
			}
			return parser.INF
		}

	case ast.Unary:
		if t.Op == token.PLUS {
			return inp.visit(t.Expr)
		}
		if t.Op == token.MINUS {
			return -inp.visit(t.Expr)
		}

	case ast.NumNode:
		num, _ := strconv.ParseFloat(t.Tok.Literal, 64)
		return num

	case ast.VarNode:
		inp.visitVar(t)
	case ast.Compound:
		inp.visitCompound(t)
	case ast.AssignStatement:
		inp.visitAssignment(t)
	default:
		fmt.Println("no match", t)
	}
	return 0
}

func (inp *Interpreter) visitCompound(t ast.Compound) {
	for _, child := range t.Children {
		switch node := child.(type) {
		case ast.AssignStatement:
			inp.visitAssignment(node)
		case ast.NoOp:
			return
		}
	}
}

func (inp *Interpreter) visitAssignment(st ast.AssignStatement) {
	varName := st.Left.Value
	rValue := inp.visit(st.Right)
	inp.VarMap[varName] = rValue
}

func (inp *Interpreter) visitVar(node ast.VarNode) float64 {
	name := node.Value
	value, ok := inp.VarMap[name]
	if !ok {
		return 0
	}
	val, _ := strconv.ParseFloat(value.(string), 64)
	return val
}
