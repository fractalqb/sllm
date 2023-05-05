package sllm

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"sync"
)

const (
	tmplEsc byte = '`'
	nameSep byte = ':'
	argErr  byte = '!'
)

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
	out := byteSlicePool.Get().([]byte)
	out, err := Expand(out[:0], tmpl, args)
	if err != nil {
		byteSlicePool.Put(out)
		return 0, err
	}
	n, err := wr.Write(out)
	byteSlicePool.Put(out)
	return n, err
}

// Sprint uses Expand to return the sllm message as a string.
func Sprint(tmpl string, args ArgPrintFunc) (string, error) {
	out := byteSlicePool.Get().([]byte)
	out, err := Expand(out[:0], tmpl, args)
	if err != nil {
		byteSlicePool.Put(out)
		return "", err
	}
	tmp := string(out)
	byteSlicePool.Put(out)
	return tmp, nil
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

var byteSlicePool = sync.Pool{New: func() any { return make([]byte, 128) }}

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
		nmSepIdx := len(buf)
		buf = append(buf, nameSep)
		if argIdx < 0 {
			lastIdx++
			argIdx = lastIdx
		}
		_, err = args((*ArgWriter)(&buf), argIdx, argName)
		if err != nil {
			buf = buf[:nmSepIdx+1]
			buf[nmSepIdx] = argErr
			buf = append(buf, '(')
			(*ArgWriter)(&buf).WriteString(err.Error())
			buf = append(buf, ')')
		}
		buf = append(buf, tmplEsc)
		tmpl = tmpl[i+pLen:]
	}
	return buf, nil
}

func parseArg(tmpl string) (argName string, argIdx, parseLen int, err error) {
	parseLen = strings.IndexByte(tmpl, tmplEsc)
	if parseLen < 0 {
		return "", -1, 0, errors.New("unterminated argument")
	}
	nms := strings.IndexByte(tmpl[:parseLen], nameSep)
	switch {
	case nms < 0:
		argName = tmpl[:parseLen]
		argIdx = -1
	case nms > 0:
		argName = tmpl[:nms]
		nms++
		if nms == parseLen {
			return argName, -1, parseLen, errors.New("empty explicit index")
		}
		argIdx = 0
		for nms < parseLen {
			b := tmpl[nms]
			if b < '0' || b > '9' {
				return argName, -1, nms, errors.New("not a digit in explicit arg index")
			}
			argIdx = 10*argIdx + int(b) - '0'
			nms++
		}
	default:
		return "", -1, 0, errors.New("empty argument name")
	}
	return argName, argIdx, parseLen + 1, nil
}
