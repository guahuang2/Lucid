package token

type TokenType string

const (
	//keywords
	EOF    = "EOF"
	IDENT  = "IDENT"
	VAR    = "VAR"
	RETURN = "RETURN"
	PRINT  = "PRINT"
	WHILE  = "WHILE"
	FOR    = "FOR"
	IF     = "IF"
	ELSE   = "ELSE"
	NULL   = "NULL"
	DOT    = "DOT"

	//Value type
	FLOAT  = "FLOAT"
	INT    = "INT"
	STRING = "STRING"
	CHAR   = "CHAR"
	TRUE   = "TRUE"
	FALSE  = "FALSE"

	//OPERATOR

	MULIPLY = "MUTIPLY"
	DEVIDE  = "DEVIDE"
	MODULUS = "MODULUS"

	PLUS  = "PLUS"
	MINUS = "MINUS"

	MOD = "MOD"

	NOT      = "NOT"
	LESS     = "LESS"
	GREATER  = "GREATER"
	LESSEQU  = "LESSEQU"
	GREATEQU = "GREATEQU"
	EQUA     = "EQUA" //==
	NOTEQU   = "NOTEQU"

	OR  = "OR"
	AND = "AND"

	ASSIGN = "ASSIGN" //=
	COM    = "COM"
	COML   = "COML"
	COMR   = "COMR"

	//PUNCTUATOR
	SEMICOLON  = "SEMICOLON"
	PUNCTUATOR = "PUNCTUATOR"
	LEFTPAR    = "LEFTPAR"
	RIGHTPAR   = "RIGHTPAR"
	LEFTBRAC   = "LEFTBRAC"
	RIGHTBRAC  = "RIGHTBRAC"
	COLON      = "COLON"

	//ERROR
	INVALID = "INVALID"
)

type Token struct {
	Type    TokenType
	Literal string
}

func New(Type TokenType, Literal string) *Token {
	return &Token{Type: Type, Literal: Literal}
}
