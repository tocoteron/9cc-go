package main

import (
	"fmt"
	"os"

	"github.com/tocoteron/9cc-go/internal/app/compiler/io"
	"github.com/tocoteron/9cc-go/internal/app/compiler/parser"
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
	code := parser.Parse()

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".globl main\n")
	fmt.Printf("main:\n")

	// Prologue
	fmt.Printf("  push rbp\n")
	fmt.Printf("  mov rbp, rsp\n")
	fmt.Printf("  sub rsp, %d\n", 26*8)

	parser.Generate(code)

	// Epilogue
	fmt.Printf("  mov rsp, rbp\n")
	fmt.Printf("  pop rbp\n")
	fmt.Printf("  ret\n")
}
