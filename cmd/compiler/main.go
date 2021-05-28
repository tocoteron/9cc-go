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
	node := parser.Parse()

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".globl main\n")
	fmt.Printf("main:\n")

	parser.Generate(node)

	fmt.Printf("  pop rax\n")
	fmt.Printf("  ret\n")
}
