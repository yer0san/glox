package main

import (
	"fmt"
	"os"
	"github.com/yer0san/glox/lox"
)

// he did a different design on the l := Lox{}, keep that in mind if this breaks
// vm := Lox{} --global, idek what that means yet :)

func main(){
	l := lox.Lox{}
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox [script]");
		os.Exit(64)
	}
	if len(os.Args) == 2 {
		l.RunFile(os.Args[1])
	} else {
		l.RunPrompt();
	}
}



