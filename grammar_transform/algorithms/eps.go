package algorithms

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"

	combinations "github.com/mxschmitt/golang-combinations"
)

const (
	UNDEFINED = iota
	PENDING
	DEFINED
)

type State struct {
	Status int
	Flag   bool
}

func DeleteEpsRules(g *Grammar) *Grammar {
	epsPart := Part{Type: T, Value: EPS}
	emptyNonTerms := getEmptyNonterminals(g.NonTerminals, g.Rules)
	emptyRules := make([]int, 0)
	newRules := make([]*Rule, 0)
	for i, rule := range g.Rules {
		if rule.Right[0] == epsPart {
			emptyRules = append(emptyRules, i)
			continue
		}

		term := false
		index := make([]string, 0)
		for j, part := range rule.Right {
			search, limit := sort.SearchStrings(emptyNonTerms, part.Value), len(emptyNonTerms)
			if !term && search == limit {
				term = true
			} else if search < limit {
				index = append(index, fmt.Sprint(j))
			}
		}
		combs := combinations.All(index)
		if !term {
			combs = combs[:len(combs)-1]
		}
		for _, comb := range combs {
			r := &Rule{Left: rule.Left, Right: []Part{}}
			start := 0
			for _, c := range comb {
				j, _ := strconv.Atoi(c)
				r.Right = append(r.Right, rule.Right[start:j]...)
				start = j + 1
			}
			r.Right = append(r.Right, rule.Right[start:]...)
			newRules = append(newRules, r)
		}
	}

	for i := range emptyRules {
		g.Rules = append(g.Rules[:emptyRules[i]-i], g.Rules[emptyRules[i]-i+1:]...)
	}
	g.Rules = append(g.Rules, newRules...)
	g.Rules = removeDuplicates(g.Rules)
	if sort.SearchStrings(emptyNonTerms, g.Start) < len(emptyNonTerms) {
		start := g.Start + NEW
		g.NonTerminals = append(g.NonTerminals, start)
		g.Rules = append(g.Rules, &Rule{
			Left: Part{Type: N, Value: start},
			Right: []Part{{
				Type:  N,
				Value: g.Start,
			}},
		}, &Rule{
			Left:  Part{Type: N, Value: start},
			Right: []Part{epsPart},
		})
		g.Start = start
	}
	return g
}

func getEmptyNonterminals(nonterms []string, rules []*Rule) []string {
	states := make(map[string]*State, len(nonterms))
	for _, nt := range nonterms {
		states[nt] = new(State)
	}
	for _, nt := range nonterms {
		states[nt].Status = PENDING
		searchEmpty(states, nt, rules)
		states[nt].Status = DEFINED
	}

	result := make([]string, 0)
	for nt, state := range states {
		if state.Flag {
			result = append(result, nt)
		}
	}
	sort.Strings(result)
	return result
}

func searchEmpty(states map[string]*State, nt string, rules []*Rule) {
	if states[nt].Status == DEFINED {
		return
	}

	index := make([]int, 0)
	part, epsPart := Part{Type: N, Value: nt}, Part{Type: T, Value: EPS}
	for i, rule := range rules {
		if rule.Left == part {
			if rule.Right[0] == epsPart {
				states[nt].Flag = true
				return
			}
			index = append(index, i)
		}
	}
	nonterms := make([]string, 0)
	for _, i := range index {
		for _, right := range rules[i].Right {
			if right.Type == N && states[right.Value].Status == UNDEFINED {
				nonterms = append(nonterms, right.Value)
			}
		}
	}

	for _, v := range nonterms {
		states[v].Status = PENDING
		searchEmpty(states, v, rules)
		states[v].Status = DEFINED
		if states[v].Flag {
			states[nt].Flag = true
		}
	}
}

func removeDuplicates(slice []*Rule) []*Rule {
	allKeys := make(map[*Rule]Rule, 0)
	list := make([]*Rule, 0)
	for _, item := range slice {
		if value, ok := allKeys[item]; !ok && !reflect.DeepEqual(*item, value) {
			allKeys[item] = *item
			list = append(list, item)
		}
	}
	return list
}
