package main

import "errors"

var ErrUnexpectedChar = errors.New("unexpected character")
var ErrunterminatedMultilineComment = errors.New("Unterminated multiline comment")
var ErrUnterminatedString = errors.New("unterminated string")
var ErrMissingRightParen = errors.New("Expect ')' after expression")
var ErrExpectExpression = errors.New("Expect expression")

var ErrExpectedLeftOpr = errors.New("Expected left operand for the operator")