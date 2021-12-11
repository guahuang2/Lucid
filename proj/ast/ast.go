package ast

import (
	"bytes"
	"fmt"
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
	Node
	GetType(symTable *st.SymbolTable) types.Type
}

type Stat interface {
	Node
	PerformSABuild(errors []string, symTable *st.SymbolTable) []string
}

type Program struct {
	Token        *token.Token
	Package      *Package
	Import       *Import
	Types        *Types
	Declarations *Declarations
	Functions    *Functions
	st           st.SymbolTable
}

func NewProgram(pac *Package, imp *Import, typ *Types, decs *Declarations, funcs *Functions, st st.SymbolTable) *Program {
	return &Program{nil, pac, imp, typ, decs, funcs, st}
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

	p.Package.TypeCheck(errors, symTable)
	p.Import.TypeCheck(errors, symTable)
	p.Types.TypeCheck(errors, symTable)
	p.Declarations.TypeCheck(errors, symTable)
	p.Functions.TypeCheck(errors, symTable)
	return errors
}

func (p *Program) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {

	p.st = *symTable
	st.SymbolTableMap["global"] = symTable
	p.Import.PerformSABuild(errors, symTable)
	p.Types.PerformSABuild(errors, symTable)
	p.Declarations.PerformSABuild(errors, symTable)
	p.Functions.PerformSABuild(errors, symTable)
	return errors
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
	return out.String()
}

func (p *Package) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
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
	out.WriteString(i.TokenLiteral())
	out.WriteString("\"")
	out.WriteString(";")
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
	out.WriteString("{\n")
	for _, typedec := range t.typedecls {
		out.WriteString(typedec.String())
		out.WriteString("\n")
	}
	out.WriteString("}\n")
	return out.String()
}

func (t *Types) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

func (t *Types) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	//Check whether typedecleartions' id exist in global symboltable
	for _, typedecl := range t.typedecls {
		if _, ext := symTable.Contain(typedecl.Ident.Id); ext {
			errors = append(errors, fmt.Sprint("Struct name #{entry} has already been used"))
		} else {
			symTable.InsertHigherinput(typedecl.Ident.Id, types.NewStructTy(typedecl.Ident.Id), typedecl.Fields.LocalST)
		}
		typedecl.PerformSABuild(errors, symTable)
	}
	return errors
}

type TypeDeclaration struct {
	Token  *token.Token
	Ident  IdentLiteral
	Fields *Fields
}

func NewTypeDeclaration(ident IdentLiteral, fields *Fields) *TypeDeclaration {
	return &TypeDeclaration{nil, ident, fields}
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
	t.Fields.PerformSABuild(errors, symTable)
	return errors
}

type Fields struct {
	Token   *token.Token
	Decls   []Decl
	LocalST st.SymbolTable
}

func NewFields(decls []Decl) *Fields {
	return &Fields{nil, decls, st.SymbolTable{}}
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
	}
	return out.String()
}

func (f *Fields) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	for _, decl := range f.Decls {
		decl.TypeCheck(errors, symTable)
	}
	return errors
}

func (f *Fields) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	//check whether declaration's name confilict with each other
	for _, decl := range f.Decls {
		decl.PerformSABuild(errors, symTable)
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
	d.Type.TypeCheck(errors, symTable)
	return errors
}

func (d *Decl) GetType() types.Type {
	//Equivelent to Decl.Type.GetType()
	return d.Type.GetType()
}

func (d *Decl) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	if _, ext := symTable.Contain(d.Ident.Id); ext {
		errors = append(errors, fmt.Sprint("#{entry} field has already been used"))
	} else {
		symTable.Insert(d.Ident.Id, d.Type.GetType())
	}
	return errors
}

type Type struct {
	Token      *token.Token
	TypeString string
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
	out.WriteString(t.Token.Literal)
	return out.String()
}

func (t *Type) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	if t.TypeString != "int" || t.TypeString != "bool" {
		return errors
	}
	if _, exist := symTable.Contain(t.TypeString); !exist {
		errors = append(errors, fmt.Sprintf("Structured named #{t.TypeString not defined}"))
	}
	return errors
}

