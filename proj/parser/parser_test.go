package parser

import (
	ct "proj/token"
	"testing"
)

func Test1(t *testing.T) {

	// The expected result struct represents the token stream for the input source
	tokens := []ct.Token{
		{ct.LET, "let"},
		{ct.ID, "a"},
		{ct.ASSIGN, "="},
		{ct.INT, "7"},
		{ct.PLUS, "+"},
		{ct.INT, "99"},
		{ct.SEMICOLON, ";"},
		{ct.EOF, "EOF"},
	}

	// Define  a new scanner for some Cal program
	parser := New(tokens)

	if got := parser.Parse(); !got {
		t.Errorf("\nParse(%v)\nExpected:%v\nGot:%v", tokens, true, got)
	}
}

func Test2(t *testing.T) {

	// The expected result struct represents the token stream for the input source
	tokens := []ct.Token{
		{ct.ID, "a"},
		{ct.ASSIGN, "="},
		{ct.INT, "7"},
		{ct.ASTERISK, "*"},
		{ct.INT, "99"},
		{ct.SEMICOLON, ";"},
		{ct.EOF, "EOF"},
	}

	// Define  a new scanner for some Cal program
	parser := New(tokens)

	if got := parser.Parse(); !got {
		t.Errorf("\nParse(%v)\nExpected:%v\nGot:%v", tokens, true, got)
	}
}

func Test3(t *testing.T) {

	// The expected result struct represents the token stream for the input source
	tokens := []ct.Token{
		{ct.PRINT, "print"},
		{ct.INT, "7"},
		{ct.ASTERISK, "*"},
		{ct.INT, "99"},
		{ct.SEMICOLON, ";"},
		{ct.EOF, "EOF"},
	}

	// Define  a new scanner for some Cal program
	parser := New(tokens)

	if got := parser.Parse(); !got {
		t.Errorf("\nParse(%v)\nExpected:%v\nGot:%v", tokens, true, got)
	}
}

func Test4(t *testing.T) {

	// The expected result struct represents the token stream for the input source
	tokens := []ct.Token{
		{ct.PRINT, "print"},
		{ct.INT, "7"},
		{ct.ASTERISK, "*"},
		{ct.INT, "99"},
		{ct.SEMICOLON, ";"},
	}

	// Define  a new scanner for some Cal program
	parser := New(tokens)

	if got := parser.Parse(); got {
		t.Errorf("\nParse(%v)\nExpected:%v\nGot:%v", tokens, false, got)
	}
}

func Test5(t *testing.T) {

	// The expected result struct represents the token stream for the input source
	tokens := []ct.Token{
		{ct.ID, "a"},
		{ct.ASSIGN, "="},
		{ct.PRINT, "print"},
		{ct.INT, "7"},
		{ct.ASTERISK, "*"},
		{ct.INT, "99"},
		{ct.SEMICOLON, ";"},
		{ct.EOF, "EOF"},
	}

	// Define  a new scanner for some Cal program
	parser := New(tokens)

	if got := parser.Parse(); got {
		t.Errorf("\nParse(%v)\nExpected:%v\nGot:%v", tokens, false, got)
	}
}

func Test6(t *testing.T) {

	// The expected result struct represents the token stream for the input source
	tokens := []ct.Token{
		{ct.ID, "a"},
		{ct.ASSIGN, "="},
		{ct.ID, "b"},
		{ct.SEMICOLON, ";"},
		{ct.EOF, "EOF"},
	}

	// Define  a new scanner for some Cal program
	parser := New(tokens)

	if got := parser.Parse(); !got {
		t.Errorf("\nParse(%v)\nExpected:%v\nGot:%v", tokens, true, got)
	}
}

func Test7(t *testing.T) {

	// The expected result struct represents the token stream for the input source
	tokens := []ct.Token{
		{ct.ID, "a"},
		{ct.ASSIGN, "="},
		{ct.ID, "b"},
		{ct.PLUS, "+"},
		{ct.INT, "3"},
		{ct.SLASH, "/"},
		{ct.INT, "4"},
		{ct.SEMICOLON, ";"},
		{ct.EOF, "EOF"},
	}

	// Define  a new scanner for some Cal program
	parser := New(tokens)

	if got := parser.Parse(); !got {
		t.Errorf("\nParse(%v)\nExpected:%v\nGot:%v", tokens, true, got)
	}
}
