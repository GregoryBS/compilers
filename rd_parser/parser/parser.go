package parser

import (
	"os"
)

type Parser struct {
	Data  string
	Index int
	File  *os.File
}

func (p *Parser) Block() {
	if p.Data[p.Index] != LBRACE {
		panic("Left brace '{' expected")
	}
	p.apply("")
	p.operatorList()
	if p.Data[p.Index] != RBRACE {
		panic("Right brace '}' expected")
	}
	p.apply("")
	if p.Index < len(p.Data) {
		panic("Something odd in the end")
	}
}

func (p *Parser) operatorList() {
	p.operator()
	p.tail()
}

func (p *Parser) operator() {
	p.identifier()
	if p.Data[p.Index] != EQ {
		panic("Equal sign '=' expected")
	}
	p.apply("")
	p.Expression()
}

func (p *Parser) tail() {
	if p.Data[p.Index] == SEP {
		p.apply("")
		p.operator()
		p.tail()
	}
}

func (p *Parser) Expression() {
	p.mathExpr()
	p.relationSign()
	p.mathExpr()
}

func (p *Parser) mathExpr() {
	defer func() {
		if msg := recover(); msg != nil {
			p.addSign()
			p.term()
			p.math()
		}
	}()
	p.term()
	p.math()
}

func (p *Parser) math() {
	defer func() {
		if msg := recover(); msg != nil {
			return
		}
	}()
	p.addSign()
	p.term()
	p.math()
}

func (p *Parser) term() {
	p.factor()
	p.term_()
}

func (p *Parser) term_() {
	defer func() {
		if msg := recover(); msg != nil {
			return
		}
	}()
	p.mulSign()
	p.factor()
	p.term_()
}

func (p *Parser) factor() {
	p.primaryExpr()
	p.factor_()
}

func (p *Parser) factor_() {
	if p.Data[p.Index] == CIRKUM {
		p.apply("")
		p.primaryExpr()
		p.factor_()
	}
}

func (p *Parser) primaryExpr() {
	defer func() {
		if msg := recover(); msg != nil {
			p.number()
		}
	}()
	defer func() {
		if msg := recover(); msg != nil {
			p.identifier()
		}
	}()
	if p.Data[p.Index] != LBRACKET {
		panic("Left bracket '(' expected")
	}
	p.apply("")
	p.mathExpr()
	if p.Data[p.Index] != RBRACKET {
		panic("Right bracket ')' expected")
	}
	p.apply("")
}
