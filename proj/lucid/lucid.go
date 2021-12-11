package lucid

import (
	"flag"
	"fmt"
	cc "proj/context"
	"proj/parser"
	"proj/sa"
	"proj/scanner"
)

func Start(ctx *cc.CompilerContext) {
	scanner := scanner.New(ctx)
	parser := parser.New(ctx, scanner)
	ast := parser.Parse()
	sa.PerformSA(ast)
	fmt.Println(ast)
}

func main() {
	/*Parse args and flags*/
	lexPtr := flag.Bool("lex", false, "Use -lex fileName to scan the tokens in the .golite file")
	flag.Parse()
	argsWithoutProg := flag.Args()
	inputFileName := argsWithoutProg[0]
	ctx := cc.New(*lexPtr, inputFileName)
	/*Print the tokens of the input file if in lex mode*/
	if *lexPtr {
		scanner := scanner.New(ctx)
		scanner.PrintAllTokens()
	}

}
