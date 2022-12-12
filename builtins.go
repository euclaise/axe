package main

import (
	"fmt"
	"os"
	"strings"
)

var builtins = map[string]Value{
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
	"readfile": {t: TypeBuiltin, s: "readfile", bu: Value.ReadFile},
	"writefile": {t: TypeBuiltin, s: "writefile", bu: Value.ReadFile},
	"strsplit": {t: TypeBuiltin, s: "strsplit", bu: Value.ReadFile},
	"nth": {t: TypeBuiltin, s: "nth", bu: Value.Nth},
	"head": {t: TypeBuiltin, s: "head", bu: Value.Head},
	"tail": {t: TypeBuiltin, s: "tail", bu: Value.Tail},
	"last": {t: TypeBuiltin, s: "last", bu: Value.Last},
	"append": {t: TypeBuiltin, s: "append", bu: Value.Append},
	"flat": {t: TypeBuiltin, s: "flat", bu: Value.Flat},
	"join": {t: TypeBuiltin, s: "join", bu: Value.Join},
}

var globals = builtins

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
		throw("%s, line %d: '==' - unhandled type (%d)",
			b.file, b.line, b.t)
		return false
	}
}

func (callee Value) Eq(args List) *Value {
	if len(args) < 2 {
		throw("%s, line %d: Too few args to '=='", callee.file, callee.line)
		return &Value{t: TypeError}
	}
	res := Value{t: TypeBool, b: true, file: callee.file, line: callee.line}
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
		throw("%s, line %d: Too few args to 'and'", callee.file, callee.line)
		return &Value{t: TypeError}
	}
	res := Value{t: TypeBool, b: true, file: callee.file, line: callee.line}
	for _, val := range args {
		if val.t != TypeBool {
			throw("%s, line %d: 'and' on non-bool", callee.file, val.line)
			return &Value{t: TypeError}
		}
		res.b = res.b && val.b
	}
	return &res
}

func (callee Value) Or(args List) *Value {
	if len(args) < 2 {
		throw("%s, line %d: Too few args to 'or'", callee.file, callee.line)
		return &Value{t: TypeError}
	}
	res := Value{t: TypeBool, b: false, file: callee.file, line: callee.line}
	for _, val := range args {
		res.b = res.b || val.Bool()
	}
	return &res
}

func (callee Value) Lt(args List) *Value { 
	if len(args) != 2 {
		throw("%s, line %d: Wrong number of args to '<'",
			callee.file, callee.line)
		return &Value{t: TypeError}
	}
	res := Value{t: TypeBool, line: callee.line}
	res.b = args[0].Float() < args[1].Float()
	return &res
}

func (callee Value) Le(args List) *Value {
	if len(args) != 2 {
		throw("%s, line %d: Wrong number of args to '<='",
			callee.file, callee.line)
		return &Value{t: TypeError}
	}
	res := Value{t: TypeBool, b: false, line: callee.line}
	res.b = args[0].Float() <= args[1].Float()
	return &res
}

func (callee Value) Gt(args List) *Value {
	if len(args) != 2 {
		throw("%s, line %d: Wrong number of args to '>'",
			callee.file, callee.line)
		return &Value{t: TypeError}
	}
	res := Value{t: TypeBool, b: false, file: callee.file, line: callee.line}
	res.b = args[0].Float() > args[1].Float()
	return &res
}

func (callee Value) Ge(args List) *Value {
	if len(args) != 2 {
		throw("%s, line %d: Wrong number of args to '>='",
			callee.file, callee.line)
		return &Value{t: TypeError}
	}
	res := Value{t: TypeBool, b: false, line: callee.line}
	res.b = args[0].Float() >= args[1].Float()
	return &res
}

func (callee Value) Add(args List) *Value {
	if len(args) < 2 {
		throw("%s, line %d: Too few args to '+'", callee.file, callee.line)
		return &Value{t: TypeError}
	}
	args[0].Float()
	res := args[0]
	for _, arg := range args[1:] {
		res.f += arg.Float()
	}
	return &res
}

func (callee Value) Sub(args List) *Value {
	if len(args) < 2 {
		throw("%s, line %d: Too few args to '-'", callee.file, callee.line)
		return &Value{t: TypeError}
	}
	args[0].Float()
	res := args[0]
	for _, arg := range args[1:] {
		res.f -= arg.Float()
	}
	return &res
}

func (callee Value) Mul(args List) *Value {
	if len(args) < 2 {
		throw("%s, line %d: Too few args to '*'", callee.file, callee.line)
		return &Value{t: TypeError}
	}
	args[0].Float()
	res := args[0]
	for _, arg := range args[1:] {
		res.f *= arg.Float()
	}
	return &res
}

func (callee Value) Div(args List) *Value {
	if len(args) < 2 {
		throw("%s, line %d: Too few args to '/'", callee.file, callee.line)
		return &Value{t: TypeError}
	}
	args[0].Float()
	res := args[0]
		for _, arg := range args {
		res.f /= arg.f
	}
	return &res
}

