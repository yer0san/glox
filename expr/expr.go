package expr

import . "github.com/yer0san/glox/token"

type Expr interface {
	Accept(v Visitor) (any, error)
}

type Binary struct {
	Left Expr
	Operator Token
	Right Expr
}

type Literal struct {
	Value any
}

type Grouping struct {
	Expression Expr
}

type Unary struct {
	Operator Token
	Right Expr
}

type Visitor interface {
	VisitBinaryExpr(expr *Binary) (any, error)
	VisitLiteralExpr(expr *Literal) (any, error)
	VisitGroupingExpr(expr *Grouping) (any, error)
	VisitUnaryExpr(expr *Unary) (any, error)
}

func (b *Binary) Accept(v Visitor) (any, error) {
	res, _ := v.VisitBinaryExpr(b)
	return res, nil
}

func (g *Grouping) Accept(v Visitor) (any, error) {
	res, _ := v.VisitGroupingExpr(g)
	return res, nil
}

func (l *Literal) Accept(v Visitor) (any, error) {
	res, _ := v.VisitLiteralExpr(l)
	return res, nil
}

func (u *Unary) Accept(v Visitor) (any, error) {
	res, _ := v.VisitUnaryExpr(u)
	return res, nil
}



// ----------------------------------------------------
// Ignore
// this is practice

type ExprForRpn interface {
	AcceptRpn(v VisitorForRpn) any
}

type BinaryRpn struct {
	Left ExprForRpn
	Operator Token
	Right ExprForRpn
}

type LiteralRpn struct {
	Value any
}

type GroupingRpn struct {
	Expression ExprForRpn
}

type UnaryRpn struct {
	Operator Token
	Right ExprForRpn
}

type VisitorForRpn interface {
	VisitBinary(expr *BinaryRpn) any
	VisitLiteral(expr *LiteralRpn) any
	VisitGrouping(expr *GroupingRpn) any
	VisitUnary(expr *UnaryRpn) any
} // this is a temporary practice code, or not

func (b *BinaryRpn) AcceptRpn(v VisitorForRpn) any {
	return v.VisitBinary(b)
}

func (g *GroupingRpn) AcceptRpn(v VisitorForRpn) any {
	return v.VisitGrouping(g)
}

func (l *LiteralRpn) AcceptRpn(v VisitorForRpn) any {
	return v.VisitLiteral(l)
}

func (u *UnaryRpn) AcceptRpn(v VisitorForRpn) any {
	return v.VisitUnary(u)
}


// ----------------------------------------------------