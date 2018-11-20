package eval

import (
	"github.com/dito/src/ast"
	"github.com/dito/src/object"
	"github.com/dito/src/parser"
	"github.com/dito/src/scanner"
	"io/ioutil"
	"os"
)

func evalImportStatement(node *ast.ImportStatement, env *object.Environment) object.Object {
	filename := node.Value
	filepath := "lib/" + filename + ".dito"
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return object.NewError("Import File %s could not be opened.", filepath)
	}
	il := scanner.Init(string(file))
	ip := parser.New(il)
	iprogram := ip.ParseProgram()

	if len(ip.Errors()) != 0 {
		ip.PrintParseErrors(os.Stderr, ip.Errors())
		return object.NewError("Could not import file due to parse errors.")
	}
	Eval(iprogram, env)
	// return object.NONE
	return nil
}
