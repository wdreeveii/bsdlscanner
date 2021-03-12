package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/scanner"
)

func main() {

	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Println("Please provide the file that you want to parse as the first argument")
		return
	}
	infile := os.Args[1]

	f, err := os.Open(infile)
	if err != nil {
		fmt.Println(err)
		return
	}

	var s scanner.Scanner
	s.Init(f)
	s.Whitespace ^= 1 << '\n' // don't skip new lines

	data := getPinMapString(&s)

	var pmapscan scanner.Scanner
	pmapscan.Init(strings.NewReader(data))

	signals := buildPinMap(&pmapscan)

	var pmap = make(map[string]string)
	var pins = make([]string, 0, len(pmap))
	for sig, plist := range signals {
		for _, pin := range plist {
			pins = append(pins, pin)
			pmap[pin] = sig
		}
	}

	sort.Sort(ByPinOrder(pins))

	for _, k := range pins {
		fmt.Println(k, pmap[k])
	}

}

func parsePin(s string) (int, int) {
	var ii int
	for ii = 0; ii < len(s); ii++ {
		if (s[ii] >= 'A' && s[ii] <= 'Z') || (s[ii] >= 'a' && s[ii] <= 'z') {

		} else {
			break
		}
	}

	pnum, err := strconv.Atoi(s[ii:])
	if err != nil {
		panic(err)
	}
	return ii, pnum
}

type ByPinOrder []string

func (a ByPinOrder) Len() int      { return len(a) }
func (a ByPinOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPinOrder) Less(i, j int) bool {
	il, ival := parsePin(a[i])
	jl, jval := parsePin(a[j])
	if il == jl {
		if a[i][:il] == a[j][:jl] {
			return ival < jval
		}
		return a[i] < a[j]
	}
	return il < jl
}

func buildPinMap(s *scanner.Scanner) map[string][]string {
	var m = make(map[string][]string)
	var key string
	var vals []string
	var state = "key"
	tok := s.Scan()
	for tok != scanner.EOF {
		switch state {
		case "key":
			key = s.TokenText()
			tok = s.Scan()
			state = "val"
		case "val":
			if tok == '(' {
				state = "group"
			} else if tok == ',' {
				key = ""
				vals = nil
				state = "key"
			} else {
				vals = append(vals, s.TokenText())
				m[key] = vals
			}
		case "group":
			if tok == ')' {
				state = "val"
			} else if tok == ',' {

			} else {
				vals = append(vals, s.TokenText())
				m[key] = vals
			}
		}
		//fmt.Println(s.TokenText())
		tok = s.Scan()
	}

	return m
}

func simpleScan(s *scanner.Scanner) rune {
	tok := s.Scan()
	for tok != scanner.EOF {
		if tok == '-' && s.Peek() == '-' {
			for tok != scanner.EOF && tok != '\n' {
				// gobble characters
				//fmt.Println("gobbling", s.TokenText())
				tok = s.Scan()
			}
		}
		if tok == '\n' {
			tok = s.Scan()
			continue
		}
		break
	}
	return tok
}

func getPinMapString(s *scanner.Scanner) string {
	var data string
	var state = "searching"

	tok := simpleScan(s)
	for tok != scanner.EOF {
		switch state {
		case "searching":
			if s.TokenText() == "PIN_MAP_STRING" {
				state = "setup"
			}
		case "setup":
			if tok == '=' {
				state = "found"
			}
		case "found":
			if tok == ';' {
				return data
			}
			if tok != '&' {
				t := s.TokenText()

				data += t[1 : len(t)-1]
			}
		}

		tok = simpleScan(s)
	}
	return data
	//fmt.Printf("%s: %s %v\n", s.Position, s.TokenText(), tok)
	//fmt.Printf("%s", s.TokenText())
	//return &Statement{}, nil
}
