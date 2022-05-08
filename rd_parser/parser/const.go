package parser

import (
	"unicode"
)

const (
	CIRKUM   = '^'
	LBRACKET = '('
	RBRACKET = ')'
	PLUS     = '+'
	MINUS    = '-'
	MUL      = '*'
	DIV      = '/'
	MOD      = '%'
	LT       = '<'
	LE       = "<="
	EQ       = '='
	GE       = ">="
	GT       = '>'
	NE       = "<>"
	LBRACE   = '{'
	RBRACE   = '}'
	SEP      = ';'
)

func (p *Parser) apply(term string) {
	if term == "" {
		term = p.Data[p.Index : p.Index+1]
	}
	p.File.WriteString(term)
	p.Index += len(term)
}

func (p *Parser) addSign() {
	switch p.Data[p.Index] {
	case PLUS, MINUS:
		p.apply("")
	default:
		panic("Add operation sign expected")
	}
}

func (p *Parser) mulSign() {
	switch p.Data[p.Index] {
	case MUL, DIV, MOD:
		p.apply("")
	default:
		panic("Multiply operation sign expected")
	}
}

func (p *Parser) relationSign() {
	sign := p.Data[p.Index : p.Index+2]
	if sign == LE || sign == GE || sign == NE {
		p.apply(sign)
	} else if sign[0] == LT || sign[0] == GT || sign[0] == EQ {
		p.apply("")
	} else {
		panic("Relation operation sign expected")
	}
}

func (p *Parser) number() {
	if unicode.IsDigit(rune(p.Data[p.Index])) {
		p.apply("")
		for unicode.IsDigit(rune(p.Data[p.Index])) {
			p.apply("")
		}
	} else {
		panic("Number expected")
	}
}

func (p *Parser) identifier() {
	if unicode.IsLetter(rune(p.Data[p.Index])) {
		p.apply("")
		for unicode.IsLetter(rune(p.Data[p.Index])) {
			p.apply("")
		}
	} else {
		panic("Symbol identifier expected")
	}
}
