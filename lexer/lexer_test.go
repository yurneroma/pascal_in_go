package lexer

import "testing"

func TestIsalpha(t *testing.T) {
	var (
		text   = "abc"
		expect = true
	)
	lexer := NewLexer(text)
	flag := lexer.isalpha()
	if flag != expect {
		t.Errorf("flag is %+v; expected  %+v, text is %s\n ", flag, expect, text)
	}
}

func TestIsnum(t *testing.T) {
	var (
		text   = "3456"
		expect = true
	)
	lexer := NewLexer(text)
	flag := lexer.isnum()
	if flag != expect {
		t.Errorf("flag is %+v; expected  %+v, text is %s\n ", flag, expect, text)
	}
}
