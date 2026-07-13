package main

import (
	"bufio"
	"fmt"
	"os"

)

type Lox struct {
}

var hadError bool = false

// he did a different design on the l := Lox{}, keep that in mind if this breaks
// vm := Lox{} --global, idek what that means yet :)

func main(){
	l := Lox{}
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox [script]");
		os.Exit(64)
	}
	if len(os.Args) == 2 {
		l.runFile(os.Args[1])
	} else {
		l.runPrompt();
	}
}

func (l *Lox) runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	l.run(string(bytes))
	if hadError {
		os.Exit(65)
	}
	return nil
}

func (l *Lox) runPrompt() error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if ok := scanner.Scan(); !ok{
			break
		}
		line := scanner.Text()

		l.run(line)
		hadError = false
	}
	return scanner.Err()
}

func (l *Lox) run(source string) {
	lexer := NewLexer(source)

	tokens := lexer.scanTokens()

	parser := Parser{tokens: tokens}
	expr, err := parser.parse()

	if err != nil {
		// idk what to do here
		// maybe get the program to go into a panic mode or something
		return
	}
	printer := &AstPrinter{}
	fmt.Println(printer.Print(expr))
}

func reportError(line int, err error) {
	report(line, "", err)
}

func reportParserError(token *Token, err error) {
	if token.token_type == EOF {
		report(token.line, " at end", err)
	} else {
		where := fmt.Sprintf("at '%s'", token.lexeme)
		report(token.line, where, err)
	}
} // check it out later

func report(line int, where string, err error) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, err)
	hadError = true;
}


