package ast

import (
	"bytes"
	"fmt"
	"proj/ir"
	st "proj/symboltable"
	"proj/token"
	"proj/types"
)

type Node interface {
	TokenLiteral() string
	String() string
	TypeCheck(errors []string, symTable *st.SymbolTable) []string
}

type Expr interface {
	// Node /*The expression type of node, implements GetType */
	Node
	GetType(symTable *st.SymbolTable) types.Type
	TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable)
	GetRegLoc() int
}

type Stat interface {
	// Node /*The statement type of node, implements PerformSABuild */
	Node
	PerformSABuild(errors []string, symTable *st.SymbolTable) []string
	TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable)
}

type Program struct {
	Token             *token.Token
	Package           *Package
	Import            *Import
	Types             *Types
	Declarations      *Declarations
	Functions         *Functions
	GlobalSymbolTable *st.SymbolTable
}

func NewProgram(pac *Package, imp *Import, typ *Types, decs *Declarations, funcs *Functions) *Program {
	return &Program{nil, pac, imp, typ, decs, funcs, nil}
}

func (p *Program) TokenLiteral() string {
	if p.Token != nil {
		return p.Token.Literal
	}
	panic("Could not determine token literal for program statement.")
}

func (p *Program) String() string {
	out := bytes.Buffer{}
	out.WriteString(p.Package.String())
	out.WriteString(p.Import.String())
	out.WriteString(p.Types.String())
	out.WriteString(p.Declarations.String())
	out.WriteString(p.Functions.String())
	return out.String()
}

func (p *Program) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = p.Package.TypeCheck(errors, symTable)
	errors = p.Import.TypeCheck(errors, symTable)
	errors = p.Types.TypeCheck(errors, symTable)
	errors = p.Declarations.TypeCheck(errors, symTable)
	errors = p.Functions.TypeCheck(errors, symTable)
	return errors
}

func (p *Program) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	p.GlobalSymbolTable = symTable
	errors = p.Package.PerformSABuild(errors, symTable)
	errors = p.Import.PerformSABuild(errors, symTable)
	errors = p.Types.PerformSABuild(errors, symTable)
	errors = p.Declarations.PerformSABuild(errors, symTable)
	errors = p.Functions.PerformSABuild(errors, symTable)
	return errors
}

func (p *Program) TranslateToILoc(symTable *st.SymbolTable) {
	p.Functions.TranslateToILoc(symTable)
}

type Package struct {
	Token *token.Token
	Ident IdentLiteral
}

func NewPackage(ident IdentLiteral) *Package {
	return &Package{nil, ident}
}

func (p *Package) TokenLiteral() string {
	if p.Token != nil {
		return p.Token.Literal
	}
	panic("Could not determine token literal for package statement")
}

func (p *Package) String() string {
	out := bytes.Buffer{}
	out.WriteString("package")
	out.WriteString(" ")
	out.WriteString(p.Ident.String())
	out.WriteString(";")
	out.WriteString("\n")
	return out.String()
}

func (p *Package) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

func (p *Package) PerformSABuild(errors []string, table *st.SymbolTable) []string {
	return errors
}

type Import struct {
	Token *token.Token
	Ident IdentLiteral
}

func NewImport(ident IdentLiteral) *Import {
	return &Import{nil, ident}
}

func (i *Import) TokenLiteral() string {
	if i.Token != nil {
		return i.Token.Literal
	}
	panic("Could not determine token literal for import statement")
}

func (i *Import) String() string {
	out := bytes.Buffer{}
	out.WriteString("import")
	out.WriteString(" ")
	out.WriteString("\"")
	out.WriteString(i.Ident.String())
	out.WriteString("\"")
	out.WriteString(";")
	out.WriteString("\n")
	return out.String()
}

func (i *Import) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

func (i *Import) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

type Types struct {
	token     *token.Token
	typedecls []TypeDeclaration
}

func NewTypes(typedecls []TypeDeclaration) *Types {
	return &Types{nil, typedecls}
}

func (t *Types) TokenLiterals() string {
	return t.token.Literal
}

func (t *Types) String() string {
	out := bytes.Buffer{}

	for _, typedec := range t.typedecls {
		out.WriteString(typedec.String())
		out.WriteString("\n")
	}

	return out.String()
}

func (t *Types) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	for _, typedecl := range t.typedecls {
		errors = typedecl.TypeCheck(errors, symTable)
	}
	return errors
}

func (t *Types) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	for _, typedecl := range t.typedecls {
		errors = typedecl.PerformSABuild(errors, symTable)
	}
	return errors
}

type TypeDeclaration struct {
	Token   *token.Token
	Ident   IdentLiteral
	Fields  *Fields
	LocalST *st.SymbolTable
}

func NewTypeDeclaration(ident IdentLiteral, fields *Fields) *TypeDeclaration {
	return &TypeDeclaration{nil, ident, fields, nil}
}

func (t *TypeDeclaration) TokenLiterals() string {
	if t.Token != nil {
		return t.Ident.String()
	}
	panic("Could not determine identfication literal for type statement")
}

func (t *TypeDeclaration) String() string {
	out := bytes.Buffer{}
	out.WriteString("type")
	out.WriteString(" ")
	out.WriteString(t.Ident.String())
	out.WriteString(" ")
	out.WriteString("struct")
	out.WriteString("{\n")
	out.WriteString(t.Fields.String())
	out.WriteString("\n}")
	out.WriteString(";")
	return out.String()
}

func (t *TypeDeclaration) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

func (t *TypeDeclaration) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	/*
		Create local st from global st and add the fields to local st
	*/
	if _, ext := symTable.Contain(t.Ident.Id); ext {
		errors = append(errors, fmt.Sprintf("Struct name: %s has already been used", t.Ident.Id))
	} else {
		t.LocalST = st.NewWithFather(symTable, "Struct:"+t.Ident.String())
		symTable.InsertStructDefinition(t.Ident.Id, types.NewStructTy(t.Ident.Id), *t.LocalST)
	}
	typeEntry, _ := symTable.Contain(t.Ident.Id)
	paraStringList := []string{}
	for _, decl := range t.Fields.Decls {
		paraStringList = append(paraStringList, decl.Type.TypeString)
	}
	typeEntry.GetValue().ParaNames = paraStringList
	errors = t.Fields.PerformSABuild(errors, t.LocalST)
	return errors
}

type Fields struct {
	Token *token.Token
	Decls []Decl
}

func NewFields(decls []Decl) *Fields {
	return &Fields{nil, decls}
}

func (f *Fields) TokenLiterals() string {
	if f.Token != nil {
		return f.Token.Literal
	}
	panic("Could not determine token literals for fields")
}

func (f *Fields) String() string {
	out := bytes.Buffer{}
	for _, decl := range f.Decls {
		out.WriteString(decl.String())
		out.WriteString(";")
		out.WriteString("\n")
	}
	return out.String()
}

func (f *Fields) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	//for _, decl := range f.Decls {
	//	decl.TypeCheck(errors, symTable)
	//}
	return errors
}

func (f *Fields) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	/*
		call PerformSABuild for all contained filed
	*/
	for _, decl := range f.Decls {
		errors = decl.PerformSABuild(errors, symTable)
	}
	return errors
}

type Decl struct {
	Token *token.Token
	Ident IdentLiteral
	Type  *Type
}

func NewDecl(ident IdentLiteral, Type *Type) *Decl {
	return &Decl{nil, ident, Type}
}

func (d *Decl) TokenLiterals() string {
	if d.Token != nil {
		return d.Token.Literal
	}
	panic("Could not determine token literals for decl")
}

func (d *Decl) String() string {
	out := bytes.Buffer{}
	out.WriteString(d.Ident.String())
	out.WriteString(" ")
	out.WriteString(d.Type.String())
	return out.String()
}

func (d *Decl) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	/*
		Check whether fields names have been used in the local symbol table
		and whether fields type has been declared
	*/
	//errors=d.Type.TypeCheck(errors, symTable)
	return errors
}

func (d *Decl) GetType(symTable *st.SymbolTable) types.Type {
	//Equivalent to Decl.Type.GetType(symTable)
	return d.Type.GetType(symTable)
}

func (d *Decl) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	/*
		Check whether decl name exist in local symbol table and the validity of its type,
		if valid, insert it into the local symbol table
	*/
	if _, ext := symTable.ContainLocally(d.Ident.Id); ext {
		errors = append(errors, fmt.Sprintf("field name:%s  has already been used", d.Ident.Id))
	} else {
		typeSig := d.Type.GetType(symTable)
		if typeSig == types.BoolTySig || typeSig == types.IntTySig {
			symTable.Insert(d.Ident.Id, typeSig)
		} else if typeSig.GetType() == types.StructTySig { //type struct
			symTable.InsertStructDefinition(typeSig.GetName(), typeSig, *symTable)
		} else { //type unknown
			symTable.Insert(d.Ident.Id, typeSig)
			errors = append(errors, fmt.Sprintf("Sturct:%s  not declared", d.Type.TypeString))
		}
		fmt.Sprintf("%s defined with type : %s hahahahaha", d.Ident.Id, typeSig.GetName())
	}
	fmt.Printf("Ident:%s defined in symboltable: %s \n", d.Ident.Id, symTable)
	return errors
}

