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
var ErrMissingRightParen = errors.New("Expected ')' after expression")
var ErrExpectExpression = errors.New("Expected expression")

var ErrExpectedLeftExpr = errors.New("Expected left expression for the operator")

var ErrOperandNotNumber = errors.New("Operand must be number.")
var ErrOperandsNotSameType = errors.New("Operands must be two numbers or two strings.")

var ErrExpectSemicolonAfterExpr = errors.New("Expected ';' after value or expression.")

var ErrExpectedVariableName = errors.New("Expected variable name.")

var ErrInvalidAssignmentTarget = errors.New("Invalid assignment target")

var ErrExpectedRightBrace = errors.New("Expected '}' after block")

var ErrExpectedLeftParen = errors.New("Expected '('")
var ErrExpectedRightParen = errors.New("Expected ')'")

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

func ErrUndefinedVariable(lexeme string) error {
	mes := fmt.Sprintf("undefined variable %s.", lexeme)
	return errors.New(mes)
}
