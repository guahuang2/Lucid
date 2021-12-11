package parser

import (
	"flag"
	"fmt"
	"proj/ast"
	cc "proj/context"
	cs "proj/scanner"
	st "proj/symboltable"
	ct "proj/token"
)

type Parser struct {
	tokens          []ct.Token
	currIdx         int
	compilerContext *cc.CompilerContext
	scanner         *cs.Scanner
	successfulBuild bool
}

func New(compilerContext *cc.CompilerContext, scanner *cs.Scanner) *Parser {
	return &Parser{[]ct.Token{}, 0, compilerContext, scanner, false}
}

func (p *Parser) currToken() ct.Token {
	return p.tokens[p.currIdx]
}

func (p *Parser) nextToken() {
	if p.currIdx >= len(p.tokens)-1 {
		fmt.Errorf("syntax error:%s", "Try reading out of token")
	} else {
		p.currIdx += 1
	}
}

func (p *Parser) parseError(msg string) {
	out := flag.CommandLine.Output()
	fmt.Fprintf(out, "syntax error: #{msg}\n")
}

func (p *Parser) match(tokenType ct.TokenType) (ct.Token, bool) {
	//Check whether the current token matches the given type
	if tokenType == p.currToken().Type {
		p.nextToken()
		return p.currToken(), true
	}
	return p.currToken(), false
}

func (p *Parser) Parse() *ast.Program {
	return program(p)
}

func program(p *Parser) *ast.Program {
	pac := packageStmt(p)
	imp := importStmt(p)
	typ := types(p)
	decs := declarations(p)
	funcs := functions(p)
	if p.currToken().Type != ct.EOF {
		p.parseError(fmt.Sprintf("Expected end of file but found :#{p.currToken.Literal}"))
	}
	if p.successfulBuild == false {
		return ast.NewProgram(pac, imp, typ, decs, funcs, *st.NewSymbolTable())
	}
	return nil
}

func packageStmt(p *Parser) *ast.Package {
	var pacTok, idTok ct.Token
	var pacMatch, idMatch bool

	if pacTok, pacMatch = p.match(ct.PACKAGE); !pacMatch {
		return nil
	}
	if idTok, idMatch = p.match(ct.IDENT); !idMatch {
		return nil
	}
	if _, scMatch := p.match(ct.SEMICOLON); !scMatch {
		return nil
	}

	node := ast.NewPackage(ast.IdentLiteral{&idTok, idTok.Literal})
	node.Token = &pacTok
	return node
}

func importStmt(p *Parser) *ast.Import {
	var impTok, fmtTok ct.Token
	var impMatch, fmtMatch bool

	if impTok, impMatch = p.match(ct.IMPORT); !impMatch {
		return nil
	}
	if _, lQtdMatch := p.match(ct.DOUQUAT); !lQtdMatch {
		return nil
	}
	if fmtTok, fmtMatch = p.match(ct.FMT); !fmtMatch {
		return nil
	}
	if _, rQtdMatch := p.match(ct.DOUQUAT); !rQtdMatch {
		return nil
	}
	if _, scMatch := p.match(ct.SEMICOLON); !scMatch {
		return nil
	}

	node := ast.NewImport(ast.IdentLiteral{&fmtTok, fmtTok.Literal})
	node.Token = &impTok
	return node
}

func types(p *Parser) *ast.Types {
	var typdecs []ast.TypeDeclaration
	typdec := typeDeclaration(p)
	if typdec != nil {
		typdecs = append(typdecs, *typdec)
		for {
			if typdec = typeDeclaration(p); typdec != nil {
				typdecs = append(typdecs, *typdec)
			} else {
				break
			}
		}
		node := ast.NewTypes(typdecs)
		return node
	}
	return nil
}

func typeDeclaration(p *Parser) *ast.TypeDeclaration {
	var typTok, idTok ct.Token
	var typMac, idMac, structMac, lbrMac, rbrMac, scMac bool
	if typTok, typMac = p.match(ct.TYPE); !typMac {
		return nil
	}
	if idTok, idMac = p.match(ct.IDENT); !idMac {
		return nil
	}
	if _, structMac = p.match(ct.STRUCT); !structMac {
		return nil
	}
	if _, lbrMac = p.match(ct.LEFTBRAC); !lbrMac {
		return nil
	}
	astFields := fields(p)
	if astFields == nil {
		return nil
	}
	if _, rbrMac = p.match(ct.RIGHTBRAC); !rbrMac {
		return nil
	}
	if _, scMac = p.match(ct.SEMICOLON); !scMac {
		return nil
	}

	node := ast.NewTypeDeclaration(ast.IdentLiteral{&idTok, idTok.Literal}, astFields)
	node.Token = &typTok
	return node
}