func (t *Type) GetType() types.Type {
	if t.TypeString == "int" {
		return types.IntTySig
	} else if t.TypeString == "bool" {
		return &types.BoolTy{}
	} else {
		return &types.StructTy{}
	}
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
	out.WriteString("{\n")
	for _, dec := range d.Declarations {
		out.WriteString("\t")
		out.WriteString(dec.String())
		out.WriteString("\n")
	}
	out.WriteString("}\n")
	return out.String()
}

func (d *Declarations) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	//Check whether the declarations already been declared in the given symbol table
	for _, decl := range d.Declarations {
		decl.PerformSABuild(errors, symTable)
	}
	return errors
}

func (d *Declarations) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	for _, decl := range d.Declarations {
		decl.TypeCheck(errors, symTable)
	}
	return errors
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
	//Check whether the declaration already been declared in the symbol table
	for _, id := range d.Ids.Idents {
		if _, ext := symTable.Contain(id.Id); ext {
			errors = append(errors, fmt.Sprint("In line #{d.Token.Rows}, #{entry} ident has already been used in line"))
		} else {
			symTable.Insert(id.Id, d.Type.GetType())
		}
	}
	return errors
}

func (d *Declaration) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
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
	out.WriteString("{\n")
	for _, fun := range funcs.functionArray {
		out.WriteString("\t")
		out.WriteString(fun.String())
		out.WriteString("\n")
	}
	out.WriteString("}\n")
	return out.String()
}

func (funcs *Functions) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	for _, function := range funcs.functionArray {
		function.PerformSABuild(errors, symTable)
	}
	return errors
}

func (funcs *Functions) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	for _, function := range funcs.functionArray {
		function.TypeCheck(errors, symTable)
	}
	return errors
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

func NewFunction(ident IdentLiteral, params *Parameters, returnType *ReturnType,
	declarations *Declarations, statements *Statements, symbolTable *st.SymbolTable) *Function {
	return &Function{nil, ident, params, returnType, declarations, statements, st.NewWithFather(symbolTable)}
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
	out.WriteString("id")
	out.WriteString(" ")
	out.WriteString(f.Parameters.String())
	out.WriteString(" ")
	out.WriteString(f.ReturnType.String())
	out.WriteString("{")
	out.WriteString(f.Declarations.String())
	out.WriteString(" ")
	out.WriteString(f.Statements.String())
	out.WriteString("}")
	return out.String()
}

func (f *Function) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	_, exist := symTable.Contain(f.Ident.Id)
	if !exist {
		errors = append(errors, fmt.Sprint("Function name #{f.Ident.Id} already defined"))
	} else {
		symTable.InsertFunctionEntry(f.Ident.Id, types.NewFuncTy(f.Ident.Id), f.localST, f.Parameters.getParameterTypeArray(), f.ReturnType.Type.GetType())
	}
	f.Parameters.PerformSABuild(errors, f.localST)
	f.ReturnType.PerformSABuild(errors, f.localST)
	f.Declarations.PerformSABuild(errors, f.localST)
	f.Statements.PerformSABuild(errors, f.localST)
	return errors
}

func (f *Function) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	f.Parameters.TypeCheck(errors, f.localST)
	f.ReturnType.TypeCheck(errors, f.localST)
	f.Declarations.TypeCheck(errors, f.localST)
	f.Statements.TypeCheck(errors, f.localST)

	return errors
}

type Parameters struct {
	Token *token.Token
	Decls []Decl
}

func NewParameters(decls []Decl) *Parameters {
	return &Parameters{nil, decls}
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
	for _, decl := range p.Decls {
		out.WriteString(",")
		out.WriteString(decl.String())
	}
	out.WriteString(")")
	return out.String()
}

func (p *Parameters) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// Write parameters into the local symboltable
	for _, decl := range p.Decls {
		decl.PerformSABuild(errors, symTable)
	}
	return errors
}

func (p *Parameters) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	for _, decl := range p.Decls {
		CheckDeclared(decl.Ident.Id, errors, symTable)
	}
	return errors
}

