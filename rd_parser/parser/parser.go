package parser

import (
	"os"

	"github.com/goccy/go-graphviz/cgraph"
)

type Parser struct {
	Data  string
	Index int
	File  *os.File
	Graph *Graph
}

func (p *Parser) Block() {
	block := p.Graph.AddNode("block")

	if p.Data[p.Index] != LBRACE {
		panic("Left brace '{' expected")
	}
	p.apply("")
	lbrace := p.Graph.AddNode(string(LBRACE))

	opList := p.Graph.AddNode("operatorList")
	p.operatorList(opList)

	if p.Data[p.Index] != RBRACE {
		panic("Right brace '}' expected")
	}
	p.apply("")
	rbrace := p.Graph.AddNode(string(RBRACE))

	if p.Index < len(p.Data) {
		panic("Something odd in the end")
	}
	p.Graph.AddEdges(block, lbrace, opList, rbrace)
}

func (p *Parser) operatorList(begin *cgraph.Node) {
	operator := p.Graph.AddNode("operator")
	p.operator(operator)

	tail := p.Graph.AddNode("tail")
	p.tail(tail)

	p.Graph.AddEdges(begin, operator, tail)
}

func (p *Parser) operator(begin *cgraph.Node) {
	id := p.Graph.AddNode("identifier")
	p.identifier(id)

	if p.Data[p.Index] != EQ {
		panic("Equal sign '=' expected")
	}
	p.apply("")
	eq := p.Graph.AddNode(string(EQ))

	expr := p.Graph.AddNode("expression")
	p.Expression(expr)

	p.Graph.AddEdges(begin, id, eq, expr)
}

func (p *Parser) tail(begin *cgraph.Node) {
	if p.Data[p.Index] == SEP {
		p.apply("")
		sep := p.Graph.AddNode(string(SEP))

		operator := p.Graph.AddNode("operator")
		p.operator(operator)

		tail := p.Graph.AddNode("tail")
		p.tail(tail)

		p.Graph.AddEdges(begin, sep, operator, tail)
		return
	}
	eps := p.Graph.AddNode(EPS)
	p.Graph.AddEdges(begin, eps)
}

func (p *Parser) Expression(begin *cgraph.Node) {
	me1 := p.Graph.AddNode("mathExpr")
	p.mathExpr(me1)

	rs := p.Graph.AddNode("relationSign")
	p.relationSign(rs)

	me2 := p.Graph.AddNode("mathExpr")
	p.mathExpr(me2)

	p.Graph.AddEdges(begin, me1, rs, me2)
}

func (p *Parser) mathExpr(begin *cgraph.Node) {
	nodes, edges := len(p.Graph.Nodes), len(p.Graph.Edges)
	defer func() {
		if msg := recover(); msg != nil {
			p.Graph.DeleteNodes(p.Graph.Nodes[nodes:])
			p.Graph.DeleteEdges(p.Graph.Edges[edges:])
			p.Graph.Nodes = p.Graph.Nodes[:nodes]
			p.Graph.Edges = p.Graph.Edges[:edges]

			as := p.Graph.AddNode("addSign")
			p.addSign(as)

			t := p.Graph.AddNode("term")
			p.term(t)

			m := p.Graph.AddNode("math")
			p.math(m)

			p.Graph.AddEdges(begin, as, t, m)
		}
	}()
	t := p.Graph.AddNode("term")
	p.term(t)

	m := p.Graph.AddNode("math")
	p.math(m)

	p.Graph.AddEdges(begin, t, m)
}

func (p *Parser) math(begin *cgraph.Node) {
	nodes, edges := len(p.Graph.Nodes), len(p.Graph.Edges)
	defer func() {
		if msg := recover(); msg != nil {
			p.Graph.DeleteNodes(p.Graph.Nodes[nodes:])
			p.Graph.DeleteEdges(p.Graph.Edges[edges:])
			p.Graph.Nodes = p.Graph.Nodes[:nodes]
			p.Graph.Edges = p.Graph.Edges[:edges]

			eps := p.Graph.AddNode(EPS)
			p.Graph.AddEdges(begin, eps)
		}
	}()
	as := p.Graph.AddNode("addSign")
	p.addSign(as)

	t := p.Graph.AddNode("term")
	p.term(t)

	m := p.Graph.AddNode("math")
	p.math(m)

	p.Graph.AddEdges(begin, as, t, m)
}

