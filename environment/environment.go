package environment

import (
	. "github.com/yer0san/glox/errors"
	. "github.com/yer0san/glox/token"
)

type Environment struct {
	Enclosing *Environment
	Values map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{Values: make(map[string]any)}
}

func NewEnclosingEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		Enclosing: enclosing,
		Values: make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) {
	e.Values[name] = value
}

func (e *Environment) Get(name *Token) (any, error) {
	value, ok := e.Values[name.Lexeme]

	if !ok {
		if e.Enclosing != nil {
			val, err := e.Enclosing.Get(name)

			if err != nil {
				return nil, err
			}
			return val, nil
		}

		err := ErrUndefinedVariable(name.Lexeme)
		ReportError(name, err)
		return nil, err
	}
	return value, nil
}

func (e *Environment) Assign(name *Token, value any) (error) {
	_, ok := e.Values[name.Lexeme]

	if ok {
		e.Values[name.Lexeme] = value
		return nil
	}
	if e.Enclosing != nil {
		err := e.Enclosing.Assign(name, value)
		if err != nil {
			return err
		}
	}else {
		err := ErrUndefinedVariable(name.Lexeme)
		ReportError(name, err)
		return err
	}
	return nil
} // this should work