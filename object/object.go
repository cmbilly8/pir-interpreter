package object

import (
	"fmt"
)

const (
	INT_OBJ         = "INT"
	BOOL_OBJ        = "BOOL"
	MT_OBJ          = "MT"
	MAYBE_OBJ       = "MAYBE"
	GIVES_VALUE_OBJ = "GIVES_VALUE"
	ERROR_OBJ       = "ERROR"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Int struct {
	Value int64
}

func (i *Int) Type() ObjectType { return INT_OBJ }
func (i *Int) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Bool struct {
	Value bool
}

func (b *Bool) Type() ObjectType { return BOOL_OBJ }
func (b *Bool) Inspect() string {
	if b.Value {
		return "ay"
	}
	return "nay"
}

type MT struct{}

func (mt *MT) Type() ObjectType { return MT_OBJ }
func (mt *MT) Inspect() string  { return "MT" }

type Maybe struct{}

func (m *Maybe) Type() ObjectType { return MAYBE_OBJ }
func (m *Maybe) Inspect() string  { return "Maybe" }

type GivesValue struct {
	Value Object
}

func (gv *GivesValue) Type() ObjectType { return GIVES_VALUE_OBJ }
func (gv *GivesValue) Inspect() string  { return gv.Value.Inspect() }

type Error struct {
	Message string
	Line    int
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return fmt.Sprintf("ERROR: %s", e.Message) }

func IsError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}
