package main

import (
	"bufio"
	"fmt"
	"lang_recognize/analysis"
	"os"
)

func Menu() {
	fmt.Println("1 - Input regexp for analysis")
	fmt.Println("2 - Input word-chain for FSM-modelling")
	fmt.Println("0 - Exit")
}

func main() {
	Menu()
	var fsm *analysis.FSM
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		switch scanner.Text() {
		case "1":
			if scanner.Scan() {
				fsm = analysis.BuildFSM(scanner.Text())
				if fsm == nil {
					fmt.Println("Error occured while building fsm")
				} else {
					fsm.GetGraph()
				}
			}
		case "2":
			if scanner.Scan() {
				if fsm != nil {
					fmt.Println("word is acceptable:", analysis.Modelling(fsm, scanner.Text()))
				} else {
					fmt.Println("Error occured while modelling fsm")
				}
			}
		case "0":
			os.Exit(0)
		default:
			fmt.Println("Input error")
		}
		Menu()
	}
}
