package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

type Parser struct {
	scanner     *bufio.Scanner
	command     string
	advanced    bool
	cachedHMC   bool
	commandType CommandType
	symbol      string
	dest        string
	comp        string
	jump        string
}

type CommandType int

const (
	ACommand CommandType = iota
	CCommand
	LCommand
)

var (
	acRegexp = regexp.MustCompile(`^@([^()]+)$`)
	ccRegexp = regexp.MustCompile(`^([AMD]*)=?([\-+!|&AMD01]+);?(\w*)$`)
	lcRegexp = regexp.MustCompile(`^\(([^()]+)\)$`)
)

func NewParser(scanner *bufio.Scanner) *Parser {
	return &Parser{
		scanner:   scanner,
		advanced:  true,
		cachedHMC: false,
	}
}

func (p *Parser) HasMoreCommands() bool {
	if !p.advanced {
		return p.cachedHMC
	}
	p.advanced = false

	for p.scanner.Scan() {
		sl := strings.Split(p.scanner.Text(), "//")
		line := strings.TrimSpace(sl[0])
		if line != "" {
			p.command = line
			p.cachedHMC = true
			return true
		}
	}

	p.cachedHMC = false
	return false
}

func (p *Parser) Advance() error {
	p.advanced = true

	if acRegexp.MatchString(p.command) {
		p.commandType = ACommand
		l := acRegexp.FindStringSubmatch(p.command)
		p.symbol = l[1]
		return nil
	}

	if ccRegexp.MatchString(p.command) {
		p.commandType = CCommand
		l := ccRegexp.FindStringSubmatch(p.command)
		p.dest = l[1]
		p.comp = l[2]
		p.jump = l[3]
		return nil
	}

	if lcRegexp.MatchString(p.command) {
		p.commandType = LCommand
		l := lcRegexp.FindStringSubmatch(p.command)
		p.symbol = l[1]
		return nil
	}

	return fmt.Errorf(`unexpected command: "%s"`, p.command)
}

func (p *Parser) CommandType() CommandType {
	return p.commandType
}

func (p *Parser) Symbol() string {
	return p.symbol
}

func (p *Parser) Dest() string {
	return p.dest
}

func (p *Parser) Comp() string {
	return p.comp
}

func (p *Parser) Jump() string {
	return p.jump
}