type Type struct {
	Token      *token.Token
	TypeString string //int,bool or id
}

func NewType(TypeString string) *Type {
	return &Type{nil, TypeString}
}

func (t *Type) TokenLiterals() string {
	if t.Token != nil {
		return t.Token.Literal
	}
	panic("Could not determine token literals for Type")
}

func (t *Type) String() string {
	out := bytes.Buffer{}
	out.WriteString(t.TypeString)
	return out.String()
}

func (t *Type) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	if t.TypeString == "int" || t.TypeString == "bool" {
		return errors
	}
	if _, exist := symTable.Contain(t.TypeString[1:]); !exist {
		errors = append(errors, fmt.Sprintf("Structured named %s not defined on line num: %d", t.TypeString, t.Token.Rows))
	}
	return errors
}

func (t *Type) GetType(symTable *st.SymbolTable) types.Type {
	/*
		Check whether the type has been defined, if not unknown is returned
	*/
	var typeSig types.Type
	if t.TypeString == "" {
		typeSig = types.NilTySig
	} else if t.TypeString == "int" {
		typeSig = types.IntTySig
	} else if t.TypeString == "bool" {
		typeSig = types.BoolTySig
	} else if _, exist := symTable.ContainStructure(t.TypeString[1:]); exist {
		typeSig = types.NewStructTy(t.TypeString[1:])
	} else {
		typeSig = types.NewUnknownTy(t.TypeString[1:])
	}
	return typeSig
}

type Declarations struct {
	Token        *token.Token
	Declarations []Declaration
}

func NewDeclarations(decls []Declaration) *Declarations {
	return &Declarations{nil, decls}
}

func (d *Declarations) TokenLiterals() string {
	if d.Token != nil {
		return d.Token.Literal
	}
	panic("Could not determine token literals for Declarations")
}

func (d *Declarations) String() string {
	out := bytes.Buffer{}
	for _, dec := range d.Declarations {
		out.WriteString(dec.String())
		out.WriteString("\n")
	}
	return out.String()
}

func (d *Declarations) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	//Check whether the declarations already been declared in the given symbol table
	for _, decl := range d.Declarations {
		errors = decl.PerformSABuild(errors, symTable)
	}
	return errors
}

func (d *Declarations) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	for _, decl := range d.Declarations {
		errors = decl.TypeCheck(errors, symTable)
	}
	return errors
}

func (d *Declarations) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	for _, decl := range d.Declarations {
		decl.TranslateToILoc(frag, table)
	}
}

type Declaration struct {
	Token *token.Token
	Ids   *Ids
	Type  *Type
}

func NewDeclaration(ids *Ids, Type *Type) *Declaration {
	return &Declaration{nil, ids, Type}
}

func (d *Declaration) TokenLiterals() string {
	if d.Token != nil {
		return d.Token.Literal
	}
	panic("Could not determine token literals for declaration")
}

func (d *Declaration) String() string {
	out := bytes.Buffer{}
	out.WriteString("var")
	out.WriteString(" ")
	out.WriteString(d.Ids.String())
	out.WriteString(" ")
	out.WriteString(d.Type.String())
	out.WriteString(";")
	return out.String()
}

func (d *Declaration) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	/*
		Check whether the Ids of the declaration already been declared in the symbol table
		If not, add it to the local symbol table
	*/
	for _, id := range d.Ids.Idents {
		if _, ext := symTable.Contain(id.Id); ext {
			errors = append(errors, fmt.Sprintf("%s ident has already been used on line num: %d", id.Id, id.Token.Rows))
		} else {
			entry := symTable.InsertWithNewReg(id.Id, d.Type.GetType(symTable))
			fmt.Printf("Ident:%s defined in symboltable: %s \n", id.Id, symTable)
			fmt.Printf("Entry %s register loc: %d \n", id.Id, entry.GetValue().RegisterLoc)
		}
	}
	return errors
}

func (d *Declaration) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = d.Type.TypeCheck(errors, symTable)
	return errors
}

func (d *Declaration) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	for _, id := range d.Ids.Idents {
		_, exist := table.Contain(id.Id)
		if !exist {
			panic("fail sa")
		}
	}
}

type Ids struct {
	Token  *token.Token
	Idents []IdentLiteral //id literal list
}

func NewIds(idents []IdentLiteral) *Ids {
	return &Ids{nil, idents}
}

func (id *Ids) TokenLiterals() string {
	if id.Token != nil {
		return id.Token.Literal
	}
	panic("Could not determine token literals for ids")
}

func (id *Ids) String() string {
	out := bytes.Buffer{}
	for i, id := range id.Idents {
		if i > 0 {
			out.WriteString(",")
		}
		out.WriteString(id.String())
	}
	return out.String()
}

func (id *Ids) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

type Functions struct {
	Token         *token.Token
	functionArray []Function
}

func NewFunctions(funcs []Function) *Functions {
	return &Functions{nil, funcs}
}

func (funcs *Functions) TokenLiterals() string {
	if funcs.Token != nil {
		return funcs.Token.Literal
	}
	panic("Could not determine token literals for the functions")
}

func (funcs *Functions) String() string {
	out := bytes.Buffer{}
	for _, fun := range funcs.functionArray {
		out.WriteString(fun.String())
		out.WriteString("\n")
	}
	return out.String()
}

func (funcs *Functions) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	for _, function := range funcs.functionArray {
		errors = function.PerformSABuild(errors, symTable)
	}
	return errors
}

func (funcs *Functions) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	for _, function := range funcs.functionArray {
		errors = function.TypeCheck(errors, symTable)
	}
	return errors
}

func (funcs *Functions) TranslateToILoc(table *st.SymbolTable) {
	for _, function := range funcs.functionArray {
		var tempFuncFrag = &ir.FuncFrag{Body: []ir.Instruction{}}
		ir.ControlFlowFrags = append(ir.ControlFlowFrags, tempFuncFrag)
		function.TranslateToILoc(tempFuncFrag, table)
	}
}

type Function struct {
	Token        *token.Token
	Ident        IdentLiteral
	Parameters   *Parameters
	ReturnType   *ReturnType
	Declarations *Declarations
	Statements   *Statements
	localST      *st.SymbolTable
}

func NewFunction(ident IdentLiteral, params *Parameters, returnType *ReturnType, declarations *Declarations, statements *Statements) *Function {
	return &Function{nil, ident, params, returnType, declarations, statements, nil}
}

func (f *Function) TokenLiterals() string {
	if f.Token != nil {
		return f.Token.Literal
	}
	panic("Could not determine token literals for functions")
}

func (f *Function) String() string {
	out := bytes.Buffer{}
	out.WriteString("func")
	out.WriteString(" ")
	out.WriteString(f.Ident.String())
	out.WriteString(" ")
	out.WriteString(f.Parameters.String())
	out.WriteString(" ")
	out.WriteString(f.ReturnType.String())
	out.WriteString("{")
	out.WriteString("\n")
	out.WriteString(f.Declarations.String())
	out.WriteString(" ")
	out.WriteString(f.Statements.String())
	out.WriteString("}")
	out.WriteString("\n")
	return out.String()
}

func (f *Function) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	/*
		Stores the functionName, parameter types and return types in the global symbol table.
	*/
	//fmt.Println("Start function PerformSA")
	_, exist := symTable.Contain(f.Ident.Id)
	if exist {
		errors = append(errors, fmt.Sprintf("Function name %s already defined", f.Ident.Id))
	}
	//fmt.Println("localST created")
	f.localST = st.NewWithFather(symTable, f.Ident.String())

	errors = f.Parameters.PerformSABuild(errors, f.localST)
	errors = f.ReturnType.PerformSABuild(errors, f.localST)
	errors = f.Declarations.PerformSABuild(errors, f.localST)
	errors = f.Statements.PerformSABuild(errors, f.localST)
	symTable.InsertFunctionEntry(f.Ident.Id, types.NewFuncTy(f.Ident.Id), f.localST, f.Parameters.getParameterTypeArray(symTable), f.ReturnType.Type.GetType(symTable))
	return errors
}

func (f *Function) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = f.Parameters.TypeCheck(errors, f.localST)
	errors = f.ReturnType.TypeCheck(errors, f.localST)
	errors = f.Declarations.TypeCheck(errors, f.localST)
	errors = f.Statements.TypeCheck(errors, f.localST)

	return errors
}

func (f *Function) TranslateToILoc(funcFrag *ir.FuncFrag, symTable *st.SymbolTable) {
	//Create a funcFrag with the statements using local symbol table
	funcFrag.Label = f.Ident.Id
	var localST *st.SymbolTable
	entry, exist := symTable.Contain(f.Ident.Id)
	if exist {
		localST = entry.GetValue().LocalSymbolTable
	} else {
		panic("Fail sa")
	}
	//Assign register to parameters
	entry.GetValue().ParametersRegisterLocList = f.Parameters.GenerateRegisterList(localST)
	f.Declarations.TranslateToILoc(funcFrag, localST)
	f.Statements.TranslateToILoc(funcFrag, localST)
}

