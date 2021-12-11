package ir

import (
	"bytes"
	"fmt"
	"proj/regDepatcher"
	"strconv"
)

type Pop struct {
	sourceReg []int
	funcName  string
}

func NewPop(sourceReg []int, funcName string) *Pop {
	return &Pop{sourceReg, funcName}
}

func (instr *Pop) GetTargets() []int { return []int{} }

func (instr *Pop) GetSources() []int {
	sources := []int{}
	for _, src := range instr.sourceReg {
		sources = append(sources, src)
	}
	return sources
}

func (instr *Pop) GetImmediate() *int { return nil }

func (instr *Pop) GetGlobal() string { return "" }

func (instr *Pop) GetLabel() string { return "" }

func (instr *Pop) SetLabel(newLabel string) {}

func (instr *Pop) String() string {
	var out bytes.Buffer
	var strSource string

	for id, src := range instr.sourceReg {
		if id != 0 {
			strSource = strSource + ","
		}
		strSource = strSource + "r" + strconv.Itoa(src)
	}

	out.WriteString(fmt.Sprintf("pop {%s} @%v", strSource, instr.funcName))

	return out.String()
}

func (instr *Pop) ToAssembly(funcVarDict map[int]int, paramRegIds map[int]int) []string {
	instruction := []string{}

	instruction = append(instruction, fmt.Sprintf("\tldr x0,[x29,#24]"))
	offset := len(instr.sourceReg) * 8
	if offset%16 != 0 {
		offset += 8
	}

	instruction = append(instruction, fmt.Sprintf("\tadd sp,sp,#%v", offset))

	iteration := 8
	if len(instr.sourceReg) <= 8 {
		iteration = len(instr.sourceReg)
	}

	for i := 1; i < iteration; i++ {
		regDepatcher.ReleaseReg(i)
	}
	return instruction
}
