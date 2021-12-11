package symboltable

import (
	"proj/types"
)

type EntryValue struct {
	EntryType        types.Type
	LocalSymbolTable *SymbolTable
	Parameters       []types.Type
	ReturnType       types.Type
	FunctionName     string
	StructName       string
}

type Entry interface {
	GetValue() *EntryValue
}

type higherLevelEntry struct {
	//Entry type for struct
	entryType        types.Type
	localSymbolTable SymbolTable
}

func NewHigherEntry(t types.Type, st SymbolTable) *higherLevelEntry {
	return &higherLevelEntry{t, st}
}

func (h *higherLevelEntry) GetValue() *EntryValue {
	return &EntryValue{EntryType: h.entryType}
}

type functionEntry struct {
	//Entry type for function
	entryType        types.Type
	localSymbolTable *SymbolTable
	parameters       []types.Type
	returnType       types.Type
}

func NewFunctionEntry(t types.Type, st *SymbolTable, parameters []types.Type, returnType types.Type) *functionEntry {
	return &functionEntry{t, st, parameters, returnType}
}

func (f *functionEntry) GetValue() *EntryValue {
	return &EntryValue{EntryType: f.entryType, LocalSymbolTable: f.localSymbolTable, Parameters: f.parameters, ReturnType: f.returnType}
}

type lowLevelEntry struct {
	//Entry type for int and bool
	entryType types.Type
}

func (l *lowLevelEntry) GetValue() *EntryValue {
	return &EntryValue{EntryType: l.entryType}
}

func NewEntry(t types.Type) *lowLevelEntry {
	return &lowLevelEntry{t}
}

type SymbolTable struct {
	typeMap           map[string]Entry
	fatherSymbolTable *SymbolTable
}

func NewSymbolTable() *SymbolTable {
	//Create a global symbol table
	return &SymbolTable{map[string]Entry{}, nil}
}

func NewWithFather(father *SymbolTable) *SymbolTable {
	//Create a symbol table with father
	return &SymbolTable{map[string]Entry{}, father}
}

func (st *SymbolTable) Insert(input string, t types.Type) {
	st.typeMap[input] = NewEntry(t)
}

func (st *SymbolTable) InsertHigherinput(input string, t types.Type, localST SymbolTable) {
	st.typeMap[input] = NewHigherEntry(t, localST)
}

func (st *SymbolTable) InsertFunctionEntry(input string, t types.Type, localST *SymbolTable, para []types.Type, returnType types.Type) {
	st.typeMap[input] = NewFunctionEntry(t, localST, para, returnType)
}

func (st *SymbolTable) Contain(input string) (Entry, bool) {
	//Look for the key, if exist in the symboltable or its ancestor
	for {
		exist, pre := st.typeMap[input]
		if pre {
			return exist, pre
		} else if st.fatherSymbolTable != nil {
			st = st.fatherSymbolTable
		} else {
			return nil, false
		}
	}
}

func (st *SymbolTable) GetFatherSymbol() (*SymbolTable, bool) {
	if st.fatherSymbolTable != nil {
		return st.fatherSymbolTable, true
	} else {
		return nil, false
	}
}

var SymbolTableMap map[string]*SymbolTable
