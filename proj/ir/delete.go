package ir

import (
	"bytes"
	"fmt"
	"proj/regDepatcher"
)

type Delete struct {
	sourceReg int
}

func NewDelete(sourceReg int) *Delete {
	return &Delete{sourceReg}
}

func (instr *Delete) GetTargets() []int { return []int{} }

func (instr *Delete) GetSources() []int {
	source := []int{}
	source = append(source, instr.sourceReg)
	return source
}

func (instr *Delete) GetImmediate() *int { return nil }

func (instr *Delete) GetGlobal() string { return "" }

func (instr *Delete) GetLabel() string { return "" }

func (instr *Delete) SetLabel(newLabel string) {}

func (instr *Delete) String() string {
	var out bytes.Buffer
	sourceRegister := fmt.Sprintf("r%v", instr.sourceReg)
	out.WriteString(fmt.Sprintf("delete %s", sourceRegister))
	return out.String()
}

func (instr *Delete) ToAssembly(funcVarDict map[int]int, paramRegIds map[int]int) []string {
	instruction := []string{}

	targetRegId := regDepatcher.NextAvailReg()
	delOffset := funcVarDict[instr.sourceReg]
	instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,#%v]", targetRegId, delOffset))
	instruction = append(instruction, fmt.Sprintf("\tmov x0,x%v", targetRegId))
	instruction = append(instruction, fmt.Sprintf("\tbl free"))
	regDepatcher.ReleaseReg(targetRegId)

	return instruction
}
