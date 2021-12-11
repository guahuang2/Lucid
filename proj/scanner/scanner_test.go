package scanner

import (
	"proj/context"
	"proj/token"
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
			return
		}

	}
}

func Test1(t *testing.T) {

	// This is a raw string in Go (aka its a multiline string). This will be easy to
	ctx := context.New(false, "../golite/test1.golite")
	// The expected result struct represents the token stream for the input source
	expected := []ExpectedResult{
		{token.PACKAGE, "package"},
		{token.IDENT, "main"},
		{token.SEMICOLON, ";"},
		{token.IMPORT, "import"},
		{token.DOUQUAT, "\""},
		{token.FMT, "fmt"},
		{token.DOUQUAT, "\""},
		{token.SEMICOLON, ";"},
		{token.FUNC, "func"},
		{token.IDENT, "main"},
		{token.LEFTPAR, "("},
		{token.RIGHTPAR, ")"},
		{token.LEFTBRAC, "{"},
		{token.VAR, "var"},
		{token.IDENT, "a"},
		{token.INT, "int"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.NUMBER, "3"},
		{token.PLUS, "+"},
		{token.NUMBER, "4"},
		{token.PLUS, "+"},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.FMT, "fmt"},
		{token.DOT, "."},
		{token.PRINT, "Print"},
		{token.LEFTPAR, "("},
		{token.IDENT, "a"},
		{token.RIGHTPAR, ")"},
		{token.SEMICOLON, ";"},
		{token.RIGHTBRAC, "}"},
	}

	// Define  a new scanner for some Cal program
	scanner := New(ctx)

	// Verify that the scanner produces the tokens in the order that you expect.
	VerifyTest(t, expected, scanner)
}
