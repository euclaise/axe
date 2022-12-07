package main

import "fmt"

const (
	TypeError = iota
	TypeInt
	TypeFloat
	TypeBool
	TypeStr
	TypeSym
	TypeFn
	TypeList
)

type Value struct {
	t int //type
	
	i int64
	f float64
	b bool
	s string //string, sym
	fn Fn
	l List

	line int
}

type List []Value

type Fn struct {
	args []string
	expr List
}

func (v Value) Print() {
	switch v.t {
	case TypeError: fmt.Print("[error]")
	case TypeInt: fmt.Printf("%d", v.i)
	case TypeFloat: fmt.Printf("%f", v.f)
	case TypeBool: fmt.Printf("%t", v.b)
	case TypeStr: fmt.Printf("\"%s\"", v.s)
	case TypeSym: fmt.Printf("%s", v.s)
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
	default:
		fmt.Println("[nil]")
	}
}
