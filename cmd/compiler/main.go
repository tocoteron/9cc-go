package main

import (
	"fmt"
	"os"

	"github.com/tocoteron/9cc-go/internal/app/compiler/io"
	"github.com/tocoteron/9cc-go/internal/app/compiler/tokenizer"
)

func main() {
	args := os.Args

	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "The number of arguments is invalid")
		os.Exit(1)
	}

	io.UserInput = args[1]
	tokenizer.CurrentToken = tokenizer.Tokenize(io.UserInput)

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".globl main\n")
	fmt.Printf("main:\n")

	fmt.Printf("  mov rax, %d\n", tokenizer.ExpectNumber())

	for !tokenizer.AtEOF() {
		if tokenizer.Consume('+') {
			fmt.Printf("  add rax, %d\n", tokenizer.ExpectNumber())
			continue
		}

		tokenizer.Expect('-')
		fmt.Printf("  sub rax, %d\n", tokenizer.ExpectNumber())
	}

	fmt.Printf("  ret\n")
}
