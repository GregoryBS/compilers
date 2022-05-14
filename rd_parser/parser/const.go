package parser

import (
	"unicode"

	"github.com/goccy/go-graphviz/cgraph"
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
	EPS      = "eps"
)

func (p *Parser) apply(term string) {
	if term == "" {
		term = p.Data[p.Index : p.Index+1]
	}
	p.File.WriteString(term)
	p.Index += len(term)
}

func (p *Parser) addSign(begin *cgraph.Node) {
	switch p.Data[p.Index] {
	case PLUS, MINUS:
		s := p.Graph.AddNode(string(p.Data[p.Index]))
		p.apply("")
		p.Graph.AddEdges(begin, s)
	default:
		panic("Add operation sign expected")
	}
}

func (p *Parser) mulSign(begin *cgraph.Node) {
	switch p.Data[p.Index] {
	case MUL, DIV, MOD:
		s := p.Graph.AddNode(string(p.Data[p.Index]))
		p.apply("")
		p.Graph.AddEdges(begin, s)
	default:
		panic("Multiply operation sign expected")
	}
}

func (p *Parser) relationSign(begin *cgraph.Node) {
	sign := p.Data[p.Index : p.Index+2]
	if sign == LE || sign == GE || sign == NE {
		s := p.Graph.AddNode(sign)
		p.apply(sign)
		p.Graph.AddEdges(begin, s)
	} else if sign[0] == LT || sign[0] == GT || sign[0] == EQ {
		s := p.Graph.AddNode(sign[:1])
		p.apply("")
		p.Graph.AddEdges(begin, s)
	} else {
		panic("Relation operation sign expected")
	}
}

func (p *Parser) number(begin *cgraph.Node) {
	if unicode.IsDigit(rune(p.Data[p.Index])) {
		id := string(p.Data[p.Index])
		p.apply("")
		for unicode.IsDigit(rune(p.Data[p.Index])) {
			id += string(p.Data[p.Index])
			p.apply("")
		}
		end := p.Graph.AddNode(id)
		p.Graph.AddEdges(begin, end)
	} else {
		panic("Number expected")
	}
}

func (p *Parser) identifier(begin *cgraph.Node) {
	if unicode.IsLetter(rune(p.Data[p.Index])) {
		id := string(p.Data[p.Index])
		p.apply("")
		for unicode.IsLetter(rune(p.Data[p.Index])) {
			id += string(p.Data[p.Index])
			p.apply("")
		}
		end := p.Graph.AddNode(id)
		p.Graph.AddEdges(begin, end)
	} else {
		panic("Symbol identifier expected")
	}
}
