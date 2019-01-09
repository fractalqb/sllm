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

type valEsc struct {
	wr io.Writer
}

func (ew valEsc) Write(p []byte) (n int, err error) {
	var i int
	for _, b := range p {
		if b == tmplEsc {
			i, err = ew.wr.Write(tEsc2)
		} else {
			i, err = ew.wr.Write([]byte{b})
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

func (e SyntaxError) Error() string {
	var sb bytes.Buffer
	ExpandArgs(&sb, "syntax error in `tmpl`:`pos`:`desc`", nil, e.Tmpl, e.Pos, e.Err)
	return sb.String()
}

func Expand(
	wr io.Writer,
	tmpl string,
	writeArg func(wr io.Writer, idx int, name string) error,
) (err error) {
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
				_, err = wr.Write([]byte(name))
				if err != nil {
					return err
				}
				_, err = wr.Write(nmSepStr)
				if err != nil {
					return err
				}
				err := writeArg(valEsc{wr}, idx, name)
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
			_, err := wr.Write([]byte{b})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type ArgIndexError struct {
	Tmpl string
	Pos  int
}

func (iior ArgIndexError) Error() string {
	var sb bytes.Buffer
	ExpandArgs(&sb, "argument index 'idx' out of range in template `tmpl`", nil, iior.Pos, iior.Tmpl)
	return sb.String()
}

func ExpandArgs(wr io.Writer, tmpl string, undef []byte, argv ...interface{}) (err error) {
	return Expand(wr, tmpl, func(wr io.Writer, idx int, name string) error {
		if idx < 0 || idx >= len(argv) {
			if undef == nil {
				return ArgIndexError{tmpl, idx}
			} else {
				_, err = wr.Write(undef)
			}
		} else {
			_, err = fmt.Fprint(wr, argv[idx])
		}
		return err
	})
}

type UndefArgError struct {
	Tmpl string
	Arg  string
}

func (ua UndefArgError) Error() string {
	var sb bytes.Buffer
	ExpandArgs(&sb, "undefined argument for `arg` in template `tmpl`", nil, ua.Arg, ua.Tmpl)
	return sb.String()
}

func ExpandMap(wr io.Writer, tmpl string, undef []byte, args map[string]interface{}) (err error) {
	return Expand(wr, tmpl, func(wr io.Writer, idx int, name string) error {
		if val, ok := args[name]; !ok {
			if undef == nil {
				return UndefArgError{tmpl, name}
			} else {
				_, err = wr.Write(undef)
			}
		} else {
			_, err = fmt.Fprint(wr, val)
		}
		return err
	})
}
