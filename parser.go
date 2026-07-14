package main 

type Parser struct {
	tokens []*Token
	current int // starts at 0 by default ig
}

func NewParser(tokens []*Token) *Parser{
	return &Parser{tokens: tokens}
}

func (p *Parser) expression() (Expr, error) {
	return p.comma()
}

func (p *Parser) comma() (Expr, error) {
	if p.match(COMMA) {
		reportError(p.previous(), ErrExpectedLeftOpr)

		_, err := p.equality()
		
		if err != nil {
			return nil, err
		}
		return nil, ErrExpectedLeftOpr
	}

	expr, err := p.equality()

	if err != nil {
		return nil, err
	}

	for p.match(COMMA) {
		operator := *p.previous()
		right, err := p.equality()

		if err != nil {
			return nil, err
		}

		expr = &Binary {
			Left: expr,
			Operator: operator,
			Right: right,
		}
	}
	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	if p.match(BANG_EQUAL, EQUAL_EQUAL) {
		reportError(p.previous(), ErrExpectedLeftOpr)

		_, err := p.comparison()
		
		if err != nil {
			return nil, err
		}
		return nil, ErrExpectedLeftOpr
	}

	expr, err := p.comparison()

	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		tokenPtr := p.previous()
		operator := *tokenPtr
		right, err := p.comparison()

		if err != nil {
			return nil, err
		}

		expr = &Binary {
			Left: expr,
			Operator: operator,
			Right: right,
		}
	}
	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	if p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		reportError(p.previous(), ErrExpectedLeftOpr)

		_, err := p.term()
		
		if err != nil {
			return nil, err
		}
		return nil, ErrExpectedLeftOpr
	}

	expr, err := p.term()

	if err != nil {
		return nil, err
	}


	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := *p.previous()
		right, err := p.term()

		if err != nil {
			return nil, err
		}

		expr = &Binary {
			Left: expr,
			Operator: operator,
			Right: right,
		}
	}
	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	if p.match(MINUS, PLUS) {
		reportError(p.previous(), ErrExpectedLeftOpr)

		_, err := p.factor()
		
		if err != nil {
			return nil, err
		}
		return nil, ErrExpectedLeftOpr
	}

	expr, err := p.factor()

	if err != nil {
		return nil, err
	}
	
	for p.match(MINUS, PLUS) {
		operator := *p.previous()
		right, err := p.factor()

		if err != nil {
			return nil, err
		}

		expr = &Binary{
			Left: expr,
			Operator: operator,
			Right: right,
		}
	}
	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	if p.match(SLASH, STAR) {
		reportError(p.previous(), ErrExpectedLeftOpr)

		_, err := p.unary()
		
		if err != nil {
			return nil, err
		}
		return nil, ErrExpectedLeftOpr
	}

	expr, err := p.unary()

	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		opr := *p.previous()
		right, err := p.unary()

		if err != nil {
			return nil, err
		}

		expr = &Binary{
			Left: expr,
			Operator: opr,
			Right: right,
		}
	}
	return expr, nil
}

func (p *Parser) unary() (Expr, error) {	
	if p.match(BANG, MINUS) {
		opr := *p.previous()
		right, err := p.unary()

		if err != nil {
			return nil, err
		}

		return &Unary{Operator: opr, Right: right}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return &Literal{Value: false}, nil
	}
	if p.match(TRUE) {
		return &Literal{Value: true}, nil
	}
	if p.match(NIL) {
		return &Literal{Value: nil}, nil
	}

	if p.match(NUMBER, STRING) {
		return &Literal{Value: p.previous().literal}, nil 
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.expression()

		if err != nil {
			return nil, err
		}

		_, err = p.consume(RIGHT_PAREN)

		if err != nil {
			return nil, err
		}

		return &Grouping{Expression: expr}, nil
	}

	reportError(p.peek(), ErrExpectExpression)
	return nil, ErrExpectExpression
}

// HELPERS
func (p *Parser) match(tokentypes ...TokenType) bool {
	for _, tokentype := range tokentypes {
		if p.check(tokentype) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokentype TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	token := p.peek()
	return token.token_type == tokentype
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	token := p.peek()
	return token.token_type == EOF
}

func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return p.tokens[p.current-1]
}


func (p *Parser) consume(tknType TokenType) (*Token, error) {
	if p.check(tknType) {
		return p.advance(), nil
	}
	reportError(p.peek(), ErrMissingRightParen)
	return nil, ErrMissingRightParen
}


func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if prev := p.previous(); prev.token_type == SEMICOLON {
			return
		}
		switch p.peek().token_type {
			case CLASS, FUN, VAR, FOR, IF, WHILE, RETURN, PRINT:
				return
		}
		p.advance()
	}
	
}

func (p *Parser) parse() (Expr, error) {
	return p.expression()  
} // entry method
