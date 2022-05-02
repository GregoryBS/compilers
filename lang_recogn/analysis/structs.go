package analysis

import "fmt"

const (
	eps       = "eps"
	union     = '|'
	concat    = '.'
	iteration = '*'
	separator = ';'
)

type Node struct {
	Name string
}

type Link struct {
	Source string
	Target string
	Value  string
}

type Graph struct {
	Nodes []*Node
	Links []*Link
}

type FSM struct {
	alphabet map[byte]bool
	graph    *Graph
	start    string
	finish   []string
}

func (m *FSM) GetGraph() {
	fmt.Println("Nodes:")
	for _, node := range m.graph.Nodes {
		fmt.Println(node.Name)
	}
	fmt.Println("Links:")
	for _, link := range m.graph.Links {
		fmt.Println(link.Source, link.Target, link.Value)
	}
}