func fields(p *Parser) *ast.Fields {
	var decls []ast.Decl

	dec := decl(p)
	if dec != nil {
		decls = append(decls, *dec)
		for {
			dec = decl(p)
			if decl != nil {
				decls = append(decls, *dec)
			} else {
				break
			}
		}
	}

	node := ast.NewFields(decls)
	return node
}

func decl(p *Parser) *ast.Decl {
	var tok ct.Token
	var match bool
	if tok, match = p.match(ct.IDENT); !match {
		return nil
	}
	astType := typeExpression(p)
	if astType == nil {
		return nil
	}
	node := ast.NewDecl(ast.IdentLiteral{&tok, tok.Literal}, astType)
	node.Token = &tok
	return node
}

func typeExpression(p *Parser) *ast.Type {
	if typeTok, match := p.match(ct.INT); match {
		return ast.NewType(typeTok.Literal)
	}
	if typeTok, match := p.match(ct.BOOL); match {
		return ast.NewType(typeTok.Literal)
	}
	if typeTok, match := p.match(ct.ASTERISK); match {
		if idTok, idMatch := p.match(ct.IDENT); idMatch {
			return ast.NewType(typeTok.Literal + idTok.Literal)
		}
		return nil
	}
	return nil
}

func declarations(p *Parser) *ast.Declarations {
	var decs []ast.Declaration
	dec := declaration(p)
	if dec != nil {
		decs = append(decs, *dec)
		for {
			if dec = declaration(p); dec != nil {
				decs = append(decs, *dec)
			} else {
				break
			}
		}
		node := ast.NewDeclarations(decs)
		return node
	}
	return nil
}

func declaration(p *Parser) *ast.Declaration {
	var varmatch, scmatch bool
	var varTok ct.Token
	if varTok, varmatch = p.match(ct.VAR); !varmatch {
		return nil
	}
	idTok := ids(p)
	if idTok == nil {
		return nil
	}
	typeTok := typeExpression(p)
	if typeTok == nil {
		return nil
	}
	if _, scmatch = p.match(ct.SEMICOLON); !scmatch {
		return nil
	}
	node := ast.NewDeclaration(idTok, typeTok)
	node.Token = &varTok
	return node
}

func ids(p *Parser) *ast.Ids {
	var idTok ct.Token
	var idMatch bool

	var ids []ast.IdentLiteral

	for {
		if idTok, idMatch = p.match(ct.IDENT); idMatch {
			ids = append(ids, ast.IdentLiteral{&idTok, idTok.Literal})
		} else {
			break
		}
	}

	node := ast.NewIds(ids)
	return node
}

func functions(p *Parser) *ast.Functions {
	var funs []ast.Function
	fun := function(p)
	if fun != nil {
		funs = append(funs, *fun)
		for {
			if fun = function(p); fun != nil {
				funs = append(funs, *fun)
			} else {
				break
			}
		}
		node := ast.NewFunctions(funs)
		return node
	}
	return nil
}

func function(p *Parser) *ast.Function {
	var funcTok, idTok ct.Token
	var funcMatch, idMatch bool
	if funcTok, funcMatch = p.match(ct.FUNC); !funcMatch {
		return nil
	}
	if idTok, idMatch = p.match(ct.IDENT); idMatch {
		return nil
	}
	paras := parameters(p)
	if paras == nil {
		return nil
	}
	retTyp := returnType(p)
	if retTyp == nil {
		return nil
	}
	if _, lbraceMatch := p.match(ct.LEFTBRAC); !lbraceMatch {
		return nil
	}
	decls := declarations(p)
	if decls == nil {
		return nil
	}
	stmts := statements(p)
	if stmts == nil {
		return nil
	}
	if _, rbraceMatch := p.match(ct.RIGHTBRAC); !rbraceMatch {
		return nil
	}

	node := ast.NewFunction(ast.IdentLiteral{&idTok, idTok.Literal}, paras, retTyp, decls, stmts, st.NewSymbolTable())
	node.Token = &funcTok
	return node
}

func parameters(p *Parser) *ast.Parameters {
	var lParenTok ct.Token
	var lParenMatch bool
	if lParenTok, lParenMatch = p.match(ct.LEFTPAR); !lParenMatch {
		return nil
	}

	var decls []ast.Decl

	for {
		if _, commaMatch := p.match(ct.PUNCTUATOR); !commaMatch {
			return nil
		}
		if decl := decl(p); decl != nil {
			decls = append(decls, *decl)
		} else {
			break
		}
	}
	if _, rParenMatch := p.match(ct.RIGHTPAR); !rParenMatch {
		return nil
	}

	node := ast.NewParameters(decls)
	node.Token = &lParenTok
	return node
}

