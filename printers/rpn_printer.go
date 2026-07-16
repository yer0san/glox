// practice program

package printers 

import (
	"fmt"
	"strings"
	. "github.com/yer0san/glox/expr"
)
type RpnPrinter struct {}

func (p *RpnPrinter) Print(expr ExprForRpn) string {
	return expr.AcceptRpn(p).(string)
}

func (p *RpnPrinter) VisitBinary(expr *BinaryRpn) any {
	return p.buildRpn(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *RpnPrinter) VisitGrouping(expr *GroupingRpn) any {
	return p.buildRpn("", expr.Expression)
}

func (p *RpnPrinter) VisitLiteral(expr *LiteralRpn) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (p *RpnPrinter) VisitUnary(expr *UnaryRpn) any {
	// TODO : not correct, try to fix
	return p.buildRpn(expr.Operator.Lexeme, expr.Right)
}

// real logic
func (p *RpnPrinter) buildRpn(name string, exprs ...ExprForRpn) string {
	var builder strings.Builder

	for _, expr := range exprs {
		builder.WriteString(expr.AcceptRpn(p).(string))
		builder.WriteRune(' ')
	}
	builder.WriteString(name)

	return builder.String()
}


// func main() {
// 	expression2 := &BinaryRpn{
// 		Left : &GroupingRpn{
// 			Expression : &BinaryRpn{
// 				Left : &LiteralRpn{
// 					Value : 1,
// 				},
// 				Operator: Token{
// 					token_type: PLUS,
// 					lexeme: "+",
// 					literal: nil,
// 					line: 1,
// 				}, // i can use the constructor too
// 				Right: &LiteralRpn{
// 					Value: 2,
// 				},
// 			},
// 		},
// 		Operator: NewToken(STAR, "*", nil, 1), // using constructor
// 		Right : & GroupingRpn{
// 			Expression : &BinaryRpn{
// 				Left : &LiteralRpn{
// 					Value : 4,
// 				},
// 				Operator: Token{
// 					token_type: MINUS,
// 					lexeme: "-",
// 					literal: nil,
// 					line: 1,
// 				},
// 				Right: &LiteralRpn{
// 					Value: 3,
// 				},
// 			},
// 		},
// 	}

// 	printer := &RpnPrinter{}
// 	fmt.Println(printer.Print(expression2))
// }

// RPN challenge DONE
