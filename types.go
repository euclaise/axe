package main

import "fmt"

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

	line int
}

type List []Value

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
		if len(v.l) == 2 {
			if v.l[0].t == TypeSym && v.l[0].s == "quote" {
				fmt.Print("'")
				v.l[1].Print()
				return
			}
		}
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
