package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"pascal_in_go/interpreter"
	"pascal_in_go/lexer"
)

func main() {

	for {

		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("read error")
		}

		lexer := lexer.NewLexer(text)
		interpreter := interpreter.NewInterpreter(lexer)
		result := interpreter.Expr()
		fmt.Printf("%+v\n", result)
	}

}
