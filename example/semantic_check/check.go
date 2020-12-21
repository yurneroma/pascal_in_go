package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"pascal_in_go/lexer"
	"pascal_in_go/parser"
	"pascal_in_go/types"
)

func main() {

	filename := os.Args[1]
	stream, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	text := string(stream)
	fmt.Println("origin code: ")
	fmt.Println("-------------------")
	fmt.Println(text)
	lexer := lexer.NewLexer(text)
	parser := parser.NewParser(lexer)
	symboltable := &types.SymbolTable{Symbols: make(map[string]types.Symbol), ErrorList: make([]error, 0)}
	symboltable.InitBuiltins()
	tree := parser.Program()
	symboltable.Visit(tree)
	errList := symboltable.ErrorList
	for _, err := range errList {
		fmt.Println("error:  ", err)
	}
}
