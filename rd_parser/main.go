package main

import (
	"bufio"
	"fmt"
	"os"
	"recursive_descent/parser"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

func parseResult(g *graphviz.Graphviz, graph *cgraph.Graph) {
	fmt.Println()
	if msg := recover(); msg != nil {
		fmt.Println(msg)
	} else {
		fmt.Println("Parsing successful")
		if err := g.RenderFilename(graph, graphviz.PNG, "graph.png"); err != nil {
			fmt.Println("Cannot save graph in file", err)
		}
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

	g := graphviz.New()
	defer g.Close()
	graph, err := g.Graph()
	if err != nil {
		fmt.Println("cannot create graphviz graph", err)
		return
	}
	defer graph.Close()

	p := &parser.Parser{
		Data:  strings.Join(strings.Fields(data), ""),
		File:  output,
		Graph: &parser.Graph{graph, []*cgraph.Node{}, []*cgraph.Edge{}},
	}
	defer parseResult(g, graph)
	p.Block()
	return
}
