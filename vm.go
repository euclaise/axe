package main

import (
	"fmt"
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

var btrace = false // trace builtins
var itrace = false // trace instructions
var strace = false // trace top of stack

var locals = []map[string]Value{}

type Stack []Value

var stack = Stack{Value{}}

func (s *Stack) Pop() Value {
	res := (*s)[len(*s) - 1]
	*s = (*s)[:len(*s) - 1]
	return res
}

func (s Stack) Top() Value {
	return s[len(s) - 1]
}

func (s *Stack) Push(v Value) {
	*s = append(*s, v)
}

func (ins Ins) Run(fn *Fn) bool {
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
			throw("%s, line %d (vm): Could not find variable %s",
				ins.imm.file, ins.imm.line, ins.imm.s)
			return false
		}
		stack.Push(val)
	case InsStoreV:
		top := stack.Top()
		if len(locals) > 0 {
			locals[len(locals)-1][ins.imm.s] = top
		} else {
			globals[ins.imm.s] = top
		}
	case InsIf:
		cond := stack.Pop()
		if cond.Bool() {
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
			if btrace {
				fmt.Printf(">>> btrace: %s\n", callee.s)
			}
			v := callee.bu(callee, args)
			if thrown {
				thrown = false
				stack.Push(Value{t: TypeError})
			} else if v != nil {
				stack.Push(*v)
			} else {
				return false
			}
		} else {
			throw("%s, line %d: Call to non-fn (%d)",
				callee.file, callee.line, callee.t)
			fmt.Print("Value: ")
			callee.Print()
			fmt.Println()
			return false
		}
	}
	return true
}

func (b Block) Run() {
	old := stack
	for _, ins := range b.body {
		if itrace {
			fmt.Print(">>> itrace: ")
			ins.Print()
			fmt.Println()
		}
		ins.Run(b.fn)
		if strace {
			fmt.Print(">>> strace: ")
			stack.Top().Print()
			fmt.Println()
		}
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
