package main


type Expr interface {
	Accept(v Visitor) any
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
	VisitBinaryExpr(expr *Binary) any
	VisitLiteralExpr(expr *Literal) any
	VisitGroupingExpr(expr *Grouping) any
	VisitUnaryExpr(expr *Unary) any
} // TODO : possible change to generics

func (b *Binary) Accept(v Visitor) any {
	return v.VisitBinaryExpr(b)
}

func (g *Grouping) Accept(v Visitor) any {
	return v.VisitGroupingExpr(g)
}

func (l *Literal) Accept(v Visitor) any {
	return v.VisitLiteralExpr(l)
}

func (u *Unary) Accept(v Visitor) any {
	return v.VisitUnaryExpr(u)
}