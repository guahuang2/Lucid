package ir

import (
	"bytes"
	"fmt"
	"proj/regDepatcher"
	"strconv"
)

type Push struct {
	sourceReg []int
	funcName  string
}

func NewPush(sourceReg []int, funcName string) *Push {
	return &Push{sourceReg, funcName}
}

func (instr *Push) GetTargets() []int { return []int{} }

func (instr *Push) GetSources() []int {
	sources := []int{}
	for _, src := range instr.sourceReg {
		sources = append(sources, src)
	}
	return sources
}

func (instr *Push) GetImmediate() *int { return nil }

func (instr *Push) GetGlobal() string { return "" }

func (instr *Push) GetLabel() string { return "" }

func (instr *Push) SetLabel(newLabel string) {}

func (instr *Push) String() string {
	var out bytes.Buffer
	var strSource string

	for id, src := range instr.sourceReg {
		if id != 0 {
			strSource = strSource + ","
		}
		strSource = strSource + "r" + strconv.Itoa(src)
	}

	out.WriteString(fmt.Sprintf("push {%s} @%v", strSource, instr.funcName))

	return out.String()
}

func (instr *Push) ToAssembly(funcVarDict map[int]int, paramRegIds map[int]int) []string {
	instruction := []string{}

	instruction = append(instruction, fmt.Sprintf("\tstr x0,[x29,#24]"))
	offset := len(instr.sourceReg) * 8
	if offset%16 != 0 {
		offset += 8
	}
	instruction = append(instruction, fmt.Sprintf("\tsub sp,sp,#%v", offset))

	iteration := 8
	if len(instr.sourceReg) <= 8 {
		iteration = len(instr.sourceReg)
	}

	for i := 0; i < iteration; i++ {
		aggRegId := regDepatcher.NextAvailReg()
		argOffset := funcVarDict[instr.sourceReg[i]]
		instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,#%v]", aggRegId, argOffset))
		instruction = append(instruction, fmt.Sprintf("\tmov x%v,x%v", i, aggRegId))
		regDepatcher.ReleaseReg(aggRegId)
		regDepatcher.OccupyReg(i)
	}
	return instruction
}
