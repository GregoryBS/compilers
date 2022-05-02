package analysis

func BuildFSM(regexp string) *FSM {
	ndfsm := new(FSM)
	ndfsm.alphabet = getAlphabet(regexp)

	tree := parseExpression(regexp, ndfsm.alphabet)
	start, finish := &Node{"n0"}, &Node{"n1"}
	ndfsm.graph = &Graph{
		[]*Node{start, finish},
		[]*Link{},
	}
	if isOperation(tree.symbol) {
		err := calcTree(tree, ndfsm)
		if err != nil {
			return nil
		}
	} else {
		ndfsm.graph.Links = append(ndfsm.graph.Links, &Link{start.Name, finish.Name, tree.symbol})
	}

	//ndfsm.GetGraph()
	dfsm := Determinate(ndfsm)
	//dfsm.GetGraph()
	return Minimize(dfsm)
}

func getAlphabet(regexp string) map[byte]bool {
	result := make(map[byte]bool, 0)
	for i := range regexp {
		switch regexp[i] {
		case iteration, concat, union, '(', ')':
			continue
		default:
			result[regexp[i]] = true
		}
	}
	return result
}

func Find(nodes []*Node, name string) int {
	for i, n := range nodes {
		if name == n.Name {
			return i
		}
	}
	return len(nodes)
}
