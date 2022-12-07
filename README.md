Axe is a basic Lisp-like language interpreter that I wrote in a few hours.

I might keep expanding it at some point, I do want another scripting lang
that I can use on Plan 9.

That said, I feel guilty writing a tree-walk interpreter, so I'll probably
rewrite with JIT or as a VM if I continue working on it.

It has lexical scoping, functions, and global and local variables.
It supports floats, bools, strings, symols, and functions.
It doesn't have much else.

Operations:
- `(exit [num])` exits with error code [num], or 0 if not present
- `(= <a:sym> <b>)` sets the variable a to b
- `(+ <a> <b> [c ...])`
- `(- <a> <b> [c ...])`
- `(* <a> <b> [c ...])`
- `(/ <a> <b> [c ...])`
- `(fn (<arg:sym ...>) <...>)`
