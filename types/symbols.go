package types

import (
	"errors"
	"fmt"
	"pascal_in_go/ast"
)

type Symbol interface {
	ShowName() string
	ShowType() string
}

type BuiltinTypeSymbol struct {
	Name string
	Type string
}
type VarSymbol struct {
	Name string
	Type string
}

func (bts BuiltinTypeSymbol) ShowName() string {
	return bts.Name
}

func (bts BuiltinTypeSymbol) ShowType() string {
	return bts.Type
}

func (vs VarSymbol) ShowName() string {
	return vs.Name
}

func (vs VarSymbol) ShowType() string {
	return vs.Type
}

type SymbolTable struct {
	Symbols   map[string]Symbol
	ErrorList []error
}

func (symtab *SymbolTable) define(symbol Symbol) {
	switch t := symbol.(type) {
	case BuiltinTypeSymbol:
		fmt.Printf("BuiltinSymbol : %+v\n", t)
	case VarSymbol:
		fmt.Printf("VarSymbol : %+v\n", t)

	}
	name := symbol.ShowName()
	symtab.Symbols[name] = symbol
}

func (symtab *SymbolTable) lookup(name string) Symbol {
	fmt.Println("lookup: ", name)
	symbol := symtab.Symbols[name]
	return symbol
}

func (symtab *SymbolTable) InitBuiltins() {
	symtab.define(BuiltinTypeSymbol{Type: "INTEGER"})
	symtab.define(BuiltinTypeSymbol{Type: "REAL"})
}

func (symtab *SymbolTable) visitProgram(t ast.Program) {
	symtab.visitBlock(t.Block)
}

func (symtab *SymbolTable) visitBlock(t ast.Block) {
	for _, vardecl := range t.Decl.VarDeclList {
		symtab.visitVarDecl(vardecl)
	}

	// for _, procedureDecl := range t.Decl.ProceDeclList {
	// 	//todo
	// }
	symtab.visitCompound(t.Compound)
}

func (symtab *SymbolTable) visitVarDecl(t ast.VarDecl) {
	typeName := string(t.Type)
	varName := t.Node.Literal
	varSymbol := VarSymbol{Type: typeName, Name: varName}
	symbol := symtab.lookup(varName)
	if symbol != nil {
		msg := fmt.Sprintf("Duplicate  identifier %s", varName)
		err := errors.New(msg)
		symtab.addError(err)
		return
	}
	symtab.define(varSymbol)
}

func (symtab *SymbolTable) visitCompound(t ast.Compound) {
	for _, child := range t.Children {
		switch node := child.(type) {
		case ast.Statement:
			symtab.visitStatement(node)
		case ast.NoOp:
			return
		}
	}
}

func (symtab *SymbolTable) visitStatement(t ast.Statement) {
	switch node := t.Statement.(type) {
	case ast.AssignStatement:
		symtab.visitAssignment(node)
	case ast.Compound:
		symtab.visitCompound(node)
	case ast.NoOp:
		return
	}
}

func (symtab *SymbolTable) visitAssignment(st ast.AssignStatement) {
	varName := st.Left.Literal
	res := symtab.lookup(varName)
	if res == nil {
		msg := fmt.Sprintf("varname %s undeclared", varName)
		err := errors.New(msg)
		symtab.addError(err)
		return
	}
	symtab.Visit(st.Right)
}

func (symtab *SymbolTable) addError(err error) {
	errList := symtab.ErrorList
	errList = append(errList, err)
	symtab.ErrorList = errList
}

func (symtab *SymbolTable) Visit(astTree ast.Expr) {
	if astTree == nil {
		return
	}

	switch t := astTree.(type) {
	case ast.Program:
		symtab.visitProgram(t)
	case ast.Block:
		symtab.visitBlock(t)
	case ast.Decl:
		symtab.visitDecl(t)
	case ast.Compound:
		symtab.visitCompound(t)
	case ast.AssignStatement:
		symtab.visitAssignment(t)
	case ast.Statement:
		symtab.visitStatement(t)
	case ast.BinNode:
		symtab.visitBinNode(t)
	case ast.Unary:
		symtab.Visit(t.Expr)

	case ast.NumNode:

	case ast.VarNode:
		symtab.visitVar(t)
	case ast.Procedure:
		symtab.visitProcedure(t)

	default:
		fmt.Println("no match", t)
	}
}

func (symtab *SymbolTable) visitBinNode(t ast.BinNode) {
	symtab.Visit(t.Left)
	symtab.Visit(t.Right)
}

func (symtab *SymbolTable) visitVar(t ast.VarNode) {
	name := t.Literal
	symbol := symtab.lookup(name)
	if symbol == nil {
		msg := fmt.Sprintf("varname %s undeclared", name)
		err := errors.New(msg)
		symtab.addError(err)
		return
	}
}

func (symtab *SymbolTable) visitProcedure(t ast.Procedure) {
	return
}

func (symtab *SymbolTable) visitDecl(t ast.Decl) {

}
