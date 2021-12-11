package ir

import (
	"bytes"
	"fmt"
	"proj/regDepatcher"
)

type Mul struct {
	target     int // The target register for the instruction
	sourceReg1 int // The first source register of the instruction
	sourceReg2 int // The second source register of the instruction
}

func NewMul(target int, sourceReg1 int, sourceReg2 int) *Mul {
	return &Mul{target, sourceReg1, sourceReg2}
}

func (instr *Mul) GetTargets() []int {
	targets := []int{}
	targets = append(targets, instr.target)
	return targets
}
func (instr *Mul) GetSources() []int {
	sources := []int{}
	sources = append(sources, instr.sourceReg1, instr.sourceReg2)
	return sources
}
func (instr *Mul) GetImmediate() *int {

	//Return nil if this instruction does not have an immediate
	return nil
}
func (instr *Mul) GetGlobal() string {
	return ""
}
func (instr *Mul) GetLabel() string {
	return ""
}
func (instr *Mul) SetLabel(newLabel string) {}

func (instr *Mul) String() string {

	var out bytes.Buffer

	targetReg := fmt.Sprintf("r%v", instr.target)
	sourceReg1 := fmt.Sprintf("r%v", instr.sourceReg1)
	sourceReg2 := fmt.Sprintf("r%v", instr.sourceReg2)

	out.WriteString(fmt.Sprintf("mul %s,%s,%s", targetReg, sourceReg1, sourceReg2))

	return out.String()

}

func (instr *Mul) ToAssembly(funcVarDict map[int]int, paramRegIds map[int]int) []string {
	instruction := []string{}

	var source1RedId, source2RedId int
	var isParam1, isParam2 bool

	// load operand 1
	if source1RedId, isParam1 = paramRegIds[instr.sourceReg1]; !isParam1 {
		source1Offset := funcVarDict[instr.sourceReg1]
		source1RedId = regDepatcher.NextAvailReg()
		instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,#%v]", source1RedId, source1Offset))
	}

	// load operand 2
	if source2RedId, isParam2 = paramRegIds[instr.sourceReg2]; !isParam2 {
		source2Offset := funcVarDict[instr.sourceReg2]
		source2RedId = regDepatcher.NextAvailReg()
		instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,#%v]", source2RedId, source2Offset))
	}

	// multiply
	targetRegId := regDepatcher.NextAvailReg()
	instruction = append(instruction, fmt.Sprintf("\tmul x%v,x%v,x%v", targetRegId, source1RedId, source2RedId))

	// store result
	targetOffset := funcVarDict[instr.target]
	instruction = append(instruction, fmt.Sprintf("\tstr x%v,[x29,#%v]", targetRegId, targetOffset))

	if !isParam1 {
		regDepatcher.ReleaseReg(source1RedId)
	}
	if !isParam2 {
		regDepatcher.ReleaseReg(source2RedId)
	}
	regDepatcher.ReleaseReg(targetRegId)

	return instruction
}
