package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"pascal_in_go/interpreter"
	"pascal_in_go/lexer"
	"pascal_in_go/parser"
)

func main() {

	for {

		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("read error")
		}

		lexer := lexer.NewLexer(text)
		parser := parser.NewParser(lexer)
		inp := interpreter.NewInterpreter(parser)
		result := inp.Expr()
		fmt.Printf("%+v\n", result)
	}

}
