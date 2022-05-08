package main

import (
	"bufio"
	"fmt"
	"os"
	"recursive_descent/parser"
	"strings"
)

func parseResult() {
	fmt.Println()
	if msg := recover(); msg != nil {
		fmt.Println(msg)
	} else {
		fmt.Println("Parsing successful")
	}
}

func main() {
	var err error
	var input, output *os.File
	switch len(os.Args) {
	case 3:
		if output, err = os.Create(os.Args[2]); err != nil {
			output = os.Stdout
		}
		defer output.Close()
		fallthrough
	case 2:
		if input, err = os.Open(os.Args[1]); err != nil {
			input = os.Stdin
		}
		defer input.Close()
		if output == nil {
			output = os.Stdout
		}
	default:
		input, output = os.Stdin, os.Stdout
	}

	var data string
	reader := bufio.NewReader(input)
	for {
		buffer, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		data += buffer
	}

	p := &parser.Parser{Data: strings.Join(strings.Fields(data), ""), File: output}
	defer parseResult()
	p.Block()
	return
}
