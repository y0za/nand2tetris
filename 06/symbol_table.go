package main

import "fmt"

type SymbolTable struct {
	table map[string]int
}

func NewSymbolTable() *SymbolTable {
	table := map[string]int{
		"SP":     0,
		"LCL":    1,
		"ARG":    2,
		"THIS":   3,
		"THAT":   4,
		"SCREEN": 16384,
		"KBD":    24576,
	}
	for i := 0; i <= 15; i += 1 {
		s := fmt.Sprintf("R%d", i)
		table[s] = i
	}
	return &SymbolTable{
		table: table,
	}
}

func (st *SymbolTable) AddEntry(symbol string, address int) {
	st.table[symbol] = address
}

func (st *SymbolTable) Contains(symbol string) bool {
	_, ok := st.table[symbol]
	return ok
}

func (st *SymbolTable) GetAddress(symbol string) int {
	return st.table[symbol]
}
