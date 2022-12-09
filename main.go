package main

import (
	"bufio"
	"fmt"
	"os"
)

var rd *bufio.Reader
var fromfile = false

func die(f string, args ...any) {
	fmt.Printf(f, args...)
	fmt.Println()
	os.Exit(1)
}

func throw(f string, args ...any) {
	fmt.Printf(f, args...)
	fmt.Println()
	if fromfile {
		os.Exit(1)
	}
}

func main() {
	root := Fn{locals: nil}
	if len(os.Args) == 1 {
		rd = bufio.NewReader(os.Stdin)
		for {
			root.first = &Block{fn: &root}
			fmt.Print("> ")
			root.first.Gen(GetValue())
			//for _, i := range root.first.body {
			//	i.Print()
			//}
			root.first.Run()
			stack.Top().Print()
			fmt.Println()
		}
	} else if len(os.Args) == 2 {
		fromfile = true
		reader, err := os.Open(os.Args[1])
		rd = bufio.NewReader(reader)
		if err != nil {
			die("Could not read file %s\nUsage: %s [file]",
				os.Args[1], os.Args[0])
		}

		v := Value{t: TypeInt}
		root.first = &Block{fn: &root}
		for v.t != TypeError {
			v = GetValue()
			root.first.Gen(v)
		}
		root.first.Run()
	} else {
		die("Usage: %s [file]", os.Args[0])
	}

}
