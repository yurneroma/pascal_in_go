package parser

import (
	"log"
	"pascal_in_go/ast"
	"pascal_in_go/lexer"
	"pascal_in_go/token"
)

/* context free grammar
program : PROGRAM Variable SEMI block DOT

block : declarations compound_statement

declarations :  VAR(variable_declaration SEMI)+  | empty

variable_declaration : ID(COMMA ID)* COLON type_spec

type_spec : INTEGER | REAL

compound_statement :  BEGIN   statement_list  END

statement_list : statement | statement SEMI  statement_list

statement :  compound_statement | assignment  | empty

assignment :  variable  ASSIGN expr

expr : term ((PLUS |  MINUS) term )*

term : factor ((MUL | INTEGER_DIV | FLOAT_DIV) factor )*

factor :  PLUS factor
		| MINUS factor
		| REAL_CONST
		| INTEGER_CONST
		| Lparenthesized expr Rparenthesized
		| variable

variable :  ID
*/

//INF represents the infinity
const INF = 0x3fffffff

//Parser struct
type Parser struct {
	Lexer    lexer.Lexer `json:"lexer"`
	CurToken token.Token `json:"curToken"`
}

// NewParser  init the parser
func NewParser(lexer lexer.Lexer) *Parser {
	return &Parser{Lexer: lexer, CurToken: lexer.NextToken()}
}

func (parser *Parser) Program() ast.Expr {
	/*
		program : PROGRAM Variable SEMI block DOT
	*/
	parser.eat(token.PROGRAM)
	varNode := parser.variable()
	name := varNode.ToStr()
	parser.eat(token.SEMI)
	block := parser.block()
	parser.eat(token.DOT)
	return ast.Program{Block: block, Name: name}
}

func (parser *Parser) block() ast.Block {
	//block : declarations compound_statement
	declarations := parser.declarations()
	compound := parser.comStatement()
	return ast.Block{Decl: declarations, Compound: compound}
}

func (parser *Parser) declarations() ast.Decl {
	/*
		declarations : VAR (variable_declaration SEMI)+
						| (PROCEDURE ID SEMI block SEMI)*
		               | empty
	*/
	decls := ast.Decl{}
	vardecls := make([]ast.VarDecl, 0)
	if parser.CurToken.Type == token.VAR {
		parser.eat(token.VAR)
		for parser.CurToken.Type == token.ID {
			varDecls := parser.varDecl()
			vardecls = append(vardecls, varDecls...)
			parser.eat(token.SEMI)
		}
	}
	decls.VarDeclList = vardecls

	procedureList := make([]ast.Procedure, 0)
	for parser.CurToken.Type == token.PROCEDURE {
		parser.eat(token.PROCEDURE)
		procedName := parser.CurToken.Literal
		parser.eat(token.ID)
		parser.eat(token.SEMI)
		block := parser.block()
		procedure := ast.Procedure{Name: procedName, Block: block}
		procedureList = append(procedureList, procedure)
		parser.eat(token.SEMI)

	}
	decls.ProceDeclList = procedureList
	return decls
}

func (parser *Parser) varDecl() []ast.VarDecl {
	/*
		variable_declaration:  ID(COMMA ID)*  COLON type_spec
	*/
	varNodes := make([]ast.VarNode, 0)
	tok := parser.CurToken
	varNodes = append(varNodes, ast.VarNode{tok, tok.Literal})
	parser.eat(token.ID)

	for parser.CurToken.Type == token.COMMA {
		parser.eat(token.COMMA)
		tok = parser.CurToken
		varNodes = append(varNodes, ast.VarNode{tok, tok.Literal})
		parser.eat(token.ID)
	}

	parser.eat(token.COLON)
	typeSpec := parser.typeSpec()
	decls := make([]ast.VarDecl, 0)
	for _, elem := range varNodes {
		decls = append(decls, ast.VarDecl{Node: elem, Type: typeSpec})
	}
	return decls
}
func (parser *Parser) typeSpec() token.Type {
	/*
		type_spec : INTEGER
					| REAL
	*/

	curType := parser.CurToken.Type
	if parser.CurToken.Type == token.REAL {
		parser.eat(token.REAL)

	}
	if parser.CurToken.Type == token.INTEGER {
		parser.eat(token.INTEGER)
	}
	return curType

}

