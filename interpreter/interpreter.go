package interpreter

import (
	"fmt"
	"log"
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
	log.Printf("tree is %+v\n", astTree)
	log.Println("-------------------")
	inp.visit(astTree)
	return inp.VarMap
}

func (inp *Interpreter) visit(astTree ast.Expr) float64 {
	if astTree == nil {
		return 0
	}

	switch t := astTree.(type) {
	case ast.Program:
		inp.visitProgram(t)
	case ast.Block:
		inp.visitBlock(t)
	case ast.Decl:
		inp.visitDecl(t)
	case ast.Compound:
		inp.visitCompound(t)
	case ast.AssignStatement:
		inp.visitAssignment(t)
	case ast.Statement:
		inp.visitStatement(t)
	case ast.BinNode:
		res := inp.visitBinNode(t)
		return res
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
		return inp.visitVar(t)
	default:
		fmt.Println("no match", t)
	}
	return 0
}

func (inp *Interpreter) visitBinNode(t ast.BinNode) float64 {
	if t.Tok.Type == token.PLUS {
		left := inp.visit(t.Left)
		right := inp.visit(t.Right)
		return left + right
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

	return 0
}

func (inp *Interpreter) visitProgram(t ast.Program) {
	inp.visitBlock(t.Block)
}

func (inp *Interpreter) visitBlock(t ast.Block) {
	for _, decl := range t.Decls {
		inp.visitDecl(decl)
	}
	inp.visitCompound(t.Compound)
}

func (inp *Interpreter) visitDecl(t ast.Decl) {}

func (inp *Interpreter) visitCompound(t ast.Compound) {
	for _, child := range t.Children {
		switch node := child.(type) {
		case ast.Statement:
			inp.visitStatement(node)
		case ast.NoOp:
			return
		}
	}
}

func (inp *Interpreter) visitStatement(t ast.Statement) {
	switch node := t.Statement.(type) {
	case ast.AssignStatement:
		inp.visitAssignment(node)
	case ast.Compound:
		inp.visitCompound(node)
	case ast.NoOp:
		return
	}
}
func (inp *Interpreter) visitAssignment(st ast.AssignStatement) {
	varName := st.Left.Literal
	rValue := inp.visit(st.Right)
	inp.VarMap[varName] = rValue
}

func (inp *Interpreter) visitVar(node ast.VarNode) float64 {
	name := node.Literal
	value, ok := inp.VarMap[name]
	if !ok {
		return 0
	}
	switch t := value.(type) {
	case string:
		val, _ := strconv.ParseFloat(t, 64)
		return val
	case float64:
		return t
	default:
		return 0
	}

}
