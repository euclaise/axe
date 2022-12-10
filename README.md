Axe is a basic Lisp-like language interpreter that I wrote in a few
 h̶o̶u̶r̶s̶ days.

Unlike lisp, it is based on vectors rather than linked lists.
If the interpreter were well optimized, this almost certainly improves
performance on average but I haven't tested it.

I might keep expanding it at some point, I do want another scripting lang
that I can use on Plan 9.

I felt guilty writing a tree-walk interpreter, so I rewrote it to compile to a
high-level VM-like linear representation.
This change significantly improved performance.

It has lexical scoping, functions, and global and local variables.
It supports floats, bools, strings, symols, and functions.

# Operations
- `(fn (&arg1 &arg2 &...) $expr)` Create a new function.
    Use `do` for multi-expr functions
- `'$expr` or `(quote $expr)` returns `(quote $expr)`, printed as `'$expr`
- `(do $expr1 $expr2 ...)` Runs $expr1, $expr2, ..., returning the final value
- `(= $var $expr)` Sets $var to $expr
- `(+ $expr1 $expr2 ...)` Adds all values, left to right
- `(- $expr1 $expr2 ...)` Subtracts all values, left to right
- `(* $expr1 $expr2 ...)` Multiplies all values, left to right
- `(/ $expr1 $expr2 ...)` Divides all values, left to right
- `(== $expr1 $expr2 ...)` Returns true if $expr1, $expr2, ... are equal
- `(!= $expr1 $expr2 ...)` Returns true if any of $expr2, ...
    are not equal to $expr1
- `(> $expr1 $expr2)` Returns true if $expr1 is greater than $expr2
- `(> $expr1 $expr2)` Returns true if $expr1 is lesser than $expr2
- `(>= $expr1 $expr2)` Returns true if $expr1 is greater than or equal to $expr2
- `(<= $expr1 $expr2)` Returns true if $expr1 is lesser than or equal to $expr2
- `(or $expr1 $expr2 ...)` Returns true if any of $expr1, $expr2, ... are true
- `(and $expr1 $expr2 ...)` Returns true if any of $expr1, $expr2, ... are false
- `(not $expr)` Inverts $expr
- `(exit $num)` Exits with error code $num (rounded to the nearest int),
    or 0 if not present
- `(print $expr)` Prints expr
- `(cond ($cond1 $expr1) ($cond2 $expr2) &...)` Runs each test until one is
    true true, returning the matching expr if so.  Returns error if none match
- `(while $cond $expr)` Runs expr while test is true

# Example:
```
(= fib (fn (n)
    (cond
        ((== n 0) 0)
        ((== n 1) 1)
        (true (+ (fib (- n 1)) (fib (- n 2)))))))
(= a 0)
(= z (mu (if (< a 25)
        (do
            (= a (+ a 1))
            (print (fib a))
            (z))
        (print "done"))))

(z)
```
