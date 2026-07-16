package errors

import (
	"errors"
	"fmt"
	. "github.com/yer0san/glox/token"
)


var HadError bool = false

var ErrUnexpectedChar = errors.New("unexpected character")
var ErrunterminatedMultilineComment = errors.New("Unterminated multiline comment")
var ErrUnterminatedString = errors.New("unterminated string")
var ErrMissingRightParen = errors.New("Expect ')' after expression")
var ErrExpectExpression = errors.New("Expect expression")

var ErrExpectedLeftOpr = errors.New("Expected left operand for the operator")

var ErrOperandNotNumber = errors.New("Operand must be number.")
var ErrOperandsNotSameType = errors.New("Operands must be two numbers or two strings.")


func ReportLexingError(line int, err error) {
	Report(line, "", err)
}

func ReportError(token *Token, err error) {
	if token.Token_type == EOF {
		Report(token.Line, " at end", err)
	} else {
		where := fmt.Sprintf("at '%s'", token.Lexeme)
		Report(token.Line, where, err)
	}
}

func Report(line int, where string, err error) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, err)
	HadError = true;
}