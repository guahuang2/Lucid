package ir

import (
	"fmt"
	"math"
)

type registerGen struct {
	count int //The current counter for the label
}

func NewRegister() int {
	retVal := rGen.count
	rGen.count += 1
	return retVal
}

type labelGen struct {
	count int //The current label number
}

func NewLabelWithPre(prefix string) string {
	retVal := fmt.Sprintf("%s_L%d", prefix, lGen.count)
	LongestLableLength = int(math.Max(float64(LongestLableLength), float64(len(retVal))))
	lGen.count += 1
	return retVal
}

func NewLabel() string {
	retVal := fmt.Sprintf("L%d", lGen.count)
	LongestLableLength = int(math.Max(float64(LongestLableLength), float64(len(retVal))))
	lGen.count += 1
	return retVal
}

var rGen *registerGen
var lGen *labelGen
var LongestLableLength int
var ControlFlowFrags []*FuncFrag

// The init() function will only be called once per package. This is where you can setup singletons for types
func init() {
	rGen = &registerGen{0}
	lGen = &labelGen{0}
	LongestLableLength = 0
}
