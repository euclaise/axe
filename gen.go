package main

func (b *Block) SetVar(s string, v Value) {
	_, found_global := globals[s]
	if b.fn.locals != nil {
		if _, ok := b.fn.locals[s]; !ok && !found_global {
			// Not found globally or locally, set local
			b.fn.locals[s] = v
		}
	} else { // No local scope, set global
		globals[s] = v
	}
}

func (b *Block) LookupVar(s string) *Value {
	if b.fn.locals != nil {
		if v, ok := b.fn.locals[s]; ok  {
			return &v
		}
	}
	if v, ok := globals[s]; ok {
		return &v
	}
	return nil
}

func (b *Block) Gen(v Value) bool {
	switch v.t {
	case TypeFloat,
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
			throw("%s, line %d: Failed to find variable %s",
				v.file, v.line, v.s)
			return false
		}
	case TypeList:
		if len(v.l) == 0 {
			throw("%s, line %d: Empty expression", v.file, v.line)
			return false
		}
		switch v.l[0].s {
		case "quote":
			if len(v.l) != 2 {
				throw("'quote' takes 1 args")
				return false
			}
			b.body = append(b.body, Ins{
				op: InsImm,
				imm: v.l[1],
			})
		case "fn":
			// (fn (&a &b &...) $expr)
			newf := Fn{locals: map[string]Value{}}

			if len(v.l) != 3 {
				throw("%s, line %d: fn takes 2 args, got %d",
					v.l[0].file, v.l[0].line, len(v.l))
				return false
			}

			for _, arg := range v.l[1].List() {
				if arg.t != TypeSym {
					throw("%s, line %d: Args should be symbols",
						arg.file, arg.line)
					return false
				}
				newf.args = append(newf.args, arg.Symbol())
				newf.locals[arg.s] = Value{}
			}
			newf.first = new(Block)
			newf.first.fn = &newf
			newf.first.Gen(v.l[2])
			b.body = append(b.body, Ins{
				op:  InsImm,
				imm: Value{t: TypeFn, fn: newf},
			})
		case "macro":
			// (fn (&a &b &...) $expr)
			newf := Fn{locals: map[string]Value{}, macro: true}

			if len(v.l) != 4 {
				throw("%s, line %d: macro takes 2 args, got %d",
						v.l[0].file, v.l[0].line, len(v.l))
				return false
			}

			for _, arg := range v.l[2].List() {
				if arg.t != TypeSym {
					throw("%s, line %d: Args should be symbols",
						arg.file, arg.line)
					return false
				}
				newf.args = append(newf.args, arg.s)
				newf.locals[arg.s] = Value{}
			}
			newf.first = new(Block)
			newf.first.fn = &newf
			newf.first.Gen(v.l[3])
			b.body = append(b.body, Ins{
				op:  InsImm,
				imm: Value{t: TypeFn, fn: newf},
			})

			b.SetVar(v.l[1].Symbol(), Value{t: TypeFn, fn: newf})
			b.body = append(b.body, Ins{
				op: InsStoreV,
				imm: Value{
					t: TypeFn,
					fn: newf,
				},
			})
		case "mu":
			// (mu $expr)
			newb := Block{[]Ins{}, b.fn}
			if len(v.l) != 2 {
				throw("%s, line %d: mu takes 1 args, got %d",
					v.l[0].file, v.l[0].line, len(v.l))
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
				throw("%s, line %d: if takes 3 args (cond, true, false), got %d",
					v.l[0].file, v.l[0].line, len(v.l))
				return false
			}

			b.Gen(v.l[1])
			bt.Gen(v.l[2])
			bf.Gen(v.l[3])
			b.body = append(b.body, Ins{op: InsIf, bt: bt, bf: bf})
		case "cond":
			// (cond ($cond1 $expr1) &($cond2 $expr2) &...)
			if len(v.l) < 2 {
				throw("%s, line %d: cond requires at least one condition",
					v.l[0].file, v.l[0].line)
				return false
			}

			cur := b
			for _, arg := range v.l[1:] {
				bt := &Block{fn: b.fn}
				bf := &Block{fn: b.fn}

				if len(arg.l) != 2 {
					throw("%s, line %d: Conditions should be of (cond expr) format",
						arg.file, arg.line)
					return false
				}

				cur.Gen(arg.List()[0])
				bt.Gen(arg.List()[1])
				cur.body = append(cur.body, Ins{op: InsIf, bt: bt, bf: bf})
				cur = bf
			}
			cur.body = append(cur.body, Ins{
				op: InsImm,
				imm: Value{t: TypeError},
			})

		// First add args (backwards), then callee, then call instruction
		case "=":
			// (= $var $expr)
			if len(v.l) != 3 {
				throw("%s, line %d: Wrong arg count for '='", v.file, v.line)
				return false
			}
			s := v.l[1].Symbol()
			if _, ok := builtins[s]; ok {
				throw("%s, line %d: '%s' is a builtin, cannot be rebound",
					v.l[1].file, v.l[1].line, s)
					return false
			}
			b.SetVar(s, Value{})

			b.Gen(v.l[2])
			b.body = append(b.body, Ins{op: InsStoreV, imm: v.l[1]})
		default:
			// ($callee &expr1 &expr2 &...)
			is_macro := false
			if v.l[0].t == TypeSym {
				lookup := b.LookupVar(v.l[0].s)
				is_macro = lookup != nil &&
						lookup.t == TypeFn &&
						lookup.fn.macro == true
			}
			if is_macro {
				for i := len(v.l) - 1; i > 0; i-- { // >= to pass callee
					b.body = append(b.body, Ins{
						op: InsImm,
						imm: Value{
							t: TypeList,
							l: List{
								Value{
									t: TypeSym,
									s: "quote",
								},
								v.l[i],
							},
						},
					})
					b.Gen(v.l[0])
				}
			} else {
				for i := len(v.l) - 1; i >= 0; i-- { // >= to pass callee
					b.Gen(v.l[i])
				}
			}
			b.body = append(b.body, Ins{
				op:   InsCall,
				imm:  Value{line: v.l[0].line},
				argn: len(v.l) - 1,
			})
		}
	}
	return true
}
