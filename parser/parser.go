package parser

import (
	"log"
	"pascal_in_go/ast"
	"pascal_in_go/lexer"
	"pascal_in_go/token"
)

//INF represents the infinity
const INF = 0x3fffffff

//Parser struct
type Parser struct {
	Lexer    lexer.Lexer `json:"lexer"`
	CurToken token.Token `json:"curToken"`
}

// NewParser  init the parser
func NewParser(lexer lexer.Lexer) *Parser {
	tok := lexer.NextToken()
	return &Parser{Lexer: lexer, CurToken: tok}
}

// eat function compare the current token type with the passed token
// type and if they match then "eat" the current token
// and assign the next token to the  parser's current_token,
// otherwise raise an exception.
func (parser *Parser) eat(tokenType token.Type) {
	if parser.CurToken.Type == tokenType {
		parser.CurToken = parser.Lexer.NextToken()
	} else {
		log.Fatal("type not match, cur and input  is ", parser.CurToken, ":", tokenType)
	}
}

func (parser *Parser) factor() ast.Expr {
	tok := parser.CurToken
	if tok.Type == token.INTEGER {
		parser.eat(token.INTEGER)
		res := ast.NumNode{
			Tok:   tok,
			Value: tok.Literal,
		}
		return res
	}

	if tok.Type == token.LPAREN {
		parser.eat(token.LPAREN)
		res := parser.AstBuild()
		parser.eat(token.RPAREN)
		return res
	}

	if tok.Type == token.MINUS {
		parser.eat(token.MINUS)
		expr := parser.factor()
		res := ast.Unary{
			Op:   token.MINUS,
			Expr: expr}

		return res
	}

	if tok.Type == token.PLUS {
		parser.eat(token.PLUS)
		expr := parser.factor()
		res := ast.Unary{
			Op:   token.PLUS,
			Expr: expr,
		}

		return res
	}
	return nil
}

func (parser *Parser) term() ast.Expr {
	// context free grammar
	// calc > 1 + 9 * 2 - 6 / 3
	// term : factor ((MUL|DIV)factor)*
	left := parser.factor()
	for parser.CurToken.Type == token.DIV || parser.CurToken.Type == token.MUL {
		tok := parser.CurToken
		if tok.Type == token.MUL {
			parser.eat(token.MUL)
			rnode := parser.factor()
			left = ast.BinNode{Left: left, Right: rnode, Tok: tok}
		}

		if tok.Type == token.DIV {
			parser.eat(token.DIV)
			rnode := parser.factor()
			left = ast.BinNode{Left: left, Right: rnode, Tok: tok}
		}
	}
	return left
}

//AstBuild implements the ast tree
func (parser *Parser) AstBuild() ast.Expr {
	/* context free grammar
	program : Compound_statement DOT

	compound_statement :  START   statement_list  END

	statement_list : statement | statement SEMI  statement_list

	statement :  compound_statement | assignment  | empty

	assignment :  variable  ASSIGN expr

	expr : term ((PLUS |  MINUS) term )*

	term : factor ((MUL | DIV) factor )*

	factor :  (PLUS | MINUS) factor  | INTEGER | Lparenthesized expr Rparenthesized | variable

	variable :  ID
	*/

	left := parser.term()
	for parser.CurToken.Type == token.PLUS || parser.CurToken.Type == token.MINUS {
		tok := parser.CurToken
		if tok.Type == token.PLUS {
			parser.eat(token.PLUS)
			rnode := parser.term()
			left = ast.BinNode{Left: left, Right: rnode, Tok: tok}
		}
		if tok.Type == token.MINUS {
			parser.eat(token.MINUS)
			rnode := parser.term()
			left = ast.BinNode{Left: left, Right: rnode, Tok: tok}
		}
	}

	return left
}

func (parser *Parser) variable() ast.Expr {
	tok := parser.CurToken
	if tok.Type == token.ID {
		parser.eat(token.ID)
		res := ast.VarNode{
			Tok:   tok,
			Value: tok.Literal}
		return res
	}
	return nil
}