type Parameters struct {
	Token   *token.Token
	Decls   []Decl
	regList []int
}

func NewParameters(decls []Decl) *Parameters {
	return &Parameters{nil, decls, nil}
}

func (p *Parameters) TokenLiterals() string {
	if p.Token != nil {
		return p.Token.Literal
	}
	panic("Could not determine token literals for parameters")
}

func (p *Parameters) String() string {
	out := bytes.Buffer{}
	out.WriteString("(")
	for idx, decl := range p.Decls {
		if idx >= 1 {
			out.WriteString(",")
		}
		out.WriteString(decl.String())
	}
	out.WriteString(")")
	return out.String()
}

func (p *Parameters) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	/*
		Write parameters into the local symbol table
	*/
	for _, decl := range p.Decls {
		errors = decl.PerformSABuild(errors, symTable)
	}
	return errors
}

func (p *Parameters) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	//for _, decl := range p.Decls {
	//	if !CheckDeclared(decl.Type.TypeString, errors, symTable){
	//		errors= append(errors, fmt.Sprintf("struct type %s not defined",decl.Type.TypeString ))
	//	}
	//}
	return errors
}

func (p *Parameters) getParameterTypeArray(symTable *st.SymbolTable) []types.Type {
	//convert the parameter types to the array
	var paraTypeArray []types.Type
	for _, decl := range p.Decls {
		paraTypeArray = append(paraTypeArray, decl.Type.GetType(symTable))
	}
	return paraTypeArray
}

func (p *Parameters) GenerateRegisterList(localST *st.SymbolTable) []int {
	/*
		Assign register to the parameter
	*/
	RegList := []int{}
	for i := 0; i < len(p.Decls); i++ {
		paraEntry, exist := localST.Contain(p.Decls[i].Ident.Id)
		if !exist {
			panic("SA fail")
		} else {
			regNum := ir.NewRegister()
			RegList = append(RegList, regNum)
			paraEntry.GetValue().RegisterLoc = regNum
		}
	}
	return RegList
}

func CheckDeclared(typeString string, errors []string, symTable *st.SymbolTable) bool {
	//check whether t is a primitive type or has been declared
	if typeString == "bool" || typeString == "int" {
		return true
	}
	for {
		_, exist := symTable.Contain(typeString)
		if exist {
			return true
		} else {
			stFather, fatherExisit := symTable.GetFatherSymbol()
			if !fatherExisit {
				return false
			} else {
				symTable = stFather
			}
		}
	}
}

type ReturnType struct {
	Token *token.Token
	Type  *Type
}

func NewReturnType(t *Type) *ReturnType {
	return &ReturnType{nil, t}
}

func (r *ReturnType) String() string {
	out := bytes.Buffer{}
	if r.Type != nil {
		out.WriteString(r.Type.String())
	}
	return out.String()
}

func (r *ReturnType) TokenLiterals() string {
	if r.Token != nil {
		return r.Token.Literal
	}
	panic("Could not determine token literals for return statement")
}

func (r *ReturnType) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	/*
		Check whether return type has been declared
	*/
	if r.Type == nil {
		//Nil return
		return errors
	}
	CheckDeclared(r.Type.TypeString, errors, symTable)
	return errors
}

func (r *ReturnType) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

func (r *ReturnType) GetType(symTable *st.SymbolTable) types.Type {
	return r.Type.GetType(symTable)
}

type Statements struct {
	Token      *token.Token
	Statements []Statement
}

func NewStatements(stmts []Statement) *Statements {
	return &Statements{nil, stmts}
}

func (s *Statements) TokenLiterals() string {
	if s.Token != nil {
		return s.Token.Literal
	}
	panic("Could not determine token literals for the statements")
}

func (s *Statements) String() string {
	out := bytes.Buffer{}
	for _, stmt := range s.Statements {
		out.WriteString(stmt.String())
		out.WriteString("\n")
	}
	return out.String()
}

func (s *Statements) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	for _, statement := range s.Statements {
		errors = statement.TypeCheck(errors, symTable)
	}
	return errors
}

func (s *Statements) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	for _, statement := range s.Statements {
		errors = statement.PerformSABuild(errors, symTable)
	}
	return errors
}

func (s *Statements) TranslateToILoc(funcFrag *ir.FuncFrag, symTable *st.SymbolTable) {
	if symTable == nil {
		panic("Nill symboltalbe in Statements TranslateToILoc")
	}
	for _, statement := range s.Statements {
		statement.TranslateToILoc(funcFrag, symTable)
	}
}

type Statement struct {
	Token    *token.Token
	statExpr Stat
}

func NewStatement(s Stat) *Statement {
	return &Statement{nil, s}
}

func (s *Statement) TokenLiterals() string {
	return s.Token.Literal
}

func (s *Statement) String() string {
	out := bytes.Buffer{}
	out.WriteString(s.statExpr.String())
	return out.String()
}

func (s *Statement) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors1 := s.statExpr.TypeCheck(errors, symTable)
	return errors1
}

func (s *Statement) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	errors1 := s.statExpr.PerformSABuild(errors, symTable)
	return errors1
}

func (s *Statement) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	s.statExpr.TranslateToILoc(frag, table)
}

type Block struct {
	Token *token.Token
	stat  *Statements
}

func (b *Block) TokenLiteral() string {
	return b.Token.Literal
}

func NewBlock(stat *Statements) *Block {
	return &Block{nil, stat}
}

func (b *Block) TokenLiterals() string {
	if b.Token != nil {
		return b.Token.Literal
	}
	panic("Could not determine token literals for block")
}

func (b *Block) String() string {
	out := bytes.Buffer{}
	out.WriteString("{")
	out.WriteString("\n")
	out.WriteString(b.stat.String())
	out.WriteString("}")
	return out.String()
}

func (b *Block) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = b.stat.TypeCheck(errors, symTable)
	return errors
}

func (b *Block) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	errors = b.stat.PerformSABuild(errors, symTable)
	return errors
}

func (b *Block) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	b.stat.TranslateToILoc(frag, table)
}

type Assignment struct {
	Token  *token.Token
	Lvalue *LValue
	Expr   *Expression
}

func NewAssignment(lvalue *LValue, expr *Expression) *Assignment {
	return &Assignment{nil, lvalue, expr}
}

func (a *Assignment) TokenLiteral() string {
	if a.Token != nil {
		return a.Token.Literal
	}
	panic("Could not determine token literals for assignment")
}

func (a *Assignment) String() string {
	out := bytes.Buffer{}
	out.WriteString(a.Lvalue.String())
	out.WriteString("=")
	out.WriteString(a.Expr.String())
	out.WriteString(";")
	out.WriteString("\n")
	return out.String()
}

func (a *Assignment) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	//Check whether type of LValue == type of Expression
	lt := a.Lvalue.GetType(symTable)
	rt := a.Expr.GetType(symTable)
	if lt.GetName() != rt.GetName() {
		errors = append(errors, fmt.Sprintf("Assignment type error: Expected: %s, Actual: %s on line num: %d", lt, rt, a.Token.Rows))
	}
	return errors
}

func (a *Assignment) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	errors1 := a.Lvalue.PerformSABuild(errors, symTable)
	return errors1
}

func (a *Assignment) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	a.Lvalue.TranslateToILoc(frag, table)
	a.Expr.TranslateToILoc(frag, table)
	regLoc := a.Expr.RegisterLoc
	var inst ir.Instruction
	if len(a.Lvalue.Idents) == 1 {
		varName := a.Lvalue.Idents[0].Id
		if entry, exist := table.ContainGlobally(varName); entry != nil && exist {
			// global variable assignment
			inst = ir.NewStr(*regLoc, -1, -1, varName, ir.GLOBALVAR)
		} else {
			lvReg := a.Lvalue.RegisterLoc
			inst = ir.NewMov(lvReg, *regLoc, ir.AL, ir.REGISTER)
		}
	} else {
		// struct field assignment
		//get struct name
		frag.Body = append(frag.Body, inst)
		var structName string
		countFields := 0
		if stuctEntry, exist := table.Contain(a.Lvalue.Idents[0].Id); !exist {
			panic("SA fail")
		} else {
			structName = stuctEntry.GetValue().EntryType.GetName()
			for _, currField := range stuctEntry.GetValue().ParaNames {
				if a.Lvalue.Idents[1].Id == currField {
					break
				}
				countFields += 1
			}

		}
		inst = ir.NewStrRef(*regLoc, a.Lvalue.RegisterLoc, a.Lvalue.Idents[1].Id, structName, countFields)
		fmt.Println(inst)
	}
	frag.Body = append(frag.Body, inst)
}

type Read struct {
	Token *token.Token
	Ident IdentLiteral
}

