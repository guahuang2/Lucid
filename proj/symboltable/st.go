package symboltable

import (
	"fmt"
	"proj/ir"
	"proj/types"
)

type EntryValue struct {
	EntryType                 types.Type
	LocalSymbolTable          *SymbolTable
	Parameters                []types.Type
	ReturnType                types.Type
	FunctionName              string
	StructName                string
	BoolValue                 bool
	IntValue                  int
	ParametersRegisterLocList []int
	RegisterLoc               int
	ParaNames                 []string
}

type Entry interface {
	GetValue() *EntryValue
}

type IntLiteralEntry struct {
	EntryType types.Type
	IntValue  int
}

//func NewIntLiteralEntry(t types.Type, st SymbolTable) *structDefinitionEntry {
//	return &structDefinitionEntry{t, st}
//}
//
//type structInstanceEntry struct {
//	//Entry type for struct instance
//	entryType        types.Type
//	localSymbolTable SymbolTable
//}
//
//func NewstructInstanceEntry(t types.Type, st SymbolTable) *structDefinitionEntry {
//	return &structDefinitionEntry{t, st}
//}

//func (h *structInstanceEntry) GetValue() *EntryValue {
//	return &EntryValue{EntryType: h.entryType}
//}

type structDefinitionEntry struct {
	//Entry type for struct definition
	entryValue *EntryValue
}

func NewStructDefinition(t types.Type, st SymbolTable) *structDefinitionEntry {
	return &structDefinitionEntry{&EntryValue{EntryType: t, LocalSymbolTable: &st}}
}

func (h *structDefinitionEntry) GetValue() *EntryValue {
	return h.entryValue
}

type functionEntry struct {
	//Entry type for function
	entryValue *EntryValue
}

func NewFunctionEntry(t types.Type, st *SymbolTable, parameters []types.Type, returnType types.Type) *functionEntry {
	return &functionEntry{&EntryValue{EntryType: t, LocalSymbolTable: st, Parameters: parameters, ReturnType: returnType}}
}

func (f *functionEntry) GetValue() *EntryValue {
	return f.entryValue
}

type lowLevelEntry struct {
	//Entry type for int, bool and Unknown

	entryValue *EntryValue
}

func (l *lowLevelEntry) GetValue() *EntryValue {
	return l.entryValue
}

func NewEntry(t types.Type) *lowLevelEntry {
	return &lowLevelEntry{&EntryValue{EntryType: t}}
}

type SymbolTable struct {
	tableName         string
	typeMap           map[string]Entry
	fatherSymbolTable *SymbolTable
}

func (st *SymbolTable) String() string {
	return st.tableName
}

func NewSymbolTable(tableName string) *SymbolTable {
	//Create a symbol table without father
	return &SymbolTable{tableName, make(map[string]Entry), nil}
}

func NewWithFather(father *SymbolTable, tableName string) *SymbolTable {
	//Create a symbol table with father
	return &SymbolTable{tableName, map[string]Entry{}, father}
}

func (st *SymbolTable) GetRegisterLoc(id string) int {
	ety, exist := st.Contain(id)
	if !exist {
		panic(fmt.Sprintf("Variable Named :%s don't exisit", id))
	} else if ety.GetValue().EntryType != types.IntTySig && ety.GetValue().EntryType != types.BoolTySig {
		panic(fmt.Sprintf("Variable %s type: %s doesn't has a Register location", id, ety.GetValue().EntryType))
	}
	return ety.GetValue().RegisterLoc
}

func (st *SymbolTable) Insert(input string, t types.Type) {
	st.typeMap[input] = NewEntry(t)
}

func (st *SymbolTable) InsertWithNewReg(input string, t types.Type) Entry {
	st.typeMap[input] = NewEntry(t)
	st.typeMap[input].GetValue().RegisterLoc = ir.NewRegister()
	return st.typeMap[input]
}

func (st *SymbolTable) InsertStructDefinition(structName string, t types.Type, localST SymbolTable) {
	st.typeMap[structName] = NewStructDefinition(t, localST)
}

func (st *SymbolTable) InsertFunctionEntry(input string, t types.Type, localST *SymbolTable, para []types.Type, returnType types.Type) {
	st.typeMap[input] = NewFunctionEntry(t, localST, para, returnType)
}

func (st *SymbolTable) Contain(input string) (Entry, bool) {
	//Check whether the key exist in the local symboltable or its ancestor
	cur := st
	//fmt.Printf("Check for ident %s in symboltable: %s\n",input,st)
	//fmt.Println(cur.typeMap)
	for {
		entry, pre := cur.typeMap[input]
		if pre {
			return entry, pre
		} else if cur.fatherSymbolTable != nil {
			cur = cur.fatherSymbolTable
		} else {
			return nil, false
		}
	}
}

func (st *SymbolTable) ContainGlobally(input string) (Entry, bool) {
	//Check whether the key exist in the global symbol table
	if st == nil {
		panic("Nill symboltalbe in ContainGlobally")
	}
	cur := *st
	for {
		entry, pre := cur.typeMap[input]
		if pre {
			return entry, cur.fatherSymbolTable == nil
		} else if cur.fatherSymbolTable != nil {
			cur = *cur.fatherSymbolTable
		} else {
			return nil, false
		}
	}
}

func (st *SymbolTable) ContainLocally(input string) (Entry, bool) {
	//Check whether the key exist in the local symboltable or its ancestor
	exist, pre := st.typeMap[input]
	if pre {
		return exist, pre
	} else {
		return nil, false
	}
}

func (st *SymbolTable) ContainStructure(input string) (Entry, bool) {
	//Check whether the structure has been declared and return its definition symboltable
	entry, exist := st.Contain(input)
	if exist && entry.GetValue().EntryType == types.StructTySig {
		return entry, exist
	} else {
		return nil, false
	}
}

func (st *SymbolTable) GetFatherSymbol() (*SymbolTable, bool) {
	if st.fatherSymbolTable != nil {
		return st.fatherSymbolTable, true
	} else {
		return nil, false
	}
}

func (st *SymbolTable) GetRegLoc(input string) int {
	entry, exist := st.Contain(input)
	if exist {
		return entry.GetValue().RegisterLoc
	} else {
		panic("SA fail")
	}
}

var SymbolTableMap map[string]*SymbolTable
