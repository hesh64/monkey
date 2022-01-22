package main

import (
	"bufio"
	"fmt"
	"io"
	"monkey/internal/evaluator"
	"monkey/internal/lexer"
	"monkey/internal/object"
	"monkey/internal/parser"
	"os"
	user "os/user"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	environment := object.NewEnv()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, environment)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errs []string) {
	for _, msg := range errs {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func main() {
	user, err := user.Current()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	fmt.Printf("Hello %s! this is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	Start(os.Stdin, os.Stdout)
}
