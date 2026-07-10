package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct {}

func (p *AstPrinter) Print(expr Expr) string {
	return expr.Accept(p).(string)
}

func (p *AstPrinter) VisitBinaryExpr(expr *Binary) any {
	return p.parenthesize(expr.Operator.lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitGroupingExpr(expr *Grouping) any {
	return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) VisitLiteralExpr(expr *Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value) 
}

func (p *AstPrinter) VisitUnaryExpr(expr *Unary) any {
	return p.parenthesize(expr.Operator.lexeme, expr.Right)
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)


	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.Accept(p).(string))
	}

	builder.WriteString(")")

	return builder.String()

}

// TEST

// func main() {
// 	expression := &Binary{
// 		Left: &Unary{
// 			Operator: Token{token_type: MINUS, lexeme: "-", literal: nil, line: 1},
// 			Right:    &Literal{Value: 123.0},
// 		},
// 		Operator: Token{token_type: STAR, lexeme: "*", literal: nil, line: 1},
// 		Right: &Unary{
// 			Operator: Token{token_type: BANG, lexeme: "!", literal: nil, line: 1},
// 			Right: &Literal{Value: 45.67},
// 		},
// 	}

// 	expression2 := &Binary{
// 		Left : &Grouping{
// 			Expression : &Binary{
// 				Left : &Literal{
// 					Value : 1,
// 				},
// 				Operator: Token{
// 					token_type: PLUS,
// 					lexeme: "+",
// 					literal: nil,
// 					line: 1,
// 				}, // i can use the constructor too
// 				Right: &Literal{
// 					Value: 2,
// 				},
// 			},
// 		},
// 		Operator: NewToken(STAR, "*", nil, 1), // using constructor
// 		Right : & Grouping{
// 			Expression : &Binary{
// 				Left : &Literal{
// 					Value : 4,
// 				},
// 				Operator: Token{
// 					token_type: MINUS,
// 					lexeme: "-",
// 					literal: nil,
// 					line: 1,
// 				},
// 				Right: &Literal{
// 					Value: 2,
// 				},
// 			},
// 		},
// 	}


// // the expression was originally --    -123 * 45.67  so the precedence works
// 	printer := &AstPrinter{}
// 	fmt.Println(printer.Print(expression))
// 	fmt.Println(printer.Print(expression2))
// 	// so it already prints Parent Left Right recursivelyf
// }

