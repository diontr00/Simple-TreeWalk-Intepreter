package object

import (
	"fmt"
)

type ObjectType string

// Every value is wrapped inside a struct , which fulfill the Object interface
type Object interface {
	Type() ObjectType
	Inspect() string
}

// After parsing associated ast literal , we turn them into internal object ,using those struct
// INTEGER
// -------------------------------------------------------------------
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// BOOLEAN
// --------------------------------------------------------------------
type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

// NIL
type Nil struct{}

func (n *Nil) Inspect() string  { return "nil" }
func (n *Nil) Type() ObjectType { return NIL_OBJ }

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NIL_OBJ     = "NIL"
	STRING_OBJ  = "STRING"
)

// STRING
// --------------------------------------------------------------------
type String struct {
	Value string
}

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRING_OBJ }
