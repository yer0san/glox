package main

import "errors"

var ErrUnexpectedChar = errors.New("unexpected character")
var ErrunterminatedMultilineComment = errors.New("Unterminated multiline comment") // a bit long, meh
var ErrUnterminatedString = errors.New("unterminated string")
var ErrMissingRightParen = errors.New("Expect ')' after expression")
var ErrExpectExpression = errors.New("Expect expression")