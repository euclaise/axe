package main

import (
	"fmt"
	"os"
	"bufio"
)

var rd *bufio.Reader
var fromfile = false

func die(f string, args ...any) {
	fmt.Printf(f, args...)
	fmt.Println()
	os.Exit(1)
}

func fatal(err error, f string, args ...any) {
	if err != nil {
		die(f, args...)
	}
}

func main() {
	if len(os.Args) == 1 {
		rd = bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")
			GetValue().Eval().Print()
			fmt.Println()
		}
	} else if len(os.Args) == 2 {
		reader, err := os.Open(os.Args[1])
		rd = bufio.NewReader(reader)
		fatal(err, "Could not read file %s\nUsage: %s [file]",
				os.Args[1], os.Args[0])

		v := Value{t: TypeInt}
		for v.t != TypeError {
			v = GetValue().Eval()
		}
	} else {
		die("Usage: %s [file]", os.Args[0])
	}

}
