package scanner

import (
	"hw2/token"
	"testing"
)

type ExpectedResult struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func VerifyTest(t *testing.T, tts []ExpectedResult, scanner *Scanner) {

	for i, tt := range tts {
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("FAILED[%d] - incorrect token.\nexpected=%v\ngot=%v\n",
				i, tt.expectedType, tok.Type)

			return
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("FAILED[%d] - incorrect token literal.\nexpected=%v\ngot=%v\n",
				i, tt.expectedLiteral, tok.Literal)

		}
	}
}

func Test1(t *testing.T) {

	// This is a raw string in Go (aka its a multiline string). This will be easy to
	input1 := `var five = 5;
print five; 
`
	// The expected result struct represents the token stream for the input source
	expected := []ExpectedResult{
		{token.VAR, "var"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.PRINT, "print"},
		{token.IDENT, "five"},
		{token.SEMICOLON, ";"},
		{token.EOF, "EOF"},
	}

	// Define  a new scanner for some Cal program
	scanner := New(input1)

	// Verify that the scanner produces the tokens in the order that you expect.
	VerifyTest(t, expected, scanner)
}

func Test2(t *testing.T) {

	// This is a raw string in Go (aka its a multiline string). This will be easy to
	input1 := `if x> 3:
	return x;
`
	// The expected result struct represents the token stream for the input source
	expected := []ExpectedResult{
		{token.IF, "if"},
		{token.IDENT, "x"},
		{token.GREATER, ">"},
		{token.INT, "3"},
		{token.COLON, ":"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.SEMICOLON, ";"},
		{token.EOF, "EOF"},
	}

	// Define  a new scanner for some Cal program
	scanner := New(input1)

	// Verify that the scanner produces the tokens in the order that you expect.
	VerifyTest(t, expected, scanner)
}

func Test3(t *testing.T) {

	// This is a raw string in Go (aka its a multiline string). This will be easy to
	input1 := `while a1<= 3:
	a=a+1;
`
	// The expected result struct represents the token stream for the input source
	expected := []ExpectedResult{
		{token.WHILE, "while"},
		{token.IDENT, "a1"},
		{token.LESSEQU, "<="},
		{token.INT, "3"},
		{token.COLON, ":"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.EOF, "EOF"},
	}

	// Define  a new scanner for some Cal program
	scanner := New(input1)

	// Verify that the scanner produces the tokens in the order that you expect.
	VerifyTest(t, expected, scanner)
}
