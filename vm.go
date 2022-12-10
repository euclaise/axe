package main

import (
	"fmt"
	"os"
)

const (
	InsImm = iota
	InsStoreV
	InsLoadV
	InsIf
	InsCall
)

type Ins struct {
	op   int
	imm  Value
	argn int
	bt   *Block // true block
	bf   *Block // false block
}

type Fn struct {
	args   []string
	locals map[string]Value
	first  *Block // main block
	macro bool
}

type Block struct {
	body []Ins
	fn   *Fn
}

var globals = map[string]Value{
	"+": {t: TypeBuiltin, s: "+" },
	"-": {t: TypeBuiltin, s: "-" },
	"*": {t: TypeBuiltin, s: "*" },
	"/": {t: TypeBuiltin, s: "/" },
	"and": {t: TypeBuiltin, s: "and" },
	"or": {t: TypeBuiltin, s: "or" },
	">": {t: TypeBuiltin, s: ">" },
	"<": {t: TypeBuiltin, s: "<" },
	">=": {t: TypeBuiltin, s: ">=" },
	"<=": {t: TypeBuiltin, s: "<=" },
	"==": {t: TypeBuiltin, s: "==" },
	"!=": {t: TypeBuiltin, s: "!=" },
	"print": {t: TypeBuiltin, s: "print" },
	"exit": {t: TypeBuiltin, s: "exit" },
}
var locals = []map[string]Value{}

type Stack []Value

var stack = Stack{Value{}}

func (s *Stack) Pop() Value {
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}

func (s Stack) Top() Value {
	return s[len(s)-1]
}

func (s *Stack) Push(v Value) {
	*s = append(*s, v)
}

func (a Value) Eq(b Value) bool {
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
			if !a.l[i].Eq(b.l[i]) {
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

func (ins Ins) Run(fn *Fn) {
	switch ins.op {
	case InsImm:
		stack.Push(ins.imm)
	case InsLoadV:
		ok := false
		var val Value
		if len(locals) > 0 {
			val, ok = locals[len(locals)-1][ins.imm.s]
		}
		if !ok {
			val, ok = globals[ins.imm.s]
		}
		if !ok {
			throw("Line %d (vm): Could not find variable %s",
				ins.imm.line, ins.imm.s)
		}
		stack.Push(val)
	case InsStoreV:
		top := stack.Top()
		if _, ok := globals[ins.imm.s]; ok {
			globals[ins.imm.s] = top
			return
		}
		if len(locals) > 0 {
			locals[len(locals)-1][ins.imm.s] = top
		} else {
			globals[ins.imm.s] = top
		}
	case InsIf:
		cond := stack.Pop()
		if cond.t != TypeBool {
			throw("Line %d: 'if' on non-bool", cond.line)
			return
		}
		if cond.b {
			ins.bt.Run()
		} else {
			ins.bf.Run()
		}
	case InsCall:
		callee := stack.Pop()
		args := []Value{}
		for i := 0; i < ins.argn; i++ {
			args = append(args, stack.Pop())
		}
		if callee.t == TypeFn {
			locals = append(locals, map[string]Value{})
			for i, arg := range args {
				locals[len(locals)-1][callee.fn.args[i]] = arg
			}
			callee.fn.first.Run()
			locals = locals[:len(locals)-1]
		} else if callee.t == TypeBlock {
			callee.bl.Run()
		} else if callee.t == TypeBuiltin {
			switch callee.s {
			case "==":
				if len(args) < 2 {
					throw("Line %d: Too few args to '=='", args[1].line)
					return
				}
				res := Value{t: TypeBool, b: true, line: callee.line}
				for _, val := range args[1:] {
					if !val.Eq(args[0]) {
						res.b = false
						break
					}
				}
				stack.Push(res)
			case "and":
				if len(args) < 2 {
					throw("Line %d: Too few args to 'and'", args[1].line)
					return
				}
				res := Value{t: TypeBool, b: true, line: callee.line}
				for _, val := range args {
					if val.t != TypeBool {
						throw("Line %d: 'and' on non-bool", val.line)
						return
					}
					res.b = res.b && val.b
				}
				stack.Push(res)
			case "or":
				if len(args) < 2 {
					throw("Line %d: Too few args to 'or'", args[1].line)
					return
				}
				res := Value{t: TypeBool, b: false, line: callee.line}
				for _, val := range args {
					if val.t != TypeBool {
						throw("Line %d: 'or' on non-bool", val.line)
						return
					}
					res.b = res.b || val.b
				}
				stack.Push(res)
			case "<":
				if len(args) != 2 {
					throw("Line %d: Wrong number of args to '<'", args[1].line)
					return
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
				stack.Push(res)
			case "<=":
				if len(args) != 2 {
					throw("Line %d: Wrong number of args to '<='", args[1].line)
					return
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
				stack.Push(res)
			case ">":
				if len(args) != 2 {
					throw("Line %d: Wrong number of args to '>'", args[1].line)
					return
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
				stack.Push(res)
			case ">=":
				if len(args) != 2 {
					throw("Line %d: Wrong number of args to '>='", args[1].line)
					return
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
				stack.Push(res)
			case "+":
				if len(args) < 2 {
					throw("Line %d: Too few args to '+'", args[1].line)
					return
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
				stack.Push(res)
			case "-":
				if len(args) < 2 {
					throw("Line %d: Too few args to '-'", args[1].line)
					return
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
				stack.Push(res)
			case "*":
				if len(args) < 2 {
					throw("Line %d: Too few args to '*'", args[1].line)
					return
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
				stack.Push(res)
			case "/":
				if len(args) < 2 {
					throw("Line %d: Too few args to '/'", args[1].line)
					return
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
				stack.Push(res)
			case "print":
				if len(args) != 1 {
					throw("'print' expects 1 args, not %d", len(args))
					return
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
					fmt.Printf("'%s\n", args[0].s)
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
				stack.Push(args[0])
			case "exit":
				if len(args) > 2 {
					throw("'exit' takes at most 1 args, not %d", len(args))
					return
				}
				if len(args) == 0 {
					os.Exit(0)
				}
				if args[0].t != TypeInt {
					throw("Line %d: Trying to exit with non-int code",
						args[0].line)
					return
				}
				os.Exit(int(args[1].i))
			default:
				throw("Line %d: Unknown builtin", callee.line)
				return
			}
		} else {
			callee.Print()
			fmt.Println()
			throw("Line %d: Call to non-fn (%d)", callee.line, callee.t)
			return
		}
	}
}

func (b Block) Run() {
	old := stack
	for _, ins := range b.body {
		ins.Run(b.fn)
	}
	old.Push(stack.Top())
	stack = old
}

func (ins Ins) Print() {
	switch ins.op {
	case InsImm:
		fmt.Print("IMM: ")
		ins.imm.Print()
		fmt.Println()
	case InsStoreV:
		fmt.Print("STOREV: ")
		ins.imm.Print()
		fmt.Println()
	case InsLoadV:
		fmt.Print("LOADV: ")
		ins.imm.Print()
		fmt.Println()
	case InsCall:
		fmt.Printf("CALL(%d)\n", ins.argn)
	case InsIf:
		fmt.Printf("IF\n")
	default:
		fmt.Println("Ins.Print: Unhandled")
	}
}
