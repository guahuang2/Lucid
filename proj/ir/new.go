package ir

import (
	"bytes"
	"fmt"
)

type NewStruct struct {
	target   int
	dataType string
	size     int
}

func GetNewStructInst(target int, dataType string, fieldsSize int) *NewStruct {
	return &NewStruct{target, dataType, fieldsSize}
}

func (instr *NewStruct) GetTargets() []int {
	target := []int{}
	target = append(target, instr.target)
	return target
}

func (instr *NewStruct) GetSources() []int { return []int{} }

func (instr *NewStruct) GetImmediate() *int { return nil }

func (instr *NewStruct) GetGlobal() string {
	return instr.dataType
}

func (instr *NewStruct) GetLabel() string { return "" }

func (instr *NewStruct) SetLabel(newLabel string) {}

func (instr *NewStruct) String() string {
	var out bytes.Buffer
	targetReg := fmt.Sprintf("r%v", instr.target)
	out.WriteString(fmt.Sprintf("new %s,%s", targetReg, instr.dataType))
	return out.String()
}

func (instr *NewStruct) ToAssembly(funcVarDict map[int]int, paramRegIds map[int]int) []string {
	instruction := []string{}

	// prepare for malloc, push x0... to stack
	offset := 16
	for i := 0; i < len(paramRegIds); i++ {
		instruction = append(instruction, fmt.Sprintf("\tstr x%v,[x29,%v]", i, offset))
		offset += 8
	}

	space := instr.size * 8
	instruction = append(instruction, fmt.Sprintf("\tmov x0,#%v", space))
	instruction = append(instruction, "\tbl malloc")
	targetOffset := funcVarDict[instr.target]
	instruction = append(instruction, fmt.Sprintf("\tstr x0,[x29,#%v]", targetOffset))

	// restore registers after malloc
	offset = 16
	for i := 0; i < len(paramRegIds); i++ {
		instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,%v]", i, offset))
		offset += 8
	}

	return instruction
}
