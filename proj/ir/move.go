package ir

import (
	"bytes"
	"fmt"
	"proj/regDepatcher"
)

type Mov struct {
	flag    ApsrFlag
	target  int
	operand int
	opty    OperandTy
	retFlag bool
}

func NewMov(target int, operand int, flag ApsrFlag, opty OperandTy) *Mov {
	return &Mov{flag, target, operand, opty, false}
}

func (instr *Mov) GetTargets() []int {
	targets := []int{}
	targets = append(targets, instr.target)
	return targets
}

func (instr *Mov) GetSources() []int {
	sources := []int{}
	if instr.opty == REGISTER {
		sources = append(sources, instr.operand)
	}
	return sources
}

func (instr *Mov) GetImmediate() *int {
	if instr.opty == IMMEDIATE {
		return &instr.operand
	}
	return nil
}

func (instr *Mov) GetGlobal() string { return "" }

func (instr *Mov) GetLabel() string { return "" }

func (instr *Mov) SetLabel(newLabel string) {}

func (instr *Mov) String() string {
	var out bytes.Buffer
	var flag string

	switch instr.flag {
	case GT:
		flag = "gt"
		break
	case LT:
		flag = "lt"
		break
	case GE:
		flag = "ge"
		break
	case LE:
		flag = "le"
		break
	case EQ:
		flag = "eq"
		break
	case NE:
		flag = "ne"
		break
	case AL:
		flag = ""
		break
	}
	operator := fmt.Sprintf("mov%v", flag)
	targetReg := fmt.Sprintf("r%v", instr.target)
	var prefix string
	if instr.opty == IMMEDIATE {
		prefix = "#"
	} else {
		prefix = "r"
	}
	operand2 := fmt.Sprintf("%v%v", prefix, instr.operand)
	out.WriteString(fmt.Sprintf("%s %s,%s", operator, targetReg, operand2))

	if instr.retFlag {
		out.WriteString(fmt.Sprintf(" @Return"))
	}

	return out.String()
}

func (instr *Mov) ToAssembly(funcVarDict map[int]int, paramRegIds map[int]int) []string {
	instruction := []string{}

	if instr.flag == AL {
		if instr.retFlag {
			tempRegId := regDepatcher.NextAvailReg()
			tempOffset := funcVarDict[instr.target]
			instruction = append(instruction, fmt.Sprintf("\tmov x%v,x0", tempRegId))
			instruction = append(instruction, fmt.Sprintf("\tstr x%v,[x29,#%v]", tempRegId, tempOffset))
			regDepatcher.ReleaseReg(tempRegId)
			return instruction
		}
		var sourceRegId, targetRegId int
		var isSourceParam, isTargetParam bool

		if targetRegId, isTargetParam = paramRegIds[instr.target]; !isTargetParam {
			targetRegId = regDepatcher.NextAvailReg()
		}

		if instr.opty == REGISTER {
			if sourceRegId, isSourceParam = paramRegIds[instr.operand]; !isSourceParam {
				sourceOffset := funcVarDict[instr.operand]
				sourceRegId = regDepatcher.NextAvailReg()
				instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,#%v]", sourceRegId, sourceOffset))
			}
		}

		if instr.opty == REGISTER {
			instruction = append(instruction, fmt.Sprintf("\tmov x%v,x%v", targetRegId, sourceRegId))
		} else {
			instruction = append(instruction, fmt.Sprintf("\tmov x%v,#%v", targetRegId, instr.operand))
		}

		if !isTargetParam {
			targetOffset := funcVarDict[instr.target]
			instruction = append(instruction, fmt.Sprintf("\tstr x%v,[x29,#%v]", targetRegId, targetOffset))
		}

		if instr.opty == REGISTER && !isSourceParam {
			regDepatcher.ReleaseReg(sourceRegId)
		}
		if !isTargetParam {
			regDepatcher.ReleaseReg(targetRegId)
		}
	} else {
		label := NewLabelWithPre("skipMov")
		cmpResReg := regDepatcher.NextAvailReg()
		cmpResOffset := funcVarDict[instr.target]
		tempReg := regDepatcher.NextAvailReg()
		if instr.flag == LT {
			instruction = append(instruction, fmt.Sprintf("\tb.ge %v", label))

			if instr.opty == IMMEDIATE {
				instruction = append(instruction, fmt.Sprintf("\tmov x%v,#%v", tempReg, instr.operand))
			} else {
				operandOffset := funcVarDict[instr.operand]
				instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,#%v]", tempReg, operandOffset))
			}
			instruction = append(instruction, fmt.Sprintf("\tmov x%v,x%v", cmpResReg, tempReg))
			instruction = append(instruction, fmt.Sprintf("\tstr x%v,[x29,#%v]", cmpResReg, cmpResOffset))

			instruction = append(instruction, fmt.Sprintf("%v:", label))
		} else if instr.flag == NE {
			instruction = append(instruction, fmt.Sprintf("\tb.eq %v", label))

			if instr.opty == IMMEDIATE {
				instruction = append(instruction, fmt.Sprintf("\tmov x%v,#%v", tempReg, instr.operand))
			} else {
				operandOffset := funcVarDict[instr.operand]
				instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,#%v]", tempReg, operandOffset))
			}
			instruction = append(instruction, fmt.Sprintf("\tmov x%v,x%v", cmpResReg, tempReg))
			instruction = append(instruction, fmt.Sprintf("\tstr x%v,[x29,#%v]", cmpResReg, cmpResOffset))

			instruction = append(instruction, fmt.Sprintf("%v:", label))
		} else if instr.flag == LE {
			instruction = append(instruction, fmt.Sprintf("\tb.gt %v", label))

			if instr.opty == IMMEDIATE {
				instruction = append(instruction, fmt.Sprintf("\tmov x%v,#%v", tempReg, instr.operand))
			} else {
				operandOffset := funcVarDict[instr.operand]
				instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,#%v]", tempReg, operandOffset))
			}
			instruction = append(instruction, fmt.Sprintf("\tmov x%v,x%v", cmpResReg, tempReg))
			instruction = append(instruction, fmt.Sprintf("\tstr x%v,[x29,#%v]", cmpResReg, cmpResOffset))

			instruction = append(instruction, fmt.Sprintf("%v:", label))
		} else if instr.flag == EQ {
			instruction = append(instruction, fmt.Sprintf("\tb.ne %v", label))

			if instr.opty == IMMEDIATE {
				instruction = append(instruction, fmt.Sprintf("\tmov x%v,#%v", tempReg, instr.operand))
			} else {
				operandOffset := funcVarDict[instr.operand]
				instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,#%v]", tempReg, operandOffset))
			}
			instruction = append(instruction, fmt.Sprintf("\tmov x%v,x%v", cmpResReg, tempReg))
			instruction = append(instruction, fmt.Sprintf("\tstr x%v,[x29,#%v]", cmpResReg, cmpResOffset))

			instruction = append(instruction, fmt.Sprintf("%v:", label))

		}
		regDepatcher.ReleaseReg(tempReg)
		regDepatcher.ReleaseReg(cmpResReg)
	}
	return instruction
}

func (instr *Mov) SetRetFlag() {
	instr.retFlag = true
}
