package sllm

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

const (
	tmplEsc = '`'
	nmSep   = ':'
)

var tEsc2 = []byte{tmplEsc, tmplEsc}
var nmSepStr = []byte{nmSep}

// ValueEsc is used by Expand to escap the argument when written as value of
// a parameter. It is assumed that a user of this package should not use this
// type directly. However the type it will be needed if one has to provide an
// own implemenetation of the writeArg parameter of the Expand function.
type ValueEsc struct {
	wr io.Writer
}

// Write escapes the content so that it can be reliably recognized in a sllm
// message, i.e. replace a '`' with '``'.
func (ew ValueEsc) Write(p []byte) (n int, err error) {
	var i int
	var b1 [1]byte
	bs := b1[:]
	for _, b := range p {
		if b == tmplEsc {
			i, err = ew.wr.Write(tEsc2)
		} else {
			bs[0] = b
			i, err = ew.wr.Write(bs)
		}
		if err != nil {
			return n, err
		}
		n += i
	}
	return n, nil
}

// SyntaxError describes errors of the sllm template syntax in a message
// template.
type SyntaxError struct {
	// Tmpl is the errornous template string
	Tmpl string
	// Pas is the byte position within the template string
	Pos int
	// Err is the description of the error
	Err string
}

func (err SyntaxError) Error() string {
	var sb bytes.Buffer
	ExpandArgs(&sb, "syntax error in `tmpl`:`pos`:`desc`",
		nil, err.Tmpl, err.Pos, err.Err)
	return sb.String()
}

type ParamWriter = func(wr ValueEsc, idx int, name string) (int, error)

func Expand(
	wr io.Writer,
	tmpl string,
	writeArg ParamWriter,
) (n int, err error) {
	var b1 [1]byte
	bs := b1[:]
	valEsc := ValueEsc{wr}
	tlen := len(tmpl)
	idx := 0
	for i := 0; i < tlen; i++ {
		if b := tmpl[i]; b == tmplEsc {
			wn, err := wr.Write(tEsc2[:1])
			n += wn
			if err != nil {
				return n, err
			}
			i++
			switch {
			case i >= tlen:
				return n, SyntaxError{tmpl, i, "unterminated argument"}
			case tmpl[i] == tmplEsc:
				wn, err := wr.Write(tEsc2[:1])
				n += wn
				if err != nil {
					return n, err
				}
			default:
				name := tmpl[i:]
				nmLen := strings.IndexByte(name, tmplEsc)
				if nmLen < 0 {
					return n, SyntaxError{tmpl, i, "unterminated argument"}
				}
				name = name[:nmLen]
				if strings.IndexByte(name, nmSep) >= 0 {
					return n, SyntaxError{tmpl, i, fmt.Sprintf("name contains '%c'", nmSep)}
				}
				wn, err = io.WriteString(wr, name)
				n += wn
				if err != nil {
					return n, err
				}
				wn, err = wr.Write(nmSepStr)
				n += wn
				if err != nil {
					return n, err
				}
				wn, err = writeArg(valEsc, idx, name)
				n += wn
				if err != nil {
					return n, err
				}
				wn, err = wr.Write(tEsc2[:1])
				n += wn
				if err != nil {
					return n, err
				}
				idx++
				i += nmLen
			}
		} else {
			bs[0] = b
			wn, err := wr.Write(bs)
			n += wn
			if err != nil {
				return n, err
			}
		}
	}
	return n, nil
}

type IllegalArgIndex struct {
	Tmpl string
	Pos  int
}

func (err IllegalArgIndex) Error() string {
	var sb bytes.Buffer
	ExpandArgs(&sb, "argument index 'idx' out of range in template `tmpl`",
		nil, err.Pos, err.Tmpl)
	return sb.String()
}

func ExpandArgs(
	wr io.Writer,
	tmpl string,
	undef []byte,
	argv ...interface{},
) (n int, err error) {
	return Expand(wr, tmpl, func(wr ValueEsc, idx int, name string) (int, error) {
		if idx < 0 || idx >= len(argv) {
			if undef == nil {
				return 0, IllegalArgIndex{tmpl, idx}
			} else {
				n, err = wr.Write(undef)
			}
		} else {
			n, err = writeVal(wr, argv[idx])
		}
		return n, err
	})
}

type UndefinedArg struct {
	Tmpl string
	Arg  string
}

func (err UndefinedArg) Error() string {
	var sb bytes.Buffer
	ExpandArgs(&sb, "undefined argument for `arg` in template `tmpl`",
		nil, err.Arg, err.Tmpl)
	return sb.String()
}

type ArgMap = map[string]interface{}

func ExpandMap(wr io.Writer, tmpl string, undef []byte, args ArgMap) (n int, err error) {
	return Expand(wr, tmpl, func(wr ValueEsc, idx int, name string) (int, error) {
		if val, ok := args[name]; !ok {
			if undef == nil {
				return 0, UndefinedArg{tmpl, name}
			} else {
				n, err = wr.Write(undef)
			}
		} else {
			n, err = writeVal(wr, val)
		}
		return n, err
	})
}

// ExtractParams extracs the parameter names from template tmpl and appends them
// to appendTo.
func ExtractParams(appendTo []string, tmpl string) ([]string, error) {
	_, err := Expand(ioutil.Discard, tmpl,
		func(wr ValueEsc, idx int, name string) (int, error) {
			appendTo = append(appendTo, name)
			return len(name), nil
		},
	)
	return appendTo, err
}

var writeVal = fmt.Fprint
