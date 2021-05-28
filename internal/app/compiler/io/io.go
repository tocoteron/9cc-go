package io

import (
	"fmt"
	"os"
	"strings"
)

var UserInput string

func Error(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func ErrorAt(loc string, format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s\n", UserInput)

	pos := len(UserInput) - len(loc)

	fmt.Fprintf(os.Stderr, strings.Repeat(" ", pos))
	fmt.Fprintf(os.Stderr, "^ ")
	fmt.Fprintf(os.Stderr, format+"\n", a...)

	os.Exit(1)
}
