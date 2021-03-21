package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "%s: invalid number of arguments\n", os.Args[0])
		os.Exit(1)
	}

	i, s := strtol(os.Args[1], 10)

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".globl main\n")
	fmt.Printf("main:\n")
	fmt.Printf("  mov rax, %d\n", i)

	for s != "" {
		switch []rune(s)[0] {
		case '+':
			s = s[1:]
			i, s = strtol(s, 10)
			fmt.Printf("  add rax, %d\n", i)
			continue
		case '-':
			s = s[1:]
			i, s = strtol(s, 10)
			fmt.Printf("  sub rax, %d\n", i)
			continue
		default:
			fmt.Fprintf(os.Stderr, "unexpected character: %c\n", []rune(s)[0])
			os.Exit(1)
		}
	}
	fmt.Printf("  ret\n")
}

func strtol(s string, b int) (int, string) {
	if !unicode.IsDigit([]rune(s)[0]) {
		return 0, s
	}
	j := len(s)
	for i, c := range s {
		if !unicode.IsDigit(c) {
			j = i
			break
		}
	}
	n, _ := strconv.ParseInt(s[:j], b, 32)
	return int(n), s[j:]
}