func returnType(p *Parser) *ast.ReturnType {
	typTok := typeExpression(p)
	if typTok != nil {
		node := ast.NewReturnType(typTok)
		node.Token = typTok.Token
		return node
	} else {
		node := ast.NewReturnType(nil)
		return node
	}
}

func statements(p *Parser) *ast.Statements {
	var stmts []ast.Statement
	stmt := statement(p)
	if stmt != nil {
		stmts = append(stmts, *stmt)
		for {
			if stmt = statement(p); stmt != nil {
				stmts = append(stmts, *stmt)
			} else {
				break
			}
		}
		node := ast.NewStatements(stmts)
		return node
	}
	return nil
}

func statement(p *Parser) *ast.Statement {
	bloc := block(p)
	if bloc != nil {
		return ast.NewStatement(bloc)
	}
	assi := assignment(p)
	if assi != nil {
		return ast.NewStatement(assi)
	}
	prin := print(p)
	if prin != nil {
		return ast.NewStatement(prin)
	}
	cond := conditional(p)
	if cond != nil {
		return ast.NewStatement(cond)
	}
	loopAst := loop(p)
	if loopAst != nil {
		return ast.NewStatement(loopAst)
	}
	ret := returnStmt(p)
	if ret != nil {
		return ast.NewStatement(ret)
	}
	readAst := read(p)
	if readAst != nil {
		return ast.NewStatement(readAst)
	}
	invoc := invocation(p)
	if invoc != nil {
		return ast.NewStatement(invoc)
	}
	return nil
}

func block(p *Parser) *ast.Block {
	if lbraceTok, lbraceMatch := p.match(ct.LEFTBRAC); lbraceMatch {
		stmtsTok := statements(p)
		if stmtsTok == nil {
			return nil
		}
		if _, rbraceMatch := p.match(ct.RIGHTBRAC); rbraceMatch {
			node := ast.NewBlock(stmtsTok)
			node.Token = &lbraceTok
			return node
		}
	}
	return nil
}

func assignment(p *Parser) *ast.Assignment {
	lval := lvalue(p)
	if lval == nil {
		return nil
	}
	if _, match := p.match(ct.ASSIGN); !match {
		return nil
	}
	expr := expression(p)
	if expr == nil {
		return nil
	}
	if _, match := p.match(ct.SEMICOLON); !match {
		return nil
	}
	node := ast.NewAssignment(lval, expr)
	return node
}

func read(p *Parser) *ast.Read {
	var fmtTok, idTok ct.Token
	var fmtMatch, idMatch bool

	if fmtTok, fmtMatch = p.match(ct.FMT); !fmtMatch {
		return nil
	}
	if _, match := p.match(ct.DOT); !match {
		return nil
	}
	if _, match := p.match(ct.SCAN); !match {
		return nil
	}
	if _, match := p.match(ct.LEFTPAR); !match {
		return nil
	}
	if _, match := p.match(ct.AMPERSAND); !match {
		return nil
	}
	if idTok, idMatch = p.match(ct.IDENT); !idMatch {
		return nil
	}
	if _, match := p.match(ct.RIGHTPAR); !match {
		return nil
	}
	if _, match := p.match(ct.SEMICOLON); !match {
		return nil
	}

	node := ast.NewRead(ast.IdentLiteral{&idTok, idTok.Literal})
	node.Token = &fmtTok
	return node
}

func print(p *Parser) *ast.Print {
	var fmtTok, printTok, idTok ct.Token
	var fmtMatch, printMatch, idMatch bool

	if fmtTok, fmtMatch = p.match(ct.FMT); !fmtMatch {
		return nil
	}
	if _, match := p.match(ct.DOT); !match {
		return nil
	}
	printTok, printMatch = p.match(ct.PRINT)
	if !printMatch {
		printTok, printMatch = p.match(ct.PRINTLN)
	}
	if !printMatch {
		return nil
	}
	if _, match := p.match(ct.LEFTPAR); !match {
		return nil
	}
	if idTok, idMatch = p.match(ct.IDENT); !idMatch {
		return nil
	}
	if _, match := p.match(ct.RIGHTPAR); !match {
		return nil
	}
	if _, match := p.match(ct.SEMICOLON); !match {
		return nil
	}

	node := ast.NewPrint(printTok.Literal, ast.IdentLiteral{&idTok, idTok.Literal})
	node.Token = &fmtTok
	return node
}