func (p *Parameters) getParameterTypeArray() []types.Type {
	//convert the parameter types to the array
	var paraTypeArray []types.Type
	for _, decl := range p.Decls {
		paraTypeArray = append(paraTypeArray, decl.Type.GetType())
	}
	return paraTypeArray
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
	return r.Type.TypeString
}

func (r *ReturnType) TokenLiterals() string {
	if r.Token != nil {
		return r.Token.Literal
	}
	panic("Could not determine token literals for return statement")
}

func (r *ReturnType) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	CheckDeclared(r.Type.TypeString, errors, symTable)
	return errors
}

func (r *ReturnType) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
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
	out.WriteString("{\n")
	for _, stmt := range s.Statements {
		out.WriteString("\t")
		out.WriteString(stmt.String())
		out.WriteString("\n")
	}
	out.WriteString("}\n")
	return out.String()
}

func (s *Statements) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	for _, expr := range s.Statements {
		errors = expr.TypeCheck(errors, symTable)
	}
	return errors
}

func (s *Statements) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	for _, expr := range s.Statements {
		errors = expr.TypeCheck(errors, symTable)
	}
	return errors
}

type Statement struct {
	Token    *token.Token
	statExpr Stat
}

//func NewStatement(	block *Block,
//	assignment *Assignment,
//	print *Print,
//	conditional *Conditional,
//	loop *Loop,
//	retn  *Return,
//	read *Read,
//	invocation *Invocation ,
//) *Statement {
//	return &Statement{nil, block,assignment,print,conditional,loop,retn,read,invocation}
//}

func NewStatement(s Stat) *Statement {
	return &Statement{nil, s}
}

func (s *Statement) TokenLiterals() string {
	return s.Token.Literal
}

func (s *Statement) String() string {
	out := bytes.Buffer{}
	out.WriteString(s.String())
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

type Block struct {
	Token *token.Token
	stat  *Statements
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
	out.WriteString(b.stat.String())
	out.WriteString("{")
	return out.String()
}

func (b *Block) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors1 := b.stat.TypeCheck(errors, symTable)
	return errors1
}

func (b *Block) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	errors1 := b.stat.PerformSABuild(errors, symTable)
	errors = append(errors, errors1...)
	return errors
}

type Assignment struct {
	Token  *token.Token
	Lvalue *LValue
	Expr   *Expression
}

func NewAssignment(lvalue *LValue, expr *Expression) *Assignment {
	return &Assignment{nil, lvalue, expr}
}

func (p *Assignment) TokenLiteral() string {
	if p.Token != nil {
		return p.Token.Literal
	}
	panic("Could not determine token literals for assignment")
}

func (p *Assignment) String() string {
	out := bytes.Buffer{}
	out.WriteString(p.Lvalue.String())
	out.WriteString("=")
	out.WriteString(p.Expr.String())
	out.WriteString(";")
	return out.String()
}

func (p *Assignment) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	//Check whether type of LValue == type of Expression
	lt := p.Lvalue.GetType()
	rt := p.Expr.GetType()
	if lt.GetName() != rt.GetName() {
		errors = append(errors, fmt.Sprintf("#{p.Token.LineNum} assignment type error: Expected: #{entry.GetType().GetName()}, Actual: #{exprType.GetName()}"))
	}
	return errors
}

func (p *Assignment) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	errors1 := p.Lvalue.PerformSABuild(errors, symTable)
	return errors1
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
	// Check whether read id has been declared
	if _, exist := symTable.Contain(r.Ident.Id); exist {
		return errors
	}
	errors = append(errors, fmt.Sprintf("#{p.Ident.Id} has not been declared #{p.Token.LineNum}"))
	return errors
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
	errors = append(errors, fmt.Sprintf("#{p.Ident.Id} has not been declared #{p.Token.LineNum}"))
	return errors
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
	errors1 := c.Block.TypeCheck(errors, symTable)
	errors = append(errors, errors1...)
	if c.ElseExists {
		errors2 := c.ElseBlock.TypeCheck(errors, symTable)
		errors = append(errors, errors2...)
	}
	return errors
}

