package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
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

	asmCode, err := ioutil.ReadAll(asm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read file %v", err)
		return
	}

	code, err := assemble(asmCode)
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

func assemble(asmCode []byte) ([]byte, error) {
	var err error

	sc := bufio.NewScanner(bytes.NewReader(asmCode))
	st := NewSymbolTable()

	err = preProcess(sc, st)
	if err != nil {
		return nil, err
	}

	sc = bufio.NewScanner(bytes.NewReader(asmCode))
	return buildProcess(sc, st)
}

func preProcess(sc *bufio.Scanner, st *SymbolTable) error {
	var err error
	parser := NewParser(sc)
	n := 0

	for parser.HasMoreCommands() {
		err = parser.Advance()
		if err != nil {
			return err
		}

		switch parser.CommandType() {
		case CCommand, ACommand:
			n += 1
		case LCommand:
			symbol := parser.Symbol()
			st.AddEntry(symbol, n)
		}
	}

	return nil
}

func buildProcess(sc *bufio.Scanner, st *SymbolTable) ([]byte, error) {
	var buf bytes.Buffer
	var err error
	parser := NewParser(sc)
	na := 16

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
		case ACommand:
			symbol := parser.Symbol()
			i, err := strconv.Atoi(symbol)
			address := i
			if err != nil {
				if !st.Contains(symbol) {
					st.AddEntry(symbol, na)
					na += 1
				}
				address = st.GetAddress(symbol)
			}
			line := fmt.Sprintf("0%015b\n", address)
			buf.WriteString(line)
		}
	}

	return buf.Bytes(), nil
}
