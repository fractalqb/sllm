package sllm

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
)

const (
	tmplEsc byte = '`'
	nmSep   byte = ':'
)

var tokenStr = string([]byte{tmplEsc, nmSep})

// ArgWriter is used by Expand to escape the argument when written as
// a value of a parameter. It is assumed that a user of this package
// should not use this type directly. However the type will be needed
// if one wants to provide an own ArgPrintFunc.
type ArgWriter []byte

// Write escapes the content so that it can be reliably recognized in a sllm
// message, i.e. replace each single backtick '`' with two backticks.
func (w *ArgWriter) Write(p []byte) (n int, err error) {
	buf := *w
	n = len(buf)
	ep := 0
	for i, b := range p {
		if b == tmplEsc {
			buf = append(buf, p[ep:i+1]...)
			ep = i
		}
	}
	buf = append(buf, p[ep:]...)
	*w = buf
	return len(buf) - n, nil
}

func (w *ArgWriter) WriteString(s string) (n int, err error) {
	return w.Write([]byte(s))
}

func (w *ArgWriter) WriteInt(i int64) (n int, err error) {
	l := len(*w)
	*w = strconv.AppendInt(*w, i, 10)
	return len(*w) - l, nil
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
	s, _ := Sprint("syntax error in `tmpl`:`pos`:`desc`",
		Argv("", err.Tmpl, err.Pos, err.Err))
	return s
}

// ArgPrintFunc is used by the Expand function to process an argument when it
// appears in the expand process of a template. Expand will pass the index idx
// and the name of the argument to expand, i.e. write into the writer wr.
// A ArgPrintFunc returns the number of bytes writen and—just in case—an error.
//
// NOTE: The writer wr of type ArgWriter will escape whatever ArgPrintFunc
// writes so that the template escape symbol '`' remains recognizable.
//
// See also Args and Map
type ArgPrintFunc = func(wr *ArgWriter, idx int, name string) (int, error)

// Fprint uses Expand to print the sllm message to wr.
func Fprint(wr io.Writer, tmpl string, args ArgPrintFunc) (int, error) {
	out, err := Expand(nil, tmpl, args)
	if err != nil {
		return 0, err
	}
	return wr.Write(out)
}

// Sprint uses Expand to return the sllm message as a string.
func Sprint(tmpl string, args ArgPrintFunc) (string, error) {
	out, err := Expand(nil, tmpl, args)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Bprint uses Expand to print the sllm message to an in-memory buffer. Bprint
// is somewhat faster than Print.
func Bprint(wr *bytes.Buffer, tmpl string, args ArgPrintFunc) (int, error) {
	out, err := Expand(wr.Bytes(), tmpl, args)
	if err != nil {
		return 0, err
	}
	wr.Reset()
	return wr.Write(out)
}

// Expand appends a message to buf by expanding all arguments of the given
// template tmpl. The actual process of expanding an argument is left to
// ArgPrintFunc args.
//
// See also Args and Map
func Expand(buf []byte, tmpl string, args ArgPrintFunc) ([]byte, error) {
	errTmpl := SyntaxError{Tmpl: tmpl}
	off, lastIdx := 0, -1
	for len(tmpl) > 0 {
		i := strings.IndexByte(tmpl, tmplEsc)
		if i < 0 {
			buf = append(buf, tmpl...)
			return buf, nil
		}
		if i++; i >= len(tmpl) {
			errTmpl.Pos = off + i
			errTmpl.Err = "unterminated argument"
			return buf, errTmpl
		}
		if tmpl[i] == tmplEsc {
			i++
			buf = append(buf, tmpl[:i]...)
			tmpl = tmpl[i:]
			off += i
			continue
		}
		argName, argIdx, pLen, err := parseArg(tmpl[i:])
		if err != nil {
			errTmpl.Pos = off + i + pLen
			errTmpl.Err = err.Error()
			return buf, errTmpl
		}
		buf = append(buf, tmpl[:i+len(argName)]...)
		buf = append(buf, nmSep)
		if argIdx < 0 {
			lastIdx++
			argIdx = lastIdx
		}
		_, err = args((*ArgWriter)(&buf), argIdx, argName)
		if err != nil {
			return buf, err
		}
		buf = append(buf, tmplEsc)
		tmpl = tmpl[i+pLen:]
	}
	return buf, nil
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
