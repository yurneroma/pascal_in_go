# pascal_in_go
implement pascal  with golang  


## Grammar

program : Compound_statement DOT

compound_statement :  BEGIN   statement_list  END

statement_list : statement | statement SEMI  statement_list

statement :  compound_statement | assignment  | empty

assignment :  variable  ASSIGN expr

expr : term ((PLUS |  MINUS) term )*

term : factor ((MUL | DIV) factor )*

factor :  (PLUS | MINUS) factor  | INTEGER | Lparenthesized expr Rparenthesized | variable

variable :  ID


## Example