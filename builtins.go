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
	"btrace!": {t: TypeBuiltin, s: "btrace!", bu: Value.Btrace},
	"strace!": {t: TypeBuiltin, s: "strace!", bu: Value.Strace},
	"itrace!": {t: TypeBuiltin, s: "itrace!", bu: Value.Itrace},
}

func (a Value) Eq2(b Value) bool {
	if b.t != a.t {
		return false
	}
	switch b.t {
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
		throw("Line %d: '==' - unhandled type (%d)",
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
	if args[0].t != TypeFloat || args[1].t != TypeFloat {
		throw("Line %d: '<' only takes numeric args", args[1].line)
		return nil
	}
	res.b = args[0].f < args[1].f
	return &res
}

func (callee Value) Le(args List) *Value {
	if len(args) != 2 {
		throw("Line %d: Wrong number of args to '<='", callee.line)
		return nil
	}
	if args[0].t != TypeFloat || args[1].t != TypeFloat {
		throw("Line %d: '<=' only takes numeric args", args[1].line)
		return nil
	}
	res := Value{t: TypeBool, b: false, line: callee.line}
	res.b = args[0].f <= args[1].f
	return &res
}

func (callee Value) Gt(args List) *Value {
	if len(args) != 2 {
		throw("Line %d: Wrong number of args to '>'", callee.line)
		return nil
	}
	if args[0].t != TypeFloat || args[1].t != TypeFloat {
		throw("Line %d: '>' only takes numeric args", args[1].line)
		return nil
	}
	res := Value{t: TypeBool, b: false, line: callee.line}
	res.b = args[0].f > args[1].f
	return &res
}

func (callee Value) Ge(args List) *Value {
	if len(args) != 2 {
		throw("Line %d: Wrong number of args to '>='", callee.line)
		return nil
	}
	if args[0].t != TypeFloat || args[1].t != TypeFloat {
		throw("Line %d: '>=' only takes numeric args", args[1].line)
		return nil
	}
	res := Value{t: TypeBool, b: false, line: callee.line}
	res.b = args[0].f >= args[1].f
	return &res
}

func (callee Value) Add(args List) *Value {
	if len(args) < 2 {
		throw("Line %d: Too few args to '+'", callee.line)
		return nil
	}
	if args[0].t != TypeFloat {
		throw("Line %d: '+' only takes numeric args", args[0].line)
		return nil
	}
	res := args[0]
	for _, arg := range args[1:] {
		if arg.t != TypeFloat {
			throw("Line %d: '+' only takes numeric args", arg.line)
			return nil
		}
		res.f += arg.f
	}
	return &res
}

func (callee Value) Sub(args List) *Value {
	if len(args) < 2 {
		throw("Line %d: Too few args to '-'", callee.line)
		return nil
	}
	if args[0].t != TypeFloat {
		throw("Line %d: '-' only takes numeric args", args[0].line)
		return nil
	}
	res := args[0]
	for _, arg := range args[1:] {
		if arg.t != TypeFloat {
			throw("Line %d: '-' only takes numeric args", arg.line)
			return nil
		}
		res.f -= arg.f
	}
	return &res
}

func (callee Value) Mul(args List) *Value {
	if len(args) < 2 {
		throw("Line %d: Too few args to '*'", callee.line)
		return nil
	}
	res := args[0]
	if args[0].t != TypeFloat {
		throw("Line %d: '*' only takes numeric args", args[0].line)
		return nil
	}
	for _, arg := range args[1:] {
		if arg.t != TypeFloat {
			throw("Line %d: '*' only takes numeric args", arg.line)
			return nil
		}
		res.f *= arg.f
	}
	return &res
}

func (callee Value) Div(args List) *Value {
	if len(args) < 2 {
		throw("Line %d: Too few args to '/'", callee.line)
		return nil
	}
	res := args[0]
	if args[0].t != TypeFloat {
		throw("Line %d: '/' only takes numeric args", args[0].line)
		return nil
	}
	for _, arg := range args {
		if arg.t != TypeFloat {
			throw("Line %d: '/' only takes numeric args", arg.line)
			return nil
		}
		res.f /= arg.f
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
	if args[0].t != TypeFloat {
		throw("Line %d: Trying to exit with non-numeric code",
			args[0].line)
		return nil
	}
	os.Exit(int(args[1].f))
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
		n := len(stack) - int(args[0].f)
		for i := len(stack) - 1; i >= n; i++ {
			stack[i].Print()
		}
	} else {
		throw("'dumps!' takes at most 1 args, not %d", len(args))
	}
	return nil
}

func (callee Value) Btrace(args List) *Value {
	if len(args) != 0 {
		throw("'btrace!' takes no args")
	} else {
		btrace = !btrace
	}
	return nil
}

func (callee Value) Itrace(args List) *Value {
	if len(args) != 0 {
		throw("'itrace!' takes no args")
	} else {
		itrace = !itrace
	}
	return nil
}

func (callee Value) Strace(args List) *Value {
	if len(args) != 0 {
		throw("'strace!' takes no args")
	} else {
		strace = !strace
	}
	return nil
}
