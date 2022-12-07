package main

import (
	"fmt"
)

var stack = []map[string]Value{
	{}, //Globals
}

func (v Value) Eval() Value {
	var res Value
	switch v.t {
	case TypeError: return v
	case TypeInt: return v
	case TypeFloat: return v
	case TypeBool: return v
	case TypeStr: return v
	case TypeSym:
		if x, ok := stack[len(stack) - 1][v.s]; ok {
			return x
		} else if x, ok := stack[0][v.s]; ok { //globals
			return x
		}
		fmt.Printf("Line %d: Could not find variable %s\n", v.line, v.s)
		return Value{}
	case TypeList:
		if len(v.l) == 0 {
			fmt.Printf("Line %d: Empty list\n", v.line)
			return Value{}
		}
		if v.l[0].t == TypeSym {
			if b, ok := map[string]L1Fn{
				"=" : Value.Set,
				"+" : Value.Add,
				"-" : Value.Sub,
				"*" : Value.Mul,
				"/" : Value.Div,
				"fn" : Value.Fn,
				"exit" : Value.Exit,
				"print" : Value.bPrint,
			}[v.l[0].s]; ok {
				if len(v.l) == 1 {
					return b(v, List{})
				} else {
					return b(v, v.l[1:])
				}
			}
		}
		fnv := v.l[0].Eval()
		if fnv.t != TypeFn {
			fmt.Printf("Line %d: Call to non-fn\n", v.line);
			return Value{}
		}
		fn := fnv.fn
		if len(v.l) != len(fn.args) + 1 {
			fmt.Printf("Line %d: Wrong number of args\n", v.line);
			return Value{}
		}
		stack = append(stack, map[string]Value{})
		for i, arg := range fn.args {
			stack[len(stack)-1][arg] = v.l[i+1].Eval()
		}
		for _, expr := range fn.expr {
			res = expr.Eval()
		}
		stack = stack[:len(stack)-1]
	}
	return res
}
