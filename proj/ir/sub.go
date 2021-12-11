package ir

import (
	"bytes"
	"fmt"
	"proj/regDepatcher"
)

type Sub struct {
	target    int       // The target register for the instruction
	sourceReg int       // The first source register of the instruction
	operand   int       // The operand either register or constant
	opty      OperandTy // The type for the operand (REGISTER, IMMEDIATE)
}

func NewSub(target int, sourceReg int, operand int, opty OperandTy) *Sub {
	return &Sub{target, sourceReg, operand, opty}
}

func (instr *Sub) GetTargets() []int {
	targets := []int{}
	targets = append(targets, instr.target)
	return targets
}

func (instr *Sub) GetSources() []int {
	sources := []int{}
	if instr.opty == REGISTER {
		sources = append(sources, instr.sourceReg, instr.operand)
	} else {
		sources = append(sources, instr.sourceReg)
	}
	return sources
}

func (instr *Sub) GetImmediate() *int {

	if instr.opty == IMMEDIATE {
		return &instr.operand
	}
	return nil
}

func (instr *Sub) GetGlobal() string {
	return ""
}

func (instr *Sub) GetLabel() string {
	return ""
}

func (instr *Sub) SetLabel(newLabel string) {}

func (instr *Sub) String() string {

	var out bytes.Buffer

	targetReg := fmt.Sprintf("r%v", instr.target)
	sourceReg := fmt.Sprintf("r%v", instr.sourceReg)

	var prefix string

	if instr.opty == IMMEDIATE {
		prefix = "#"
	} else {
		prefix = "r"
	}
	operand2 := fmt.Sprintf("%v%v", prefix, instr.operand)

	out.WriteString(fmt.Sprintf("sub %s,%s,%s", targetReg, sourceReg, operand2))

	return out.String()

}

func (instr *Sub) ToAssembly(funcVarDict map[int]int, paramRegIds map[int]int) []string {
	instruction := []string{}
	var source1RegId int
	var source2RegId int
	var isParam1 bool
	var isParam2 bool

	// load operand 1
	if source1RegId, isParam1 = paramRegIds[instr.sourceReg]; !isParam1 {
		source1Offset := funcVarDict[instr.sourceReg]
		source1RegId = regDepatcher.NextAvailReg()
		instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,#%v]", source1RegId, source1Offset))
	}

	// load operand 2
	if source2RegId, isParam2 = paramRegIds[instr.operand]; !isParam2 {
		source2RegId = regDepatcher.NextAvailReg()
		if instr.opty == REGISTER {
			source2Offset := funcVarDict[instr.operand]
			instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,#%v]", source2RegId, source2Offset))
		} else {
			instruction = append(instruction, fmt.Sprintf("\tmov x%v,#%v", source2RegId, instr.operand))
		}
	}

	// sub
	targetRegId := regDepatcher.NextAvailReg()
	instruction = append(instruction, fmt.Sprintf("\tsubs x%v,x%v,x%v", targetRegId, source1RegId, source2RegId))

	// store result
	targetOffset := funcVarDict[instr.target]
	instruction = append(instruction, fmt.Sprintf("\tstr x%v,[x29,#%v]", targetRegId, targetOffset))

	regDepatcher.ReleaseReg(targetRegId)
	if !isParam1 {
		regDepatcher.ReleaseReg(source1RegId)
	}
	if !isParam2 {
		regDepatcher.ReleaseReg(source2RegId)
	}

	return instruction
}