func NewRead(ident IdentLiteral) *Read {
	return &Read{nil, ident}
}

func (r *Read) TokenLiteral() string {
	if r.Token != nil {
		return r.Token.Literal
	}
	panic("Could not determine token literals for read")
}

func (r *Read) String() string {
	out := bytes.Buffer{}
	out.WriteString("fmt")
	out.WriteString(".")
	out.WriteString("Scan")
	out.WriteString("(")
	out.WriteString("&")
	out.WriteString(r.Ident.Id)
	out.WriteString(")")
	out.WriteString(";")
	return out.String()
}

func (r *Read) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

func (r *Read) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	/*
		Check whether read id has been declared
	*/
	if _, exist := symTable.Contain(r.Ident.Id); exist {
		return errors
	}
	errors = append(errors, fmt.Sprintf("%s has not been declared on line num:%d ", r.Ident.Id, r.Token.Rows))
	return errors
}

func (r *Read) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	frag.Body = append(frag.Body, ir.NewRead(table.GetRegisterLoc(r.Ident.Id), r.Ident.Id))
}

type Print struct {
	Token       *token.Token
	printMethod string // "Print" | "Println"
	Ident       IdentLiteral
}

func NewPrint(printMethod string, ident IdentLiteral) *Print {
	return &Print{nil, printMethod, ident}
}

func (p *Print) TokenLiteral() string {
	if p.Token != nil {
		return p.Token.Literal
	}
	panic("Could not determine token literals for print")
}

func (p *Print) String() string {
	out := bytes.Buffer{}
	out.WriteString("fmt")
	out.WriteString(".")
	out.WriteString(p.printMethod)
	out.WriteString("(")
	out.WriteString(p.Ident.String())
	out.WriteString(")")
	out.WriteString(";")
	return out.String()
}

func (p *Print) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

func (p *Print) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// Check whether print id has been declared
	if _, exist := symTable.Contain(p.Ident.Id); exist {
		return errors
	}
	errors = append(errors, fmt.Sprintf("%s has not been declared on line num: %d", p.Ident.Id, p.Token.Rows))
	return errors
}

func (p *Print) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	reg := table.GetRegisterLoc(p.Ident.Id)
	if p.printMethod == "Print" {
		frag.Body = append(frag.Body, ir.NewPrint(reg))
	} else { //Println
		frag.Body = append(frag.Body, ir.NewPrintln(reg))
	}
}

type Conditional struct {
	Token      *token.Token
	Expr       *Expression
	Block      *Block
	ElseExists bool
	ElseBlock  *Block
}

func NewConditional(expr *Expression, block *Block, elseExists bool, elseBlock *Block) *Conditional {
	return &Conditional{nil, expr, block, elseExists, elseBlock}
}

func (c *Conditional) TokenLiteral() string {
	if c.Token != nil {
		return c.Token.Literal
	}
	panic("Could not determine token literals for conditional")
}

func (c *Conditional) String() string {
	out := bytes.Buffer{}
	out.WriteString("if")
	out.WriteString(" ")
	out.WriteString("(")
	out.WriteString(c.Expr.String())
	out.WriteString(")")
	out.WriteString(" ")
	out.WriteString(c.Block.String())
	if c.ElseExists {
		out.WriteString("else")
		out.WriteString(c.ElseBlock.String())
	}
	return out.String()
}

func (c *Conditional) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = c.Block.TypeCheck(errors, symTable)
	if c.ElseExists {
		errors = c.ElseBlock.TypeCheck(errors, symTable)
	}
	exprType := c.Expr.GetType(symTable)
	if exprType != types.BoolTySig {
		errors = append(errors, fmt.Sprintf("Conditional expression type: %s ,expected: bool on line num: %d", exprType.GetName(), c.Expr.Token.Rows))
		return errors
	}
	return errors
}

func (c *Conditional) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	errors = c.Block.PerformSABuild(errors, symTable)
	if c.ElseExists {
		errors = c.ElseBlock.PerformSABuild(errors, symTable)
	}
	return errors
}

func (c *Conditional) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	/*

	 */
	elseLabel := ir.NewLabelWithPre("else")
	doneLabel := ir.NewLabelWithPre("done")
	// conditional expression
	c.Expr.TranslateToILoc(frag, table)
	// translate cmp
	cmpInstruct := ir.NewCmp(*c.Expr.RegisterLoc, 0, ir.IMMEDIATE)
	frag.Body = append(frag.Body, cmpInstruct)
	var brFalseInst ir.Instruction
	if c.ElseExists {
		brFalseInst = ir.NewBranch(ir.EQ, elseLabel)
	} else {
		brFalseInst = ir.NewBranch(ir.EQ, doneLabel)
	}
	frag.Body = append(frag.Body, brFalseInst)
	// translate if clause
	c.Block.TranslateToILoc(frag, table)
	// translate else clause
	if c.ElseExists {
		ir.ControlFlowFrags = append(ir.ControlFlowFrags, &ir.FuncFrag{})
		ir.ControlFlowFrags[len(ir.ControlFlowFrags)-1].Label = elseLabel
		//brEndInst := ir.NewBranch(ir.AL, doneLabel)
		//elsLabelInst := ir.NewLabelStmt(elseLabel)
		//frag.Body = append(frag.Body, elsLabelInst)
		c.ElseBlock.TranslateToILoc(ir.ControlFlowFrags[len(ir.ControlFlowFrags)-1], table)
	}
	//Add finish branch
	frag.Body = append(frag.Body, ir.NewBranch(ir.AL, doneLabel))
	ir.ControlFlowFrags = append(ir.ControlFlowFrags, &ir.FuncFrag{})
	ir.ControlFlowFrags[len(ir.ControlFlowFrags)-1].Label = doneLabel
}

type Loop struct {
	Token *token.Token
	Expr  *Expression
	Block *Block
}

func NewLoop(expr *Expression, block *Block) *Loop {
	return &Loop{nil, expr, block}
}

func (p *Loop) TokenLiteral() string {
	if p.Token != nil {
		return p.Token.Literal
	}
	panic("Could not determine token literals for loop")
}

func (p *Loop) String() string {
	out := bytes.Buffer{}
	out.WriteString("for")
	out.WriteString(" ")
	out.WriteString("(")
	out.WriteString(p.Expr.String())
	out.WriteString(")")
	out.WriteString(" ")
	out.WriteString(p.Block.String())
	return out.String()
}

func (p *Loop) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = p.Block.TypeCheck(errors, symTable)
	//Check whether the expression is the bool type
	exprType := p.Expr.GetType(symTable)
	if exprType != types.BoolTySig {
		errors = append(errors, fmt.Sprintf("conditional expression is not a boolean value on line: %d", p.Token.Rows))
		return errors
	}
	return errors
}

func (p *Loop) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	errors = p.Expr.PerformSABuild(errors, symTable)
	errors = p.Block.PerformSABuild(errors, symTable)
	return errors
}

func (p *Loop) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	condLabel := ir.NewLabelWithPre("condLabel")
	bodyLabel := ir.NewLabelWithPre("loopBody")
	// b condLabel1
	frag.Body = append(frag.Body, ir.NewBranch(ir.AL, condLabel))

	// loop body
	loopFrag := ir.FuncFrag{}
	loopFrag.Label = bodyLabel
	ir.ControlFlowFrags = append(ir.ControlFlowFrags, &loopFrag)
	p.Block.TranslateToILoc(&loopFrag, table)

	// conditional expression
	conditionalFrag := ir.FuncFrag{}
	conditionalFrag.Label = condLabel
	ir.ControlFlowFrags = append(ir.ControlFlowFrags, &conditionalFrag)
	p.Expr.TranslateToILoc(&conditionalFrag, table)
	conditionalFrag.Body = append(conditionalFrag.Body, ir.NewCmp(*p.Expr.RegisterLoc, 1, ir.IMMEDIATE))
	conditionalFrag.Body = append(conditionalFrag.Body, ir.NewBranch(ir.EQ, bodyLabel))
}

type Return struct {
	Token *token.Token
	Expr  *Expression
}

func (r *Return) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	var retInst ir.Instruction
	r.Expr.TranslateToILoc(frag, table)
	if r.Expr == nil {
		retInst = ir.NewRet(-1, ir.VOID)
	} else {
		retInst = ir.NewRet(*r.Expr.RegisterLoc, ir.REGISTER)
	}
	lastFrag := ir.ControlFlowFrags[len(ir.ControlFlowFrags)-1]
	lastFrag.Body = append(lastFrag.Body, retInst)
}

func NewReturn(exprExists bool, expr *Expression) *Return {
	return &Return{nil, expr}
}

func (r *Return) TokenLiteral() string {
	if r.Token != nil {
		return r.Token.Literal
	}
	panic("Could not determine token literals for return")
}

func (r *Return) String() string {
	out := bytes.Buffer{}
	out.WriteString("return")

	if r.Expr != nil {
		out.WriteString(" ")
		out.WriteString(r.Expr.String())
		out.WriteString(" ")
	}
	out.WriteString(";")
	return out.String()
}

