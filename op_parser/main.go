package main

import (
	"bufio"
	"fmt"
	"operator_precedence/parser"
	"os"
	"strings"
)

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
		if err != nil || buffer == "\n" {
			break
		}
		data += buffer
	}
	data = strings.Join(strings.Fields(data), "")

	if result, err := parser.Parse(data); err == nil {
		output.WriteString(result + "\n")
	} else {
		fmt.Println("Error while parsing:", err.Error())
	}
}
