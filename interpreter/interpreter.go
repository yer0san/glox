package interpreter

import (
	"fmt"

	. "github.com/yer0san/glox/errors"
	. "github.com/yer0san/glox/expr"
	"github.com/yer0san/glox/stmt"
	. "github.com/yer0san/glox/token"
	. "github.com/yer0san/glox/environment"
)

type Interpreter struct {
	Environment *Environment
}

func (i *Interpreter) VisitLiteralExpr(expr *Literal) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitGroupingExpr(expr *Grouping) (any, error) {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitUnaryExpr(expr *Unary) (any, error) {
	right, err := i.evaluate(expr.Right)

	if err != nil {
		ReportError(expr.Operator, err)
		return nil, err
	}

	switch expr.Operator.Token_type {
		case MINUS:
			err := i.checkNumberOperand(right)

			if err != nil {
				ReportError(expr.Operator, err)
				return nil, err
			}

			return -(right.(float64)), nil
		case BANG:
			return !(i.isTruthy(right)), nil
	}
	return nil, nil // unreachable he said
}

func (i *Interpreter) VisitBinaryExpr(expr *Binary) (any, error) {
	left, err := i.evaluate(expr.Left)

	if err != nil {
		ReportError(expr.Operator, err)
		return nil, err
	}

	right, err := i.evaluate(expr.Right)

	if err != nil {
		ReportError(expr.Operator, err)
		return nil, err
	}

	switch expr.Operator.Token_type {
		case MINUS:
			err := i.checkNumberOperands(left, right)

			if err != nil {
				ReportError(expr.Operator, err)
				return nil, err
			}
			return left.(float64) - right.(float64), nil

		case PLUS:
			if leftNum, ok := left.(float64); ok {
				if rightNum, ok := right.(float64); ok {
					return leftNum + rightNum, nil
				}
			}

			if leftStr, ok := left.(string); ok {
				if rightStr, ok := right.(string); ok {
					return leftStr + rightStr, nil
				}
			}
			ReportError(expr.Operator, ErrOperandsNotSameType)
			return nil, ErrOperandsNotSameType
		case SLASH:
			err := i.checkNumberOperands(left, right)

			if err != nil {
				ReportError(expr.Operator, err)
				return nil, err
			}

			return left.(float64) / right.(float64), nil
		case STAR:
			err := i.checkNumberOperands(left, right)

			if err != nil {
				return nil, err
			}

			return left.(float64) * right.(float64), nil
		
		case GREATER:
			err := i.checkNumberOperands(left, right)

			if err != nil {
				ReportError(expr.Operator, err)
				return nil, err
			}

			return left.(float64) > right.(float64), nil
		case GREATER_EQUAL:
			err := i.checkNumberOperands(left, right)

			if err != nil {
				ReportError(expr.Operator, err)
				return nil, err
			}

			return left.(float64) >= right.(float64), nil
		case LESS:
			err := i.checkNumberOperands(left, right)

			if err != nil {
				ReportError(expr.Operator, err)
				return nil, err
			}

			return left.(float64) < right.(float64), nil
		case LESS_EQUAL:
			err := i.checkNumberOperands(left, right)

			if err != nil {
				ReportError(expr.Operator, err)
				return nil, err
			}

			return left.(float64) <= right.(float64), nil

		case BANG_EQUAL:
			return left != right, nil
		case EQUAL_EQUAL:
			return left == right, nil
	}
	return nil, nil
}

func (i *Interpreter) VisitVariableExpr(expr *Variable) (any, error) {
	return i.Environment.Get(expr.Name)
}

func (i *Interpreter) VisitAssignExpr(expr *Assign) (any, error) {
	value, err := i.evaluate(expr.Value)

	if err != nil {
		//??
		return nil, err
	}
	i.Environment.Assign(expr.Name, value)
	return value, nil
}


func (i *Interpreter) VisitExprStmt(statement *stmt.Expression) {
	i.evaluate(statement.Expr)
}

func (i *Interpreter) VisitPrintStmt(statement *stmt.Print) {
	value, err := i.evaluate(statement.Expr)

	if err != nil {
		//???
		return
	}
	switch v := value.(type) {
		default:
			fmt.Println(v)
	}
}

func (i *Interpreter) VisitVarStmt(statement *stmt.Var) {
	var value any
	if statement.Initializer != nil {
		val, err := i.evaluate(statement.Initializer)

		if err != nil {
			return 
		}
		value = val // check if this is correct
	}
	i.Environment.Define(statement.Name.Lexeme, value)
}

func (i *Interpreter) VisitBlockStmt(stmt *stmt.Block) {
	i.executeBlock(stmt.Statements, NewEnclosingEnvironment(i.Environment))
}

// HELPERS
func (i *Interpreter) executeBlock(statements []stmt.Stmt, environment *Environment) {
	prev := i.Environment

	i.Environment = environment
	for _, statement := range statements {
		i.execute(statement)
	}
	i.Environment = prev
}

func (i *Interpreter) evaluate(expr Expr) (any, error) {
	val, err := expr.Accept(i)
	return val, err
}

func (i *Interpreter) isTruthy(obj any) bool {
	switch v := obj.(type) {
		case nil:
			return false
		case bool:
			return v
		default:
			return true
	}
}

func (i *Interpreter) checkNumberOperand(operand any) error {
	if _, ok := operand.(float64); ok {
		return nil
	}
	return ErrOperandNotNumber
}

func (i *Interpreter) checkNumberOperands(left any, right any) error {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return nil
		}
	}
	return ErrOperandNotNumber
}

func (i *Interpreter) execute(statement stmt.Stmt) {
	statement.Accept(i)
}

// entry
func (i *Interpreter) Interpret(statements []stmt.Stmt) {
	for _, statement := range statements {
		i.execute(statement)
	}
} // shit is rly complicated
