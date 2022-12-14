package main

import (
	"io"
	"strconv"
	"unicode"
)

const EOF = 0

var line = 1
var filename string

func PeekRune() rune {
	r, _, err := rd.ReadRune()
	if err == io.EOF {
		return 0
	} else if err != nil {
		panic(err)
	}
	if rd.UnreadRune() != nil {
		panic(err)
	}
	return r
}

func GetRune() rune {
	r, _, err := rd.ReadRune()
	if err == io.EOF {
		die("%s: Unexpected early EOF", filename)
	} else if err != nil {
		panic(err)
	}
	return r
}

func SkipWS() rune {
	r := PeekRune()
	for unicode.IsSpace(r) {
		GetRune()
		if r == '\n' {
			line++
		}
		r = PeekRune()
	}
	if r == ';' {
		for r != '\n' {
			r = GetRune()
		}
		line++
		SkipWS()
	}
	return PeekRune()
}

func GetValue() Value {
	v := Value{line: line, file: filename}
	r := SkipWS()

	if r == 0 {
		return v
	} else if r == '\'' {
		GetRune()
		v.t = TypeList
		v.l = List{Value{t: TypeSym, s: "quote"}}
		v.l = append(v.l, GetValue())
	} else if r == '(' {
		v.t = TypeList
		v.l = GetList()
	} else if r == '"' {
		v.t = TypeStr
		v.s = ""
		GetRune() // First '"'
		r = GetRune()
		for r != '"' {
			v.s += string(r)
			r = GetRune()
		}
	} else if unicode.IsDigit(r) || r == '-' {
		tmp := ""
		neg := r == '-'
		for unicode.IsDigit(r) || neg {
			GetRune()
			tmp += string(r)
			r = PeekRune()
			neg = false
		}
		if tmp == "-" {
			v.t = TypeSym
			v.s = "-"
			for !unicode.IsSpace(r) && r != ')' {
				GetRune()
				v.s += string(r)
				r = PeekRune()
			}
			return v
		} 
		if r == '.' {
			v.t = TypeFloat
			GetRune()
			tmp = tmp + string(r)
			if !unicode.IsDigit(r) {
				die("%s, line %d: Started float, but didn't end, \"%s\"",
					filename, line, tmp)
			}
			for unicode.IsDigit(r) {
				GetRune()
				tmp += string(r)
				r = PeekRune()
			}
		}
		v.t = TypeFloat
		v.f, _ = strconv.ParseFloat(tmp, 64)
	} else {
		v.t = TypeSym
		v.s = ""
		for !unicode.IsSpace(r) && r != ')' && r != ';' {
			GetRune()
			v.s += string(r)
			r = PeekRune()
		}
		if v.s == "" {
			GetRune()
			throw("%s, line %d: Unbalanced rparen", filename, line)
			return Value{}
		}

		if v.s == "true" {
			v.t = TypeBool
			v.b = true
			v.s = ""
		} else if v.s == "false" {
			v.t = TypeBool
			v.b = false
			v.s = ""
		}
	}

	return v
}

func GetList() List {
	var res List
	var r rune

	r = GetRune()
	if r != '(' {
		die("%s, line %d: Tried to start list without '('", filename, line)
	}
	for PeekRune() != ')' {
		res = append(res, GetValue())
		SkipWS()
	}
	GetRune() // ')'
	return res
}
