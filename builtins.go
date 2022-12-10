package main

import (
	"fmt"
	"os"
)

var globals = map[string]Value{
	"+": {t: TypeBuiltin, s: "+", bu: Value.Add},
	"-": {t: TypeBuiltin, s: "-", bu: Value.Sub},
	"*": {t: TypeBuiltin, s: "*", bu: Value.Mul},
	"/": {t: TypeBuiltin, s: "/", bu: Value.Div},
	"and": {t: TypeBuiltin, s: "and", bu: Value.And},
	"or": {t: TypeBuiltin, s: "or", bu: Value.Or},
	">": {t: TypeBuiltin, s: ">", bu: Value.Gt},
	"<": {t: TypeBuiltin, s: "<", bu: Value.Lt},
	">=": {t: TypeBuiltin, s: ">=", bu: Value.Ge},
	"<=": {t: TypeBuiltin, s: "<=", bu: Value.Le},
	"==": {t: TypeBuiltin, s: "==", bu: Value.Eq},
	"!=": {t: TypeBuiltin, s: "!=", bu: Value.Ne},
	"print": {t: TypeBuiltin, s: "print", bu: Value.bPrint},
	"exit": {t: TypeBuiltin, s: "exit", bu: Value.Exit},
	"dumps!": {t: TypeBuiltin, s: "dumps!", bu: Value.Dumps},
}

func (a Value) Eq2(b Value) bool {
	if b.t != a.t {
		return false
	}
	switch b.t {
	case TypeInt:
		return b.i == a.i
	case TypeFloat:
		return b.f == a.f
	case TypeBool:
		return b.b == a.b
	case TypeStr, TypeSym:
		return b.s == a.s
	case TypeFn:
		return b.fn.first == a.fn.first
	case TypeList:
		if len(a.l) != len(b.l) {
			return false
		}
		for i := range a.l {
			if !a.l[i].Eq2(b.l[i]) {
				return false
			}
		}
		return true
	default:
		throw("Line %d: ==, unhandled type (%d)",
			b.line, b.t)
		return false
	}
}

func (callee Value) Eq(args List) *Value {
	if len(args) < 2 {
		throw("Line %d: Too few args to '=='", callee.line)
		return nil
	}
	res := Value{t: TypeBool, b: true, line: callee.line}
	for _, val := range args[1:] {
		if !val.Eq2(args[0]) {
			res.b = false
			break
		}
	}
	return &res
}

func (callee Value) Ne(args List) *Value {
	res := callee.Eq(args)
	res.b = !res.b
	return res
}

func (callee Value) And(args List) *Value {
	if len(args) < 2 {
		throw("Line %d: Too few args to 'and'", callee.line)
		return nil
	}
	res := Value{t: TypeBool, b: true, line: callee.line}
	for _, val := range args {
		if val.t != TypeBool {
			throw("Line %d: 'and' on non-bool", val.line)
			return nil
		}
		res.b = res.b && val.b
	}
	return &res
}

func (callee Value) Or(args List) *Value {
	if len(args) < 2 {
		throw("Line %d: Too few args to 'or'", callee.line)
		return nil
	}
	res := Value{t: TypeBool, b: false, line: callee.line}
	for _, val := range args {
		if val.t != TypeBool {
			throw("Line %d: 'or' on non-bool", val.line)
			return nil
		}
		res.b = res.b || val.b
	}
	return &res
}

func (callee Value) Lt(args List) *Value { 
	if len(args) != 2 {
		throw("Line %d: Wrong number of args to '<'", callee.line)
		return nil
	}
	res := Value{t: TypeBool, line: callee.line}
	switch {
	case args[0].t == TypeInt && args[1].t == TypeInt:
		res.b = args[0].i < args[1].i
	case args[0].t == TypeInt && args[1].t == TypeFloat:
		res.b = float64(args[0].i) < args[1].f
	case args[0].t == TypeFloat && args[1].t == TypeInt:
		res.b = args[0].f < float64(args[1].i)
	case args[0].t == TypeFloat && args[1].t == TypeFloat:
		res.b = args[0].f < args[1].f
	}
	return &res
}

func (callee Value) Le(args List) *Value {
	if len(args) != 2 {
		throw("Line %d: Wrong number of args to '<='", callee.line)
		return nil
	}
	res := Value{t: TypeBool, b: false, line: callee.line}
	switch {
	case args[0].t == TypeInt && args[1].t == TypeInt:
		res.b = args[0].i <= args[1].i
	case args[0].t == TypeInt && args[1].t == TypeFloat:
		res.b = float64(args[0].i) <= args[1].f
	case args[0].t == TypeFloat && args[1].t == TypeInt:
		res.b = args[0].f <= float64(args[1].i)
	case args[0].t == TypeFloat && args[1].t == TypeFloat:
		res.b = args[0].f <= args[1].f
	}
	return &res
}

func (callee Value) Gt(args List) *Value {
	if len(args) != 2 {
		throw("Line %d: Wrong number of args to '>'", callee.line)
		return nil
	}
	res := Value{t: TypeBool, b: false, line: callee.line}
	switch {
	case args[0].t == TypeInt && args[1].t == TypeInt:
		res.b = args[0].i > args[1].i
	case args[0].t == TypeInt && args[1].t == TypeFloat:
		res.b = float64(args[0].i) > args[1].f
	case args[0].t == TypeFloat && args[1].t == TypeInt:
		res.b = args[0].f > float64(args[1].i)
	case args[0].t == TypeFloat && args[1].t == TypeFloat:
		res.b = args[0].f > args[1].f
	}
	return &res
}