func (c *Conditional) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	exprType := c.Expr.GetType() // Singleton
	if exprType != types.BoolTySig {
		errors = append(errors, fmt.Sprintf("#{p.Token.LineNum} conditional expression is not a boolean value"))
		return errors
	}
	errors1 := c.Block.PerformSABuild(errors, symTable)
	errors = append(errors, errors1...)
	if c.ElseExists {
		errors2 := c.ElseBlock.PerformSABuild(errors, symTable)
		errors = append(errors, errors2...)
	}
	return errors
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
	errors1 := p.Block.TypeCheck(errors, symTable)
	errors = append(errors, errors1...)
	return errors
}

func (p *Loop) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	//Check whether the expression is the bool type
	exprType := p.Expr.GetType()
	if exprType != types.BoolTySig {
		errors = append(errors, fmt.Sprintf("#{p.Token.LineNum} conditional expression is not a boolean value"))
		return errors
	}
	errors1 := p.Expr.PerformSABuild(errors, symTable)
	errors = append(errors, errors1...)
	errors2 := p.Block.PerformSABuild(errors, symTable)
	errors = append(errors, errors2...)
	return errors
}

type Return struct {
	Token *token.Token
	Expr  *Expression
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
	//TODO:perform return type check
	return errors
}

func (r *Return) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

type Invocation struct {
	Token *token.Token
	Ident IdentLiteral
	Args  *Arguments
}

func NewInvocation(ident IdentLiteral, args *Arguments) *Invocation {
	return &Invocation{nil, ident, args}
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
	//TODO:check
	errors1 := i.Args.TypeCheck(errors, symTable)
	return errors1
}

func (i *Invocation) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// check whether the id is a function name
	entry, exist := symTable.Contain(i.Ident.TokenLiteral())
	if exist && entry.GetValue().EntryType != types.FuncTySig {
		errors = append(errors, fmt.Sprintf("Line:#{p.Token.LineNum}, #{p.Ident.TokenLiteral()} is not a function name"))
	}
	errors1 := i.Args.PerformSABuild(errors, symTable)
	return errors1
}

func (i *Invocation) GetType() types.Type {
	//Return the return type of the invocation function
	f, _ := st.SymbolTableMap["global"].Contain(i.Ident.Id)
	return f.GetValue().ReturnType
}

type Arguments struct {
	Token *token.Token
	Exprs []Expression
}

func NewArguments(exprs []Expression) *Arguments {
	return &Arguments{nil, exprs}
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
		for _, exp := range a.Exprs {
			out.WriteString(",")
			out.WriteString(exp.String())
		}
	}
	out.WriteString(")")
	return out.String()
}

func (a *Arguments) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

func (a *Arguments) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

type LValue struct {
	Token  *token.Token
	Idents []IdentLiteral
}

func NewLvalue(idents []IdentLiteral) *LValue {
	return &LValue{nil, idents}
}

func (l *LValue) TokenLiteral() string {
	if l.Token != nil {
		return l.Token.Literal
	}
	panic("Could not determine token literals for lvalue")
}

func (l *LValue) String() string {
	out := bytes.Buffer{}
	for _, id := range l.Idents {
		out.WriteString(".")
		out.WriteString(id.String())
	}
	return out.String()
}

func (l *LValue) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

func (l *LValue) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// check whether the declared idents exist in the structs
	stCopy := *symTable

	var entry st.Entry
	var exist bool
	for _, id := range l.Idents {
		varName := id.TokenLiteral()
		if entry, exist = stCopy.Contain(varName); exist {
			errors = append(errors, fmt.Sprintf("#{varName} has not been declared #{id.Token.LineNum}"))
			break
		}
		stCopy = *entry.GetValue().LocalSymbolTable
	}
	return errors
}

func (l *LValue) GetType(symTable *st.SymbolTable) types.Type {
	return l.Idents[len(l.Idents)-1].IdType
}

type Expression struct {
	Token *token.Token
	Left  *BoolTerm
	//RightExists bool
	Rights []BoolTerm
}

