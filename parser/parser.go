package parser

import (
	"fmt"
	. "github.com/yer0san/glox/errors"
	. "github.com/yer0san/glox/expr"
	"github.com/yer0san/glox/stmt"
	. "github.com/yer0san/glox/token"
)

type Parser struct {
	Tokens []*Token
	current int // starts at 0 by default ig
}

func NewParser(tokens []*Token) *Parser{
	return &Parser{Tokens: tokens}
}

func (p *Parser) declaration() (stmt.Stmt, error) {
	if p.match(VAR) {
		val, err := p.varDecl()

		if err != nil {
			p.synchronize()
			return nil, err
		}

		return val, nil
	}
	val, err :=  p.statement()

	if err != nil {
		p.synchronize()
		return nil, err
	}

	return val, nil
}

func (p *Parser) varDecl() (stmt.Stmt, error) {
	name, err := p.consume(IDENTIFIER)

	if err != nil {
		return nil, err
	}
	
	var initializer Expr

	if (p.match(EQUAL)) {
		initializer, err = p.expression();

		if err != nil {
			return nil, err
		}
	}

	p.consume(SEMICOLON)
	return &stmt.Var{Name: name, Initializer: initializer}, nil
}

func (p *Parser) statement() (stmt.Stmt, error) {
	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(LEFT_BRACE) {
		statements, err := p.block()

		if err != nil {
			return nil, err
		}
		return &stmt.Block{Statements:statements}, nil
	}

	if p.match(IF) {
		return p.ifStatement()
	}

	if p.match(WHILE) {
		return p.whileStatement()
	}

	if p.match(FOR) {
		return p.forStatement()
	}
	return p.exprStatement()
}

