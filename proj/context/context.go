package context

import (
	"fmt"
)

type CompilerContext struct {
	lexOut     bool
	sourcePath string
}

func New(lexOut bool, sourcePath string) *CompilerContext {
	return &CompilerContext{lexOut, sourcePath}
}

func (ctx *CompilerContext) OutputLex() bool { return ctx.lexOut }

func (ctx *CompilerContext) SourcePath() string { return ctx.sourcePath }

func (ctx *CompilerContext) RuntimeError(msg string, e error) {
	if e != nil {
		fmt.Println(msg)
	}
}
