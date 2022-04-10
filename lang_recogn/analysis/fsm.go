package analysis

import (
	"errors"
	"fmt"
	"strings"
)

type treeNode struct {
	symbol      string
	left, right *treeNode
}

func (node *treeNode) isLeaf() bool {
	return node.left == nil && node.right == nil
}

func BuildFSM(regexp string) *FSM {
	ndfsm := new(FSM)
	ndfsm.alphabet = getAlphabet(regexp)
	tree := parseExpression(regexp, ndfsm.alphabet)
	printTree(tree)
	graph := &Graph{
		[]*Node{{"n0"}, {"n1"}},
		[]*Link{},
	}
	if isOperation(tree.symbol) {
		err := calcTree(tree, graph)
		if err != nil {
			return nil
		}
	} else {
		graph.Links = append(graph.Links, &Link{"n0", "n1", tree.symbol})
	}
	return ndfsm
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

func isOperation(symbol string) bool {
	switch symbol[0] {
	case union, concat, iteration:
		return true
	default:
		return false
	}
}

func isBlank(symbol string) bool {
	if symbol == "" || symbol == eps {
		return true
	}
	return false
}

func findExprLen(expr string) (int, error) {
	result := 0
	braceCounter := 1
	for i := range expr {
		switch expr[i] {
		case '(':
			braceCounter++
		case ')':
			braceCounter--
		}
		if braceCounter == 0 {
			break
		}
		result++
	}
	if braceCounter != 0 {
		return 0, errors.New("Bad expression: brace mistake")
	}
	return result, nil
}

func parseExpression(expr string, alphabet map[byte]bool) *treeNode {
	treeHead := &treeNode{symbol: eps}
	curNode, curNodeParent := &treeHead.left, treeHead
	var i int
	for i < len(expr) {
		if expr[i] == iteration {
			if *curNode == nil {
				return nil
			}
			*curNode = &treeNode{
				left: *curNode,
				right: &treeNode{
					symbol: eps,
				},
				symbol: string(iteration),
			}
		} else {
			if curNodeParent.right != nil {
				*curNode = &treeNode{
					left: *curNode,
				}
				curNodeParent = *curNode
				curNode = &curNodeParent.right
			}
			if *curNode != nil {
				curNode = &curNodeParent.right
			}
			if _, ok := alphabet[expr[i]]; ok {
				*curNode = &treeNode{
					symbol: string(expr[i]),
				}
			} else if expr[i] == union || expr[i] == concat {
				if curNodeParent.left == nil {
					return nil
				}
				curNodeParent.symbol = string(expr[i])
			} else if expr[i] == '(' {
				bracedExprLen, err := findExprLen(expr[i+1:])
				if err != nil {
					return nil
				}
				subTree := parseExpression(expr[i+1:i+bracedExprLen+1], alphabet)
				if subTree != nil {
					*curNode = subTree
					i += bracedExprLen + 1
				} else {
					return nil
				}
			}
			if isBlank(curNodeParent.symbol) && curNodeParent.right != nil {
				curNodeParent.symbol = string(concat)
			}
		}
		i++
	}

	if isBlank(treeHead.symbol) && treeHead.left != nil {
		treeHead = treeHead.left
	}
	return treeHead
}

func Find(nodes []*Node, name string) int {
	for i, n := range nodes {
		if name == n.Name {
			return i
		}
	}
	return len(nodes)
}

func getSubGraph(graph *Graph, node string) (string, string) {
	if index := strings.IndexByte(node, separator); index >= 0 {
		return node[:index], node[index+1:]
	}
	start := fmt.Sprintf("n%d", len(graph.Nodes))
	end := fmt.Sprintf("n%d", len(graph.Nodes)+1)
	graph.Nodes = append(graph.Nodes, &Node{start}, &Node{end})
	graph.Links = append(graph.Links, &Link{Source: start, Target: end, Value: node})
	return start, end
}

func applyOperation(graph *Graph, operation, left, right string) (string, error) {
	leftStart, leftEnd := getSubGraph(graph, left)
	rightStart, rightEnd := getSubGraph(graph, right)
	switch operation[0] {
	case iteration:
		graph.Links = append(graph.Links, &Link{rightStart, leftStart, eps},
			&Link{leftEnd, rightEnd, eps},
			&Link{leftEnd, leftStart, eps})
		return fmt.Sprintf("%s%c%s", rightStart, separator, rightEnd), nil
	case concat:
		for i := range graph.Links {
			if graph.Links[i].Source == rightStart {
				graph.Links[i].Source = leftEnd
			}
			if graph.Links[i].Target == rightStart {
				graph.Links[i].Target = leftEnd
			}
		}
		return fmt.Sprintf("%s%c%s", leftStart, separator, rightEnd), nil
	case union:
		start := fmt.Sprintf("n%d", len(graph.Nodes))
		end := fmt.Sprintf("n%d", len(graph.Nodes)+1)
		graph.Nodes = append(graph.Nodes, &Node{start}, &Node{end})
		graph.Links = append(graph.Links, &Link{start, leftStart, eps},
			&Link{start, rightStart, eps},
			&Link{leftEnd, end, eps},
			&Link{rightEnd, end, eps},
		)
		return fmt.Sprintf("%s%c%s", start, separator, end), nil
	}
	return "", errors.New("Invalid operation")
}

func calcTree(head *treeNode, graph *Graph) error {
	stack := [](*treeNode){head}
	for len(stack) > 0 {
		curNode := stack[len(stack)-1]
		if curNode.left.isLeaf() && curNode.right.isLeaf() {
			nameStartEnd, err := applyOperation(graph, curNode.symbol, curNode.left.symbol, curNode.right.symbol)
			if err != nil {
				return err
			}
			*curNode = treeNode{
				symbol: nameStartEnd,
			}
			stack = stack[:len(stack)-1]
		} else {
			if isOperation(curNode.left.symbol) {
				stack = append(stack, curNode.left)
			}
			if isOperation(curNode.right.symbol) {
				stack = append(stack, curNode.right)
			}
		}
	}
	return nil
}

func printTree(head *treeNode) {
	queue := make([]*treeNode, 0)
	queue = append(queue, head)
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node.left != nil {
			queue = append(queue, node.left)
		}
		if node.right != nil {
			queue = append(queue, node.right)
		}
		fmt.Println(node.symbol)
	}
}