func (callee Value) bPrint(args List) *Value {
	if len(args) != 1 {
		throw("'print' expects 1 args, not %d", len(args))
		return &Value{t: TypeError}
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
		fmt.Printf("'%s\n", args[0].s)
	case TypeFn:
		fmt.Println("[fn]")
	case TypeList:
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
	if len(args) == 0 {
		os.Exit(0)
	}
	if len(args) > 2 {
		throw("%s, line %d: 'exit' takes at most 1 args, not %d",
			callee.file, callee.line, len(args))
		return &Value{t: TypeError}
	}
	os.Exit(args[0].Int())
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
		throw("%s, line %d: 'dumps!' takes at most 1 args, not %d",
			callee.file, callee.line, len(args))
		return &Value{t: TypeError}
	}
	return nil
}

func (callee Value) Btrace(args List) *Value {
	if len(args) != 0 {
		throw("%s, line %d: 'btrace!' takes no args", callee.file, callee.line)
		return &Value{t: TypeError}
	}
	btrace = !btrace
	return nil
}

func (callee Value) Itrace(args List) *Value {
	if len(args) != 0 {
		throw("%s, line %d: 'itrace!' takes no args", callee.file, callee.line)
		return &Value{t: TypeError}
	}
	itrace = !itrace
	return nil
}

func (callee Value) Strace(args List) *Value {
	if len(args) != 0 {
		throw("%s, line %d: 'strace!' takes no args", callee.file, callee.line)
		return &Value{t: TypeError}
	} 
	strace = !strace
	return nil
}

func (callee Value) ReadFile(args List) *Value {
	if len(args) != 1 {
		throw("%s, line %d: 'readfile' takes 1 arg (filename), got %d",
			callee.file, callee.line, len(args))
		return &Value{t: TypeError}
	}

	f, err := os.ReadFile(args[0].String())
	if err != nil {
		fmt.Println(err)
		return &Value{t: TypeError}
	}

	return &Value{
		t: TypeStr,
		s: string(f),
	}
}

func (callee Value) WriteFile(args List) *Value {
	if len(args) != 2 {
		throw("%s, line %d: 'writefile' takes 2 args (filename data), got %d",
			callee.file, callee.line, len(args))
		return &Value{t: TypeError}
	}

	err := os.WriteFile(args[0].String(), []byte(args[1].String()), 0666)
	if err != nil {
		fmt.Println(err)
		return &Value{t: TypeError}
	}
	return nil
}

func (callee Value) StrSplit(args List) *Value {
	if len(args) != 2 {
		throw("%s, line %d: 'strsplit' takes 2 args (str delim), got %d",
			callee.file, callee.line, len(args))
		return &Value{t: TypeError}
	}

	strs := strings.Split(args[0].String(), args[1].String())
	
	l := List{}
	for _, str := range strs {
		l = append(l, Value{
			t: TypeStr,
			s: str,
		})
	}
	return &Value{
		t: TypeList,
		l: l,
	}
}

func (callee Value) Nth(args List) *Value {
	if len(args) != 2 {
		throw("%s, line %d: 'nth' takes 2 args (n list), got %d",
			callee.file, callee.line, len(args))
		return &Value{t: TypeError}
	}
	return &args[1].List()[args[0].Int()]
}

func (callee Value) Head(args List) *Value {
	if len(args) != 1 {
		throw("%s, line %d: 'head' takes 1 args (list), got %d",
			callee.file, callee.line, len(args))
		return &Value{t: TypeError}
	}
	return &args[0].List()[0]
}

func (callee Value) Tail(args List) *Value {
	if len(args) != 1 {
		throw("%s, line %d: 'tail' takes 1 args (list), got %d",
			callee.file, callee.line, len(args))
		return &Value{t: TypeError}
	}
	v := args[0]
	v.l = v.List()[1:]
	return &v
}

func (callee Value) Last(args List) *Value {
	if len(args) != 1 {
		throw("%s, line %d: 'last' takes 1 args (list), got %d",
			callee.file, callee.line, len(args))
		return &Value{t: TypeError}
	}
	l := args[0].List()
	return &l[len(l) - 1]
}

func (callee Value) Append(args List) *Value {
	if len(args) != 2 {
		throw("%s, line %d: 'append' takes 2 args (list ...), got %d",
			callee.file, callee.line, len(args))
	}
	res := args[0]
	res.l = append(res.l, args[1])
	return &res
}

func (callee Value) Join(args List) *Value {
	if len(args) != 2 {
		throw("%s, line %d: 'cons' takes 2 args (val val), got %d",
			callee.file, callee.line, len(args))
	}

	res := callee
	res.t = TypeList
	switch {
	case args[0].t == args[1].t:
		res.l = List{args[0], args[1]}
	case args[0].t == TypeList && args[1].t != TypeList:
		res.l = args[0].List()
		res.l = append(res.l, args[1])
	case args[0].t != TypeList && args[1].t == TypeList:
		res.l = append(List{args[0]}, args[1].List()...)
	}
	return &res
}

func (callee Value) Flat(args List) *Value {
	if len(args) != 1 {
		throw("%s, line %d: 'flat' takes 1 args (list ...), got %d",
			callee.file, callee.line, len(args))
	}

	var rl = List{}
	for _, l := range args[0].List() {
		rl = append(rl, l)
	}
	return &Value{
		t: TypeList,
		l: rl,
		line: callee.line,
		file: callee.file,
	}
}
