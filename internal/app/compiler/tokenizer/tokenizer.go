package tokenizer

import (
	"strings"

	"github.com/tocoteron/9cc-go/internal/app/compiler/io"
)

type TokenKind int

const (
	TOKEN_RESERVED TokenKind = iota
	TOKEN_NUM
	TOKEN_EOF
)

type Token struct {
	kind TokenKind
	next *Token
	val  int
	str  string
}

var CurrentToken *Token

func Consume(op byte) bool {
	if CurrentToken.kind != TOKEN_RESERVED || CurrentToken.str[0] != op {
		return false
	}

	CurrentToken = CurrentToken.next

	return true
}

func Expect(op byte) {
	if CurrentToken.kind != TOKEN_RESERVED || CurrentToken.str[0] != op {
		io.ErrorAt(CurrentToken.str, "It is not '%c'", op)
	}

	CurrentToken = CurrentToken.next
}

func ExpectNumber() int {
	if CurrentToken.kind != TOKEN_NUM {
		io.ErrorAt(CurrentToken.str, "It is not a number")
	}

	val := CurrentToken.val

	CurrentToken = CurrentToken.next

	return val
}

func atEOF() bool {
	return CurrentToken.kind == TOKEN_EOF
}

func newToken(cur *Token, kind TokenKind, str string) *Token {
	token := &Token{}
	token.kind = kind
	token.str = str
	cur.next = token

	return token
}

func Tokenize(s string) *Token {
	var head Token
	head.next = nil

	cur := &head

	for s != "" {
		if s[0] == ' ' {
			s = s[1:]
			continue
		}

		if strings.IndexByte("+-*/()", s[0]) != -1 {
			cur = newToken(cur, TOKEN_RESERVED, s)
			s = s[1:]
			continue
		}

		if s[0] >= '0' && s[0] <= '9' {
			cur = newToken(cur, TOKEN_NUM, s)
			cur.val, s = strToInt(s)
			continue
		}

		io.ErrorAt(s, "Can't tokenize")
	}

	newToken(cur, TOKEN_EOF, s)

	return head.next
}

func strToInt(s string) (int, string) {
	num := 0
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return num, s[i:]
		}

		num = num*10 + int(s[i]-'0')
	}

	return num, s[len(s):]
}