func (r *Return) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	//get Function return type
	if funcEntry, exist := symTable.Contain(symTable.String()); !exist {
		panic(fmt.Sprintf("Function named %s not defined on row: %d", symTable, r.Token.Rows))
	} else {
		lt := r.Expr.GetType(symTable).GetName()
		rt := funcEntry.GetValue().ReturnType.GetName()
		if lt != rt {
			errors = append(errors, fmt.Sprintf("Function named %s unmatched, expected: %s, got: %s on row: %d", symTable, lt, rt, r.Token.Rows))
		}
	}
	return errors
}

func (r *Return) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

type Invocation struct {
	Token        *token.Token
	Ident        IdentLiteral
	Args         *Arguments
	ReturnRegLoc int
}

func NewInvocation(ident IdentLiteral, args *Arguments) *Invocation {
	return &Invocation{nil, ident, args, -1}
}

func (i *Invocation) TokenLiteral() string {
	if i.Token != nil {
		return i.Token.Literal
	}
	panic("Could not determine token literals for invocation")
}

func (i *Invocation) String() string {
	out := bytes.Buffer{}
	out.WriteString(i.Ident.String())
	out.WriteString(i.Args.String())
	out.WriteString(";")
	return out.String()
}

func (i *Invocation) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = i.Args.TypeCheck(errors, symTable)
	return errors
}

func (i *Invocation) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// check whether the id is a function name
	entry, exist := symTable.Contain(i.Ident.TokenLiteral())
	if exist && entry.GetValue().EntryType != types.FuncTySig {
		errors = append(errors, fmt.Sprintf("Function named: %s is not defiend on row num: %d", i.Ident.Id, i.Ident.Token.Rows))
	}
	errors1 := i.Args.PerformSABuild(errors, symTable)
	return errors1
}

func (i *Invocation) GetType(symTable *st.SymbolTable) types.Type {
	//Return the return type of the invocation function
	f, _ := st.SymbolTableMap["global"].Contain(i.Ident.Id)
	return f.GetValue().ReturnType
}

func (invo *Invocation) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	//Handle the new case
	if invo.Ident.TokenLiteral() == "new" {
		if entry, exist := table.Contain(invo.Args.Exprs[0].Token.Literal); !exist {
			panic("fail sa")
		} else {
			fmt.Println(invo.Args.Exprs[0].Token.Literal)
			frag.Body = append(frag.Body, ir.GetNewStructInst(entry.GetValue().RegisterLoc, invo.Args.Exprs[0].Token.Literal, len(entry.GetValue().ParaNames)))
		}
		return
	}
	if invo.Ident.TokenLiteral() == "delete" {
		if entry, exist := table.Contain(invo.Args.Exprs[0].Token.Literal); !exist {
			panic("fail sa")
		} else {
			frag.Body = append(frag.Body, ir.NewDelete(entry.GetValue().RegisterLoc))
		}
		return
	}

	//Mov
	//Get function name
	funcName := table.String()
	funcEntry, exist := table.Contain(funcName)
	if !exist {
		panic("Fail sa")
	}
	frag.Body = append(frag.Body, ir.NewPush(funcEntry.GetValue().ParametersRegisterLocList, invo.Ident.Id))

	// bl
	frag.Body = append(frag.Body, ir.NewBl(invo.Ident.TokenLiteral()))

	// mov retrun result to tmp
	invo.ReturnRegLoc = ir.NewRegister()
	frag.Body = append(frag.Body, ir.NewMov(invo.ReturnRegLoc, 0, ir.MARG, ir.REGISTER))

	//pop

	frag.Body = append(frag.Body, ir.NewPop(funcEntry.GetValue().ParametersRegisterLocList, invo.Ident.Id))
}

type Arguments struct {
	Token       *token.Token
	Exprs       []Expression
	RegisterLoc *int
}

func NewArguments(exprs []Expression) *Arguments {
	return &Arguments{nil, exprs, nil}
}

func (a *Arguments) TokenLiteral() string {
	if a.Token != nil {
		return a.Token.Literal
	}
	panic("Could not determine token literals for arguments")
}

func (a *Arguments) String() string {
	out := bytes.Buffer{}
	out.WriteString("(")
	if len(a.Exprs) > 0 {
		for idx, exp := range a.Exprs {
			if idx >= 1 {
				out.WriteString(",")
			}
			out.WriteString(exp.String())
		}
	}
	out.WriteString(")")
	return out.String()
}

func (a *Arguments) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	for _, expr := range a.Exprs {
		errors = expr.TypeCheck(errors, symTable)
	}
	if funcEntry, exist := symTable.Contain(symTable.String()); !exist {
		panic(fmt.Sprintf("Function named %s not defined on row: %d", symTable, a.Token.Rows))
	} else {
		funcParaTypeList := funcEntry.GetValue().Parameters
		for idx, funcParaType := range funcParaTypeList {
			lt := funcParaType.GetName()
			rt := a.Exprs[idx].GetType(symTable).GetName()
			if lt != rt {
				errors = append(errors, fmt.Sprintf("Unmatched funcion parameter type, expected %s, got %s on line num: %d", lt, rt, a.Exprs[idx].Token.Rows))
			}
		}
	}
	return errors
}

func (a *Arguments) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	for _, expr := range a.Exprs {
		errors = expr.PerformSABuild(errors, symTable)
	}
	return errors
}

func (a *Arguments) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	for _, exp := range a.Exprs {
		exp.TranslateToILoc(frag, table)
	}
}

type LValue struct {
	Token       *token.Token
	Idents      []IdentLiteral
	RegisterLoc int
}

func NewLvalue(idents []IdentLiteral) *LValue {
	return &LValue{nil, idents, -1}
}

func (l *LValue) TokenLiteral() string {
	if l.Token != nil {
		return l.Token.Literal
	}
	panic("Could not determine token literals for lvalue")
}

func (l *LValue) String() string {
	out := bytes.Buffer{}
	for idx, id := range l.Idents {
		if idx >= 1 {
			out.WriteString(".")
		}
		out.WriteString(id.String())
	}
	return out.String()
}

func (l *LValue) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

func (l *LValue) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// check whether the declared idents exist in the structs
	//stCopy := *symTable
	//
	//var entry st.Entry
	//var exist bool
	//for idx, id := range l.Idents {
	//	varName := id.TokenLiteral()
	//	if idx>=1{
	//		fmt.Println("yo")
	//		fmt.Println(stCopy)
	//		stCopy = *entry.GetValue().LocalSymbolTable
	//
	//	}
	//	if entry, exist = stCopy.Contain(varName); !exist {
	//		errors = append(errors, fmt.Sprintf("%s has not been declared in symbol talble: %s on line num:%d",varName,symTable,id.Token.Rows))
	//		break
	//	}
	//}
	return errors
}

func (l *LValue) GetType(symTable *st.SymbolTable) types.Type {
	return l.Idents[len(l.Idents)-1].GetType(symTable)
}

func (l *LValue) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	/*
		Set the regisloc as the register loc of the first id
	*/
	if table == nil {
		panic("Nill symboltalbe in LValue")
	}
	l.Idents[0].TranslateToILoc(frag, table)
	l.RegisterLoc = l.Idents[0].RegisterLoc
	if len(l.Idents) == 1 {
		l.RegisterLoc = table.GetRegLoc(l.Idents[0].Id)
		return
	}
}

type Expression struct {
	Token *token.Token
	Left  *BoolTerm
	//RightExists bool
	Rights      []BoolTerm
	RegisterLoc *int
}

func NewExpression(l *BoolTerm, rs []BoolTerm) *Expression {
	return &Expression{nil, l, rs, nil}
}

func (p *Expression) TokenLiteral() string {
	if p.Token != nil {
		return p.Token.Literal
	}
	panic("Could not determine token literals for expression")
}

func (p *Expression) String() string {
	out := bytes.Buffer{}
	out.WriteString(p.Left.String())
	for _, boolTerm := range p.Rights {
		out.WriteString("||")
		out.WriteString(boolTerm.String())
	}
	return out.String()
}

func (p *Expression) GetType(symTable *st.SymbolTable) types.Type {
	lefType := p.Left.GetType(symTable)
	//for _, rTerm := range p.Rights {
	//	rigType := rTerm.GetType(symTable)
	//	if rigType != lefType {
	//		return types.UnknownTySig
	//	}
	//}
	return lefType
}

func (p *Expression) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	lefType := p.Left.GetType(symTable)
	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType(symTable)
		if lefType != rigType {
			errors = append(errors, fmt.Sprintf("#{p.Token.LineNum}: expected #{lefType.GetName()}, #{rigType.GetName()} found."))
			break
		}
	}

	return errors
}

func (p *Expression) PerformSABuild(errors []string, table *st.SymbolTable) []string {
	return errors
}

