package ir

import (
	"bytes"
	"fmt"
	"proj/regDepatcher"
)

type Ldr struct {
	target    int
	sourceReg int
	operand   int
	globalVar string
	opty      OperandTy
}

func NewLdr(target int, sourceReg int, operand int, globalVar string, opty OperandTy) *Ldr {
	return &Ldr{target, sourceReg, operand, globalVar, opty}
}

func (instr *Ldr) GetTargets() []int {
	targets := []int{}
	targets = append(targets, instr.target)
	return targets
}

func (instr *Ldr) GetSources() []int {
	sources := []int{}
	if instr.opty == REGISTER {
		sources = append(sources, instr.sourceReg, instr.operand)
	} else if instr.opty == IMMEDIATE || instr.opty == ONEOPERAND {
		sources = append(sources, instr.sourceReg)
	}
	return sources
}

func (instr *Ldr) GetImmediate() *int {
	if instr.opty == IMMEDIATE {
		return &instr.operand
	}
	return nil
}

func (instr *Ldr) GetGlobal() string {
	if instr.opty == GLOBALVAR {
		return instr.globalVar
	}
	return ""
}

func (instr *Ldr) GetLabel() string {
	return ""
}

func (instr *Ldr) SetLabel(newLabel string) {}

func (instr *Ldr) String() string {
	var out bytes.Buffer

	targetReg := fmt.Sprintf("r%v", instr.target)

	if instr.opty == REGISTER {
		sourceReg1 := fmt.Sprintf("r%v", instr.sourceReg)
		sourceReg2 := fmt.Sprintf("r%v", instr.operand)
		out.WriteString(fmt.Sprintf("ldr %s,%s,%s", targetReg, sourceReg1, sourceReg2))
	} else if instr.opty == IMMEDIATE {
		sourceReg := fmt.Sprintf("r%v", instr.sourceReg)
		operand2 := fmt.Sprintf("#%v", instr.operand)
		out.WriteString(fmt.Sprintf("ldr %s,%s,%s", targetReg, sourceReg, operand2))
	} else if instr.opty == ONEOPERAND {
		sourceReg := fmt.Sprintf("r%v", instr.sourceReg)
		out.WriteString(fmt.Sprintf("ldr %s,%s", targetReg, sourceReg))
	} else if instr.opty == GLOBALVAR {
		globVarName := fmt.Sprintf("%v", instr.globalVar)
		out.WriteString(fmt.Sprintf("ldr %s,%s", targetReg, globVarName))
	}

	return out.String()
}

func (instr *Ldr) ToAssembly(funcVarDict map[int]int, paramRegIds map[int]int) []string {
	instruction := []string{}
	if instr.opty == GLOBALVAR {
		addrRegId := regDepatcher.NextAvailReg()
		instruction = append(instruction, fmt.Sprintf("\tadrp x%v,%v", addrRegId, instr.globalVar))
		instruction = append(instruction, fmt.Sprintf("\tadd x%v,x%v, :lo12:%v", addrRegId, addrRegId, instr.globalVar))
		instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x%v]", addrRegId, addrRegId))
		targetOffset := funcVarDict[instr.target]
		instruction = append(instruction, fmt.Sprintf("\tstr x%v,[x29,#%v]", addrRegId, targetOffset))
		regDepatcher.ReleaseReg(addrRegId)
	}

	return instruction
}
