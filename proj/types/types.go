package types

type Type interface {
	GetName() string
}

type IntTy struct{}

func (intTy *IntTy) GetName() string { return "int" }

type BoolTy struct{}

func (boolTy *BoolTy) GetName() string {
	return "bool"
}

type UnknownTy struct{}

func (unknownTy *UnknownTy) GetName() string {
	return "unknown"
}

type StructTy struct {
	structName string
}

func NewStructTy(structName string) *StructTy {
	return &StructTy{structName: structName}
}

func (structTy *StructTy) GetName() string {
	return structTy.structName
}

type FunctTy struct {
	funcName string
}

func NewFuncTy(functName string) *FunctTy {
	return &FunctTy{funcName: functName}
}

func (functTy *FunctTy) GetName() string {
	return "function"
}

type NilTy struct {
	funcName string
}

func (n *NilTy) GetName() string {
	return "Nil"
}

type ImportTy struct{}

func (importTy *ImportTy) GetName() string {
	return "import"
}

var IntTySig *IntTy
var BoolTySig *BoolTy
var UnknownTySig *UnknownTy
var ImportTySig *ImportTy
var NilTySig *NilTy
var FuncTySig *FunctTy
var StructTySig *StructTy

func init() {
	IntTySig = &IntTy{}
	BoolTySig = &BoolTy{}
	UnknownTySig = &UnknownTy{}
	ImportTySig = &ImportTy{}
	NilTySig = &NilTy{}
	FuncTySig = &FunctTy{}
	StructTySig = &StructTy{}
}
