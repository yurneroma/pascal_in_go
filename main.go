package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"pascal_in_go/interpreter"
	"pascal_in_go/lexer"
	"pascal_in_go/parser"
)

func main() {

	filename := os.Args[1]
	stream, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	text := string(stream)
	fmt.Println(text)
	lexer := lexer.NewLexer(text)
	parser := parser.NewParser(lexer)
	inp := interpreter.NewInterpreter(parser)
	result := inp.Expr()
	fmt.Printf("%+v\n", result)

}
