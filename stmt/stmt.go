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

type If struct {
	Condition expr.Expr // Logical??
	ThenBranch Stmt
	ElseBranch Stmt
}

type While struct {
	Condition expr.Expr
	Body Stmt
}

type For struct {
	Init Stmt
	Condition expr.Expr
	Increment expr.Expr
	Body Stmt
}

type Visitor interface {
	VisitExprStmt(stmt *Expression)
	VisitPrintStmt(stmt *Print)
	VisitVarStmt(stmt *Var)
	VisitBlockStmt(stmt *Block)
	VisitIfStmt(stmt *If)
	VisitWhileStmt(stmt *While)
	VisitForStmt(stmt *For)
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

func (i *If) Accept(v Visitor) {
	v.VisitIfStmt(i)
}

func (w *While) Accept(v Visitor) {
	v.VisitWhileStmt(w)
}

func (f *For) Accept(v Visitor) {
	v.VisitForStmt(f)
}