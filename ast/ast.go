package ast

import (
	"fmt"
	"pascal_in_go/token"
)

//BinNode represents the binary expr
type BinNode struct {
	Left  Expr
	Right Expr
	Tok   token.Token
}

//NumNode holds the number of token
type NumNode struct {
	Tok   token.Token `json:"tok"`
	Value string      `json:"value"`
}

//VarNode represents the Variable node
type VarNode struct {
	Tok   token.Token
	Value string
}

//Expr interface represent the expr,  expr is the unit for program
type Expr interface {
	ToStr() string
}

//Unary Node
type Unary struct {
	Op   string
	Expr Expr
}

//Statement Node
type Statement struct {
	Expr Expr
}

//ToStr for BinNode
func (binNode BinNode) ToStr() string {
	left := fmt.Sprint(binNode.Left)
	right := fmt.Sprint(binNode.Right)
	op := fmt.Sprint(binNode.Tok)
	return fmt.Sprint(left, right, op)
}

//ToStr for NumNode
func (numNode NumNode) ToStr() string {
	return fmt.Sprint(numNode.Tok)
}

//ToStr for UnaryNode
func (unary Unary) ToStr() string {
	return fmt.Sprint(unary)
}

//ToStr for VarNode
func (va VarNode) ToStr() string {
	return fmt.Sprint(va)
}

//ToStr for statement
func (st Statement) ToStr() string {
	return fmt.Sprint(st)
}

type StatementList []Statement

func (stList StatementList) ToStr() string {
	return fmt.Sprint(stList)
}

type AssignStatement struct {
	Left token.ID
}
