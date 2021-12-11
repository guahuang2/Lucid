package ir

import (
	"bytes"
	"fmt"
	"proj/regDepatcher"
)

// Add represents a ADD instruction in ILOC
type Add struct {
	target    int       // The target register for the instruction
	sourceReg int       // The first source register of the instruction
	operand   int       // The operand either register or constant
	opty      OperandTy // The type for the operand (REGISTER, IMMEDIATE)
}

func (instr *Add) GetGlobal() string {
	return ""
}

//NewAdd is a constructor and initialization function for a new Add instruction
func NewAdd(target int, sourceReg int, operand int, opty OperandTy) *Add {
	return &Add{target, sourceReg, operand, opty}
}

func (instr *Add) GetTargets() []int {
	targets := make([]int, 1)
	targets = append(targets, instr.target)
	return targets
}
func (instr *Add) GetSources() []int {
	sources := make([]int, 1)
	sources = append(sources, instr.sourceReg)
	return sources
}
func (instr *Add) GetImmediate() *int {

	//Add instruction has two forms for the second operand: register, and immediate (constant)
	//make sure to check for that.
	if instr.opty == IMMEDIATE {
		return &instr.operand
	}
	//Return nil if this instruction does not have an immediate
	return nil
}
func (instr *Add) GetLabel() string {
	// Add does not work with labels so we can just return a default value
	return ""
}
func (instr *Add) SetLabel(newLabel string) {
	// Add does not work with labels can we can skip implementing this method.

}
func (instr *Add) String() string {

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

	out.WriteString(fmt.Sprintf("add %s,%s,%s", targetReg, sourceReg, operand2))

	return out.String()

}

func (instr *Add) ToAssembly(funcVarDict map[int]int, paramRegIds map[int]int) []string {
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

	// add
	targetRegId := regDepatcher.NextAvailReg()
	instruction = append(instruction, fmt.Sprintf("\tadd x%v,x%v,x%v", targetRegId, source1RegId, source2RegId))

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
