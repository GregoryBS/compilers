package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"grammar_transform/algorithms"
	"io/ioutil"
	"os"
)

func main() {
	epsFlag := flag.Bool("e", false, "eps-rules deleting algorithm")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 || len(args) > 2 {
		flag.Usage()
		fmt.Println("in_file [out_file]")
		return
	}

	var err error
	g := new(algorithms.Grammar)
	output := os.Stdout
	switch len(args) {
	case 2:
		if output, err = os.Create(args[1]); err != nil {
			fmt.Println("Unable to create the output file. Result will be returned to stdout.")
			output = os.Stdout
		}
		defer output.Close()
		fallthrough
	case 1:
		data, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Println("Unable to read the input file.")
			return
		}
		err = json.Unmarshal(data, g)
		if err != nil {
			fmt.Println("Invalid json.")
			return
		}
	}

	var result *algorithms.Grammar
	if *epsFlag {
		result = algorithms.DeleteEpsRules(g)
	} else {
		result = algorithms.DeleteLeftRecursion(g)
	}
	if result != nil {
		data, _ := json.Marshal(result)
		output.Write(data)
	} else {
		fmt.Println("Error occured while grammar transforming.")
	}
	return
}
