package main

var fid int

func (b *Block) Gen(v Value) {
	switch v.t {
	case TypeInt,
		TypeFloat,
		TypeBool,
		TypeStr:
		b.body = append(b.body, Ins{op: InsImm, imm: v})
	case TypeSym:
		found := false
		if b.fn.locals != nil {
			_, found = b.fn.locals[v.s]
		}

		if !found {
			_, found = globals[v.s]
		}

		if found {
			b.body = append(b.body, Ins{op: InsLoadV, imm: v})
		} else {
			throw("Line %d: Failed to find variable %s", v.line, v.s)
			return
		}
	case TypeList:
		if len(v.l) == 0 {
			throw("Line %d: Empty expression", v.line)
			return
		}
		switch v.l[0].s {
		case "fn":
			// (fn (&a &b &...) $expr)
			newf := Fn{locals: map[string]Value{}}
			fid++
			newf.id = fid

			if len(v.l) != 3 {
				throw("Line %d: fn takes 2 args, got %d", v.l[0].line, len(v.l))
				return
			}
			if v.l[1].t != TypeList {
				throw("Line %d: Args should be a list", v.l[1].line)
				return
			}

			for _, arg := range v.l[1].l {
				if arg.t != TypeSym {
					throw("Line %d: Args should be symbols", arg.line)
					return
				}
				newf.args = append(newf.args, arg.s)
				newf.locals[arg.s] = Value{}
			}
			newf.first = new(Block)
			newf.first.fn = &newf
			newf.first.Gen(v.l[2])
			b.body = append(b.body, Ins{
				op:  InsImm,
				imm: Value{t: TypeFn, fn: newf},
			})
		case "mu":
			// (mu $expr)
			newb := Block{[]Ins{}, b.fn}
			if len(v.l) != 2 {
				throw("Line %d: mu takes 1 args, got %d",
					v.l[0].line, len(v.l))
			}
			newb.Gen(v.l[1])
			b.body = append(b.body, Ins{
				op:  InsImm,
				imm: Value{t: TypeBlock, bl: &newb},
			})
		case "do":
			// (do $expr1 $expr2 &...)
			for _, expr := range v.l[1:] {
				b.Gen(expr)
			}
		case "if":
			// (if ($cond) $exprt $exprf)
			bt := &Block{fn: b.fn}
			bf := &Block{fn: b.fn}

			if len(v.l) != 4 {
				throw("Line %d: if takes 3 args (cond, true, false), got %d",
					v.l[0].line, len(v.l))
				return
			}

			b.Gen(v.l[1])
			bt.Gen(v.l[2])
			bf.Gen(v.l[3])
			b.body = append(b.body, Ins{op: InsIf, bt: bt, bf: bf})
		case "cond":
			// (cond ($cond1 $expr1) &($cond2 $expr2) &...)
			if len(v.l) < 2 {
				throw("Line %d: cond requires at least one condition",
					v.l[0].line)
				return
			}

			cur := b
			for _, arg := range v.l[1:] {
				bt := &Block{fn: b.fn}
				bf := &Block{fn: b.fn}

				if arg.t != TypeList || len(arg.l) != 2 {
					throw("Line %d: Conditions should be of (cond expr) format",
						arg.line)
					return
				}

				cur.Gen(arg.l[0])
				bt.Gen(arg.l[1])
				cur.body = append(cur.body, Ins{op: InsIf, bt: bt, bf: bf})
				cur = bf
			}

		// First add args (backwards), then callee, then call instruction
		case "=":
			// (= $var $expr)
			if len(v.l) != 3 {
				throw("Line %d: Wrong arg count for '='", v.line)
				return
			}
			if b.fn.locals != nil {
				if _, ok := b.fn.locals[v.l[1].s]; !ok {
					if _, ok := globals[v.l[1].s]; !ok {
						// Not found globally or locally, set local
						b.fn.locals[v.l[1].s] = Value{}
					}
				}
			} else { // No local scope, set global
				if _, ok := globals[v.l[1].s]; !ok {
					globals[v.l[1].s] = Value{}
				}
			}
			b.Gen(v.l[2])
			b.body = append(b.body, Ins{op: InsStoreV, imm: v.l[1]})
		case "<", "<=", ">", ">=":
			// ($op $expr1 $expr2)
			if len(v.l) != 3 {
				throw("Line %d: Wrong arg count to binary op", v.line)
				return
			}
			b.Gen(v.l[2])
			b.Gen(v.l[1])
			b.body = append(b.body, Ins{
				op:  InsImm,
				imm: Value{t: TypeBuiltin, s: v.l[0].s},
			})
			b.body = append(b.body, Ins{
				op:   InsCall,
				imm:  Value{line: v.l[0].line},
				argn: 2,
			})
		case "+", "-", "*", "/", "==", "and", "or":
			// ($op $expr1 $expr2 &expr3 &...)
			if len(v.l) < 3 {
				throw("Line %d: Not enough args to arithmetic op", v.line)
				return
			}
			for i := len(v.l) - 1; i > 0; i-- {
				b.Gen(v.l[i])
			}
			b.body = append(b.body, Ins{
				op:  InsImm,
				imm: Value{t: TypeBuiltin, s: v.l[0].s},
			})
			b.body = append(b.body, Ins{
				op:   InsCall,
				imm:  Value{line: v.l[0].line},
				argn: len(v.l) - 1,
			})
		case "quote":
			// (quote $expr)
			b.body = append(b.body, Ins{op: InsImm, imm: v})
		case "print":
			// (print $expr)
			if len(v.l) != 2 {
				throw("Line %d: Wrong  number of args to unary op", v.line)
				return
			}
			b.Gen(v.l[1]) // Push arg
			b.body = append(b.body, Ins{
				op:  InsImm,
				imm: Value{t: TypeBuiltin, s: v.l[0].s},
			})
			b.body = append(b.body, Ins{
				op:   InsCall,
				imm:  Value{line: v.line},
				argn: 1,
			})
		default:
			// ($callee &expr1 &expr2 &...)
			for i := len(v.l) - 1; i >= 0; i-- { // >= to pass callee
				b.Gen(v.l[i])
			}
			b.body = append(b.body, Ins{
				op:   InsCall,
				imm:  Value{line: v.l[0].line},
				argn: len(v.l) - 1,
			})
		}
	}
}
