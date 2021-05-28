package tokenizer

import (
	"strings"

	"github.com/tocoteron/9cc-go/internal/app/compiler/io"
)

type TokenKind int

const (
	TOKEN_RESERVED TokenKind = iota
	TOKEN_IDENT
	TOKEN_NUM
	TOKEN_EOF
)

type Token struct {
	kind TokenKind
	next *Token
	val  int
	Str  string
	len  int
}

var CurrentToken *Token

func Consume(op string) bool {
	if CurrentToken.kind != TOKEN_RESERVED || CurrentToken.Str[:CurrentToken.len] != op {
		return false
	}

	CurrentToken = CurrentToken.next

	return true
}

func ConsumeIdent() *Token {
	if CurrentToken.kind != TOKEN_IDENT {
		return nil
	}

	token := CurrentToken

	CurrentToken = CurrentToken.next

	return token
}

func Expect(op string) {
	if CurrentToken.kind != TOKEN_RESERVED || CurrentToken.Str[:CurrentToken.len] != op {
		io.ErrorAt(CurrentToken.Str, "It is not '%s'", op)
	}

	CurrentToken = CurrentToken.next
}

func ExpectNumber() int {
	if CurrentToken.kind != TOKEN_NUM {
		io.ErrorAt(CurrentToken.Str, "It is not a number")
	}

	val := CurrentToken.val

	CurrentToken = CurrentToken.next

	return val
}

func AtEOF() bool {
	return CurrentToken.kind == TOKEN_EOF
}

func newToken(cur *Token, kind TokenKind, str string, len int) *Token {
	token := &Token{}
	token.kind = kind
	token.Str = str
	token.len = len
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

		if strings.HasPrefix(s, "==") || strings.HasPrefix(s, "!=") ||
			strings.HasPrefix(s, "<=") || strings.HasPrefix(s, ">=") {
			cur = newToken(cur, TOKEN_RESERVED, s, 2)
			s = s[2:]
			continue
		}

		if strings.IndexByte("+-*/()<>=;", s[0]) != -1 {
			cur = newToken(cur, TOKEN_RESERVED, s, 1)
			s = s[1:]
			continue
		}

		if s[0] >= 'a' && s[0] <= 'z' {
			cur = newToken(cur, TOKEN_IDENT, s, 0)
			s = s[1:]
			continue
		}

		if s[0] >= '0' && s[0] <= '9' {
			cur = newToken(cur, TOKEN_NUM, s, 0)
			cur.val, s = strToInt(s)
			continue
		}

		io.ErrorAt(s, "Can't tokenize")
	}

	newToken(cur, TOKEN_EOF, s, 0)

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
