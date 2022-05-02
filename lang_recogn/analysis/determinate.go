package analysis

import (
	"fmt"
	"reflect"
	"sort"
)

func Determinate(ndfsm *FSM) *FSM {
	visit := make(map[string][]*Node, 0)
	start := fmt.Sprintf("n%d", len(visit))
	idx := Find(ndfsm.graph.Nodes, ndfsm.start)
	visit[start] = closure(ndfsm.graph.Nodes[idx:idx+1], ndfsm.graph.Links)
	links := make([]*Link, 0)
	queue := make([]string, 0)
	queue = append(queue, start)
	for len(queue) > 0 {
		src := queue[0]
		queue = queue[1:]
		nodes := visit[src]
		for s := range ndfsm.alphabet {
			states := closure(move(nodes, ndfsm.graph.Links, string(s)), ndfsm.graph.Links)
			var node string
			for name, n := range visit {
				if reflect.DeepEqual(n, states) {
					node = name
					break
				}
			}
			if node == "" {
				node = fmt.Sprintf("n%d", len(visit))
				visit[node] = states
				queue = append(queue, node)
			}
			links = append(links, &Link{Source: src, Target: node, Value: string(s)})
		}
	}
	result := &FSM{graph: &Graph{Nodes: setToSlice(visit), Links: links}, alphabet: ndfsm.alphabet}
	result.start = result.graph.Nodes[0].Name
	for f, nodes := range visit {
		for _, name := range ndfsm.finish {
			if Find(nodes, name) < len(nodes) {
				result.finish = append(result.finish, f)
			}
		}
	}
	return result
}

func setToSlice(set interface{}) []*Node {
	keys := make([]string, 0)
	switch m := set.(type) {
	case map[string]bool:
		for k := range m {
			keys = append(keys, k)
		}
	case map[string][]*Node:
		for k := range m {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	result := make([]*Node, 0)
	for _, k := range keys {
		result = append(result, &Node{k})
	}
	return result
}

func closure(nodes []*Node, links []*Link) []*Node {
	set := make(map[string]bool, len(nodes))
	for _, n := range nodes {
		set[n.Name] = true
	}
	for len(nodes) > 0 {
		node := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]
		for _, l := range links {
			if l.Source == node.Name && l.Value == eps {
				if _, ok := set[l.Target]; !ok {
					set[l.Target] = true
					nodes = append(nodes, &Node{l.Target})
				}
			}
		}
	}
	return setToSlice(set)
}

func move(nodes []*Node, links []*Link, symbol string) []*Node {
	set := make(map[string]bool, 0)
	for _, node := range nodes {
		for _, l := range links {
			if l.Source == node.Name && l.Value == symbol {
				if _, ok := set[l.Target]; !ok {
					set[l.Target] = true
				}
			}
		}
	}
	return setToSlice(set)
}