func conditional(p *Parser) *ast.Conditional {
	var ifTok ct.Token
	var ifMatch bool

	if ifTok, ifMatch = p.match(ct.IF); !ifMatch {
		return nil
	}
	if _, match := p.match(ct.LEFTPAR); !match {
		return nil
	}
	expr := expression(p)
	if expr == nil {
		return nil
	}
	if _, match := p.match(ct.RIGHTPAR); !match {
		return nil
	}
	bloc := block(p)
	if bloc == nil {
		return nil
	}

	var node *ast.Conditional
	_, match := p.match(ct.ELSE)
	if match {
		elsBloc := block(p)
		node = ast.NewConditional(expr, bloc, true, elsBloc)
	} else {
		node = ast.NewConditional(expr, bloc, false, nil)
	}
	node.Token = &ifTok

	return node
}

func loop(p *Parser) *ast.Loop {
	var forTok ct.Token
	var forMatch bool

	if forTok, forMatch = p.match(ct.FOR); !forMatch {
		return nil
	}
	if _, match := p.match(ct.LEFTPAR); !match {
		return nil
	}
	expr := expression(p)
	if expr == nil {
		return nil
	}
	if _, match := p.match(ct.RIGHTPAR); !match {
		return nil
	}
	bloc := block(p)
	if bloc == nil {
		return nil
	}

	node := ast.NewLoop(expr, bloc)
	node.Token = &forTok
	return node
}

func returnStmt(p *Parser) *ast.Return {
	var retTok ct.Token
	var retMatch bool
	if retTok, retMatch = p.match(ct.RETURN); !retMatch {
		return nil
	}

	var node *ast.Return
	expr := expression(p)
	if _, match := p.match(ct.SEMICOLON); !match {
		return nil
	}
	if expr == nil {
		node = ast.NewReturn(false, nil)
	} else {
		node = ast.NewReturn(true, expr)
	}
	node.Token = &retTok

	return node
}

func invocation(p *Parser) *ast.Invocation {
	var idTok ct.Token
	var idMatch bool
	if idTok, idMatch = p.match(ct.IDENT); !idMatch {
		return nil
	}
	arg := arguments(p)
	if arg == nil {
		return nil
	}
	if _, match := p.match(ct.SEMICOLON); !match {
		return nil
	}

	node := ast.NewInvocation(ast.IdentLiteral{&idTok, idTok.Literal}, arg)
	node.Token = &idTok
	return node
}

func arguments(p *Parser) *ast.Arguments {
	var lParentok ct.Token
	var lParenMatch bool
	var exprs []ast.Expression

	if lParentok, lParenMatch = p.match(ct.LEFTPAR); !lParenMatch {
		return nil
	}

	for {
		expr := expression(p)
		if expr != nil {
			exprs = append(exprs, *expr)
		} else {
			if _, match := p.match(ct.PUNCTUATOR); !match {
				break
			} else {
				return nil
			}
		}
	}

	if _, match := p.match(ct.RIGHTPAR); !match {
		return nil
	}

	node := ast.NewArguments(exprs)
	node.Token = &lParentok

	return node
}

func lvalue(p *Parser) *ast.LValue {
	var idTok ct.Token
	var ids []ast.IdentLiteral

	for {
		if _, match := p.match(ct.DOT); !match {
			break
		}
		if id, match := p.match(ct.IDENT); match {
			ids = append(ids, ast.IdentLiteral{&id, id.Literal})
		} else {
			return nil
		}
	}
	node := ast.NewLvalue(ids)
	node.Token = &idTok
	return node
}

func expression(p *Parser) *ast.Expression {
	var bts []ast.BoolTerm
	btLeft := boolterm(p)
	if btLeft == nil {
		return nil
	}
	for {
		if _, match := p.match(ct.OR); !match {
			break
		}
		btRight := boolterm(p)
		if btRight != nil {
			bts = append(bts, *btRight)
		} else {
			return nil
		}
	}

	node := ast.NewExpression(btLeft, bts)
	node.Token = btLeft.Token
	return node
}

func boolterm(p *Parser) *ast.BoolTerm {
	var ets []ast.EqualTerm
	etLeft := equalterm(p)
	if etLeft == nil {
		return nil
	}
	for {
		if _, match := p.match(ct.AND); !match {
			break
		}
		etRight := equalterm(p)
		if etRight != nil {
			ets = append(ets, *etRight)
		} else {
			return nil
		}
	}

	node := ast.NewBoolTerm(etLeft, ets)
	node.Token = etLeft.Token
	return node
}

