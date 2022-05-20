package parser

import "unicode"

const (
	CIRKUM   = "^"
	LBRACKET = "("
	RBRACKET = ")"
	PLUS     = "+"
	MINUS    = "-"
	MUL      = "*"
	DIV      = "/"
	MOD      = "%"
	LT       = "<"
	LE       = "<="
	EQ       = "="
	GE       = ">="
	GT       = ">"
	NE       = "<>"
	MARKER   = "$"
	ATOM     = "atom"
)

const (
	LESS  = '<'
	EQUAL = '='
	MORE  = '>'
)

//Marker must be with the lowest and atom with the biggest priority
var terminals = []string{MARKER, LBRACKET, EQ, LT, GT, NE, LE, GE, PLUS, MINUS, MUL, DIV, MOD, CIRKUM, RBRACKET, ATOM}
var priority = []int{0, 1, 2, 2, 2, 2, 2, 2, 3, 3, 4, 4, 4, 5, 6, 7}
var leftAssociative = []bool{false, false, false, false, false, false, false, false, true, true, true, true, true, false, false, false}
var parseErrors = map[byte]string{
	1: "Missing operator",
	2: "Missing expression",
	3: "Unexpected right bracket",
	4: "Missing right bracket",
}

func IsAtom(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func getTable() map[string]map[string]byte {
	result := make(map[string]map[string]byte, len(terminals))
	for _, t := range terminals {
		result[t] = make(map[string]byte, len(terminals))
	}

	for i, ti := range terminals {
		for j, tj := range terminals {
			if (ti == ATOM || ti == RBRACKET) && (tj == LBRACKET || tj == ATOM) {
				result[ti][tj] = 1
			} else if ti == MARKER && tj == MARKER {
				result[ti][tj] = 2
			} else if ti == MARKER && tj == RBRACKET {
				result[ti][tj] = 3
			} else if ti == LBRACKET && tj == MARKER {
				result[ti][tj] = 4
			} else if ti == LBRACKET && tj == LBRACKET {
				result[ti][tj] = LESS
			} else if ti == RBRACKET && tj == RBRACKET {
				result[ti][tj] = MORE
			} else if ti == LBRACKET && tj == RBRACKET {
				result[ti][tj] = EQUAL
			} else if priority[i] > priority[j] {
				result[ti][tj] = MORE
			} else if priority[i] < priority[j] {
				result[ti][tj] = LESS
			} else if priority[i] == priority[j] {
				if leftAssociative[i] {
					result[ti][tj] = MORE
				} else {
					result[ti][tj] = LESS
				}
			}
		}
	}
	return result
}
