package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"proj/assembly"
	cc "proj/context"
	"proj/ir"
	"proj/parser"
	"proj/sa"
	"proj/scanner"
	"strings"
)

func StartCompiling(ctx *cc.CompilerContext) {
	scanner := scanner.New(ctx)
	parser := parser.New(ctx, scanner)
	fmt.Println("Start parsing")
	ast := parser.Parse()
	fmt.Println("Parse successful")
	fmt.Println("Printing AST:")
	fmt.Println(ast)
	fmt.Println("Start perform SA")
	sa.PerformSA(ast)
}

func GenerateIlocInstructions(ctx *cc.CompilerContext) {
	scanner := scanner.New(ctx)
	parser := parser.New(ctx, scanner)
	fmt.Println("Start parsing")
	programAst := parser.Parse()
	fmt.Println("Parse successful")
	fmt.Println("Start perform SA")
	sa.PerformSA(programAst)
	fmt.Println("Start Translatating ast into iloc")
	ir.ControlFlowFrags = make([]*ir.FuncFrag, 0)
	programAst.TranslateToILoc(programAst.GlobalSymbolTable)
	fmt.Println("Printing ILOC instructions:")
	PrintIlocInstructions(ir.ControlFlowFrags)
}

func PrintIlocInstructions(FuncFrags []*ir.FuncFrag) {
	//Find the longest label
	for _, funcFrag := range FuncFrags {
		ir.LongestLableLength = int(math.Max(float64(ir.LongestLableLength), float64(len(funcFrag.Label)+1)))
	}
	//Print instructions
	for _, funcFrag := range FuncFrags {
		fmt.Println(funcFrag.Label + ":")
		for _, instruction := range funcFrag.Body {
			if instruction != nil {
				fmt.Printf("%s%s\n", strings.Repeat(" ", ir.LongestLableLength), instruction.String())
			}

		}
	}
}

func getAssembly(ctx *cc.CompilerContext) []string {
	scanner := scanner.New(ctx)
	parser := parser.New(ctx, scanner)
	ast := parser.Parse()

	sa.PerformSA(ast)
	ir.ControlFlowFrags = make([]*ir.FuncFrag, 0)
	ast.TranslateToILoc(ast.GlobalSymbolTable)
	armInstructString := assembly.ToAssembly(ir.ControlFlowFrags, ast.GlobalSymbolTable)
	return armInstructString
}

func main() {
	/*Parse args and flags*/
	lexPtr := flag.Bool("lex", false, "Use -lex fileName to print the scanned tokens in the specified file")
	astPtr := flag.Bool("ast", false, "Use -ast fileName to print the ast for the specified file")
	ilocPtr := flag.Bool("iloc", false, "Use -iloc fileName to print the iloc instructions for the specified file")
	armPtr := flag.Bool("S", false, "Use -s to print out arm code")
	flag.Parse()
	argsWithoutProg := flag.Args()
	inputFileName := argsWithoutProg[0]
	ctx := cc.New(*lexPtr, inputFileName)
	/*Print the tokens of the input file if in lex mode*/
	if *lexPtr {
		scanner := scanner.New(ctx)
		scanner.PrintAllTokens()
	} else if *astPtr {
		StartCompiling(ctx)
	} else if *ilocPtr {
		GenerateIlocInstructions(ctx)
	} else if *armPtr {
		fileName := filepath.Base(ctx.SourcePath())
		fileType := filepath.Ext(ctx.SourcePath())
		fileName = strings.TrimSuffix(fileName, fileType) + ".s"

		armInstList := getAssembly(ctx)
		outStr := "" // dump arm code into a string
		for _, line := range armInstList {
			outStr = outStr + line + "\n"
		}

		f, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		_, err2 := f.WriteString(outStr)
		if err2 != nil {
			log.Fatal(err2)
		}
	}

}
