package ir

import (
	"bytes"
	"fmt"
	"proj/regDepatcher"
)

type Not struct {
	target  int       // The target register for the instruction
	operand int       // The operand either register or constant
	opty    OperandTy // The type for the operand (REGISTER, IMMEDIATE)
}

func NewNot(target int, operand int, opty OperandTy) *Not {
	return &Not{target, operand, opty}
}

func (instr *Not) GetTargets() []int {
	targets := []int{}
	targets = append(targets, instr.target)
	return targets
}

func (instr *Not) GetSources() []int {

	sources := []int{}
	if instr.opty != IMMEDIATE {
		sources = append(sources, instr.operand)
		return sources
	}
	return sources
}

func (instr *Not) GetImmediate() *int {

	if instr.opty == IMMEDIATE {
		return &instr.operand
	}
	return nil
}

func (instr *Not) GetGlobal() string {
	return ""
}

func (instr *Not) GetLabel() string {
	return ""
}

func (instr *Not) SetLabel(newLabel string) {}

func (instr *Not) String() string {

	var out bytes.Buffer

	targetReg := fmt.Sprintf("r%v", instr.target)

	var prefix string
	var operand2 string

	if instr.opty == IMMEDIATE {
		prefix = "#"
	} else {
		prefix = "r"
	}
	operand2 = fmt.Sprintf("%v%v", prefix, instr.operand)

	out.WriteString(fmt.Sprintf("not %s,%s", targetReg, operand2))

	return out.String()

}

func (instr *Not) ToAssembly(funcVarDict map[int]int, paramRegIds map[int]int) []string {
	instruction := []string{}

	// load operand
	sourceRegId := regDepatcher.NextAvailReg()
	if instr.opty == REGISTER {
		source2Offset := funcVarDict[instr.operand]
		instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,#%v]", sourceRegId, source2Offset))
	} else {
		instruction = append(instruction, fmt.Sprintf("\tmov x%v,#%v", sourceRegId, instr.operand))
	}

	targetRegId := regDepatcher.NextAvailReg()
	//instruction = append(instruction, fmt.Sprintf("neg x%v, x%v", targetRegId, sourceRegId))
	tempRedId := regDepatcher.NextAvailReg()
	instruction = append(instruction, fmt.Sprintf("\tmov x%v,#1", tempRedId))
	instruction = append(instruction, fmt.Sprintf("\tsubs x%v,x%v,x%v", targetRegId, tempRedId, targetRegId))
	regDepatcher.ReleaseReg(tempRedId)

	// store result
	targetOffset := funcVarDict[instr.target]
	instruction = append(instruction, fmt.Sprintf("\tstr x%v,[x29,#%v]", targetRegId, targetOffset))

	regDepatcher.ReleaseReg(sourceRegId)
	regDepatcher.ReleaseReg(targetRegId)

	return instruction
}
