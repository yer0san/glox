package main

import "strconv"

type Lexer struct {
	source 		string
	tokens 		[]*Token
	start 		int
	current 	int
	line 		int
}

func NewLexer(source string) *Lexer{
	tokens := []*Token{}

	return &Lexer{
		source:  source, 
		tokens:  tokens,
		start:   0,
		current: 0,
		line:    1,
	}
}

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	// "super":  SUPER,
	// "this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

func (l *Lexer) scanTokens() []*Token {
	for !l.isAtEnd() {
		l.start = l.current
		l.scanToken()
	}
	tok := NewToken(EOF, "", nil, l.line)
	l.tokens = append(l.tokens, &tok)
	return l.tokens
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) scanToken() {

	var c rune = l.advance()
	switch c {
		case '(':
			l.addToken(LEFT_PAREN)
		case ')':
			l.addToken(RIGHT_PAREN)
		case '{':
			l.addToken(LEFT_BRACE)
		case '}':
			l.addToken(RIGHT_BRACE)
		case ',':
			l.addToken(COMMA)
		case '.':
			l.addToken(DOT)
		case '-':
			l.addToken(MINUS)
		case '+':
			l.addToken(PLUS)
		case ';':
			l.addToken(SEMICOLON)
		case '*':
			l.addToken(STAR)
		
		case '!':
			b := l.match('=')
			if b {
				l.addToken(BANG_EQUAL)
			} else {
				l.addToken(BANG)
			}
		case '=':
			b := l.match('=')
			if b {
				l.addToken(EQUAL_EQUAL)
			} else {
				l.addToken(EQUAL)
			}
		case '<':
			b := l.match('=')
			if b {
				l.addToken(LESS_EQUAL)
			} else {
				l.addToken(LESS)
			}
		case '>':
			b := l.match('=')
			if b {
				l.addToken(GREATER_EQUAL)
			} else {
				l.addToken(GREATER)
			}
		case '/':
			if l.match('/'){
				// FIX: peek return - should be newline but idk how to do that
				for l.peek() != '\n' && !l.isAtEnd() {
					l.advance()
				}
			} else {
				l.addToken(SLASH)
			}
		case ' ':
		case '\r': // FIX: what is \r
		case '\t': //FIX: what is \t
			break
		case '\n': // FIX: what is \n
			l.line++
		case '"':
			l.str()

		default:
			if l.isDigit(c) {
				l.number()
			} else if l.isAlpha(c) {
				l.identifier();
			} else {
				reportError(l.line, "unexpected character.")
			}
	
	}
}

func (l *Lexer) identifier() {
	for l.isAlphaNumberic(l.peek()) {
		l.advance()
	}
	text := l.source[l.start:l.current]
	tokentype, ok := keywords[text]

	if !ok {
		tokentype = IDENTIFIER
	}
	l.addToken(tokentype)
}

func (l *Lexer) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
}

func (l *Lexer) isAlphaNumberic(c rune) bool {
	return l.isAlpha(c) || l.isDigit(c)
}

func (l *Lexer) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (l *Lexer) number() {
	for l.isDigit(l.peek()) {
		l.advance()
	}

	if l.peek() == '.' && l.isDigit(l.peekNext()) {
		l.advance()
		for l.isDigit(l.peek()) {
			l.advance()
		}
	}
	
	value, err := strconv.ParseFloat(l.source[l.start:l.current], 64)

	if err != nil {
		panic(err)
	}

	l.addT(NUMBER, value)
}

func (l *Lexer) peekNext() rune {
	if l.current + 1 >= len(l.source) {
		return '0'
	}
	return rune(l.source[l.current+1])
}

func (l *Lexer) match(expected rune) bool {
	if l.isAtEnd() {
		return false
	}

	if rune(l.source[l.current]) != expected{
		return false
	}

	l.current++
	return true
}

func (l *Lexer) peek() rune {
	if l.isAtEnd(){
		// FIX: what do we return here?
		return '\n'
	}
	return rune(l.source[l.current])
}

func (l *Lexer) advance() rune {
	pos := l.source[l.current]
	l.current++
	return rune(pos)
}

func (l *Lexer) addToken(tokenType TokenType) {
	l.addT(tokenType, nil)
}

func (l *Lexer) addT(tokenType TokenType, literal any) {
	text := l.source[l.start: l.current]
	token := &Token{
				token_type: tokenType, 
				lexeme: text,
				literal: literal, 
				line: l.line,
			}

	l.tokens = append(l.tokens, token)
}

func (l *Lexer) str() {
	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\n' { // FIX: again, new line??
			l.line++
			l.advance()
		}
	}

	if l.isAtEnd() {
		reportError(l.line, "unterminated string")
		return
	}

	// The closing "
	l.advance()

	value := l.source[l.start+1:l.current-1]
	l.addT(STRING, value)
}

