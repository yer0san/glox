package token

import "fmt"

type Token struct {
	Token_type 	TokenType
	Lexeme 		string
	Literal 	any
	Line 		int
}

func NewToken(token_type TokenType, lexeme string, literal any, line int) Token{
	return Token{token_type, lexeme, literal, line}
}

func (t *Token) String() string{
	mes := fmt.Sprintf("%v %s %v", t.Token_type, t.Lexeme, t.Literal)
	return mes
}