func (p *Expression) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	p.Left.TranslateToILoc(frag, table)
	if p.Rights == nil || len(p.Rights) == 0 {
		p.RegisterLoc = p.Left.RegisterLoc
		return
	}

	leftSource := p.Left.RegisterLoc
	for _, rTerm := range p.Rights {
		rTerm.TranslateToILoc(frag, table)
		target := ir.NewRegister()
		// in this way, OperandTy is always REGISTER
		instruction := ir.NewOr(target, *leftSource, *rTerm.RegisterLoc, ir.REGISTER)
		frag.Body = append(frag.Body, instruction)
		leftSource = &target
	}
	p.RegisterLoc = leftSource
}

type BoolTerm struct {
	Token         *token.Token
	EqualTermList []EqualTerm
	RegisterLoc   *int
}

func NewBoolTerm(rs []EqualTerm) *BoolTerm {
	return &BoolTerm{nil, rs, nil}
}

func (p *BoolTerm) TokenLiteral() string {
	if p.Token != nil {
		return p.Token.Literal
	}
	panic("Could not determine token literals for boolTerm")
}

func (p *BoolTerm) String() string {
	out := bytes.Buffer{}
	for i, equalTerm := range p.EqualTermList {
		if i >= 1 {
			out.WriteString("&&")
		}
		out.WriteString(equalTerm.String())
	}
	return out.String()
}

func (p *BoolTerm) GetType(symTable *st.SymbolTable) types.Type {
	return p.EqualTermList[0].GetType(symTable)
}

func (p *BoolTerm) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	for i, _ := range p.EqualTermList {
		if i >= 1 && p.EqualTermList[i].GetType(symTable) != p.EqualTermList[i-1].GetType(symTable) {
			errors = append(errors, fmt.Sprintf("#{p.Token.LineNum}: expected #{p. EqualTermList[i-1].GetType(symTable)} type, #{p.EqualTermList[i].GetType(symTable)} type found."))
			break
		}
	}
	return errors
}

func (p *BoolTerm) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	p.EqualTermList[0].TranslateToILoc(frag, table)
	if len(p.EqualTermList) == 1 {
		p.RegisterLoc = p.EqualTermList[0].RegisterLoc
	}

	leftSource := *p.EqualTermList[0].RegisterLoc
	for _, rTerm := range p.EqualTermList[1:] {
		rTerm.TranslateToILoc(frag, table)
		target := ir.NewRegister()
		// OperandTy is always REGISTER
		instruction := ir.NewAnd(target, leftSource, *rTerm.RegisterLoc, ir.REGISTER)
		frag.Body = append(frag.Body, instruction)
		leftSource = target
	}
	p.RegisterLoc = &leftSource
}

type EqualTerm struct {
	Token            *token.Token
	EqualOperator    []string
	RelationTermList []RelationTerm
	RegisterLoc      *int
}

func NewEqualTerm(operators []string, RelationTermList []RelationTerm) *EqualTerm {
	return &EqualTerm{nil, operators, RelationTermList, nil}
}

func (p *EqualTerm) TokenLiteral() string {
	if p.Token != nil {
		return p.Token.Literal
	}
	panic("Could not determine token literals for equalTerm")
}

func (p *EqualTerm) String() string {
	out := bytes.Buffer{}
	for i, relationTerm := range p.RelationTermList {
		if i >= 1 {
			out.WriteString(p.EqualOperator[i-1])
		}
		out.WriteString(relationTerm.String())
	}
	return out.String()
}

func (p *EqualTerm) GetType(symTable *st.SymbolTable) types.Type {
	return p.RelationTermList[0].GetType(symTable)
}

func (p *EqualTerm) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	for i, _ := range p.RelationTermList {
		if i >= 1 && p.RelationTermList[i].GetType(symTable) != p.RelationTermList[i-1].GetType(symTable) {
			errors = append(errors, fmt.Sprintf("#{p.Token.LineNum}: expected #{p. RelationTermList[i-1].GetType(symTable)} type, #{p.RelationTermList[i].GetType(symTable)} type found."))
			break
		}
	}
	return errors
}

func (p *EqualTerm) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	p.RelationTermList[0].TranslateToILoc(frag, table)
	if len(p.RelationTermList) == 1 {
		p.RegisterLoc = &p.RelationTermList[0].RegisterLoc
	}

	leftSource := p.RelationTermList[0].RegisterLoc
	for idx, rTerm := range p.RelationTermList[1:] {
		// Put into a new register the "false" value ("false" = 0) before the cmp
		rTerm.TranslateToILoc(frag, table)
		target := ir.NewRegister()
		instruction1 := ir.NewMov(target, 0, ir.AL, ir.IMMEDIATE)
		instruction2 := ir.NewCmp(leftSource, rTerm.RegisterLoc, ir.REGISTER)
		var instruction3 ir.Instruction
		if p.EqualOperator[idx] == "==" {
			instruction3 = ir.NewMov(target, 1, ir.EQ, ir.IMMEDIATE)
		} else { // "!="
			instruction3 = ir.NewMov(target, 1, ir.NE, ir.IMMEDIATE)
		}

		frag.Body = append(frag.Body, instruction1, instruction2, instruction3)
		leftSource = target
	}
	p.RegisterLoc = &leftSource
}

type RelationTerm struct {
	Token *token.Token
	Left  *SimpleTerm
	//RightExists bool
	RelationOperators []string // '>'| '<' | '<=' | '>='
	Rights            []SimpleTerm
	RegisterLoc       int
}

func NewRelationTerm(l *SimpleTerm, operators []string, rs []SimpleTerm) *RelationTerm {
	return &RelationTerm{nil, l, operators, rs, -1}
}

func (p *RelationTerm) TokenLiteral() string {
	if p.Token != nil {
		return p.Token.Literal
	}
	panic("Could not determine token literals for relationTerm")
}

func (p *RelationTerm) String() string {
	out := bytes.Buffer{}
	out.WriteString(p.Left.String())
	for i, operator := range p.RelationOperators {
		simpleTerm := p.Rights[i]
		out.WriteString(operator)
		out.WriteString(simpleTerm.String())
	}
	return out.String()
}

func (p *RelationTerm) GetType(symTable *st.SymbolTable) types.Type {
	lefType := p.Left.GetType(symTable)

	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType(symTable)
		if rigType != lefType {
			return types.UnknownTySig
		}
	}
	if len(p.Rights) != 0 {
		return types.BoolTySig
	} else {
		return lefType
	}
}

func (p *RelationTerm) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	lefType := p.Left.GetType(symTable)
	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType(symTable)
		if lefType != rigType {
			errors = append(errors, fmt.Sprintf("#{p.Token.LineNum}: expected #{lefType.GetName()}, #{rigType.GetName()} found."))
			break
		}
	}

	return errors
}

func (p *RelationTerm) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	p.Left.TranslateToILoc(frag, table)
	if p.Rights == nil || len(p.Rights) == 0 {
		p.RegisterLoc = p.Left.RegisterLoc
	}

	leftSource := p.Left.RegisterLoc
	for idx, rTerm := range p.Rights {
		rTerm.TranslateToILoc(frag, table)
		relationOperator := p.RelationOperators[idx]
		// Put into a new register the "false" value ("false" = 0) before the cmp
		target := ir.NewRegister()
		instruction1 := ir.NewMov(target, 0, ir.AL, ir.IMMEDIATE)
		instruction2 := ir.NewCmp(leftSource, rTerm.RegisterLoc, ir.REGISTER)
		var instruction3 ir.Instruction
		if relationOperator == ">" {
			instruction3 = ir.NewMov(target, 1, ir.GT, ir.IMMEDIATE)
		} else if relationOperator == "<" {
			instruction3 = ir.NewMov(target, 1, ir.LT, ir.IMMEDIATE)
		} else if relationOperator == "<=" {
			instruction3 = ir.NewMov(target, 1, ir.LE, ir.IMMEDIATE)
		} else { // >=
			instruction3 = ir.NewMov(target, 1, ir.GE, ir.IMMEDIATE)
		}

		frag.Body = append(frag.Body, instruction1, instruction2, instruction3)
		leftSource = target
	}
	p.RegisterLoc = leftSource
}

type SimpleTerm struct {
	Token *token.Token
	Left  *Term
	//RightExists bool
	SimpleTermOperators []string // '+' | '-'
	Rights              []Term
	RegisterLoc         int
}

func NewSimpleTerm(l *Term, operators []string, rs []Term) *SimpleTerm {
	return &SimpleTerm{nil, l, operators, rs, -1}
}

func (p *SimpleTerm) TokenLiteral() string {
	if p.Token != nil {
		return p.Token.Literal
	}
	panic("Could not determine token literals for simpleTerm")
}

func (p *SimpleTerm) String() string {
	out := bytes.Buffer{}
	out.WriteString(p.Left.String())
	for i, operator := range p.SimpleTermOperators {
		term := p.Rights[i]
		out.WriteString(operator)
		out.WriteString(term.String())
	}
	return out.String()
}

