package stmt

import (
	"github.com/yer0san/glox/expr"
	. "github.com/yer0san/glox/token"
) 

type Stmt interface {
	Accept(v Visitor)
}

type Expression struct {
	Expr expr.Expr
}

type Print struct {
	Expr expr.Expr
}

type Var struct {
	Name *Token
	Initializer expr.Expr
}

type Block struct {
	Statements []Stmt
}

type Visitor interface {
	VisitExprStmt(stmt *Expression)
	VisitPrintStmt(stmt *Print)
	VisitVarStmt(stmt *Var)
	VisitBlockStmt(stmt *Block)
}

func (e *Expression) Accept(v Visitor) {
	v.VisitExprStmt(e)
}

func (p *Print) Accept(v Visitor) {
	v.VisitPrintStmt(p)
}

func (p *Var) Accept(v Visitor) {
	v.VisitVarStmt(p)
}

func (b *Block) Accept(v Visitor) {
	v.VisitBlockStmt(b)
}