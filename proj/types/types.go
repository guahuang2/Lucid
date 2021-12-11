package types

type Type interface {
	GetName() string
	GetType() Type
}

type IntTy struct{}

func (intTy *IntTy) GetName() string { return "int" }
func (intTy *IntTy) GetType() Type   { return IntTySig }

type BoolTy struct{}

func (boolTy *BoolTy) GetName() string {
	return "bool"
}
func (boolTy *BoolTy) GetType() Type { return BoolTySig }

type UnknownTy struct {
	typeStr string
}

func (unknownTy *UnknownTy) GetName() string {
	return unknownTy.typeStr
}
func (unknownTy *UnknownTy) GetType() Type {
	return UnknownTySig
}

func NewUnknownTy(nameStr string) *UnknownTy {
	return &UnknownTy{typeStr: nameStr}
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

func (structTy *StructTy) GetType() Type {
	return StructTySig
}

type FunctTy struct {
	funcName string
}

func NewFuncTy(functName string) *FunctTy {
	return &FunctTy{funcName: functName}
}

func (functTy *FunctTy) GetName() string {
	return functTy.funcName
}

func (functTy *FunctTy) GetType() Type {
	return FuncTySig
}

type NilTy struct {
	funcName string
}

func (n *NilTy) GetName() string {
	return "Nil"
}
func (n *NilTy) GetType() Type {
	return NilTySig
}

type ImportTy struct{}

func (importTy *ImportTy) GetName() string {
	return "import"
}

func (importTy *ImportTy) GetType() Type {
	return ImportTySig
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
