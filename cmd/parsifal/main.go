package main

import (
	"fmt"
	"os"

	"github.com/boxy-pug/parsifal/jsonparser"
)

// simple rudimentary cli for reading a json file and outputting the raw go node
func main() {
	args := os.Args

	if len(args) != 2 {
		fmt.Println("please provide filepath to json file as second arg")
		os.Exit(1)
	}

	file, err := os.Open(args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	lexer := jsonparser.NewLexer(file)
	parser := jsonparser.NewParser(lexer)

	node, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%+v", node)
}