func (p *SimpleTerm) GetType(symTable *st.SymbolTable) types.Type {
	lefType := p.Left.GetType(symTable)

	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType(symTable)
		if rigType != lefType {
			return types.UnknownTySig
		}
	}
	if len(p.Rights) != 0 {
		return types.IntTySig
	} else {
		return lefType
	}
}

func (p *SimpleTerm) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	lefType := p.Left.GetType(symTable)
	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType(symTable)
		if lefType != rigType {
			errors = append(errors, fmt.Sprintf("#{p.Token.LineNum}: expected #{lefType.GetName()}, #{rigType.GetName()} found."))
			break
		}
	}

	return errors
}

func (p *SimpleTerm) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	p.Left.TranslateToILoc(frag, table)
	if p.Rights == nil {
		p.RegisterLoc = p.Left.RegisterLoc
		return
	}

	leftSource := p.Left.RegisterLoc
	for idx, rTerm := range p.Rights {
		rTerm.TranslateToILoc(frag, table)
		target := ir.NewRegister()
		var instruction ir.Instruction
		if p.SimpleTermOperators[idx] == "+" {
			instruction = ir.NewAdd(target, leftSource, rTerm.RegisterLoc, ir.REGISTER)
		} else { // "-"
			instruction = ir.NewSub(target, leftSource, rTerm.RegisterLoc, ir.REGISTER)
		}
		frag.Body = append(frag.Body, instruction)
		leftSource = target
	}
	p.RegisterLoc = leftSource
	//fmt.Println(p.RegisterLoc)
}

type Term struct {
	Token *token.Token
	Left  *UnaryTerm
	//RightExists bool
	TermOperators []string // '*' | '/'
	Rights        []UnaryTerm
	RegisterLoc   int
}

func NewTerm(l *UnaryTerm, operators []string, rs []UnaryTerm) *Term {
	return &Term{nil, l, operators, rs, -1}
}

func (p *Term) TokenLiteral() string {
	if p.Token != nil {
		return p.Token.Literal
	}
	panic("Could not determine token literals for term")
}

func (p *Term) String() string {
	out := bytes.Buffer{}
	out.WriteString(p.Left.String())
	for i, operator := range p.TermOperators {
		unaryTerm := p.Rights[i]
		out.WriteString(operator)
		out.WriteString(unaryTerm.String())
	}
	return out.String()
}

func (p *Term) GetType(symTable *st.SymbolTable) types.Type {
	lefType := p.Left.GetType(symTable)

	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType(symTable)
		if rigType != lefType {
			return types.UnknownTySig
		}
	}
	if len(p.Rights) != 0 {
		return types.IntTySig
	} else {
		return lefType
	}
}

func (p *Term) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	lefType := p.Left.GetType(symTable)
	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType(symTable)
		if lefType != rigType {
			errors = append(errors, fmt.Sprintf("#{p.Token.LineNum}: expected #{lefType.GetName()}, #{rigType.GetName()} found."))
			break
		}
	}

	return errors
}

func (p *Term) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	p.Left.TranslateToILoc(frag, table)
	if p.Rights == nil || len(p.Rights) == 0 {
		p.RegisterLoc = p.Left.RegisterLoc
		return
	}

	leftSource := p.Left.RegisterLoc
	for idx, rTerm := range p.Rights {
		rTerm.TranslateToILoc(frag, table)
		target := ir.NewRegister()
		var instruction ir.Instruction
		if p.TermOperators[idx] == "*" {
			instruction = ir.NewMul(target, leftSource, rTerm.RegisterLoc)
		} else { // "/"
			instruction = ir.NewDiv(target, leftSource, rTerm.RegisterLoc)
		}
		frag.Body = append(frag.Body, instruction)
		leftSource = target
	}
	p.RegisterLoc = leftSource
}

type UnaryTerm struct {
	Token         *token.Token
	UnaryOperator string // '!' | '-' | '' <- default
	SelectorTerm  *SelectorTerm
	RegisterLoc   int
}

func NewUnaryTerm(operator string, selectorTerm *SelectorTerm) *UnaryTerm {
	return &UnaryTerm{nil, operator, selectorTerm, -1}
}

func (p *UnaryTerm) TokenLiteral() string {
	if p.Token != nil {
		return p.Token.Literal
	}
	panic("Could not determine token literals for unaryTerm")
}

func (p *UnaryTerm) String() string {
	out := bytes.Buffer{}
	out.WriteString(p.UnaryOperator)
	out.WriteString(p.SelectorTerm.String())
	return out.String()
}

func (p *UnaryTerm) GetType(symTable *st.SymbolTable) types.Type {
	return p.SelectorTerm.GetType(symTable)
}

func (p *UnaryTerm) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	seleType := p.SelectorTerm.GetType(symTable)
	if p.UnaryOperator == "!" && seleType != types.BoolTySig {
		errors = append(errors, fmt.Sprintf("#{p.Token.LineNum}: expected #{types.BoolTySig.GetName()}, #{seleType.GetName()} found."))
	} else if p.UnaryOperator == "-" && seleType != types.IntTySig {
		errors = append(errors, fmt.Sprintf("#{p.Token.LineNum}: expected #{types.IntTySig.GetName()}, #{seleType.GetName()} found."))
	}
	return errors
}

func (p *UnaryTerm) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	p.SelectorTerm.TranslateToILoc(frag, table)
	if p.UnaryOperator == "" {
		p.RegisterLoc = p.SelectorTerm.RegisterLoc
	} else if p.UnaryOperator == "!" {
		target := ir.NewRegister()
		instruction := ir.NewNot(target, p.SelectorTerm.RegisterLoc, ir.REGISTER)
		frag.Body = append(frag.Body, instruction)
		p.RegisterLoc = target
	} else { // "-"
		target1 := ir.NewRegister()
		instruction1 := ir.NewMov(target1, 0, ir.AL, ir.IMMEDIATE) // mov r_x,#0
		target2 := ir.NewRegister()
		instruction2 := ir.NewSub(target2, target1, p.SelectorTerm.RegisterLoc, ir.REGISTER)
		frag.Body = append(frag.Body, instruction1, instruction2)
		p.RegisterLoc = target2
	}
}

type SelectorTerm struct {
	Token       *token.Token
	Fact        *Factor
	Idents      []IdentLiteral
	RegisterLoc int
}

func NewSelectorTerm(factor *Factor, idents []IdentLiteral) *SelectorTerm {
	return &SelectorTerm{nil, factor, idents, -1}
}

func (s *SelectorTerm) TokenLiteral() string {
	if s.Token != nil {
		return s.Token.Literal
	}
	panic("Could not determine token literals for selectorTerm")
}

func (s *SelectorTerm) String() string {
	out := bytes.Buffer{}
	out.WriteString(s.Fact.String())
	for _, id := range s.Idents {
		out.WriteString(".")
		out.WriteString(id.String())
	}
	return out.String()
}

func (s *SelectorTerm) GetType(symTable *st.SymbolTable) types.Type {
	facType := s.Fact.GetType(symTable)
	if len(s.Idents) == 0 {
		return facType
	} else {
		return types.StructTySig
	}
}

func (s *SelectorTerm) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	stCopy := *symTable

	var entry st.Entry
	var exist bool
	for _, id := range s.Idents {
		varName := id.TokenLiteral()
		if entry, exist = stCopy.Contain(varName); exist {
			errors = append(errors, fmt.Sprintf("%s has not been declared on line num:%d", varName, id.Token.Rows))
			break
		}
		stCopy = *entry.GetValue().LocalSymbolTable
	}
	return errors
}

func (s *SelectorTerm) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	//TODO: support nested structure
	s.Fact.TranslateToILoc(frag, table)
	s.RegisterLoc = s.Fact.RegisterLoc
	if s.Idents == nil {
		return
	}
	newLoc := ir.NewRegister()
	var structName string
	var paraOffSet int
	countFields := 0
	if stuctEntry, exist := table.Contain(s.Fact.Expr.TokenLiteral()); !exist {
		panic("SA fail")
	} else {
		structName = stuctEntry.GetValue().EntryType.GetName()
		for _, currField := range stuctEntry.GetValue().ParaNames {
			if s.Idents[0].Id == currField {
				break
			}
			countFields += 1
		}
	}

	frag.Body = append(frag.Body, ir.NewLoadRef(newLoc, s.RegisterLoc, s.Idents[0].Id, structName, paraOffSet))
	s.RegisterLoc = newLoc
	return
}

type Factor struct {
	Token       *token.Token
	Expr        Expr
	RegisterLoc int
}

func NewFactor(expr *Expr) *Factor {
	return &Factor{nil, *expr, -1}
}

func (p *Factor) TokenLiteral() string {
	if p.Token != nil {
		return p.Token.Literal
	}
	panic("Could not determine token literals for factor")
}

func (p *Factor) String() string {
	out := bytes.Buffer{}
	out.WriteString(p.Expr.String())
	return out.String()
}

func (p *Factor) GetType(symTable *st.SymbolTable) types.Type {
	return p.Expr.GetType(symTable)
}

func (p *Factor) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

