package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "%s: invalid number of arguments\n", os.Args[0])
		os.Exit(1)
	}

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf("  .globl main\n")
	fmt.Printf("main:\n")
	fmt.Printf("  mov rax, %d\n", atoi(os.Args[1]))
	fmt.Printf("  ret\n")
}

func atoi(s string) int {
	num, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	return num
}
