package token

type TokenType string

const (
	//keywords
	EOF     = "eof"
	IDENT   = "identification"
	VAR     = "variable"
	RETURN  = "return"
	PRINT   = "print"
	FOR     = "for"
	IF      = "if"
	ELSE    = "else"
	NIL     = "nil"
	FMT     = "fmt"
	SCAN    = "Scan"
	PACKAGE = "package"
	TYPE    = "type"
	IMPORT  = "import"
	STRUCT  = "struct"
	PRINTLN = "Println"
	FUNC    = "function"

	//Value type
	INT  = "int"
	BOOL = "bool"

	//constant
	NUMBER = "number"
	TRUE   = "true"
	FALSE  = "false"

	//OPERATOR
	DOT = "dot"

	ASTERISK = "*"
	DEVIDE   = "/"
	MODULUS  = "%"

	PLUS  = "+"
	MINUS = "-"

	NOT      = "!"
	LESS     = "<"
	GREATER  = ">"
	LESSEQU  = "<="
	GREATEQU = ">="
	EQUA     = "==" //==
	NOTEQU   = "!="

	OR        = "||"
	AND       = "&&"
	AMPERSAND = "&"

	ASSIGN  = "="
	COMMENT = "//"

	//PUNCTUATOR
	SEMICOLON  = "semicolon"
	PUNCTUATOR = ","
	LEFTPAR    = "left parenthesis"
	RIGHTPAR   = "right parenthesis"
	LEFTBRAC   = "left bracket"
	RIGHTBRAC  = "right bracket"
	COLON      = "colon"
	DOUQUAT    = "double quotation"
	SIGQUAT    = "single quotation"

	//ERROR
	INVALID = "error"
)

type Token struct {
	Type    TokenType
	Literal string
	Rows    int
}

func New(Type TokenType, Literal string, Rows int) *Token {
	return &Token{Type: Type, Literal: Literal, Rows: Rows}
}