func NewExpression(l *BoolTerm, rs []BoolTerm) *Expression {
	return &Expression{nil, l, rs}
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

func (p *Expression) GetType() types.Type {
	lefType := p.Left.GetType()

	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType()
		if rigType != lefType {
			return types.UnknownTySig
		}
	}
	return lefType
}

func (p *Expression) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	lefType := p.Left.GetType()
	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType()
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

type BoolTerm struct {
	Token *token.Token
	Left  *EqualTerm
	//RightExists bool
	Rights []EqualTerm
}

func NewBoolTerm(l *EqualTerm, rs []EqualTerm) *BoolTerm {
	return &BoolTerm{nil, l, rs}
}

func (p *BoolTerm) TokenLiteral() string {
	if p.Token != nil {
		return p.Token.Literal
	}
	panic("Could not determine token literals for boolTerm")
}

func (p *BoolTerm) String() string {
	out := bytes.Buffer{}
	out.WriteString(p.Left.String())
	for _, equalTerm := range p.Rights {
		out.WriteString("&&")
		out.WriteString(equalTerm.String())
	}
	return out.String()
}

func (p *BoolTerm) GetType() types.Type {
	lefType := p.Left.GetType()

	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType()
		if rigType != lefType {
			return types.UnknownTySig
		}
	}
	return lefType
}

func (p *BoolTerm) TypeCheck(errors []string, symTable st.SymbolTable) []string {
	lefType := p.Left.GetType()
	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType()
		if lefType != rigType {
			errors = append(errors, fmt.Sprintf("#{p.Token.LineNum}: expected #{lefType.GetName()}, #{rigType.GetName()} found."))
			break
		}
	}

	return errors
}

type EqualTerm struct {
	Token *token.Token
	Left  *RelationTerm
	//RightExists bool
	EqualOperator []string // '=='|'!='
	Rights        []RelationTerm
}

func NewEqualTerm(l *RelationTerm, operators []string, rs []RelationTerm) *EqualTerm {
	return &EqualTerm{nil, l, operators, rs}
}

func (p *EqualTerm) TokenLiteral() string {
	if p.Token != nil {
		return p.Token.Literal
	}
	panic("Could not determine token literals for equalTerm")
}

func (p *EqualTerm) String() string {
	out := bytes.Buffer{}
	out.WriteString(p.Left.String())
	for i, operator := range p.EqualOperator {
		relationTerm := p.Rights[i]
		out.WriteString(operator)
		out.WriteString(relationTerm.String())
	}
	return out.String()
}

func (p *EqualTerm) GetType() types.Type {
	lefType := p.Left.GetType()

	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType()
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

func (p *EqualTerm) TypeCheck(errors []string, symTable st.SymbolTable) []string {
	lefType := p.Left.GetType()
	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType()
		if lefType != rigType {
			errors = append(errors, fmt.Sprintf("#{p.Token.LineNum}: expected #{lefType.GetName()}, #{rigType.GetName()} found."))
			break
		}
	}

	return errors
}

type RelationTerm struct {
	Token *token.Token
	Left  *SimpleTerm
	//RightExists bool
	RelationOperators []string // '>'| '<' | '<=' | '>='
	Rights            []SimpleTerm
}

