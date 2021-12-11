package parser

import (
	"fmt"
	"proj/ast"
	cc "proj/context"
	cs "proj/scanner"
	ct "proj/token"
	"strconv"
)

type Parser struct {
	tokens          []ct.Token
	currIdx         int
	psuedoIdx       int
	compilerContext *cc.CompilerContext
	scanner         *cs.Scanner
	successfulBuild bool
}

func New(compilerContext *cc.CompilerContext, scanner *cs.Scanner) *Parser {
	/*
		Store the scanned tokens in the token list to allow back tracking in the parsing process
	*/
	p := Parser{tokens: []ct.Token{}, currIdx: 0, compilerContext: compilerContext, scanner: scanner, successfulBuild: true}
	for {
		tok, successful := p.scanner.NextToken()
		if successful {
			if tok.Type == ct.EOF {
				p.tokens = append(p.tokens, *tok)
				break
			}
			p.tokens = append(p.tokens, *tok)
		} else {
			continue
		}
	}
	p.currIdx = 0
	p.psuedoIdx = 0
	return &p
}

func (p *Parser) currPsuedoToken() ct.Token {
	return p.tokens[p.psuedoIdx]
}

func (p *Parser) currToken() ct.Token {
	return p.tokens[p.currIdx]
}

func (p *Parser) nextToken() {
	if p.currIdx >= len(p.tokens)-1 {
		fmt.Println("syntax error:Try reading out of token")
	} else {
		p.currIdx += 1
		p.psuedoIdx += 1
	}
}

func (p *Parser) parseError(msg string) {
	// out := flag.CommandLine.Output()
	panic(fmt.Sprintf("syntax error: %s\n", msg))
	p.successfulBuild = false
}

func (p *Parser) match(tokenType ct.TokenType) (ct.Token, bool) {
	//Check whether the current token matches the given type
	if tokenType == p.currToken().Type {
		curToken := p.currToken()
		p.nextToken()
		return curToken, true
	}
	return p.currToken(), false
}

func (p *Parser) PseudoMatch(tokenType ct.TokenType, rollback bool) (ct.Token, bool) {
	//Check whether the current token matches the given type
	if tokenType == p.currPsuedoToken().Type {
		curToken := p.currPsuedoToken()
		p.psuedoIdx += 1
		return curToken, true
	}
	if rollback {
		p.psuedoIdx = p.currIdx
	}
	return p.currToken(), false
}

func (p *Parser) RollForward() {
	//Check whether the current token matches the given type
	p.currIdx = p.psuedoIdx
}

func (p *Parser) expectedTypeErrorMessage(curToken ct.TokenType, expectToken ct.TokenType) string {
	return "unexpected token type error. Found: #{curToken.Type}, Expected: #{expectToken.Type}"
}

func (p *Parser) Parse() *ast.Program {
	prog := program(p)
	if prog == nil {
		p.parseError("Nil program error")
	}
	return prog
}

func program(p *Parser) *ast.Program {
	pac := packageStmt(p)
	if pac == nil {
		return nil
	}
	imp := importStmt(p)
	if imp == nil {
		return nil
	}
	tps := typesStmt(p)
	if tps == nil {
		return nil
	}
	decs := declarations(p)
	if decs == nil {
		return nil
	}
	funcs := functions(p)
	if funcs == nil {
		p.parseError(fmt.Sprintf("Nill calls"))
		return nil
	}
	if p.currToken().Type != ct.EOF {
		p.parseError(fmt.Sprintf(" Expected end of file but found:%s on line num: %d", p.currToken().Literal, p.currToken().Rows))
	}
	if p.successfulBuild {
		return ast.NewProgram(pac, imp, tps, decs, funcs)
	}
	return nil
}

func packageStmt(p *Parser) *ast.Package {
	var pac, id ct.Token
	var pacMatch, idMatch bool

	if pac, pacMatch = p.match(ct.PACKAGE); !pacMatch {
		return nil
	}
	if id, idMatch = p.match(ct.IDENT); !idMatch {
		return nil
	}
	if _, scMatch := p.match(ct.SEMICOLON); !scMatch {
		return nil
	}
	node := ast.NewPackage(ast.IdentLiteral{Token: &id, Id: id.Literal})
	node.Token = &pac
	return node
}