func (parser *Parser) comStatement() ast.Compound {
	/*
		compound_statement: BEGIN statement_list END
	*/
	parser.eat(token.BEGIN)
	comStatement := parser.statementList()
	parser.eat(token.END)

	root := ast.Compound{}
	for _, st := range comStatement {
		root.Children = append(root.Children, st)
	}
	return root
}

func (parser *Parser) statementList() []ast.Expr {
	/*
		statement_list : statement
						 | statement SEMI  statement_list
	*/

	stList := make([]ast.Expr, 0)
	st := parser.statement()
	stList = append(stList, st)
	for parser.CurToken.Type == token.SEMI {
		parser.eat(token.SEMI)
		res := parser.statement()
		stList = append(stList, res)
	}

	return stList
}

func (parser *Parser) statement() ast.Expr {
	/*
	    statement : compound_statement
	   				| assignment_statement
	   		 		| empty
	*/
	var st ast.Statement
	if parser.CurToken.Type == token.BEGIN {
		st.Statement = parser.comStatement()
	} else if parser.CurToken.Type == token.ID {
		st.Statement = parser.assignmentStatement()
	} else {
		st.Statement = parser.empty()
	}
	return st
}

//implements assignmentStatement
func (parser *Parser) assignmentStatement() ast.Expr {
	left := parser.variable()
	op := parser.CurToken
	parser.eat(token.ASSIGN)
	right := parser.expr()
	return ast.AssignStatement{
		Left:  left.(ast.VarNode),
		Op:    op,
		Right: right,
	}
}

func (parser *Parser) empty() ast.Expr {
	// tok := parser.CurToken
	// for tok.Type ==
	return ast.NoOp{}
}

//expr
func (parser *Parser) expr() ast.Expr {
	/*
		expr:  term((PLUS|MINUS)term)*
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

// eat function compare the current token type with the passed token
// type and if they match then "eat" the current token
// and assign the next token to the  parser's current_token,
// otherwise raise an exception.
func (parser *Parser) eat(tokenType token.Type) {
	if parser.CurToken.Type == tokenType {
		parser.CurToken = parser.Lexer.NextToken()
	} else {
		log.Fatalf("type not match, cur is %+v and input  is %+v , position is %+v, curchar is %+v",
			parser.CurToken, tokenType, parser.Lexer.Pos, string(parser.Lexer.CurChar))
	}
}

func (parser *Parser) factor() ast.Expr {
	/*
		factor :  PLUS factor
				| MINUS factor
				| REAL
				| INTEGER
				| Lparenthesized expr Rparenthesized
				| variable

	*/
	tok := parser.CurToken
	if tok.Type == token.INTEGER {
		parser.eat(token.INTEGER)
		res := ast.NumNode{
			Tok:   tok,
			Value: tok.Literal,
		}
		return res
	}

	if tok.Type == token.REAL {
		parser.eat(token.REAL)
		res := ast.NumNode{
			Tok:   tok,
			Value: tok.Literal,
		}
		return res
	}

	if tok.Type == token.LPAREN {
		parser.eat(token.LPAREN)
		res := parser.expr()
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
	if tok.Type == token.ID {
		res := parser.variable()
		return res
	}
	return nil
}

func isInSlice(a token.Type, list []token.Type) bool {
	for _, b := range list {
		if a == b {
			return true
		}
	}
	return false
}
func (parser *Parser) term() ast.Expr {
	// context free grammar
	// term : factor ((MUL | DIV)factor)*
	left := parser.factor()
	for parser.CurToken.Type == token.MUL || parser.CurToken.Type == token.DIV {
		tok := parser.CurToken
		if parser.CurToken.Type == token.DIV {
			parser.eat(token.DIV)
			right := parser.factor()
			left = ast.BinNode{Left: left, Right: right, Tok: tok}
		}
		if parser.CurToken.Type == token.MUL {
			parser.eat(token.MUL)
			right := parser.factor()
			left = ast.BinNode{Left: left, Right: right, Tok: tok}
		}
	}
	return left
}

func (parser *Parser) variable() ast.Expr {
	tok := parser.CurToken
	if tok.Type == token.ID {
		parser.eat(token.ID)
		res := ast.VarNode{
			Tok:     tok,
			Literal: tok.Literal}
		return res
	}
	return nil
}