func (p *Parser) forStatement() (stmt.Stmt, error) {
	p.consume(LEFT_PAREN)

	var init stmt.Stmt
	var condition Expr
	var increment Expr
	var err error

	if p.match(SEMICOLON) {
		// nothing i guess
	} else if p.match(VAR) {
		fmt.Println("varDecl...")
		init, err = p.varDecl()

		if err != nil {
			return nil, err
		} 
	} else {
		init, err = p.exprStatement()

		if err != nil {
			return nil, err
		}
	}
	if !p.check(SEMICOLON) {
		fmt.Println("condition...")
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	p.consume(SEMICOLON)
	if !p.check(RIGHT_PAREN) {
		fmt.Println("increment...")
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	p.consume(RIGHT_PAREN)

	if p.match(LEFT_BRACE) {
		body, err :=  p.block()

		if err != nil {
			return nil, err
		}

		return &stmt.For{
			Init: init,
			Condition: condition,
			Increment: increment,
			Body: &stmt.Block{Statements: body},
		}, nil
	}

	body, err := p.statement()

	if err != nil {
		return nil, err
	}
	return &stmt.For{
			Init: init,
			Condition: condition,
			Increment: increment,
			Body: body,
		}, nil
}


func (p *Parser) whileStatement() (stmt.Stmt, error) {
	p.consume(LEFT_PAREN)
	condition, err := p.expression()

	if err != nil {
		return nil, err
	}

	p.consume(RIGHT_PAREN)

	if p.match(LEFT_BRACE) {
		body, err :=  p.block()

		if err != nil {
			return nil, err
		}

		return &stmt.While{
			Condition: condition,
			Body: &stmt.Block{Statements: body},
		}, nil
	}

	body, err := p.statement()

	if err != nil {
		return nil, err
	}
	return &stmt.While{
		Condition: condition,
		Body: body,
	}, nil
}

func (p *Parser) ifStatement() (stmt.Stmt, error) {
	p.consume(LEFT_PAREN)
	condition, err := p.expression()

	if err != nil {
		return nil, err
	}

	p.consume(RIGHT_PAREN)

	thenBranch, err := p.statement()

	if err != nil {
		return nil, err
	}

	var elseBranch stmt.Stmt

	if p.match(ELSE) {
		elseBranch, err = p.statement()

		if err != nil {
			return nil, err
		}
	}
	return &stmt.If{
		Condition: condition, 
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (p *Parser) block() ([]stmt.Stmt, error) {
	var statements []stmt.Stmt

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() { 
		statement, err := p.declaration()

		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}
	_, err := p.consume(RIGHT_BRACE)

	if err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser) printStatement() (stmt.Stmt, error) {
	if p.match(VAR) {
		// the variable name is a Literal
		// ??
	}
	value, err := p.expression()

	if err != nil {
		return nil, err
	}

	p.consume(SEMICOLON)
	return &stmt.Print{Expr: value}, nil
}

func (p *Parser) exprStatement() (stmt.Stmt, error) {
	value, err := p.expression()

	if err != nil {
		return nil, err
	}

	p.consume(SEMICOLON)
	return &stmt.Expression{Expr: value}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.logic_or()

	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		equals := p.previous()
		value, err := p.assignment()

		if err != nil {
			return nil, err
		}

		if variable, ok := expr.(*Variable); ok {
			return &Assign{Name: variable.Name, Value: value}, nil
		}
		ReportError(equals, ErrInvalidAssignmentTarget)
		// return nil, ErrInvalidAssignmentTarget
	}

	return expr, nil
}

func (p *Parser) logic_or() (Expr, error) {
	if p.match(OR) {
		ReportError(p.previous(), ErrExpectedLeftExpr)
		_, err := p.logic_and()
		
		if err != nil {
			return nil, err
		}

		return nil, ErrExpectedLeftExpr
	}

	expr, err := p.logic_and()

	if err != nil {
		return nil, err
	}
	for p.match(OR) {
		opr := p.previous()

		right, err := p.logic_and()

		if err != nil {
			return nil, err
		}

		expr = &Logical{
			Left: expr,
			Operator: opr,
			Right: right,
		}
	}

	return expr, nil

}

func (p *Parser) logic_and() (Expr, error) {
	if p.match(AND) {
		ReportError(p.previous(), ErrExpectedLeftExpr)
		_, err := p.logic_and()
		
		if err != nil {
			return nil, err
		}

		return nil, ErrExpectedLeftExpr
	}
	
	expr, err := p.comma() // short circuiting is the interpreter's job

	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		opr := p.previous()
		right, err := p.comma()

		if err != nil {
			return nil, err
		}

		expr = &Logical{
			Left: expr,
			Operator: opr,
			Right: right,
		}
	}
	return expr, nil
}

func (p *Parser) comma() (Expr, error) {
	if p.match(COMMA) {
		ReportError(p.previous(), ErrExpectedLeftExpr)

		_, err := p.equality()
		
		if err != nil {
			return nil, err
		}
		return nil, ErrExpectedLeftExpr
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
			Operator: &operator,
			Right: right,
		}
	}
	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	if p.match(BANG_EQUAL, EQUAL_EQUAL) {
		ReportError(p.previous(), ErrExpectedLeftExpr)

		_, err := p.comparison()
		
		if err != nil {
			return nil, err
		}
		return nil, ErrExpectedLeftExpr
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
			Operator: &operator,
			Right: right,
		}
	}
	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	if p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		ReportError(p.previous(), ErrExpectedLeftExpr)

		_, err := p.term()
		
		if err != nil {
			return nil, err
		}
		return nil, ErrExpectedLeftExpr
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
			Operator: &operator,
			Right: right,
		}
	}
	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	if p.match(MINUS, PLUS) {
		ReportError(p.previous(), ErrExpectedLeftExpr)

		_, err := p.factor()
		
		if err != nil {
			return nil, err
		}
		return nil, ErrExpectedLeftExpr
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
			Operator: &operator,
			Right: right,
		}
	}
	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	if p.match(SLASH, STAR) {
		ReportError(p.previous(), ErrExpectedLeftExpr)

		_, err := p.unary()
		
		if err != nil {
			return nil, err
		}
		return nil, ErrExpectedLeftExpr
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
			Operator: &opr,
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

		return &Unary{Operator: &opr, Right: right}, nil
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

	if p.match(IDENTIFIER) {
		return &Variable{Name: p.previous()}, nil
	}

	if p.match(NUMBER, STRING) {
		return &Literal{Value: p.previous().Literal}, nil 
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

	ReportError(p.peek(), ErrExpectExpression)
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
	return token.Token_type == tokentype
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	token := p.peek()
	return token.Token_type == EOF
}

func (p *Parser) peek() *Token {
	return p.Tokens[p.current]
}

func (p *Parser) previous() *Token {
	return p.Tokens[p.current-1]
}

func (p *Parser) consume(tknType TokenType) (*Token, error) {
	if p.check(tknType) {
		return p.advance(), nil
	}

	if tknType == RIGHT_BRACE {
		ReportError(p.peek(), ErrExpectedRightBrace)
		return nil, ErrExpectedRightBrace
	}

	if tknType == LEFT_PAREN {
		ReportError(p.peek(), ErrExpectedLeftParen)
		return nil, ErrExpectedLeftParen
	}
	if tknType == RIGHT_PAREN {
		ReportError(p.peek(), ErrExpectedRightParen)
		return nil, ErrExpectedRightParen
	}

	if tknType == IDENTIFIER {
		ReportError(p.peek(), ErrExpectedVariableName)
		return nil, ErrExpectedVariableName
	}

	if tknType == SEMICOLON {
		ReportError(p.peek(), ErrExpectSemicolonAfterExpr)
		return nil, ErrExpectSemicolonAfterExpr
	}
	ReportError(p.peek(), ErrMissingRightParen)
	return nil, ErrMissingRightParen
}


func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if prev := p.previous(); prev.Token_type == SEMICOLON {
			return
		}
		switch p.peek().Token_type {
			case CLASS, FUN, VAR, FOR, IF, WHILE, RETURN, PRINT:
				return
		}
		p.advance()
	}
	
}

func (p *Parser) Parse() ([]stmt.Stmt, error) {
	var statements []stmt.Stmt
	for !p.isAtEnd() {
		stmt, err := p.declaration()

		if err != nil{
			return nil, err
		}
		statements = append(statements, stmt)
	}
	return statements, nil
} // entry method
