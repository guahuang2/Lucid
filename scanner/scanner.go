package scanner

import (
	"hw2/token"
)

var keywordsMap map[string]token.TokenType = map[string]token.TokenType{
	"var":    token.VAR,
	"return": token.RETURN,
	"print":  token.PRINT,
	"while":  token.WHILE,
	"for":    token.FOR,
	"if":     token.IF,
	"else":   token.ELSE,
	"or":     token.OR,
	"and":    token.AND,
	"true":   token.TRUE,
	"false":  token.FALSE,
}

type Scanner struct {
	tokenList []token.Token
	input     string
	idx       int
}

func calTokenList(input string) []token.Token {
	tokenList := []token.Token{}
	var curToken *token.Token
	size := len(input) - 1
	idx := 0
	for ; idx <= size; idx++ {
		c := input[idx]
		//skip space or newline
		if c == ' ' || c == '\n' || c == '\t' {
			continue
		}
		switch c {
		case '+':
			curToken = token.New(token.PLUS, "+")
		case '-':
			curToken = token.New(token.MINUS, "-")
		case '*':
			if nextChar(input, idx, size) == '/' {
				curToken = token.New(token.COMR, "*/")
				idx += 1
			} else {
				curToken = token.New(token.MULIPLY, "*")
			}
		case '/':
			if nextChar(input, idx, size) == '/' {
				curToken = token.New(token.COM, "//")
				idx += 1
			} else if nextChar(input, idx, size) == '*' {
				curToken = token.New(token.COML, "/*")
				idx += 1
			} else {
				curToken = token.New(token.DEVIDE, "/")
			}
		case '%':
			curToken = token.New(token.MODULUS, "%")
		case '<':
			if nextChar(input, idx, size) == '=' {
				curToken = token.New(token.LESSEQU, "<=")
				idx += 1
			} else {
				curToken = token.New(token.LESS, "<")
			}
		case '>':
			if nextChar(input, idx, size) == '=' {
				curToken = token.New(token.GREATEQU, ">=")
				idx += 1
			} else {
				curToken = token.New(token.GREATER, ">")
			}
		case '=':
			if nextChar(input, idx, size) == '=' {
				curToken = token.New(token.EQUA, "==")
				idx += 1
			} else {
				curToken = token.New(token.ASSIGN, "=")
			}
		case '!':
			if nextChar(input, idx, size) == '=' {
				curToken = token.New(token.GREATEQU, ">=")
				idx += 1
			} else {
				curToken = token.New(token.GREATER, ">")
			}
		case '(':
			curToken = token.New(token.LEFTPAR, "(")
		case ')':
			curToken = token.New(token.RIGHTBRAC, ")")
		case '{':
			curToken = token.New(token.LEFTBRAC, "{")
		case '}':
			curToken = token.New(token.RIGHTBRAC, "}")
		case ';':
			curToken = token.New(token.SEMICOLON, ";")
		case ',':
			curToken = token.New(token.PUNCTUATOR, ",")
		case '.':
			curToken = token.New(token.DOT, ".")
		case ':':
			curToken = token.New(token.COLON, ":")
		default:
			if isDigit(input[idx]) {
				num, step := getNum(input, idx, size)
				idx += step - 1
				curToken = token.New(token.INT, num)
			} else if isChar(input[idx]) {
				word, step := getWord(input, idx, size)
				idx += step - 1
				tokenType, ok := keywordsMap[word]
				if ok {
					// keyword
					curToken = token.New(tokenType, word)
				} else {
					//identifier
					curToken = token.New(token.IDENT, word)
				}
			} else {
				curToken = token.New(token.INVALID, string(c))
			}
		}

		tokenList = append(tokenList, *curToken)
	}
	if tokenList[len(tokenList)-1].Type != token.EOF {
		tokenList = append(tokenList, *token.New(token.EOF, "EOF"))
	}
	return tokenList
}

func isChar(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_'
}

func getWord(input string, idx int, size int) (string, int) {
	i := 0
	for ; i+idx <= size; i++ {
		if !(isChar(input[idx+i]) || (i >= 1 && isDigit(input[idx+i]))) {
			break
		}
	}
	return input[idx : idx+i], i
}

func isDigit(c byte) bool {
	if c >= '0' && c <= '9' {
		return true
	}
	return false
}

func getNum(input string, curIdx int, size int) (string, int) {
	i := 0
	for ; i+curIdx <= size; i++ {
		if !isDigit(input[curIdx+i]) {
			break
		}
	}
	return input[curIdx : curIdx+i], i
}

func nextChar(input string, curIdx int, size int) byte {
	if curIdx < size {
		return input[curIdx+1]
	}
	return ' '
}

func New(input string) *Scanner {
	tokenList := calTokenList(input)
	return &Scanner{tokenList: tokenList, input: input, idx: -1}
}

func (l *Scanner) NextToken() token.Token {
	l.idx += 1
	return l.tokenList[l.idx]
}
