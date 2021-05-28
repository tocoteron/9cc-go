package tokenizer

import "github.com/tocoteron/9cc-go/internal/app/compiler/io"

type TokenKind int

const (
	TK_RESERVED TokenKind = iota
	TK_NUM
	TK_EOF
)

type Token struct {
	kind TokenKind
	next *Token
	val  int
	str  string
}

var CurrentToken *Token

func Consume(op byte) bool {
	if CurrentToken.kind != TK_RESERVED || CurrentToken.str[0] != op {
		return false
	}

	CurrentToken = CurrentToken.next

	return true
}

func Expect(op byte) {
	if CurrentToken.kind != TK_RESERVED || CurrentToken.str[0] != op {
		io.ErrorAt(CurrentToken.str, "It is not '%c'", op)
	}

	CurrentToken = CurrentToken.next
}

func ExpectNumber() int {
	if CurrentToken.kind != TK_NUM {
		io.ErrorAt(CurrentToken.str, "It is not a number")
	}

	val := CurrentToken.val

	CurrentToken = CurrentToken.next

	return val
}

func AtEOF() bool {
	return CurrentToken.kind == TK_EOF
}

func NewToken(cur *Token, kind TokenKind, str string) *Token {
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

		if s[0] == '+' || s[0] == '-' {
			cur = NewToken(cur, TK_RESERVED, s)
			s = s[1:]
			continue
		}

		if s[0] >= '0' && s[0] <= '9' {
			cur = NewToken(cur, TK_NUM, s)
			cur.val, s = strToInt(s)
			continue
		}

		io.ErrorAt(s, "Can't tokenize")
	}

	NewToken(cur, TK_EOF, s)

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
