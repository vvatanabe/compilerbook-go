package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

const (
	TK_RESERVED TokenKind = iota // Keywords or punctuators
	TK_NUM                       // Integer literals
	TK_EOF                       // End-of-file markers
)

type TokenKind int

// Token type
type Token struct {
	Kind TokenKind `json:"kind"` // Token kind
	Next *Token    `json:"next"` // Next token
	Val  int       `json:"val"`  // If kind is TK_NUM, its value
	Str  []rune    `json:"str"`  // Token string
}

func (t *Token) String() string {
	b, _ := json.Marshal(t)
	return string(b)
}

// Current token
var token *Token

// Reports an error and exit.
func printError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

// Consumes the current token if it matches `op`.
//
// 次のトークンが期待している記号のときには、トークンを1つ読み進めて
// 真を返す。それ以外の場合には偽を返す。
func consume(op rune) bool {
	if token.Kind != TK_RESERVED || token.Str[0] != op {
		return false
	}
	token = token.Next
	return true
}

// Ensure that the current token is `op`.
//
// 次のトークンが期待している記号のときには、トークンを1つ読み進める。
// それ以外の場合にはエラーを報告する。
func expect(op rune) {
	if token.Kind != TK_RESERVED || token.Str[0] != op {
		printError("expected '%s'", string(op))
	}
	token = token.Next
}

// // Ensure that the current token is TK_NUM.
func expectNumber() int {
	if token.Kind != TK_NUM {
		printError("expected a number")
	}
	val := token.Val
	token = token.Next
	return val
}

func atEof() bool {
	return token.Kind == TK_EOF
}

// Create a new token and add it as the next token of `cur`.
//
// 新しいトークンを作成してcurに繋げる
func newToken(kind TokenKind, cur *Token, str []rune) *Token {
	tok := &Token{
		Kind: kind,
		Str:  str,
	}
	cur.Next = tok
	return tok
}

// Tokenize `p` and returns new tokens.
func tokenize(p []rune) *Token {
	var head Token
	cur := &head

	for len(p) > 0 {
		c := p[0]
		// Skip whitespace characters.
		if unicode.IsSpace(c) {
			p = p[1:]
			continue
		}

		// Punctuator
		if c == '+' || c == '-' {
			cur = newToken(TK_RESERVED, cur, p)
			p = p[1:]
			continue
		}

		// Integer literal
		if unicode.IsDigit(c) {
			cur = newToken(TK_NUM, cur, p)
			cur.Val, p = strtol(p, 10)
			continue
		}

		printError("invalid token")
	}
	newToken(TK_EOF, cur, []rune{0})
	return head.Next
}

func main() {
	if len(os.Args) != 2 {
		printError("%s: invalid number of arguments\n", os.Args[0])
	}

	// トークナイズする
	token = tokenize([]rune(os.Args[1]))

	// アセンブリの前半部分を出力
	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".globl main\n")
	fmt.Printf("main:\n")

	// The first token must be a number
	// 式の最初は数でなければならないので、それをチェックして
	// 最初のmov命令を出力
	fmt.Printf("  mov rax, %d\n", expectNumber())

	// ... followed by either `+ <number>` or `- <number>`.
	// `+ <数>`あるいは`- <数>`というトークンの並びを消費しつつ
	// アセンブリを出力
	for !atEof() {
		if consume('+') {
			fmt.Printf("  add rax, %d\n", expectNumber())
			continue
		}

		expect('-')
		fmt.Printf("  sub rax, %d\n", expectNumber())
	}

	fmt.Printf("  ret\n")
}

func strtol(s []rune, b int) (int, []rune) {
	if !unicode.IsDigit(s[0]) {
		return 0, s
	}
	j := len(s)
	for i, c := range s {
		if !unicode.IsDigit(c) {
			j = i
			break
		}
	}
	n, _ := strconv.ParseInt(string(s[:j]), b, 32)
	return int(n), s[j:]
}
