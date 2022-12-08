Axe is a basic Lisp-like language interpreter that I wrote in a few hours.

Unlike lisp, it is based on vectors rather than linked lists.
If the interpreter were well optimized, this would likely lead to improved
performance over other lisps.  However, it's just a basic tree-walk interpreter.

I might keep expanding it at some point, I do want another scripting lang
that I can use on Plan 9.

That said, I feel guilty writing a tree-walk interpreter, so I'll probably
rewrite with JIT or as a VM if I continue working on it.

It has lexical scoping, functions, and global and local variables.
It supports floats, bools, strings, symols, and functions.
It doesn't have much else.

Operations:
- `(fn ([arg:sym ...]) [expr])`
- `(do [a] [b] ...)` runs a, b, ..., returning the final value
- `(= [a:sym] [b])` sets the variable a to b
- `(+ [a] [b] ...)`
- `(- [a] [b] ...)`
- `(* [a] [b] ...)`
- `(/ [a] [b] ...)`
- `(== [a] [b] ...)` Returns true if a, b, ... are equal
- `(!= [a] [b] ...)` Returns true if any of b, ... are not equal to 1
- `(> [a] [b])` Returns true if a is greater than b
- `(> [a] [b])` Returns true if a is lesser than b
- `(>= [a] [b])` Returns true if a is greater than or equal to b
- `(<= [a] [b])` Returns true if a is lesser than or equal to b
- `(or [a] [b] ...)` Returns true if any of a, b, ... are true
- `(and [a] [b] ...)` Returns true if any of a, b, ... are false
- `(not [a])` Inverts a
- `(exit [num])` exits with error code [num], or 0 if not present
- `(print [val])` prints val
- `(cond ([test1] [expr1]) ([test2] [expr2] ...))` runs each test until one is true true, returning the matching expr if so.  Returns error if none match
- `(while [test] [expr])` Runs expr while test is true
