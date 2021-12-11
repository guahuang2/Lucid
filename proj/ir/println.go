package ir

import (
	"bytes"
	"fmt"
	"proj/regDepatcher"
)

type Println struct {
	sourceReg int
}

func NewPrintln(sourceReg int) *Println {
	return &Println{sourceReg}
}

func (instr *Println) GetTargets() []int { return []int{} }

func (instr *Println) GetSources() []int {
	source := []int{}
	source = append(source, instr.sourceReg)
	return source
}

func (instr *Println) GetImmediate() *int { return nil }

func (instr *Println) GetGlobal() string { return "" }

func (instr *Println) GetLabel() string { return "" }

func (instr *Println) SetLabel(newLabel string) {}

func (instr *Println) String() string {
	var out bytes.Buffer
	sourceRegister := fmt.Sprintf("r%v", instr.sourceReg)
	out.WriteString(fmt.Sprintf("println %s", sourceRegister))
	return out.String()
}

func (instr *Println) ToAssembly(funcVarDict map[int]int, paramRegIds map[int]int) []string {
	regDepatcher.SetPrintln()
	instruction := []string{}

	var targetRegId int
	var isTargetParam bool
	if targetRegId, isTargetParam = paramRegIds[instr.sourceReg]; !isTargetParam {
		targetRegId = regDepatcher.NextAvailReg()
		targetOffset := funcVarDict[instr.sourceReg]
		instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,#%v]", targetRegId, targetOffset))
	}

	sourceRegId := regDepatcher.NextAvailReg()
	instruction = append(instruction, fmt.Sprintf("\tadrp x%v, .PRINT_LN", sourceRegId))
	instruction = append(instruction, fmt.Sprintf("\tadd x%v,x%v, :lo12:.PRINT_LN", sourceRegId, sourceRegId))
	instruction = append(instruction, fmt.Sprintf("\tmov x1,x%v", targetRegId))
	instruction = append(instruction, fmt.Sprintf("\tmov x0,x%v", sourceRegId))
	instruction = append(instruction, "\tbl printf")

	if !isTargetParam {
		regDepatcher.ReleaseReg(targetRegId)
	}
	regDepatcher.ReleaseReg(sourceRegId)
	return instruction
}
