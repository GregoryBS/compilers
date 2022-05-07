package algorithms

import "fmt"

const (
	N   = "nonterminal"
	T   = "terminal"
	EPS = "eps"
	NEW = "_"
)

type Part struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Rule struct {
	Left  Part   `json:"left"`
	Right []Part `json:"right"`
}

type Grammar struct {
	NonTerminals []string `json:"nonterminals"`
	Terminals    []string `json:"terminals"`
	Start        string   `json:"start"`
	Rules        []*Rule  `json:"rules"`
}

func (g *Grammar) DeleteIndirect(ai, aj string) {
	partI, partJ := Part{Type: N, Value: ai}, Part{Type: N, Value: aj}
	with := make([]*Rule, 0)
	without := make([]*Rule, 0)
	for _, rule := range g.Rules {
		if rule.Left == partI {
			if rule.Right[0] == partJ {
				with = append(with, rule)
			}
		} else if rule.Left == partJ {
			without = append(without, rule)
		}
	}
	for _, rule := range with {
		replacePart := rule.Right[1:]
		rule.Right = append(without[0].Right, replacePart...)
		for i := 1; i < len(without); i++ {
			g.Rules = append(g.Rules, &Rule{Left: rule.Left, Right: append(without[i].Right, replacePart...)})
		}
	}
}

func (g *Grammar) DeleteImmediate(ai string) {
	part := Part{Type: N, Value: ai}
	with := make([]*Rule, 0)
	without := make([]*Rule, 0)
	for _, rule := range g.Rules {
		if rule.Left == part {
			if rule.Right[0] == part {
				with = append(with, rule)
			} else {
				without = append(without, rule)
			}
		}
	}
	if len(with) == 0 {
		return
	}

	pt := ai + NEW
	g.NonTerminals = append(g.NonTerminals, pt)
	part = Part{Type: N, Value: pt}
	epsPart := Part{Type: T, Value: EPS}
	for _, rule := range without {
		if rule.Right[0] != epsPart {
			rule.Right = append(rule.Right, part)
		} else {
			rule.Right[0] = part
		}
	}
	for _, rule := range with {
		rule.Left = part
		rule.Right = append(rule.Right[1:], part)
	}
	g.Rules = append(g.Rules, &Rule{Left: part, Right: []Part{epsPart}})
}

func DeleteLeftRecursion(g *Grammar) *Grammar {
	for i := range g.NonTerminals {
		for j := 0; j < i; j++ {
			g.DeleteIndirect(g.NonTerminals[i], g.NonTerminals[j])
		}
		g.DeleteImmediate(g.NonTerminals[i])
	}
	g.LeftFactor()
	return g
}

func (g *Grammar) LeftFactor() {
	groups := make(map[Part][]*Rule, 0)
	for _, rule := range g.Rules {
		groups[rule.Left] = append(groups[rule.Left], rule)
	}
	for nt, group := range groups {
		nonterms, rules := factorize(nt.Value, group)
		g.NonTerminals = append(g.NonTerminals, nonterms...)
		g.Rules = append(g.Rules, rules...)
	}
}

func factorize(nt string, group []*Rule) ([]string, []*Rule) {
	counter := 0
	epsPart := Part{Type: T, Value: EPS}
	nonterms := make([]string, 0)
	rules := make([]*Rule, 0)
	for {
		n, flag := len(group), false
		var i, j, k int
		for i = 0; i < n-1 && !flag; i++ {
			for j = i + 1; j < n && !flag; j++ {
				if group[i].Right[0] != group[j].Right[0] || group[i].Right[0] == epsPart {
					continue
				}
				flag = true
				limit := len(group[i].Right)
				if limit > len(group[j].Right) {
					limit = len(group[j].Right)
				}
				for k = 1; k < limit; k++ {
					if group[i].Right[k] != group[j].Right[k] {
						break
					}
				}
			}
		}
		i, j = i-1, j-1
		if flag {
			ntNew := fmt.Sprintf("%s%d", nt, counter)
			counter++
			nonterms = append(nonterms, ntNew)
			part := Part{Type: N, Value: ntNew}
			rule := &Rule{
				Left:  Part{Type: N, Value: nt},
				Right: append(group[i].Right[:k], part),
			}
			group = append(group, rule)
			rules = append(rules, rule)
			if len(group[i].Right[k:]) > 0 {
				group[i].Left, group[i].Right = part, group[i].Right[k:]
			} else {
				group[i].Left, group[i].Right = part, []Part{epsPart}
			}
			if len(group[j].Right[k:]) > 0 {
				group[j].Left, group[j].Right = part, group[j].Right[k:]
			} else {
				group[j].Left, group[j].Right = part, []Part{epsPart}
			}
		} else {
			break
		}
	}
	return nonterms, rules
}