func NewRelationTerm(l *SimpleTerm, operators []string, rs []SimpleTerm) *RelationTerm {
	return &RelationTerm{nil, l, operators, rs}
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

func (p *RelationTerm) GetType() types.Type {
	lefType := p.Left.GetType()

	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType()
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

func (p *RelationTerm) TypeCheck(errors []string, symTable st.SymbolTable) []string {
	lefType := p.Left.GetType()
	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType()
		if lefType != rigType {
			errors = append(errors, fmt.Sprintf("#{p.Token.LineNum}: expected #{lefType.GetName()}, #{rigType.GetName()} found."))
			break
		}
	}

	return errors
}

type SimpleTerm struct {
	Token *token.Token
	Left  *Term
	//RightExists bool
	SimpleTermOperators []string // '+' | '-'
	Rights              []Term
}

func NewSimpleTerm(l *Term, operators []string, rs []Term) *SimpleTerm {
	return &SimpleTerm{nil, l, operators, rs}
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

func (p *SimpleTerm) GetType() types.Type {
	lefType := p.Left.GetType()

	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType()
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

func (p *SimpleTerm) TypeCheck(errors []string, symTable st.SymbolTable) []string {
	lefType := p.Left.GetType()
	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType()
		if lefType != rigType {
			errors = append(errors, fmt.Sprintf("#{p.Token.LineNum}: expected #{lefType.GetName()}, #{rigType.GetName()} found."))
			break
		}
	}

	return errors
}

type Term struct {
	Token *token.Token
	Left  *UnaryTerm
	//RightExists bool
	TermOperators []string // '*' | '/'
	Rights        []UnaryTerm
}

func NewTerm(l *UnaryTerm, operators []string, rs []UnaryTerm) *Term {
	return &Term{nil, l, operators, rs}
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

func (p *Term) GetType() types.Type {
	lefType := p.Left.GetType()

	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType()
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

func (p *Term) TypeCheck(errors []string, symTable st.SymbolTable) []string {
	lefType := p.Left.GetType()
	for _, rTerm := range p.Rights {
		rigType := rTerm.GetType()
		if lefType != rigType {
			errors = append(errors, fmt.Sprintf("#{p.Token.LineNum}: expected #{lefType.GetName()}, #{rigType.GetName()} found."))
			break
		}
	}

	return errors
}

type UnaryTerm struct {
	Token         *token.Token
	UnaryOperator string // '!' | '-' | '' <- default
	SelectorTerm  *SelectorTerm
}

func NewUnaryTerm(operator string, selectorTerm *SelectorTerm) *UnaryTerm {
	return &UnaryTerm{nil, operator, selectorTerm}
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

func (p *UnaryTerm) GetType() types.Type {
	return p.SelectorTerm.GetType()
}

func (p *UnaryTerm) TypeCheck(errors []string, symTable st.SymbolTable) []string {
	seleType := p.SelectorTerm.GetType()
	if p.UnaryOperator == "!" && seleType != types.BoolTySig {
		errors = append(errors, fmt.Sprintf("#{p.Token.LineNum}: expected #{types.BoolTySig.GetName()}, #{seleType.GetName()} found."))
	} else if p.UnaryOperator == "-" && seleType != types.IntTySig {
		errors = append(errors, fmt.Sprintf("#{p.Token.LineNum}: expected #{types.IntTySig.GetName()}, #{seleType.GetName()} found."))
	}
	return errors
}

type SelectorTerm struct {
	Token  *token.Token
	Fact   *Factor
	Idents []IdentLiteral
}

func NewSelectorTerm(factor *Factor, idents []IdentLiteral) *SelectorTerm {
	return &SelectorTerm{nil, factor, idents}
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

func (s *SelectorTerm) GetType() types.Type {
	facType := s.Fact.GetType()
	if len(s.Idents) == 0 {
		return facType
	} else {
		//TODO
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
			errors = append(errors, fmt.Sprintf("#{varName} has not been declared #{id.Token.LineNum}"))
			break
		}
		stCopy = *entry.GetValue().LocalSymbolTable
	}
	return errors
}

type Factor struct {
	Token *token.Token
	Expr  *Expression
}

func NewFactor(expr *Expression) *Factor {
	return &Factor{nil, expr}
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

func (p *Factor) GetType() types.Type {
	return p.Expr.GetType()
}

func (p *Factor) TypeCheck(errors []string, symTable st.SymbolTable) []string {
	return errors
}

type IntLiteral struct {
	Token *token.Token
	Value int64
}

func (il *IntLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntLiteral) String() string       { return il.Token.Literal }

type IdentLiteral struct {
	Token *token.Token
	Id    string
}

func (idl *IdentLiteral) TokenLiteral() string { return idl.Token.Literal }
func (idl *IdentLiteral) String() string       { return idl.Token.Literal }