func (callee Value) Ge(args List) *Value {
	if len(args) != 2 {
		throw("Line %d: Wrong number of args to '>='", callee.line)
		return nil
	}
	res := Value{t: TypeBool, b: false, line: callee.line}
	switch {
	case args[0].t == TypeInt && args[1].t == TypeInt:
		res.b = args[0].i >= args[1].i
	case args[0].t == TypeInt && args[1].t == TypeFloat:
		res.b = float64(args[0].i) >= args[1].f
	case args[0].t == TypeFloat && args[1].t == TypeInt:
		res.b = args[0].f >= float64(args[1].i)
	case args[0].t == TypeFloat && args[1].t == TypeFloat:
		res.b = args[0].f >= args[1].f
	}
	return &res
}

func (callee Value) Add(args List) *Value {
	if len(args) < 2 {
		throw("Line %d: Too few args to '+'", callee.line)
		return nil
	}
	isfloat := false
	for _, arg := range args {
		if arg.t == TypeFloat {
			isfloat = true
		}
	}
	res := Value{line: callee.line}
	if isfloat {
		if args[0].t == TypeInt {
			res.f = float64(args[0].i)
		}
		res.t = TypeFloat
		for _, arg := range args {
			if arg.t == TypeFloat {
				res.f += arg.f
			} else {
				res.f += float64(arg.i)
			}
		}
	} else {
		res.t = TypeInt
		for _, arg := range args {
			res.i += arg.i
		}
	}
	return &res
}

func (callee Value) Sub(args List) *Value {
	if len(args) < 2 {
		throw("Line %d: Too few args to '-'", callee.line)
		return nil
	}
	isfloat := false
	for _, arg := range args {
		if arg.t == TypeFloat {
			isfloat = true
		}
	}
	res := args[0]
	if isfloat {
		if args[0].t == TypeInt {
			res.f = float64(args[0].i)
		}
		for _, arg := range args[1:] {
			if arg.t == TypeFloat {
				res.f -= arg.f
			} else {
				res.f -= float64(arg.i)
			}
		}
	} else {
		for _, arg := range args[1:] {
			res.i -= arg.i
		}
	}
	return &res
}

func (callee Value) Mul(args List) *Value {
	if len(args) < 2 {
		throw("Line %d: Too few args to '*'", callee.line)
		return nil
	}
	isfloat := false
	for _, arg := range args {
		if arg.t == TypeFloat {
			isfloat = true
		}
	}
	res := args[0]
	if isfloat {
		if args[0].t == TypeInt {
			res.f = float64(args[0].i)
		}
		for _, arg := range args[1:] {
			if arg.t == TypeFloat {
				res.f *= arg.f
			} else {
				res.f *= float64(arg.i)
			}
		}
	} else {
		for _, arg := range args[1:] {
			res.i *= arg.i
		}
	}
	return &res
}

func (callee Value) Div(args List) *Value {
	if len(args) < 2 {
		throw("Line %d: Too few args to '/'", callee.line)
		return nil
	}
	res := args[0]
	if args[0].t == TypeInt {
		res.f = float64(args[0].i)
	}
	res.t = TypeFloat
	for _, arg := range args[1:] {
		if arg.t == TypeFloat {
			res.f /= arg.f
		} else {
			res.f /= float64(arg.i)
		}
	}
	return &res
}

func (callee Value) bPrint(args List) *Value {
	if len(args) != 1 {
		throw("'print' expects 1 args, not %d", len(args))
		return nil
	}

	switch args[0].t {
	case TypeError:
		fmt.Println("[error]")
	case TypeInt:
		fmt.Printf("%d\n", args[0].i)
	case TypeFloat:
		fmt.Printf("%f\n", args[0].f)
	case TypeBool:
		fmt.Printf("%t\n", args[0].b)
	case TypeStr:
		fmt.Printf("%s\n", args[0].s)
	case TypeSym:
		fmt.Printf("[sym %s]\n", args[0].s)
	case TypeFn:
		fmt.Println("[fn]")
	case TypeList:
		if len(args[0].l) == 2 &&
			args[0].l[0].t == TypeSym &&
			args[0].l[0].s == "quote" {
			fmt.Print("'")
			args[0].l[1].Print()
			fmt.Println()
			break
		}
		fmt.Print("(")
		for i := range args[0].l {
			if i != 0 {
				fmt.Print(" ")
			}
			args[0].l[i].Print()
		}
		fmt.Print(")")
		fmt.Println()
	case TypeBuiltin:
		fmt.Printf("[builtin (%s)]\n", args[0].s)
	default:
		fmt.Printf("[unknown (%d)]\n", args[0].t)
	}
	return &args[0]
}

func (callee Value) Exit(args List) *Value {
	if len(args) > 2 {
		throw("'exit' takes at most 1 args, not %d", len(args))
		return nil
	}
	if len(args) == 0 {
		os.Exit(0)
	}
	if args[0].t != TypeInt {
		throw("Line %d: Trying to exit with non-int code",
			args[0].line)
		return nil
	}
	os.Exit(int(args[1].i))
	return nil
}

func (callee Value) Dumps(args List) *Value {
	if len(args) == 0 {
		for i := 1; i < len(stack); i++ {
			fmt.Printf("[%d] ", i)
			stack[i].Print()
			fmt.Println()
		}
	} else if len(args) == 1 {
		n := len(stack) - int(args[0].i)
		for i := len(stack) - 1; i >= n; i++ {
			stack[i].Print()
		}
	} else {
		throw("'dumps!' takes at most 1 args, not %d", len(args))
	}
	return nil
}
