package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args

	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "The number of arguments is invalid")
		os.Exit(1)
	}

	num, err := strconv.Atoi(args[1])
	if err != nil {
		panic(err)
	}

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".globl main\n")
	fmt.Printf("main:\n")
	fmt.Printf("  mov rax, %d\n", num)
	fmt.Printf("  ret\n")
}
