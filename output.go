package sllm

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

const (
	tmplEsc = '`'
	nmSep   = ':'
)

var (
	tEsc1    = []byte{tmplEsc}
	nmSepStr = []byte{nmSep}
	tokenStr = string([]byte{tmplEsc, nmSep})
)

// ValueEsc is used by Expand to escape the argument when written as value of
// a parameter. It is assumed that a user of this package should not use this
// type directly. However the type it will be needed if one has to provide an
// own implemenetation of the writeArg parameter of the Expand function.
type valEsc struct {
	wr  io.Writer
	buf []byte
}

// Write escapes the content so that it can be reliably recognized in a sllm
// message, i.e. replace a backtick '`' with two backticks '``'.
func (ew *valEsc) Write(p []byte) (n int, err error) {
	tmp := ew.buf[:0]
	for _, b := range p {
		if b == tmplEsc {
			tmp = append(tmp, tmplEsc)
		}
		tmp = append(tmp, b)
	}
	ew.buf = tmp
	return ew.wr.Write(tmp)
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
type ParamWriter = func(wr io.Writer, idx int, name string) (int, error)

// Expand writes a message to the io.Writer wr by expanding all arguments of
// the given template tmpl. The actual process of expanding an argument is
// left to the given ParamWriter writeArg.
func Expand(wr io.Writer, tmpl string, writeArg ParamWriter) (n int, err error) {
	errTmpl := SyntaxError{Tmpl: tmpl}
	vEsc := valEsc{wr: wr}
	off, lastIdx := 0, -1
	for len(tmpl) > 0 {
		i := strings.IndexByte(tmpl, tmplEsc)
		if i < 0 {
			m, err := io.WriteString(wr, tmpl)
			return n + m, err
		}
		if i++; i >= len(tmpl) {
			errTmpl.Pos = off + i
			errTmpl.Err = "unterminated argument"
			return n, errTmpl
		}
		if tmpl[i] == tmplEsc {
			i++
			m, err := io.WriteString(wr, tmpl[:i])
			if err != nil {
				return n + m, err
			}
			tmpl = tmpl[i:]
			off += i
			continue
		}
		argName, argIdx, pLen, err := parseArg(tmpl[i:])
		if err != nil {
			errTmpl.Pos = off + i + pLen
			errTmpl.Err = err.Error()
			return n, errTmpl
		}
		m, err := io.WriteString(wr, tmpl[:i+len(argName)])
		n += m
		if err != nil {
			return n, err
		}
		m, err = wr.Write(nmSepStr)
		n += m
		if err != nil {
			return n, err
		}
		if argIdx < 0 {
			lastIdx++
			argIdx = lastIdx
		}
		m, err = writeArg(&vEsc, argIdx, argName)
		n += m
		if err != nil {
			return n, err
		}
		m, err = wr.Write(tEsc1)
		n += m
		if err != nil {
			return n, err
		}
		tmpl = tmpl[i+pLen:]
	}
	return n, err
}

// Expands uses Expand to return the expanded temaplate as a string.
func Expands(tmpl string, writeArg ParamWriter) (string, error) {
	var buf bytes.Buffer
	_, err := Expand(&buf, tmpl, writeArg)
	return buf.String(), err
}

func parseArg(tmpl string) (argName string, argIdx, parseLen int, err error) {
	sep := strings.IndexAny(tmpl, tokenStr)
	if sep < 0 {
		return "", -1, 0, errors.New("unterminated argument")
	}
	if tmpl[sep] == tmplEsc {
		return tmpl[:sep], -1, sep + 1, nil
	}
	argName = tmpl[:sep]
	if sep++; sep >= len(tmpl) {
		return "", -1, 0, errors.New("unterminated argument")
	}
	b := tmpl[sep]
	if b == tmplEsc {
		return argName, -1, sep, errors.New("empty explicit index")
	}
	if b < '0' || b > '9' {
		return argName, -1, sep, errors.New("not a digit in explicit arg index")
	}
	argIdx = int(b) - '0'
	for sep++; sep < len(tmpl); sep++ {
		b := tmpl[sep]
		switch {
		case b == tmplEsc:
			return argName, argIdx, sep + 1, nil
		case b < '0' || b > '9':
			return argName, -1, sep, errors.New("not a digit in explicit arg index")
		}
		argIdx = 10*argIdx + int(b) - '0'
	}
	return argName, -1, 0, errors.New("unterminated argument")
}
