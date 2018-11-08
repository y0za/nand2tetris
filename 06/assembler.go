package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "need asm file")
		return
	}

	if !strings.HasSuffix(os.Args[1], ".asm") {
		fmt.Fprintln(os.Stderr, "need asm file")
		return
	}

	asm, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file %v", err)
		return
	}

	sc := bufio.NewScanner(asm)
	code, err := assemble(sc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to assemble %v", err)
		return
	}

	hfn := strings.Replace(os.Args[1], ".asm", ".hack", 1)
	hack, err := os.Create(hfn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create file %v", err)
		return
	}
	defer hack.Close()

	hack.Write(code)
}

func assemble(sc *bufio.Scanner) ([]byte, error) {
	var buf bytes.Buffer
	var err error

	parser := NewParser(sc)
	for parser.HasMoreCommands() {
		err = parser.Advance()
		if err != nil {
			return nil, err
		}

		switch parser.CommandType() {
		case CCommand:
			comp := compCode(parser.Comp())
			dest := destCode(parser.Dest())
			jump := jumpCode(parser.Jump())
			line := fmt.Sprintf("111%s%s%s\n", comp, dest, jump)
			buf.WriteString(line)
		}
	}

	return buf.Bytes(), nil
}
