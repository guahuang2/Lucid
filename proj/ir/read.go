package ir

import (
	"bytes"
	"fmt"
	"proj/regDepatcher"
)

type Read struct {
	targetReg int
	variable  string
}

func NewRead(targetReg int, variable string) *Read {
	return &Read{targetReg, variable}
}

func (instr *Read) GetTargets() []int {
	target := []int{}
	target = append(target, instr.targetReg)
	return target
}

func (instr *Read) GetSources() []int { return []int{} }

func (instr *Read) GetImmediate() *int { return nil }

func (instr *Read) GetGlobal() string { return "" }

func (instr *Read) GetLabel() string { return "" }

func (instr *Read) SetLabel(newLabel string) {}

func (instr *Read) String() string {
	var out bytes.Buffer

	targetRegister := fmt.Sprintf("r%v", instr.targetReg)
	out.WriteString(fmt.Sprintf("read %s", targetRegister))
	return out.String()
}

func (instr *Read) ToAssembly(funcVarDict map[int]int, paramRegIds map[int]int) []string {
	regDepatcher.SetScan()
	instruction := []string{}

	varTargetRegId := regDepatcher.NextAvailReg()
	varTargetOffset := funcVarDict[instr.targetReg]
	instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,#%v]", varTargetRegId, varTargetOffset))

	sourceReg := regDepatcher.NextAvailReg()
	instruction = append(instruction, fmt.Sprintf("\tadrp x%v, .READ", sourceReg))
	instruction = append(instruction, fmt.Sprintf("\tadd x%v,x%v,:lo12:.READ", sourceReg, sourceReg))

	instruction = append(instruction, fmt.Sprintf("\tadd x%v,x29,#%v", varTargetRegId, varTargetOffset))
	instruction = append(instruction, fmt.Sprintf("\tmov x0,x%v", sourceReg))
	instruction = append(instruction, "\tbl scanf")

	regDepatcher.ReleaseReg(varTargetRegId)
	regDepatcher.ReleaseReg(sourceReg)

	return instruction
}
