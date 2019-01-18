package sllm

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const (
	tmplEsc = '`'
	nmSep   = ':'
)

var tEsc2 = []byte{tmplEsc, tmplEsc}
var nmSepStr = []byte{nmSep}

type ValueEsc struct {
	wr io.Writer
}

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

type SyntaxError struct {
	Tmpl string
	Pos  int
	Err  string
}

func (err SyntaxError) Error() string {
	var sb bytes.Buffer
	ExpandArgs(&sb, "syntax error in `tmpl`:`pos`:`desc`", nil, err.Tmpl, err.Pos, err.Err)
	return sb.String()
}

func Expand(
	wr io.Writer,
	tmpl string,
	writeArg func(wr ValueEsc, idx int, name string) error,
) (err error) {
	var b1 [1]byte
	bs := b1[:]
	valEsc := ValueEsc{wr}
	tlen := len(tmpl)
	idx := 0
	for i := 0; i < tlen; i++ {
		if b := tmpl[i]; b == tmplEsc {
			_, err := wr.Write(tEsc2[:1])
			if err != nil {
				return err
			}
			i++
			switch {
			case i >= tlen:
				return SyntaxError{tmpl, i, "unterminated argument"}
			case tmpl[i] == tmplEsc:
				_, err := wr.Write(tEsc2[:1])
				if err != nil {
					return err
				}
			default:
				name := tmpl[i:]
				nmLen := strings.IndexByte(name, tmplEsc)
				if nmLen < 0 {
					return SyntaxError{tmpl, i, "unterminated argument"}
				}
				name = name[:nmLen]
				if strings.IndexByte(name, nmSep) >= 0 {
					return SyntaxError{tmpl, i, fmt.Sprintf("name contains '%c'", nmSep)}
				}
				_, err = wr.Write([]byte(name))
				if err != nil {
					return err
				}
				_, err = wr.Write(nmSepStr)
				if err != nil {
					return err
				}
				err := writeArg(valEsc, idx, name)
				if err != nil {
					return err
				}
				_, err = wr.Write(tEsc2[:1])
				if err != nil {
					return err
				}
				idx++
				i += nmLen
			}
		} else {
			bs[0] = b
			_, err := wr.Write(bs)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type IllegalArgIndex struct {
	Tmpl string
	Pos  int
}

func (err IllegalArgIndex) Error() string {
	var sb bytes.Buffer
	ExpandArgs(&sb, "argument index 'idx' out of range in template `tmpl`", nil, err.Pos, err.Tmpl)
	return sb.String()
}

func ExpandArgs(
	wr io.Writer,
	tmpl string,
	undef []byte,
	argv ...interface{},
) (err error) {
	return Expand(wr, tmpl, func(wr ValueEsc, idx int, name string) error {
		if idx < 0 || idx >= len(argv) {
			if undef == nil {
				return IllegalArgIndex{tmpl, idx}
			} else {
				_, err = wr.Write(undef)
			}
		} else {
			_, err = writeVal(wr, argv[idx])
		}
		return err
	})
}

type UndefinedArg struct {
	Tmpl string
	Arg  string
}

func (err UndefinedArg) Error() string {
	var sb bytes.Buffer
	ExpandArgs(&sb, "undefined argument for `arg` in template `tmpl`", nil, err.Arg, err.Tmpl)
	return sb.String()
}

type ArgMap = map[string]interface{}

func ExpandMap(wr io.Writer, tmpl string, undef []byte, args ArgMap) (err error) {
	return Expand(wr, tmpl, func(wr ValueEsc, idx int, name string) error {
		if val, ok := args[name]; !ok {
			if undef == nil {
				return UndefinedArg{tmpl, name}
			} else {
				_, err = wr.Write(undef)
			}
		} else {
			_, err = writeVal(wr, val)
		}
		return err
	})
}

// func ExpandData(
// 	wr io.Writer,
// 	tmpl string,
// 	undef []byte,
// 	args interface{},
// ) (err error) {
// 	switch reflect.TypeOf(args).Kind() {
// 	case reflect.Struct:
// 	case reflect.Slice:
// 	case reflect.Map:
// 	case reflect.Array:
// 	default:
// 	}

// 	// return Expand(wr, tmpl, func(wr io.Writer, idx int, name string) error {

// 	// })
// }

var writeVal = fmt.Fprint
