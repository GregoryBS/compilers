package analysis

func Modelling(fsm *FSM, word string) bool {
	links := fsm.graph.Links
	idx := Find(fsm.graph.Nodes, fsm.start)
	states := closure(fsm.graph.Nodes[idx:idx+1], links)
	for i := range word {
		states = closure(move(states, links, word[i:i+1]), links)
	}
	for _, finish := range fsm.finish {
		if Find(states, finish) < len(states) {
			return true
		}
	}
	return false
}
