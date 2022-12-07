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
	case TypeBool: fmt.Printf("%t\n", r.b)
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
	res := Value{line: v.line}
	
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
	res := Value{line: v.line}
	
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
	res := Value{line: v.line}
	
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
	res := Value{line: v.line}
	
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

func (v Value) Eq(l List) Value {
	res := Value{t: TypeBool, b: true, line: v.line}
	if len(l) < 2 {
		fmt.Printf("Line %d: Too few args to '=='\n", v.line)
		return Value{}
	}

	first := l[0].Eval()
	for _, val := range l[1:] {
		val = val.Eval()
		if val.t != first.t {
			res.b = false
			return res
		}
		switch val.t {
		case TypeInt: res.b = res.b && (val.i == first.i)
		case TypeFloat: res.b = res.b && (val.f == first.f)
		case TypeBool: res.b = res.b && (val.b == first.b)
		case TypeStr: fallthrough
		case TypeSym: res.b = res.b && (val.s == first.s)
		case TypeFn:
			if len(val.fn.args) != len(first.fn.args) ||
				len(val.fn.expr) != len(first.fn.expr) {
				res.b = false
				return res
			}
			for i, arg := range val.fn.args {
				if arg != val.fn.args[i] {
					res.b = false
					return res
				}
			}
			for i, expr := range v.fn.expr {
				res = v.Eq(List{expr, val.fn.expr[i]})
			}
		case TypeList:
			if len(val.l) != len(first.l) {
				res.b = false
				return res
			}
			for i, item := range val.l {
				res = v.Eq(List{item, first.l[i]})
			}
		}

		if !res.b {
			return res
		}
	}
	return res
}

func (v Value) Ne(l List) Value {
	res := v.Eq(l)
	res.b = !res.b
	return res
}

func (v Value) Lt(l List) Value {
	res := Value{t: TypeBool, line: v.line}
	if len(l) != 2 {
		fmt.Printf("Line %d: Wrong number of args to '<'\n", v.line)
		return Value{}
	}

	first := l[0].Eval()
	second := l[1].Eval()

	switch {
		case first.t == TypeInt && second.t == TypeInt:
			res.b = first.i < second.i
		case first.t == TypeInt && second.t == TypeFloat:
			res.b = float64(first.i) < second.f
		case first.t == TypeFloat && second.t == TypeInt:
			res.b = first.f < float64(second.i)
		case first.t == TypeFloat && second.t == TypeFloat:
			res.b = first.f < second.f
		default:
			fmt.Printf("Line %d: Type mismatch\n", v.line)
			return Value{}
	}
	return res
}

func (v Value) Gt(l List) Value {
	res := Value{t: TypeBool, line: v.line}
	if len(l) != 2 {
		fmt.Printf("Line %d: Wrong number of args to '>'\n", v.line)
		return Value{}
	}

	first := l[0].Eval()
	second := l[1].Eval()

	switch {
		case first.t == TypeInt && second.t == TypeInt:
			res.b = first.i > second.i
		case first.t == TypeInt && second.t == TypeFloat:
			res.b = float64(first.i) > second.f
		case first.t == TypeFloat && second.t == TypeInt:
			res.b = first.f > float64(second.i)
		case first.t == TypeFloat && second.t == TypeFloat:
			res.b = first.f > second.f
		default:
			fmt.Printf("Line %d: Type mismatch\n", v.line)
			return Value{}
	}
	return res
}


func (v Value) Lte(l List) Value {
	res := Value{t: TypeBool, line: v.line}
	if len(l) != 2 {
		fmt.Printf("Line %d: Wrong number of args to '<='\n", v.line)
		return Value{}
	}

	first := l[0].Eval()
	second := l[1].Eval()

	switch {
		case first.t == TypeInt && second.t == TypeInt:
			res.b = first.i <= second.i
		case first.t == TypeInt && second.t == TypeFloat:
			res.b = float64(first.i) <= second.f
		case first.t == TypeFloat && second.t == TypeInt:
			res.b = first.f <= float64(second.i)
		case first.t == TypeFloat && second.t == TypeFloat:
			res.b = first.f <= second.f
		default:
			fmt.Printf("Line %d: Type mismatch\n", v.line)
			return Value{}
	}
	return res
}

func (v Value) Gte(l List) Value {
	res := Value{t: TypeBool, line: v.line}
	if len(l) != 2 {
		fmt.Printf("Line %d: Wrong number of args to '>='\n", v.line)
		return Value{}
	}

	first := l[0].Eval()
	second := l[1].Eval()

	switch {
		case first.t == TypeInt && second.t == TypeInt:
			res.b = first.i >= second.i
		case first.t == TypeInt && second.t == TypeFloat:
			res.b = float64(first.i) >= second.f
		case first.t == TypeFloat && second.t == TypeInt:
			res.b = first.f >= float64(second.i)
		case first.t == TypeFloat && second.t == TypeFloat:
			res.b = first.f >= second.f
		default:
			fmt.Printf("Line %d: Type mismatch\n", v.line)
			return Value{}
	}
	return res
}

func (v Value) And(l List) Value {
	res := Value{t: TypeBool, b: true, line: v.line}
	
	if len(l) < 2 {
		fmt.Printf("Line %d: Too few args to 'and'", v.line)
		return Value{}
	}
	
	for _, arg := range l {
		eval := arg.Eval()
		if eval.t != TypeBool {
			fmt.Printf("Line %d: 'and' on non-bool", arg.line)
			return Value{}
		}
		if eval.b == false {
			res.b = false
			return res
		}
	}
	return res
}

func (v Value) Or(l List) Value {
	res := Value{t: TypeBool, b: false, line: v.line}
	
	if len(l) < 2 {
		fmt.Printf("Line %d: Too few args to 'or'", v.line)
		return Value{}
	}
	
	for _, arg := range l {
		eval := arg.Eval()
		if eval.t != TypeBool {
			fmt.Printf("Line %d: 'or' on non-bool", arg.line)
			return Value{}
		}
		if eval.b == true {
			res.b = true
			return res
		}
	}
	return res
}

func (v Value) Not(l List) Value {
	if len(l) != 1 {
		fmt.Printf("Line %d: Wrong number of args to 'not'", v.line)
		return Value{}
	}
	res := l[0].Eval()
	if res.t != TypeBool {
		fmt.Printf("Line %d: 'not' on non-bool", l[0].line)
		return Value{}
	}
	res.b = !res.b
	return res
}