func (p *Factor) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	p.Expr.TranslateToILoc(frag, table)
	p.RegisterLoc = p.Expr.GetRegLoc()
}

type IntLiteral struct {
	Token       *token.Token
	Value       int64
	RegisterLoc int
}

func (il *IntLiteral) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	il.RegisterLoc = ir.NewRegister()
	val := int(il.Value)
	frag.Body = append(frag.Body, ir.NewMov(il.RegisterLoc, val, ir.AL, ir.IMMEDIATE))
}

func (il *IntLiteral) GetRegLoc() int {
	return il.RegisterLoc
}

func (il *IntLiteral) TokenLiteral() string                        { return il.Token.Literal }
func (il *IntLiteral) String() string                              { return il.Token.Literal }
func (il *IntLiteral) GetType(symTable *st.SymbolTable) types.Type { return types.IntTySig }
func (il *IntLiteral) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
}
func (il *IntLiteral) GetTargetReg() int {
	return il.RegisterLoc
}

type IdentLiteral struct {
	Token       *token.Token
	Id          string //Token.Literal
	RegisterLoc int
}

func (idl *IdentLiteral) GetRegLoc() int {
	return idl.RegisterLoc
}

func (idl *IdentLiteral) TokenLiteral() string { return idl.Token.Literal }
func (idl *IdentLiteral) String() string       { return idl.Token.Literal }
func (idl *IdentLiteral) GetType(symTable *st.SymbolTable) types.Type {
	if idEntry, find := symTable.Contain(idl.TokenLiteral()); !find {
		return types.UnknownTySig
	} else {
		return idEntry.GetValue().EntryType
	}
}

func (idl *IdentLiteral) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	idlTy := idl.GetType(symTable)
	if idlTy == types.UnknownTySig {
		errors = append(errors, fmt.Sprintf("[#{idl.Token.LineNum}]: #{idl.Id} has not been defined."))
	}
	return errors
}

func (idl *IdentLiteral) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	if _, exist := table.ContainGlobally(idl.Id); exist { // if the ident is a global variable
		instruction := ir.NewLdr(idl.RegisterLoc, -1, -1, idl.Id, ir.GLOBALVAR)
		frag.Body = append(frag.Body, instruction)
	} else {
		sourceReg, exist := table.Contain(idl.Id)
		if !exist {
			panic("Try converting non-exist ident into iloc")
		}
		idl.RegisterLoc = sourceReg.GetValue().RegisterLoc
	}
}

type BoolLiteral struct {
	Token       *token.Token
	BoolValue   bool
	RegisterLoc int
}

func (bl *BoolLiteral) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	bl.RegisterLoc = ir.NewRegister()
	operand := 0
	if bl.BoolValue {
		operand = 1
	}
	frag.Body = append(frag.Body, ir.NewMov(bl.RegisterLoc, operand, ir.AL, ir.IMMEDIATE))
}

func (bl *BoolLiteral) GetRegLoc() int {
	return bl.RegisterLoc
}

func (bl *BoolLiteral) TokenLiteral() string                                         { return bl.Token.Literal }
func (bl *BoolLiteral) String() string                                               { return bl.Token.Literal }
func (bl *BoolLiteral) TypeCheck(errors []string, symTable *st.SymbolTable) []string { return errors }
func (bl *BoolLiteral) GetType(symTable *st.SymbolTable) types.Type                  { return types.BoolTySig }

type NilLiteral struct {
	Token       *token.Token
	RegisterLoc int
}

func (n *NilLiteral) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	n.RegisterLoc = ir.NewRegister()
}

func (n *NilLiteral) GetRegLoc() int {
	panic("implement me")
}

func (n *NilLiteral) TokenLiteral() string                                         { return n.Token.Literal }
func (n *NilLiteral) String() string                                               { return n.Token.Literal }
func (n *NilLiteral) GetType(symTable *st.SymbolTable) types.Type                  { return types.NilTySig }
func (n *NilLiteral) TypeCheck(errors []string, symTable *st.SymbolTable) []string { return errors }

type InvocExpr struct {
	Token       *token.Token
	Ident       IdentLiteral
	InnerArgs   *Arguments
	RegisterLoc int
}

func (ie *InvocExpr) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	//Handle the new case
	if ie.Ident.TokenLiteral() == "new" {
		if entry, exist := table.Contain(ie.InnerArgs.Exprs[0].Token.Literal); !exist {
			panic("fail sa")
		} else {
			frag.Body = append(frag.Body, ir.GetNewStructInst(entry.GetValue().RegisterLoc, ie.InnerArgs.Exprs[0].Token.Literal, len(entry.GetValue().ParaNames)))
		}
		return
	}

	//Mov
	//Get function name
	funcName := ie.Ident.Id
	_, exist := table.Contain(funcName)
	if !exist {
		panic("Fail sa")
	}
	//Get arg int list
	argIntList := []int{}
	for _, arg := range ie.InnerArgs.Exprs {
		arg.TranslateToILoc(frag, table)
		argIntList = append(argIntList, *arg.RegisterLoc)
	}
	frag.Body = append(frag.Body, ir.NewPush(argIntList, ie.Ident.Id))

	// bl
	frag.Body = append(frag.Body, ir.NewBl(ie.Ident.TokenLiteral()))

	// mov retrun result to tmp
	ie.RegisterLoc = ir.NewRegister()
	movInst := ir.NewMov(ie.RegisterLoc, 0, ir.MARG, ir.REGISTER)
	movInst.SetRetFlag()
	frag.Body = append(frag.Body, movInst)

	//pop
	frag.Body = append(frag.Body, ir.NewPop(argIntList, ie.Ident.Id))
}

func (ie *InvocExpr) GetRegLoc() int {
	return ie.RegisterLoc
}

func (ie *InvocExpr) TokenLiteral() string {
	if ie.Token != nil {
		return ie.Token.Literal
	}
	panic("Could not determine token literal for invocation expression inside Factor")
}

func (ie *InvocExpr) String() string {
	out := bytes.Buffer{}
	out.WriteString(ie.Ident.String())
	out.WriteString(ie.InnerArgs.String())
	return out.String()
}

func (ie *InvocExpr) GetType(symTable *st.SymbolTable) types.Type {
	if funcEntry, find := symTable.Contain(ie.Ident.Id); find {
		return funcEntry.GetValue().EntryType
	}
	return types.UnknownTySig
}
func (ie *InvocExpr) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// refer from Invocation.TypeCheck
	funcName := ie.Ident.TokenLiteral()
	_, find := symTable.Contain(funcName)
	if !find {
		errors = append(errors, fmt.Sprintf("#{ie.Token.LineNum}: function #{funcName} has not been defined"))
	}
	return errors
}

type PriorityExpression struct {
	Token           *token.Token
	InnerExpression *Expression
	RegisterLoc     int
}

func (pe *PriorityExpression) TranslateToILoc(frag *ir.FuncFrag, table *st.SymbolTable) {
	pe.InnerExpression.TranslateToILoc(frag, table)
	pe.RegisterLoc = *pe.InnerExpression.RegisterLoc
}

func (pe *PriorityExpression) GetRegLoc() int {
	return pe.RegisterLoc
}

func (pe *PriorityExpression) TokenLiteral() string {
	if pe.Token != nil {
		return pe.Token.Literal
	}
	panic("Could not determine token literal for expression inside Factor")
}
func (pe *PriorityExpression) String() string {
	out := bytes.Buffer{}
	out.WriteString("(")
	out.WriteString(pe.InnerExpression.String())
	out.WriteString(")")
	return out.String()
}
func (pe *PriorityExpression) GetType(symTable *st.SymbolTable) types.Type {
	return pe.InnerExpression.GetType(symTable)
}

func (pe *PriorityExpression) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = pe.InnerExpression.TypeCheck(errors, symTable)
	return errors
}

func (t *Type) GetTargetReg() int {
	return -1
}

func (rt *ReturnType) GetTargetReg() int {
	return -1
}

func (args *Arguments) GetTargetReg() int {
	return *args.RegisterLoc
}

func (lv *LValue) GetRegLoc() int {
	return lv.RegisterLoc
}

func (exp *Expression) GetRegLoc() int {
	return *exp.RegisterLoc
}

func (bt *BoolTerm) GetRegLoc() int {
	return *bt.RegisterLoc
}

func (et *EqualTerm) GetRegLoc() int {
	return *et.RegisterLoc
}

func (rt *RelationTerm) GetRegLoc() int {
	return rt.RegisterLoc
}

func (st *SimpleTerm) GetRegLoc() int {
	return st.RegisterLoc
}

func (t *Term) GetRegLoc() int {
	return t.RegisterLoc
}

func (ut *UnaryTerm) GetRegLoc() int {
	return ut.RegisterLoc
}

func (selt *SelectorTerm) GetRegLoc() int {
	return selt.RegisterLoc
}

func (f *Factor) GetRegLoc() int {
	return f.RegisterLoc
}

//func (n *NilNode) GetRegLoc() int {
//	return n.RegisterLoc
//}
