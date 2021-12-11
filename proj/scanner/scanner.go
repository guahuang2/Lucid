package scanner

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"proj/context"
	"proj/token"
)

var keywordsMap map[string]token.TokenType = map[string]token.TokenType{
	"var":     token.VAR,
	"return":  token.RETURN,
	"for":     token.FOR,
	"if":      token.IF,
	"else":    token.ELSE,
	"true":    token.TRUE,
	"false":   token.FALSE,
	"id":      token.IDENT,
	"Print":   token.PRINT,
	"nil":     token.NIL,
	"fmt":     token.FMT,
	"Scan":    token.SCAN,
	"package": token.PACKAGE,
	"type":    token.TYPE,
	"import":  token.IMPORT,
	"struct":  token.STRUCT,
	"Println": token.PRINTLN,
	"int":     token.INT,
	"bool":    token.BOOL,
	"func":    token.FUNC,
}

func calTokenList(l *Scanner, input string) []token.Token {
	tokenList := []token.Token{}
	var curToken *token.Token
	size := len(input) - 1
	idx := 0

	for ; idx <= size; idx++ {
		c := input[idx]

		//skip space, tab and newline
		if c == ' ' || c == '\n' || c == '\t' || c == '\r' {
			if c == '\n' {
				l.curRow += 1
				l.commentLine = false
			}
			continue
		}
		//fmt.Println(strconv.Itoa(int(c)))
		switch c {
		case '+':
			curToken = token.New(token.PLUS, "+", l.curRow)
		case '-':
			curToken = token.New(token.MINUS, "-", l.curRow)
		case '*':
			curToken = token.New(token.ASTERISK, "*", l.curRow)

		case '/':
			if nextChar(input, idx, size) == '/' {
				l.commentLine = true
				curToken = token.New(token.COMMENT, "//", l.curRow)
				idx += 1
			} else {
				curToken = token.New(token.DEVIDE, "/", l.curRow)
			}
		case '&':
			if nextChar(input, idx, size) == '&' {
				curToken = token.New(token.AND, "&&", l.curRow)
				idx += 1
			} else {
				curToken = token.New(token.AMPERSAND, "&", l.curRow)
			}
		case '|':
			if nextChar(input, idx, size) == '|' {
				curToken = token.New(token.OR, "||", l.curRow)
				idx += 1
			} else {
				curToken = token.New(token.INVALID, "|", l.curRow)
			}
		case '%':
			curToken = token.New(token.MODULUS, "%", l.curRow)
		case '<':
			if nextChar(input, idx, size) == '=' {
				curToken = token.New(token.LESSEQU, "<=", l.curRow)
				idx += 1
			} else {
				curToken = token.New(token.LESS, "<", l.curRow)
			}
		case '>':
			if nextChar(input, idx, size) == '=' {
				curToken = token.New(token.GREATEQU, ">=", l.curRow)
				idx += 1
			} else {
				curToken = token.New(token.GREATER, ">", l.curRow)
			}
		case '=':
			if nextChar(input, idx, size) == '=' {
				curToken = token.New(token.EQUA, "==", l.curRow)
				idx += 1
			} else {
				curToken = token.New(token.ASSIGN, "=", l.curRow)
			}
		case '!':
			if nextChar(input, idx, size) == '=' {
				curToken = token.New(token.NOTEQU, "!=", l.curRow)
				idx += 1
			} else {
				curToken = token.New(token.GREATER, ">", l.curRow)
			}
		case '(':
			curToken = token.New(token.LEFTPAR, "(", l.curRow)
		case ')':
			curToken = token.New(token.RIGHTPAR, ")", l.curRow)
		case '{':
			curToken = token.New(token.LEFTBRAC, "{", l.curRow)
		case '}':
			curToken = token.New(token.RIGHTBRAC, "}", l.curRow)
		case ';':
			curToken = token.New(token.SEMICOLON, ";", l.curRow)
		case ',':
			curToken = token.New(token.PUNCTUATOR, ",", l.curRow)
		case '.':
			curToken = token.New(token.DOT, ".", l.curRow)
		case ':':
			curToken = token.New(token.COLON, ":", l.curRow)
		case '"':
			curToken = token.New(token.DOUQUAT, "\"", l.curRow)
		case '\'':
			curToken = token.New(token.SIGQUAT, "'", l.curRow)
		default:
			if isDigit(input[idx]) {
				num, step := getNum(input, idx, size)
				idx += step - 1
				curToken = token.New(token.NUMBER, num, l.curRow)
			} else if isChar(input[idx]) {
				word, step := getWord(input, idx, size)
				idx += step - 1
				tokenType, ok := keywordsMap[word]
				if ok {
					// keyword
					curToken = token.New(tokenType, word, l.curRow)
				} else {
					//identifier
					curToken = token.New(token.IDENT, word, l.curRow)
				}
			} else {
				curToken = token.New(token.INVALID, string(c), l.curRow)
			}
		}
		if !l.commentLine {
			tokenList = append(tokenList, *curToken)
		}
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

func check(e error) {
	if e != nil {
		panic(e.Error())
	}
}

type Scanner struct {
	finalTokenList []token.Token
	curTokenliST   []token.Token
	reader         *bufio.Reader
	idx            int
	curRow         int
	commentLine    bool
}

func New(inputContext *context.CompilerContext) *Scanner {
	/*Create a Scanner according to the given context*/
	inputFile, err := os.Open(inputContext.SourcePath())
	check(err)

	reader := bufio.NewReader(inputFile)
	return &Scanner{finalTokenList: make([]token.Token, 0),
		curTokenliST: make([]token.Token, 0),
		reader:       reader,
		idx:          0,
		curRow:       1,
		commentLine:  false,
	}
}

func (l *Scanner) NextToken() (*token.Token, bool) {
	if l.idx+1 > len(l.finalTokenList) { //Check whether the pointer has exceeded the end of the finalTokenList
		inputString, err := l.reader.ReadString(' ')
		l.finalTokenList = append(l.finalTokenList, calTokenList(l, inputString)...)
		if err != nil {
			if err == io.EOF {
				l.finalTokenList = append(l.finalTokenList, *token.New(token.EOF, "eof", l.curRow))
			} else {
				check(err)
			}
		}
	}
	if l.idx+1 <= len(l.finalTokenList) {
		l.idx += 1
		return &l.finalTokenList[l.idx-1], true
	} else {
		return nil, false
	}
}

func PrintToken(t token.Token) {
	fmt.Printf("|%-20v|%-20v|%-20v|\n", t.Type, t.Literal, t.Rows)
}

func (l *Scanner) PrintAllTokens() {
	fmt.Println("Start printing tokens")
	fmt.Printf("|%-20v|%-20v|%-20v|\n", "Token Type", "Token Literal", "Rows")
	for {
		if nextToken, readSuccess := l.NextToken(); readSuccess {
			PrintToken(*nextToken)
			if nextToken.Type == token.EOF {
				return
			}
		} else {
			continue
		}

	}
	fmt.Println("Finish printing tokens")
}
