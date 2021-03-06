// Dito is a Toy interpreted programming language for fun.
// : )
package main

import (
	"github.com/dito/src/eval"
	"github.com/dito/src/object"
	"github.com/dito/src/parser"
	"github.com/dito/src/repl"
	"github.com/dito/src/scanner"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
)

func main() {
	args := os.Args[1:] // args without program.
	// https://gobyexample.com/command-line-flags
	if len(args) > 0 {
		filepath := args[0]
		file, err := ioutil.ReadFile(filepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
		execFile(string(file), os.Stdout)
		return
	}
	welcomeMsg(repl.QUIT)
	repl.Start(os.Stdin, os.Stdout)
}

func execFile(file string, out io.Writer) {
	env := object.NewEnvironment()
	l := scanner.Init(file)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		p.PrintParseErrors(out, p.Errors())
		return
	}
	evaluated := eval.Eval(program, env)
	if evaluated != nil && evaluated != object.NONE {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func welcomeMsg(quit string) {
	fmt.Printf("Dito Interactive Shell V0.01 on %s\n", runtime.GOOS)
	fmt.Printf("Enter '%s' to quit. Help is coming soon...\n", quit)
}
