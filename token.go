package main

import "fmt"

type Token struct {
	token_type 	TokenType
	lexeme 		string
	literal 	any
	line 		int
}

// How do we define a constructor method ??

func NewToken(token_type TokenType, lexeme string, literal any, line int) Token{
	return Token{token_type, lexeme, literal, line}
}

func (t *Token) String() string{
	mes := fmt.Sprintf("%v %s %v", t.token_type, t.lexeme, t.literal)
	return mes
}