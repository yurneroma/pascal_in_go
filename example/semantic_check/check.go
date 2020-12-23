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
	fmt.Println("Origin Code: ")
	fmt.Println("-------------------")
	fmt.Println(text)
	fmt.Println("-------------------")
	fmt.Println("Semantic Analyzing: ")
	lexer := lexer.NewLexer(text)
	parser := parser.NewParser(lexer)
	symboltable := &types.SymbolTable{Symbols: make(map[string]types.Symbol), ErrorList: make([]error, 0)}
	symboltable.InitBuiltins()
	tree := parser.Program()
	symboltable.Visit(tree)
	fmt.Println("-------------------")
	fmt.Println("Error Reporting : ")
	errList := symboltable.ErrorList
	for _, err := range errList {
		fmt.Println("error:  ", err)
	}
}
