package parser

import (
	"fmt"

	"github.com/goccy/go-graphviz/cgraph"
)

var counters = map[string]int{}

type Graph struct {
	Graph *cgraph.Graph
	Nodes []*cgraph.Node
	Edges []*cgraph.Edge
}

func (g *Graph) AddNode(name string) *cgraph.Node {
	node, err := g.Graph.CreateNode(fmt.Sprintf("%s (%d)", name, counters[name]))
	if err != nil {
		panic("Cannot create graph node" + err.Error())
	}
	counters[name] += 1
	g.Nodes = append(g.Nodes, node)
	return node
}

func (g *Graph) AddEdges(begin *cgraph.Node, ends ...*cgraph.Node) {
	for i := 0; i < len(ends); i++ {
		e, err := g.Graph.CreateEdge(ends[i].Name(), begin, ends[i])
		if err != nil {
			panic("Cannot create graph edge" + err.Error())
		}
		g.Edges = append(g.Edges, e)
	}
}

func (g *Graph) DeleteNodes(nodes []*cgraph.Node) {
	for _, node := range nodes {
		g.Graph.DeleteNode(node)
	}
}

func (g *Graph) DeleteEdges(edges []*cgraph.Edge) {
	for _, edge := range edges {
		g.Graph.DeleteEdge(edge)
	}
}
