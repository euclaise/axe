package main

import (
	"fmt"
	"io"
)

const (
	TypeDummy = iota
	TypeError
	TypeFloat
	TypeBool
	TypeStr
	TypeSym
	TypeFn
	TypeList
	TypeBlock
	TypeBuiltin
)

type Value struct {
	t int //type

	f  float64
	b  bool
	s  string //string, sym, builtin
	fn Fn
	bl *Block
	l  List
	n  int
	bu func(Value, List) *Value
	st io.ReadSeekCloser
	wr io.WriteCloser

	from *Block
	line int
	file string
}

type List []Value

type Stream io.ReadSeekCloser

var thrown = false

func (v Value) Float() float64 {
	if v.t != TypeFloat {
		throw("%s line %d: Type mismatch, Expected numeric", v.file, v.line)
	}
	return v.f
}

func (v Value) Int() int {
	if v.t != TypeFloat {
		throw("%s line %d: Type mismatch, Expected numeric", v.file, v.line)
	}
	return int(v.f)
}

func (v Value) Bool() bool {
	if v.t != TypeBool {
		throw("%s line %d: Type mismatch, Expected bool", v.file, v.line)
	}
	return v.b
}

func (v Value) String() string {
	if v.t != TypeStr {
		throw("%s line %d: Type mismatch, Expected str", v.file, v.line)
	}
	return v.s
}

func (v Value) Symbol() string {
	if v.t != TypeSym {
		throw("%s line %d: Type mismatch, Expected sym", v.file, v.line)
	}
	return v.s
}

func (v Value) Fn() Fn {
	if v.t != TypeFn {
		throw("%s line %d: Type mismatch, Expected fn", v.file, v.line)
	}
	return v.fn
}

func (v Value) List() List {
	if v.t != TypeList {
		throw("%s line %d: Type mismatch, Expected list", v.file, v.line)
	}
	return v.l
}

func (v Value) Builtin() func(Value, List) *Value {
	if v.t != TypeBuiltin {
		throw("%s line %d: Type mismatch, Expected builtin", v.file, v.line)
	}
	return v.bu
}


func (v Value) Print() {
	switch v.t {
	case TypeError:
		fmt.Print("[error]")
	case TypeFloat:
		fmt.Printf("%f", v.f)
	case TypeBool:
		fmt.Printf("%t", v.b)
	case TypeStr:
		fmt.Printf("\"%s\"", v.s)
	case TypeSym:
		fmt.Printf("%s", v.s)
	case TypeFn:
		fmt.Print("[fn (")
		for i := range v.fn.args {
			fmt.Print(v.fn.args[i])
			if i != len(v.fn.args)-1 {
				fmt.Print(", ")
			}
		}
		fmt.Print(")]")
	case TypeList:
		fmt.Print("(")
		for i := range v.l {
			if i != 0 {
				fmt.Print(" ")
			}
			v.l[i].Print()
		}
		fmt.Print(")")
	case TypeBlock:
		fmt.Print("[block]")
	case TypeBuiltin:
		fmt.Print("[builtin]")
	default:
		fmt.Println("[nil]")
	}
}
