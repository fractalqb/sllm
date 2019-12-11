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

// ValueEsc is used by Expand to escape the argument when written as value of
// a parameter. It is assumed that a user of this package should not use this
// type directly. However the type it will be needed if one has to provide an
// own implemenetation of the writeArg parameter of the Expand function.
type ValueEsc struct {
	wr io.Writer
}

// Write escapes the content so that it can be reliably recognized in a sllm
// message, i.e. replace a backtick '`' with two backticks '``'.
func (ew ValueEsc) Write(p []byte) (n int, err error) {
	tmp := make([]byte, 2*len(p))
	wp := 0
	for _, b := range p {
		if b == tmplEsc {
			tmp[wp] = tmplEsc
			wp++
			tmp[wp] = tmplEsc
		} else {
			tmp[wp] = b
		}
		wp++
	}
	return ew.wr.Write(tmp[:wp])
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
	Expand(&sb, "syntax error in `tmpl`:`pos`:`desc`",
		Args(nil, err.Tmpl, err.Pos, err.Err))
	return sb.String()
}

// ParamWriter is used by the Expand function to process an argument when it
// appears in the expand process of a template. Expand will pass the index idx
// and the name of the argument to expand, i.e. write into the writer wr.
// A ParamWriter returns the number of bytes writen and—just in case—an error.
//
// NOTE The writer wr of type ValueEsc will escape whatever ParamWriter
//      writes to wr so that the template escape symbol '`' remains
//      recognizable.
type ParamWriter = func(wr ValueEsc, idx int, name string) (int, error)

// Expand writes a message to the io.Writer wr by expanding all arguments of
// the given template tmpl. The actual process of expanding an argument is
// left to the given ParamWriter writeArg.
func Expand(wr io.Writer, tmpl string, writeArg ParamWriter) (n int, err error) {
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

// Expands uses Expand to return the expanded temaplate as a string.
func Expands(tmpl string, writeArg ParamWriter) (string, error) {
	var buf bytes.Buffer
	_, err := Expand(&buf, tmpl, writeArg)
	return buf.String(), err
}