func equalterm(p *Parser) *ast.EqualTerm {
	var eqOps []string
	var rts []ast.RelationTerm
	var eqTok ct.Token
	var match bool
	rtLeft := relationterm(p)
	if rtLeft == nil {
		return nil
	}
	for {
		if eqTok, match = p.match(ct.EQUA); !match {
			if eqTok, match = p.match(ct.NOTEQU); !match {
				break
			}
		}
		eqOps = append(eqOps, eqTok.Literal)
		rtRight := relationterm(p)
		if rtRight != nil {
			rts = append(rts, *rtRight)
		} else {
			return nil
		}
	}

	node := ast.NewEqualTerm(rtLeft, eqOps, rts)
	node.Token = rtLeft.Token
	return node
}

func relationterm(p *Parser) *ast.RelationTerm {
	var rlOps []string
	var sts []ast.SimpleTerm
	var rlTok ct.Token
	var match bool
	stLeft := simpleterm(p)
	if stLeft == nil {
		return nil
	}
	for {
		if rlTok, match = p.match(ct.GREATER); !match {
			if rlTok, match = p.match(ct.LESS); !match {
				if rlTok, match = p.match(ct.GREATER); !match {
					if rlTok, match = p.match(ct.LESS); !match {
						break
					}
				}
			}
		}
		rlOps = append(rlOps, rlTok.Literal)
		stRight := simpleterm(p)
		if stRight != nil {
			sts = append(sts, *stRight)
		} else {
			return nil
		}
	}
	node := ast.NewRelationTerm(stLeft, rlOps, sts)
	node.Token = stLeft.Token
	return node
}

func simpleterm(p *Parser) *ast.SimpleTerm {
	var stOps []string
	var tms []ast.Term
	var stTok ct.Token
	var match bool
	termLeft := term(p)
	if termLeft == nil {
		return nil
	}
	for {
		if stTok, match = p.match(ct.PLUS); !match {
			if stTok, match = p.match(ct.MINUS); !match {
				break
			}
		}
		stOps = append(stOps, stTok.Literal)
		tmRight := term(p)
		if tmRight != nil {
			tms = append(tms, *tmRight)
		} else {
			return nil
		}
	}

	node := ast.NewSimpleTerm(termLeft, stOps, tms)
	node.Token = termLeft.Token
	return node
}

func term(p *Parser) *ast.Term {
	var tmOps []string
	var uts []ast.UnaryTerm
	var tmTok ct.Token
	var match bool
	utLeft := unaryterm(p)
	if utLeft == nil {
		return nil
	}
	for {
		if tmTok, match = p.match(ct.ASTERISK); !match {
			if tmTok, match = p.match(ct.DEVIDE); !match {
				break
			}
		}
		tmOps = append(tmOps, tmTok.Literal)
		utRight := unaryterm(p)
		if utRight != nil {
			uts = append(uts, *utRight)
		} else {
			return nil
		}
	}

	node := ast.NewTerm(utLeft, tmOps, uts)
	node.Token = utLeft.Token
	return node
}

func unaryterm(p *Parser) *ast.UnaryTerm {
	op := ""
	var uniOp ct.Token
	var match bool
	if uniOp, match = p.match(ct.NOT); match {
		op = uniOp.Literal
	} else if uniOp, match = p.match(ct.MINUS); match {
		op = uniOp.Literal
	}
	selTok := selectorterm(p)
	if selTok == nil {
		return nil
	}
	node := ast.NewUnaryTerm(op, selTok)
	node.Token = &uniOp
	return node
}

func selectorterm(p *Parser) *ast.SelectorTerm {
	var ids []ast.IdentLiteral
	var idTok ct.Token
	var match bool
	facTok := factor(p)
	if facTok == nil {
		return nil
	}
	for {
		if _, match := p.match(ct.DOT); !match {
			break
		}
		if idTok, match = p.match(ct.IDENT); !match {
			return nil
		}
		ids = append(ids, ast.IdentLiteral{&idTok, idTok.Literal})
	}

	node := ast.NewSelectorTerm(facTok, ids)
	node.Token = facTok.Token
	return node
}

func factor(p *Parser) *ast.Factor {
	var tok ct.Token
	var match bool
	var node *ast.Factor
	if tok, match = p.match(ct.LEFTPAR); match {
		expr := expression(p)
		if expr != nil {
			if _, match := p.match(ct.RIGHTPAR); match {
				node = ast.NewFactor(expr)
				node.Token = &tok
			}
		}
	}
	return node
}
