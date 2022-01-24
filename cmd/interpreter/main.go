package main

import (
	"io"
	"monkey/internal/evaluator"
	"monkey/internal/lexer"
	"monkey/internal/object"
	"monkey/internal/parser"
	"os"
)

func readFirstArg() string {
	if len(os.Args) <= 1 {
		panic("call the repel main")
	}
	return os.Args[1]
}

func readFile(filename string) (string, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return string(file), nil
}

func printParserErrors(out io.Writer, errs []string) {
	for _, msg := range errs {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func main() {
	environment := object.NewEnv()

	filename := readFirstArg()
	fileContent, err := readFile(filename)
	if err != nil {
		panic(err)
	}

	l := lexer.New(fileContent)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		return
	}

	evaluated := evaluator.Eval(program, environment)
	if evaluated != nil {
		io.WriteString(os.Stdout, evaluated.Inspect())
		io.WriteString(os.Stdout, "\n")
	}
}
