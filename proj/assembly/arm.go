package assembly

import (
	"fmt"
	"proj/ir"
	"proj/regDepatcher"
	st "proj/symboltable"
	"strings"
)

func ToAssembly(funcfrags []*ir.FuncFrag, symTable *st.SymbolTable) []string {

	armInsList := []string{}
	regDepatcher.RegInit()
	regDepatcher.IOInit()

	armInsList = append(armInsList, "\t.arch armv8-a")
	// Append global variables
	if len(funcfrags) > 0 && strings.Contains(funcfrags[0].Label, "Global Variable") {
		for _, instruction := range funcfrags[0].Body {
			if instruction.GetGlobal() != "" {
				varName := instruction.GetGlobal()
				armInsList = append(armInsList, fmt.Sprintf("\t.comm %v,8,8", varName))
			}
		}
	}

	// function
	armInsList = append(armInsList, "\t.text")

	rems := funcfrags[1:]
	for _, funcfrag := range rems {
		offset := 0
		funcVarDict := make(map[int]int)
		for _, instruction := range funcfrag.Body {
			if instruction.GetTargets() != nil && len(instruction.GetTargets()) > 0 {
				offset -= 8
				funcVarDict[instruction.GetTargets()[0]] = offset
			}
		}

		entry, _ := symTable.Contain(funcfrag.Label)
		paraRegList := entry.GetValue().ParametersRegisterLocList
		paramRegIds := make(map[int]int)
		for id, regLoc := range paraRegList {
			paramRegIds[regLoc] = id
			regDepatcher.OccupyReg(id)
		}
		if funcfrag.Label == "main" {
			regDepatcher.OccupyReg(0)
		}

		armInsList = append(armInsList, "\t.type "+funcfrag.Label+",%function")
		armInsList = append(armInsList, "\t.global "+funcfrag.Label)
		armInsList = append(armInsList, "\t.p2align\t\t2")

		funcSize := offset
		if funcSize%16 != 0 {
			funcSize -= 8
		}
		armInsList = append(armInsList, fmt.Sprintf("%v:", funcfrag.Label))
		armInsList = append(armInsList, prologue(-funcSize)...)

		// find regID of parameters from symtable
		// save matching in a map
		remainingInstruction := funcfrag.Body[1:]
		for _, instruction := range remainingInstruction {
			//armInsList = append(armInsList, "ILOC: " + instruction.String())
			armInsList = append(armInsList, instruction.ToAssembly(funcVarDict, paramRegIds)...)
		}

		armInsList = append(armInsList, epilogue(-funcSize)...)
		armInsList = append(armInsList, "\t.size "+funcfrag.Label+",(.-"+funcfrag.Label+")")

		for id, _ := range paraRegList {
			regDepatcher.ReleaseReg(id)
		}
		if funcfrag.Label == "main" {
			regDepatcher.ReleaseReg(0)
		}
	}

	if regDepatcher.GetPrint() {
		printInst := []string{}
		printInst = append(printInst, ".PRINT:")
		printInst = append(printInst, "\t.asciz\t\"%ld\"")
		printInst = append(printInst, "\t.size\t.PRINT, 4")
		armInsList = append(armInsList, printInst...)
	}
	if regDepatcher.GetPrintln() {
		printInst := []string{}
		printInst = append(printInst, ".PRINT_LN:")
		printInst = append(printInst, "\t.asciz\t\"%ld\\n\"")
		printInst = append(printInst, "\t.size\t.PRINT_LN, 5")
		armInsList = append(armInsList, printInst...)
	}
	if regDepatcher.GetScan() {
		readInsts := []string{}
		readInsts = append(readInsts, ".READ:")
		readInsts = append(readInsts, "\t.asciz\t\"%ld\"")
		readInsts = append(readInsts, "\t.size\t.READ, 4")
		armInsList = append(armInsList, readInsts...)
	}

	return armInsList
}

func prologue(size int) []string {
	proInst := []string{}
	proInst = append(proInst, "\tsub sp,sp,16")
	proInst = append(proInst, "\tstp x29,x30,[sp]")
	proInst = append(proInst, "\tmov x29,sp")
	proInst = append(proInst, fmt.Sprintf("\tsub sp,sp,#%v", size))
	return proInst
}

func epilogue(size int) []string {
	epiInst := []string{}
	epiInst = append(epiInst, fmt.Sprintf("\tadd sp,sp,#%v", size))
	epiInst = append(epiInst, "\tldp x29,x30,[sp]")
	epiInst = append(epiInst, "\tadd sp,sp,16")
	epiInst = append(epiInst, "\tret")
	return epiInst
}
