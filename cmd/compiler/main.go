package main

import (
	"fmt"
	"os"
	"strings"
)

type TokenKind int

const (
	TK_RESERVED = iota
	TK_NUM
	TK_EOF
)

type Token struct {
	kind TokenKind
	next *Token
	val  int
	str  string
}

var token *Token
var userInput string

func errorAt(loc string, format string, a ...interface{}) {
	fmt.Println(userInput)

	pos := len(userInput) - len(loc)

	fmt.Fprintf(os.Stderr, strings.Repeat(" ", pos))
	fmt.Fprintf(os.Stderr, "^ ")
	fmt.Fprintf(os.Stderr, format+"\n", a...)

	os.Exit(1)
}

func consume(op byte) bool {
	if token.kind != TK_RESERVED || token.str[0] != op {
		return false
	}

	token = token.next

	return true
}

func expect(op byte) {
	if token.kind != TK_RESERVED || token.str[0] != op {
		errorAt(token.str, "It is not '%c'", op)
	}

	token = token.next
}

func expectNumber() int {
	if token.kind != TK_NUM {
		errorAt(token.str, "It is not a number")
	}

	val := token.val

	token = token.next

	return val
}

func atEOF() bool {
	return token.kind == TK_EOF
}

func newToken(cur *Token, kind TokenKind, str string) *Token {
	token := &Token{}
	token.kind = kind
	token.str = str
	cur.next = token

	return token
}

func tokenize(s string) *Token {
	var head Token
	head.next = nil

	cur := &head

	for s != "" {
		if s[0] == ' ' {
			s = s[1:]
			continue
		}

		if s[0] == '+' || s[0] == '-' {
			cur = newToken(cur, TK_RESERVED, s)
			s = s[1:]
			continue
		}

		if s[0] >= '0' && s[0] <= '9' {
			cur = newToken(cur, TK_NUM, s)
			cur.val, s = strToInt(s)
			continue
		}

		errorAt(s, "Can't tokenize")
	}

	newToken(cur, TK_EOF, s)

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

func main() {
	args := os.Args

	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "The number of arguments is invalid")
		os.Exit(1)
	}

	userInput = args[1]
	token = tokenize(userInput)

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".globl main\n")
	fmt.Printf("main:\n")

	fmt.Printf("  mov rax, %d\n", expectNumber())

	for !atEOF() {
		if consume('+') {
			fmt.Printf("  add rax, %d\n", expectNumber())
			continue
		}

		expect('-')
		fmt.Printf("  sub rax, %d\n", expectNumber())
	}

	fmt.Printf("  ret\n")
}