func importStmt(p *Parser) *ast.Import {
	var imp, imppck ct.Token
	var impMatch, imppckMatch bool

	if imp, impMatch = p.match(ct.IMPORT); !impMatch {
		return nil
	}
	if _, lQtdMatch := p.match(ct.DOUQUAT); !lQtdMatch {
		return nil
	}
	if imppck, imppckMatch = p.match(ct.FMT); !imppckMatch { //Currently only support fmt for the import package
		return nil
	}
	if _, rQtdMatch := p.match(ct.DOUQUAT); !rQtdMatch {
		return nil
	}
	if _, scMatch := p.match(ct.SEMICOLON); !scMatch {
		return nil
	}

	node := ast.NewImport(ast.IdentLiteral{&imppck, imppck.Literal, -1})
	node.Token = &imp
	return node
}

func typesStmt(p *Parser) *ast.Types {
	var typeDeclarations []ast.TypeDeclaration
	for {
		if typeDeclaration := typeDeclaration(p); typeDeclaration != nil {
			typeDeclarations = append(typeDeclarations, *typeDeclaration)
		} else {
			break
		}
	}
	return ast.NewTypes(typeDeclarations)
}

func typeDeclaration(p *Parser) *ast.TypeDeclaration {
	var typeToken, idToken ct.Token
	var typeMatch, idMactch, structMatch, lbrMactch, rightBracketMactch, semicolonMatch bool
	if typeToken, typeMatch = p.match(ct.TYPE); !typeMatch {
		return nil
	}
	if idToken, idMactch = p.match(ct.IDENT); !idMactch {
		return nil
	}
	if _, structMatch = p.match(ct.STRUCT); !structMatch {
		return nil
	}
	if _, lbrMactch = p.match(ct.LEFTBRAC); !lbrMactch {
		return nil
	}
	astFields := fields(p)
	if astFields == nil {
		return nil
	}
	if _, rightBracketMactch = p.match(ct.RIGHTBRAC); !rightBracketMactch {
		return nil
	}
	if _, semicolonMatch = p.match(ct.SEMICOLON); !semicolonMatch {
		return nil
	}

	node := ast.NewTypeDeclaration(ast.IdentLiteral{&idToken, idToken.Literal, -1}, astFields)
	node.Token = &typeToken
	return node
}

func fields(p *Parser) *ast.Fields {
	var fieldsDeclarationList = make([]ast.Decl, 0)
	for {
		fieldsDeclaration := decl(p)
		//fmt.Println(fieldsDeclaration)
		if fieldsDeclaration != nil {
			fieldsDeclarationList = append(fieldsDeclarationList, *fieldsDeclaration)
		} else {
			break
		}
		if _, match := p.match(ct.SEMICOLON); !match {
			return nil
		}
	}
	if len(fieldsDeclarationList) < 1 {
		return nil
	}
	node := ast.NewFields(fieldsDeclarationList)
	return node
}

func decl(p *Parser) *ast.Decl {
	var IDToken ct.Token
	var match bool
	if IDToken, match = p.match(ct.IDENT); !match {
		return nil
	}
	astType := typeExpression(p)
	if astType == nil {
		return nil
	}
	node := ast.NewDecl(ast.IdentLiteral{Token: &IDToken, Id: IDToken.Literal}, astType)
	node.Token = &IDToken
	return node
}

func typeExpression(p *Parser) *ast.Type {
	if typeTok, match := p.match(ct.INT); match {
		node := ast.NewType("int")
		node.Token = &typeTok
		return node
	}
	if typeTok, match := p.match(ct.BOOL); match {
		node := ast.NewType("bool")
		node.Token = &typeTok
		return node
	}
	if typeTok, match := p.match(ct.ASTERISK); match {
		if idToken, idMatch := p.match(ct.IDENT); idMatch {
			node := ast.NewType(typeTok.Literal + idToken.Literal)
			node.Token = &idToken
			return node
		}
		return nil
	}
	return nil
}

func declarations(p *Parser) *ast.Declarations {
	var declarationList []ast.Declaration
	for {
		if dec := declaration(p); dec != nil {
			declarationList = append(declarationList, *dec)
		} else {
			break
		}
	}

	return ast.NewDeclarations(declarationList)
}

func declaration(p *Parser) *ast.Declaration {
	var varMatch, semicolonMatch bool
	var varToken ct.Token
	if varToken, varMatch = p.match(ct.VAR); !varMatch {
		return nil
	}
	idToken := ids(p)
	if idToken == nil {
		return nil
	}
	typeToken := typeExpression(p)
	if typeToken == nil {
		return nil
	}
	if _, semicolonMatch = p.match(ct.SEMICOLON); !semicolonMatch {
		return nil
	}
	node := ast.NewDeclaration(idToken, typeToken)
	node.Token = &varToken
	return node
}

