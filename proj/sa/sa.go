package sa

import (
	"flag"
	"fmt"
	"proj/ast"
	st "proj/symboltable"
)

func reportErrors(errors []string) bool {
	// return true if there exists any error
	if len(errors) > 0 {
		for _, err := range errors {
			out := flag.CommandLine.Output()
			fmt.Fprintf(out, "semantic  error:%s\n", err)
		}
		return true
	}
	return false
}

func PerformSA(program *ast.Program) bool {
	// Define a new global table
	globalST := st.NewSymbolTable("Global")
	errors := make([]string, 0)

	// First Build the Symbol Table(s) for all declarations
	errors = program.PerformSABuild(errors, globalST)

	// Report errors
	if !reportErrors(errors) {
		// second perform type checking
		errors := make([]string, 0)
		//errors = program.TypeCheck(errors, globalST)
		if reportErrors(errors) { // finally no error
			return false
		}
	}
	return true
}
