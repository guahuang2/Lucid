package ast

import (
	"fmt"
	"m4/token"
	"testing"
)

func Test1(t *testing.T) {

	t1 := token.Token{token.PLUS, "+"}
	t2 := token.Token{token.INT, "4"}
	t3 := token.Token{token.INT, "5"}

	intLit1 := &IntLiteral{t2, 4}
	intLit2 := &IntLiteral{t3, 5}

	bExpr := &BinOpExpr{t1, ADD, intLit1, intLit2}

	fmt.Println(bExpr)
}
