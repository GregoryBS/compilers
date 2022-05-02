package analysis

func Minimize(dfsm *FSM) *FSM {
	groups := make([][]*Node, 2)
	splitFinish(groups, dfsm.graph.Nodes, dfsm.finish)

	queueIndex := make([]int, 0)
	queueSymbol := make([]string, 0)
	for s := range dfsm.alphabet {
		for i := range groups {
			queueIndex = append(queueIndex, i)
			queueSymbol = append(queueSymbol, string(s))
		}
	}

	for len(queueIndex) > 0 {
		class, symbol := groups[queueIndex[0]], queueSymbol[0]
		queueIndex, queueSymbol = queueIndex[1:], queueSymbol[1:]
		groupCount := len(groups)
		for i := 0; i < groupCount; i += 1 {
			left, right := split(groups[i], class, dfsm.graph.Links, symbol)
			if len(left) > 0 && len(right) > 0 {
				groups[i] = left
				groups = append(groups, right)
				for s := range dfsm.alphabet {
					queueIndex = append(queueIndex, len(groups)-1)
					queueSymbol = append(queueSymbol, string(s))
				}
			}
		}
	}

	nodes := make([]*Node, len(groups))
	links := make([]*Link, 0)
	result := &FSM{alphabet: dfsm.alphabet, finish: make([]string, 0)}
	for _, group := range groups {
		if Find(group, dfsm.start) < len(group) {
			result.start = group[0].Name
			break
		}
	}
	for i, group := range groups {
		nodes[i] = group[0]
		for _, name := range dfsm.finish {
			if Find(group, name) < len(group) {
				result.finish = append(result.finish, group[0].Name)
			}
		}
	}
	for _, link := range dfsm.graph.Links {
		if Find(nodes, link.Source) < len(nodes) && Find(nodes, link.Target) < len(nodes) {
			links = append(links, link)
		}
	}
	result.graph = &Graph{Nodes: nodes, Links: links}
	return result
}

func split(group, class []*Node, links []*Link, symbol string) ([]*Node, []*Node) {
	left := make([]*Node, 0)
	right := make([]*Node, 0)
	for _, node := range group {
		flag := true
		for _, link := range links {
			if link.Source == node.Name && link.Value == symbol && Find(class, link.Target) == len(class) {
				flag = false
				break
			}
		}
		if flag {
			left = append(left, node)
		} else {
			right = append(right, node)
		}
	}
	return left, right
}

func splitFinish(groups [][]*Node, nodes []*Node, finish []string) {
	for _, node := range nodes {
		flag := false
		for _, name := range finish {
			if node.Name == name {
				flag = true
				break
			}
		}
		if flag {
			groups[0] = append(groups[0], node)
		} else {
			groups[1] = append(groups[1], node)
		}
	}
}
