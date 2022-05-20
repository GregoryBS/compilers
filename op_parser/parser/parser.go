package parser

import (
	"errors"
)

func Parse(expr string) (string, error) {
	table := getTable()
	stack := make([]string, 1)
	stack[0] = MARKER
	var result string
	op, idx := readExpr(expr, 0)
	for len(stack) > 1 || op != MARKER {
		op1, op2 := stack[len(stack)-1], op
		if IsAtom(op1) {
			op1 = ATOM
		}
		if IsAtom(op2) {
			op2 = ATOM
		}
		switch table[op1][op2] {
		case LESS, EQUAL:
			stack = append(stack, op)
			op, idx = readExpr(expr, idx)
		case MORE:
			for {
				result += prn(stack[len(stack)-1])
				op1, op2 := stack[len(stack)-2], stack[len(stack)-1]
				if IsAtom(op1) {
					op1 = ATOM
				}
				if IsAtom(op2) {
					op2 = ATOM
				}
				stack = stack[:len(stack)-1]
				if table[op1][op2] == LESS {
					break
				}
			}
		default:
			return "", errors.New(parseErrors[table[op1][op2]])
		}
	}
	return result, nil
}

func prn(s string) string {
	if s != LBRACKET && s != RBRACKET {
		return s + " "
	}
	return ""
}

func readExpr(expr string, idx int) (string, int) {
	begin := expr[idx:]
	if IsAtom(begin[:1]) {
		end := 1
		for len(begin) > end && IsAtom(begin[end:end+1]) {
			end++
		}
		return begin[:end], idx + end
	} else if len(begin) >= 2 && (begin[:2] == LE || begin[:2] == GE || begin[:2] == NE) {
		return begin[:2], idx + 2
	}
	return begin[:1], idx + 1
}
