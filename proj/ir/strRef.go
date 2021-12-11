package ir

import (
	"bytes"
	"fmt"
	"proj/regDepatcher"
)

// to access fields of a struct
type StrRef struct {
	target     int
	source     int
	field      string
	structName string
	fieldIdx   int
}

func NewStrRef(target int, source int, field string, structName string, fieldIdx int) *StrRef {
	return &StrRef{target, source, field, structName, fieldIdx}
}

func (instr *StrRef) GetTargets() []int {
	targets := []int{}
	targets = append(targets, instr.target)
	return targets
}

func (instr *StrRef) GetSources() []int {
	sources := []int{}
	sources = append(sources, instr.source)
	return sources
}

func (instr *StrRef) GetImmediate() *int { return nil }

func (instr *StrRef) GetGlobal() string {
	return instr.field
}

func (instr *StrRef) GetLabel() string { return "" }

func (instr *StrRef) SetLabel(newLabel string) {}

func (instr *StrRef) String() string {
	var out bytes.Buffer

	targetReg := fmt.Sprintf("r%v", instr.target)
	sourceReg := fmt.Sprintf("r%v", instr.source)
	strField := fmt.Sprintf("@%v", instr.field)

	out.WriteString(fmt.Sprintf("strRef %s,%s,%s[%s]", targetReg, sourceReg, strField, instr.structName))

	return out.String()
}

func (instr *StrRef) ToAssembly(funcVarDict map[int]int, paramRegIds map[int]int) []string {
	instruction := []string{}

	var targetRegId int
	var istargetParam bool
	if targetRegId, istargetParam = paramRegIds[instr.target]; !istargetParam {
		targetRegId = regDepatcher.NextAvailReg()
		targetOffSet := funcVarDict[instr.target]
		instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,#%v]", targetRegId, targetOffSet))
	}

	sourceRegId := regDepatcher.NextAvailReg()
	sourceOffSet := funcVarDict[instr.source]
	instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,#%v]", sourceRegId, sourceOffSet))

	fieldOffset := instr.fieldIdx * 8
	instruction = append(instruction, fmt.Sprintf("\tstr x%v,[x%v,#%v]", targetRegId, sourceRegId, fieldOffset))

	if !istargetParam {
		regDepatcher.ReleaseReg(targetRegId)
	}
	regDepatcher.ReleaseReg(sourceRegId)
	return instruction
}