func ids(p *Parser) *ast.Ids {
	var ids []ast.IdentLiteral
	if idToken, idMatch := p.match(ct.IDENT); idMatch {
		ids = append(ids, ast.IdentLiteral{Token: &idToken, Id: idToken.Literal})
	}
	for {
		if _, semiMatch := p.match(ct.PUNCTUATOR); semiMatch {
			if idToken, idMatch := p.match(ct.IDENT); idMatch {
				ids = append(ids, ast.IdentLiteral{Token: &idToken, Id: idToken.Literal})
			} else {
				return nil
			}
		} else {
			break
		}
	}
	if len(ids) < 1 {
		return nil
	}
	node := ast.NewIds(ids)
	return node
}

func functions(p *Parser) *ast.Functions {
	var functionList []ast.Function

	for {
		if function := function(p); function != nil {
			functionList = append(functionList, *function)
		} else {
			break
		}
	}

	return ast.NewFunctions(functionList)
}

func function(p *Parser) *ast.Function {
	var functionToken, idToken ct.Token
	var funcMatch, idMatch bool
	if functionToken, funcMatch = p.match(ct.FUNC); !funcMatch {
		return nil
	}
	if idToken, idMatch = p.match(ct.IDENT); !idMatch {
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
	node := ast.NewFunction(ast.IdentLiteral{Token: &idToken, Id: idToken.Literal}, paras, retTyp, decls, stmts)
	node.Token = &functionToken
	return node
}

func parameters(p *Parser) *ast.Parameters {
	var declarationList []ast.Decl
	var leftParenToken ct.Token
	var leftParenMatch bool

	if leftParenToken, leftParenMatch = p.match(ct.LEFTPAR); !leftParenMatch {
		return nil
	}

	declFirst := decl(p)
	if declFirst != nil {
		declarationList = append(declarationList, *declFirst)
		for {
			if _, commaMatch := p.match(ct.PUNCTUATOR); !commaMatch {
				break
			}
			if decl := decl(p); decl != nil {
				declarationList = append(declarationList, *decl)
			} else {
				return nil
			}
		}
	}

	if _, rParenMatch := p.match(ct.RIGHTPAR); !rParenMatch {
		return nil
	}

	node := ast.NewParameters(declarationList)
	node.Token = &leftParenToken
	return node
}

func returnType(p *Parser) *ast.ReturnType {
	typTok := typeExpression(p)
	if typTok != nil {
		node := ast.NewReturnType(typTok)
		node.Token = typTok.Token
		return node
	} else {
		return ast.NewReturnType(ast.NewType(""))
	}
}

func statements(p *Parser) *ast.Statements {
	var statementsList []ast.Statement

	for {
		if stmt := statement(p); stmt != nil {
			statementsList = append(statementsList, *stmt)
		} else {
			break
		}
	}
	return ast.NewStatements(statementsList)
}

func statement(p *Parser) *ast.Statement {
	blck := block(p)
	if blck != nil {
		return ast.NewStatement(blck)
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
	var leftBraceketToken ct.Token
	var leftBraceketMatch bool
	if leftBraceketToken, leftBraceketMatch = p.match(ct.LEFTBRAC); !leftBraceketMatch {
		return nil
	}
	stmts := statements(p)
	if stmts == nil {
		return nil
	}
	if _, rightBraceketMatch := p.match(ct.RIGHTBRAC); !rightBraceketMatch {
		return nil
	}
	blockExpr := ast.NewBlock(stmts)
	blockExpr.Token = &leftBraceketToken
	return blockExpr
}

func assignment(p *Parser) *ast.Assignment {
	leftVal := lvalue(p)
	if leftVal == nil {
		return nil
	}
	if _, match := p.PseudoMatch(ct.ASSIGN, true); !match {
		return nil
	}
	p.RollForward()
	expr := expression(p)
	if expr == nil {
		return nil
	}
	if _, match := p.match(ct.SEMICOLON); !match {
		return nil
	}
	node := ast.NewAssignment(leftVal, expr)
	return node
}

func read(p *Parser) *ast.Read {
	var fmtTok, idToken ct.Token
	var fmtMatch, idMatch bool

	if fmtTok, fmtMatch = p.PseudoMatch(ct.FMT, true); !fmtMatch {
		return nil
	}
	if _, match := p.PseudoMatch(ct.DOT, true); !match {
		return nil
	}
	if _, match := p.PseudoMatch(ct.SCAN, true); !match {
		return nil
	}
	if _, match := p.PseudoMatch(ct.LEFTPAR, true); !match {
		return nil
	}
	if _, match := p.PseudoMatch(ct.AMPERSAND, true); !match {
		return nil
	}
	if idToken, idMatch = p.PseudoMatch(ct.IDENT, true); !idMatch {
		return nil
	}
	if _, match := p.PseudoMatch(ct.RIGHTPAR, true); !match {
		return nil
	}
	if _, match := p.PseudoMatch(ct.SEMICOLON, true); !match {
		return nil
	}
	p.RollForward()
	node := ast.NewRead(ast.IdentLiteral{Token: &idToken, Id: idToken.Literal})
	node.Token = &fmtTok
	return node
}

func print(p *Parser) *ast.Print {
	var fmtToken, printToken, idToken ct.Token
	var fmtMatch, printMatch, idMatch bool

	if fmtToken, fmtMatch = p.PseudoMatch(ct.FMT, true); !fmtMatch {
		return nil
	}
	if _, match := p.PseudoMatch(ct.DOT, true); !match {
		return nil
	}
	printToken, printMatch = p.PseudoMatch(ct.PRINT, false)
	if !printMatch {
		printToken, printMatch = p.PseudoMatch(ct.PRINTLN, true)
	}
	if !printMatch {
		return nil
	}
	if _, match := p.PseudoMatch(ct.LEFTPAR, true); !match {
		return nil
	}
	if idToken, idMatch = p.PseudoMatch(ct.IDENT, true); !idMatch {
		return nil
	}
	if _, match := p.PseudoMatch(ct.RIGHTPAR, true); !match {
		return nil
	}
	if _, match := p.PseudoMatch(ct.SEMICOLON, true); !match {
		return nil
	}

	p.RollForward()
	node := ast.NewPrint(printToken.Literal, ast.IdentLiteral{Token: &idToken, Id: idToken.Literal, RegisterLoc: -1})
	node.Token = &fmtToken
	return node
}

func conditional(p *Parser) *ast.Conditional {
	var ifToken ct.Token
	var ifMatch bool

	if ifToken, ifMatch = p.match(ct.IF); !ifMatch {
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
	node.Token = &ifToken

	return node
}

func loop(p *Parser) *ast.Loop {
	var forToken ct.Token
	var forMatch bool

	if forToken, forMatch = p.match(ct.FOR); !forMatch {
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
	node.Token = &forToken
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
	var idToken ct.Token
	var idMatch bool
	if idToken, idMatch = p.match(ct.IDENT); !idMatch {
		return nil
	}
	arg := arguments(p)
	if arg == nil {
		return nil
	}
	if _, match := p.match(ct.SEMICOLON); !match {
		return nil
	}

	node := ast.NewInvocation(ast.IdentLiteral{Token: &idToken, Id: idToken.Literal}, arg)
	node.Token = &idToken
	return node
}

func arguments(p *Parser) *ast.Arguments {
	var lParentoken ct.Token
	var lParenMatch bool
	var exprs []ast.Expression

	if lParentoken, lParenMatch = p.match(ct.LEFTPAR); !lParenMatch {
		return nil
	}
	expr := expression(p)
	if expr != nil {
		exprs = append(exprs, *expr)
	}
	for {
		if _, match := p.match(ct.PUNCTUATOR); !match {
			break
		}
		expr := expression(p)
		if expr != nil {
			exprs = append(exprs, *expr)
		} else {
			return nil
		}
	}

	if _, match := p.match(ct.RIGHTPAR); !match {
		return nil
	}

	node := ast.NewArguments(exprs)
	node.Token = &lParentoken
	return node
}

func lvalue(p *Parser) *ast.LValue {
	var idToken ct.Token
	var idList []ast.IdentLiteral
	if id, match := p.PseudoMatch(ct.IDENT, true); match {
		idList = append(idList, ast.IdentLiteral{Token: &id, Id: id.Literal})
	} else {
		return nil
	}
	for {
		if _, match := p.PseudoMatch(ct.DOT, false); !match {
			break
		}
		if id, match := p.PseudoMatch(ct.IDENT, true); match {
			idList = append(idList, ast.IdentLiteral{Token: &id, Id: id.Literal})
		} else {
			return nil
		}
	}
	node := ast.NewLvalue(idList)
	node.Token = &idToken
	return node
}

func expression(p *Parser) *ast.Expression {
	var bts []ast.BoolTerm
	btLeft := boolTerm(p)
	if btLeft == nil {
		return nil
	}
	for {
		if _, match := p.match(ct.OR); !match {
			break
		}
		btRight := boolTerm(p)
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

func boolTerm(p *Parser) *ast.BoolTerm {
	var ets []ast.EqualTerm
	feq := equalTerm(p)
	if feq == nil {
		return nil
	} else {
		ets = append(ets, *feq)
	}
	for {
		if _, match := p.match(ct.AND); !match {
			break
		}
		eqt := equalTerm(p)
		if equalTerm != nil {
			ets = append(ets, *eqt)
		} else {
			return nil
		}
	}

	node := ast.NewBoolTerm(ets)
	node.Token = feq.Token
	return node
}

func equalTerm(p *Parser) *ast.EqualTerm {
	var eqOps []string
	var relationTermList []ast.RelationTerm
	frt := relationTerm(p)
	if frt == nil {
		return nil
	} else {
		relationTermList = append(relationTermList, *frt)
	}
	var eqToken ct.Token
	var match bool
	for {
		if eqToken, match = p.match(ct.EQUA); !match {
			if eqToken, match = p.match(ct.NOTEQU); !match {
				break
			}
		}
		eqOps = append(eqOps, eqToken.Literal)
		rt := relationTerm(p)
		if rt != nil {
			relationTermList = append(relationTermList, *rt)
		} else {
			return nil
		}
	}

	node := ast.NewEqualTerm(eqOps, relationTermList)
	node.Token = frt.Token
	return node

}

func relationTerm(p *Parser) *ast.RelationTerm {
	var rlOps []string
	var sts []ast.SimpleTerm
	var rlTok ct.Token
	var match bool
	stLeft := simpleTerm(p)
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
		stRight := simpleTerm(p)
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

func simpleTerm(p *Parser) *ast.SimpleTerm {
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
	utLeft := unaryTerm(p)
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
		utRight := unaryTerm(p)
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

func unaryTerm(p *Parser) *ast.UnaryTerm {
	op := ""
	var uniOp ct.Token
	var match bool
	if uniOp, match = p.match(ct.NOT); match {
		op = uniOp.Literal
	} else if uniOp, match = p.match(ct.MINUS); match {
		op = uniOp.Literal
	}
	selTok := selectorTerm(p)
	if selTok == nil {
		return nil
	}
	node := ast.NewUnaryTerm(op, selTok)
	node.Token = &uniOp
	return node
}

func selectorTerm(p *Parser) *ast.SelectorTerm {
	var ids []ast.IdentLiteral
	var idToken ct.Token
	var match bool
	facTok := factor(p)
	if facTok == nil {
		return nil
	}
	for {
		if _, match := p.match(ct.DOT); !match {
			break
		}
		if idToken, match = p.match(ct.IDENT); !match {
			return nil
		}
		ids = append(ids, ast.IdentLiteral{Token: &idToken, Id: idToken.Literal})
	}

	node := ast.NewSelectorTerm(facTok, ids)
	node.Token = facTok.Token
	return node
}

func factor(p *Parser) *ast.Factor {
	var node ast.Expr

	if numTok, match := p.match(ct.NUMBER); match {
		val, _ := strconv.ParseInt(numTok.Literal, 10, 64)
		node = &ast.IntLiteral{Token: &numTok, Value: val, RegisterLoc: -1}
	} else if truTok, match := p.match(ct.TRUE); match {
		node = &ast.BoolLiteral{Token: &truTok, BoolValue: true}
	} else if flsTok, match := p.match(ct.FALSE); match {
		node = &ast.BoolLiteral{Token: &flsTok, BoolValue: false}
	} else if nilTok, match := p.match(ct.NIL); match {
		node = &ast.NilLiteral{Token: &nilTok}
	} else if identTok, match := p.match(ct.IDENT); match {
		//" 'id' [Arguments] "
		argu := arguments(p)
		idl := &ast.IdentLiteral{Token: &identTok, Id: identTok.Literal}
		if argu == nil {
			node = idl
		} else {
			node = &ast.InvocExpr{Token: &identTok, Ident: *idl, InnerArgs: argu}
		}
	} else if lpTok, match := p.match(ct.LEFTPAR); match {
		//"'(' Expression ')'"
		expr := expression(p)
		if expr != nil {
			if _, match := p.match(ct.RIGHTPAR); match {
				node = &ast.PriorityExpression{Token: &lpTok, InnerExpression: expr}
			}
		}
	}
	if node != nil {
		return ast.NewFactor(&node)
	} else {
		return nil
	}
}
