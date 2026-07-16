package lox

import (
	"os"
	"bufio"
	"fmt"
	. "github.com/yer0san/glox/errors"
	. "github.com/yer0san/glox/printers"
	. "github.com/yer0san/glox/interpreter"
	. "github.com/yer0san/glox/parser"
	. "github.com/yer0san/glox/lexer"
)

type Lox struct {
}

func (l *Lox) RunFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	l.run(string(bytes))
	if HadError {
		os.Exit(65)
	}
	return nil
}

func (l *Lox) RunPrompt() error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if ok := scanner.Scan(); !ok{
			break
		}
		line := scanner.Text()

		l.run(line)
		HadError = false
	}
	return scanner.Err()
}

func (l *Lox) run(source string) {
	lexer := NewLexer(source)

	tokens := lexer.ScanTokens()

	parser := Parser{Tokens: tokens}
	expr, err := parser.Parse()

	if err != nil {
		// idk what to do here
		// maybe get the program to go into a panic mode or something
		return
	}
	printer := &AstPrinter{}
	fmt.Println(printer.Print(expr))
	
	interpreter := Interpreter{}

	interpreter.Interpret(expr)
}
