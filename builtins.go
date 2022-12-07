package main

import (
	"fmt"
	"os"
)

type L1Fn func (v Value, l List) Value

func (v Value) bPrint(l List) Value {
	if len(l) != 1 {
		fmt.Printf("Line %d: Wrong format for 'puts'\n", v.line)
	}

	r := l[0].Eval()
	switch r.t {
	case TypeError: fmt.Println("[error]")
	case TypeInt: fmt.Printf("%d\n", r.i)
	case TypeFloat: fmt.Printf("%f\n", r.f)
	case TypeStr: fmt.Printf("%s\n", r.s)
	case TypeSym: fmt.Printf("'%s\n", r.s)
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
		fmt.Print("'(")
		for i := range v.l {
			if i != 0 {
				fmt.Print(" ")
			}
			v.l[i].Print()
		}
	}
	return r
}


func (v Value) Set(l List) Value {
	if len(l) != 2 || l[0].t != TypeSym {
		fmt.Printf("Line %d: Wrong format for '::'\n", v.line)
		return Value{}
	}

	res := l[1].Eval()
	stack[len(stack)-1][l[0].s] = res
	return res
}

func (v Value) Fn(l List) Value {
	res := Value{t: TypeFn}

	if len(l) < 2 || l[0].t != TypeList {
		fmt.Printf("Line %d: Wrong format for 'fn'\n", v.line)
		return Value{}
	}

	res.fn.expr = l[1:]
	res.fn.args = []string{}
	for _, arg := range l[0].l {
		if arg.t != TypeSym {
			fmt.Printf("Line %d: Wrong format for 'fn'\n", v.line)
			return Value{}
		}
		res.fn.args = append(res.fn.args, arg.s)
	}
	return res
}

func (v Value) Exit(l List) Value {
	if len(l) == 1 {
		r := l[0].Eval()
		if r.t != TypeInt {
			fmt.Printf("Line %d: Type mismatch\n", v.line)
		} else {
			os.Exit(int(r.i))
		}
	} else if len(l) > 1 {
		fmt.Printf("Line %d: Too many args to 'exit'\n", v.line)
	} else {
		os.Exit(0)
	}
	return Value{}
}


func (v Value) Add(l List) Value {
	var res Value
	res.line = v.line
	if len(l) == 0 {
		fmt.Printf("Line %d: Not enough args to '+'\n", v.line)
		return Value{}
	}

	first := l[0].Eval()
	if first.t == TypeError {
		return first
	} else if first.t != TypeInt && first.t != TypeFloat {
		fmt.Printf("Line %d: '+' on non-number\n", v.line)
		return Value{}
	}

	if first.t == TypeInt {
		res.t = TypeInt
		res.i = first.i
	} else {
		res.t = TypeFloat
		res.f = first.f
	}
	for _, x := range l[1:] {
		x = x.Eval()
		switch {
		case first.t == TypeInt && x.t == TypeInt:
			res.i += x.i
		case first.t == TypeInt && x.t == TypeFloat:
			res.i += int64(x.f)
		case first.t == TypeFloat && x.t == TypeFloat:
			res.f += x.f
		case first.t == TypeFloat && x.t == TypeInt:
			res.f += float64(x.i)
		default:
			fmt.Printf("Line %d: Type mismatch\n", x.line)
			return Value{}
		}
	}
	return res
}


func (v Value) Sub(l List) Value {
	var res Value
	res.line = v.line
	if len(l) == 0 {
		fmt.Printf("Line %d: Not enough args to '-'\n", v.line)
		return Value{}
	}

	first := l[0].Eval()
	if first.t == TypeError {
		return first
	} else if first.t != TypeInt && first.t != TypeFloat {
		fmt.Printf("Line %d: '-' on non-number\n", v.line)
		return Value{}
	}

	if first.t == TypeInt {
		res.t = TypeInt
		res.i = first.i
	} else {
		res.t = TypeFloat
		res.f = first.f
	}
	for _, x := range l[1:] {
		x = x.Eval()
		switch {
		case first.t == TypeInt && x.t == TypeInt:
			res.i -= x.i
		case first.t == TypeInt && x.t == TypeFloat:
			res.i -= int64(x.f)
		case first.t == TypeFloat && x.t == TypeFloat:
			res.f -= x.f
		case first.t == TypeFloat && x.t == TypeInt:
			res.f -= float64(x.i)
		default:
			fmt.Printf("Line %d: Type mismatch\n", x.line)
			return Value{}
		}
	}
	return res
}
func (v Value) Mul(l List) Value {
	var res Value
	res.line = v.line
	if len(l) == 0 {
		fmt.Printf("Line %d: Not enough args to '*'\n", v.line)
		return Value{}
	}

	first := l[0].Eval()
	if first.t == TypeError {
		return first
	} else if first.t != TypeInt && first.t != TypeFloat {
		fmt.Printf("Line %d: '*' on non-number\n", v.line)
		return Value{}
	}

	if first.t == TypeInt {
		res.t = TypeInt
		res.i = first.i
	} else {
		res.t = TypeFloat
		res.f = first.f
	}
	for _, x := range l[1:] {
		x = x.Eval()
		switch {
		case first.t == TypeInt && x.t == TypeInt:
			res.i *= x.i
		case first.t == TypeInt && x.t == TypeFloat:
			res.i *= int64(x.f)
		case first.t == TypeFloat && x.t == TypeFloat:
			res.f *= x.f
		case first.t == TypeFloat && x.t == TypeInt:
			res.f *= float64(x.i)
		default:
			fmt.Printf("Line %d: Type mismatch\n", x.line)
			return Value{}
		}
	}
	return res
}

func (v Value) Div(l List) Value {
	var res Value
	res.line = v.line
	if len(l) == 0 {
		fmt.Printf("Line %d: Not enough args to '/'\n", v.line)
		return Value{}
	}

	first := l[0].Eval()
	if first.t == TypeError {
		return first
	} else if first.t != TypeInt && first.t != TypeFloat {
		fmt.Printf("Line %d: '/' on non-number\n", v.line)
		return Value{}
	}

	if first.t == TypeInt {
		res.t = TypeInt
		res.i = first.i
	} else {
		res.t = TypeFloat
		res.f = first.f
	}
	for _, x := range l[1:] {
		x = x.Eval()
		switch {
		case first.t == TypeInt && x.t == TypeInt:
			res.i /= x.i
		case first.t == TypeInt && x.t == TypeFloat:
			res.i /= int64(x.f)
		case first.t == TypeFloat && x.t == TypeFloat:
			res.f /= x.f
		case first.t == TypeFloat && x.t == TypeInt:
			res.f /= float64(x.i)
		default:
			fmt.Printf("Line %d: Type mismatch\n", x.line)
			return Value{}
		}
	}
	return res
}
