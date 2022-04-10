package analysis

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
	start    int
	finish   []int
}
