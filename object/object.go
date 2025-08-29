package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"pir-interpreter/ast"
	"strings"
)

const (
	INT_OBJ         = "INT"
	BOOL_OBJ        = "BOOL"
	MT_OBJ          = "MT"
	MAYBE_OBJ       = "MAYBE"
	GIVES_VALUE_OBJ = "GIVES_VALUE"
	ERROR_OBJ       = "ERROR"
	FUNCTION_OBJ    = "FUNCTION"
	STRING_OBJ      = "STRING"
	BUILTIN_OBJ     = "BUILTIN"
	ARRAY_OBJ       = "ARRAY"
	HASHMAP_OBJ     = "HASHMAP"
	BREAK_OBJ       = "BREAK"
	CHEST_TYPE_OBJ  = "CHEST_TYPE"
	CHEST_OBJ       = "CHEST"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	AsString() string
}

type Int struct {
	Value int64
}

func (i *Int) Type() ObjectType { return INT_OBJ }
func (i *Int) AsString() string { return fmt.Sprintf("%d", i.Value) }

type Bool struct {
	Value bool
}

func (b *Bool) Type() ObjectType { return BOOL_OBJ }
func (b *Bool) AsString() string {
	if b.Value {
		return "ay"
	}
	return "nay"
}

type MT struct{}

func (mt *MT) Type() ObjectType { return MT_OBJ }
func (mt *MT) AsString() string { return "MT" }

type Maybe struct{}

func (m *Maybe) Type() ObjectType { return MAYBE_OBJ }
func (m *Maybe) AsString() string { return "Maybe" }

type GivesValue struct {
	Value Object
}

func (gv *GivesValue) Type() ObjectType { return GIVES_VALUE_OBJ }
func (gv *GivesValue) AsString() string { return gv.Value.AsString() }

type Break struct {
}

func (b *Break) Type() ObjectType { return BREAK_OBJ }
func (b *Break) AsString() string { return "break" }

type Error struct {
	Message string
	Line    int
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) AsString() string { return fmt.Sprintf("ERROR: %s", e.Message) }

func IsError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

type Function struct {
	Params []*ast.Identifier
	Body   *ast.BlockStatement
	NS     *Namespace
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) AsString() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Params {
		params = append(params, p.String())
	}
	out.WriteString("f")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") :\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n.")
	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) AsString() string { return s.Value }

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) AsString() string { return "builtin func" }

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) AsString() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.AsString())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type Hashable interface {
	Hash() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Bool) Hash() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

func (i *Int) Hash() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) Hash() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type KVP struct {
	Key   Object
	Value Object
}

type HashMap struct {
	MP map[HashKey]KVP
}

func (h *HashMap) Type() ObjectType { return HASHMAP_OBJ }

func (h *HashMap) AsString() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.MP {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.AsString(), pair.Value.AsString()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type Chest struct {
	Items map[string]Object
}

func (t *Chest) Type() ObjectType { return CHEST_OBJ }

func (t *Chest) AsString() string {
	var out bytes.Buffer
	pairs := []string{}
	for id, obj := range t.Items {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			id, obj.AsString()))
	}
	out.WriteString("|")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("|")
	return out.String()
}

type ChestType struct {
	Fields []string
}

func (ct *ChestType) Type() ObjectType { return CHEST_TYPE_OBJ }

func (ct *ChestType) AsString() string {
	var out bytes.Buffer
	out.WriteString("chest|")
	out.WriteString(strings.Join(ct.Fields, ", "))
	out.WriteString("|")
	return out.String()
}
