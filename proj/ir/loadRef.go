package ir

import (
	"bytes"
	"fmt"
	"proj/regDepatcher"
)

// to access fields of a struct
type LoadRef struct {
	target     int
	source     int
	field      string
	structName string
	offset     int
}

func NewLoadRef(target int, source int, field string, structName string, offset int) *LoadRef {
	return &LoadRef{target, source, field, structName, offset}
}

func (instr *LoadRef) GetTargets() []int {
	targets := []int{}
	targets = append(targets, instr.target)
	return targets
}

func (instr *LoadRef) GetSources() []int {
	sources := []int{}
	sources = append(sources, instr.source)
	return sources
}

func (instr *LoadRef) GetImmediate() *int { return nil }

func (instr *LoadRef) GetGlobal() string {
	return instr.field
}

func (instr *LoadRef) GetLabel() string { return "" }

func (instr *LoadRef) SetLabel(newLabel string) {}

func (instr *LoadRef) String() string {
	var out bytes.Buffer

	targetReg := fmt.Sprintf("r%v", instr.target)
	sourceReg := fmt.Sprintf("r%v", instr.source)
	strField := fmt.Sprintf("@%v", instr.field)

	out.WriteString(fmt.Sprintf("loadRef %s,%s,%s[%s]", targetReg, sourceReg, strField, instr.structName))

	return out.String()
}

func (instr *LoadRef) ToAssembly(funcVarDict map[int]int, paramRegIds map[int]int) []string {
	instruction := []string{}

	loadToRegId := regDepatcher.NextAvailReg()
	loadToOffset := funcVarDict[instr.target]
	structRegId := regDepatcher.NextAvailReg()
	structOffset := funcVarDict[instr.source]
	fieldOffset := instr.offset * 8

	instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,#%v]", structRegId, structOffset))
	instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x%v,#%v]", loadToRegId, structRegId, fieldOffset))
	instruction = append(instruction, fmt.Sprintf("\tstr x%v,[x29,#%v]", loadToRegId, loadToOffset))

	regDepatcher.ReleaseReg(loadToRegId)
	regDepatcher.ReleaseReg(structRegId)

	return instruction
}