func (p *Parser) term(begin *cgraph.Node) {
	f := p.Graph.AddNode("factor")
	p.factor(f)

	t := p.Graph.AddNode("term_")
	p.term_(t)

	p.Graph.AddEdges(begin, f, t)
}

func (p *Parser) term_(begin *cgraph.Node) {
	nodes, edges := len(p.Graph.Nodes), len(p.Graph.Edges)
	defer func() {
		if msg := recover(); msg != nil {
			p.Graph.DeleteNodes(p.Graph.Nodes[nodes:])
			p.Graph.DeleteEdges(p.Graph.Edges[edges:])
			p.Graph.Nodes = p.Graph.Nodes[:nodes]
			p.Graph.Edges = p.Graph.Edges[:edges]

			eps := p.Graph.AddNode(EPS)
			p.Graph.AddEdges(begin, eps)
		}
	}()
	ms := p.Graph.AddNode("mulSign")
	p.mulSign(ms)

	f := p.Graph.AddNode("factor")
	p.factor(f)

	t := p.Graph.AddNode("term_")
	p.term_(t)

	p.Graph.AddEdges(begin, ms, f, t)
}

func (p *Parser) factor(begin *cgraph.Node) {
	pe := p.Graph.AddNode("primaryExpr")
	p.primaryExpr(pe)

	f := p.Graph.AddNode("factor_")
	p.factor_(f)

	p.Graph.AddEdges(begin, pe, f)
}

func (p *Parser) factor_(begin *cgraph.Node) {
	if p.Data[p.Index] == CIRKUM {
		p.apply("")
		c := p.Graph.AddNode(string(CIRKUM))

		pe := p.Graph.AddNode("primaryExpr")
		p.primaryExpr(pe)

		f := p.Graph.AddNode("factor_")
		p.factor_(f)

		p.Graph.AddEdges(begin, c, pe, f)
		return
	}
	eps := p.Graph.AddNode(EPS)
	p.Graph.AddEdges(begin, eps)
}

func (p *Parser) primaryExpr(begin *cgraph.Node) {
	nodes, edges := len(p.Graph.Nodes), len(p.Graph.Edges)
	defer func() {
		if msg := recover(); msg != nil {
			p.Graph.DeleteNodes(p.Graph.Nodes[nodes:])
			p.Graph.DeleteEdges(p.Graph.Edges[edges:])
			p.Graph.Nodes = p.Graph.Nodes[:nodes]
			p.Graph.Edges = p.Graph.Edges[:edges]

			n := p.Graph.AddNode("number")
			p.number(n)

			p.Graph.AddEdges(begin, n)
		}
	}()
	defer func() {
		if msg := recover(); msg != nil {
			p.Graph.DeleteNodes(p.Graph.Nodes[nodes:])
			p.Graph.DeleteEdges(p.Graph.Edges[edges:])
			p.Graph.Nodes = p.Graph.Nodes[:nodes]
			p.Graph.Edges = p.Graph.Edges[:edges]

			id := p.Graph.AddNode("identifier")
			p.identifier(id)

			p.Graph.AddEdges(begin, id)
		}
	}()
	if p.Data[p.Index] != LBRACKET {
		panic("Left bracket '(' expected")
	}
	p.apply("")
	lb := p.Graph.AddNode(string(LBRACKET))

	me := p.Graph.AddNode("mathExpr")
	p.mathExpr(me)
	if p.Data[p.Index] != RBRACKET {
		panic("Right bracket ')' expected")
	}
	p.apply("")
	rb := p.Graph.AddNode(string(RBRACKET))

	p.Graph.AddEdges(begin, lb, me, rb)
}
