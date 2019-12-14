package sllm

import (
	"bytes"
	"io"
)

const (
	tmplEsc = '`'
	nmSep   = ':'
)

var tEsc2 = []byte{tmplEsc, tmplEsc}
var tEsc1 = tEsc2[:1]
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
	var (
		b1 [1]byte
		tn int
		xi bool
	)
	bs := b1[:]
	idx, i := 0, 0
	for i < len(tmpl) {
		if b := tmpl[i]; b == tmplEsc {
			i++
			if i >= len(tmpl) {
				return n, SyntaxError{tmpl, len(tmpl), "unterminated argument"}
			}
			if tmpl[i] == tmplEsc {
				if tn, err = wr.Write(tEsc2); err != nil {
					return n + tn, err
				}
				i++
			} else if i, tn, xi, err = xpandNm(wr, idx, tmpl, i, writeArg); err != nil {
				return n + tn, err
			} else if !xi {
				idx++
			}
			n += tn
		} else {
			bs[0] = b
			tn, err = wr.Write(bs)
			n += tn
			if err != nil {
				return n, err
			}
			i++
		}
	}
	return n, nil
}

func argIdx(tmpl string, ip int) (idx, res int, err error) {
	start := ip
	for ip < len(tmpl) {
		b := tmpl[ip]
		if b == tmplEsc {
			if ip == start {
				return -1, -1, SyntaxError{tmpl, ip, "empty explicit index"}
			}
			return ip + 1, res, nil
		}
		if b < '0' || b > '9' {
			return -1, -1, SyntaxError{tmpl, ip, "not a digit in explicit arg index"}
		}
		res = 10*res + int(b-'0')
		ip++
	}
	return -1, -1, SyntaxError{tmpl, len(tmpl), "unterminated argument"}
}

func xpandNm(wr io.Writer, idx int, tmpl string, np int, writeArg ParamWriter) (
	rp, wn int,
	xidx bool,
	err error,
) {
	var tn int
	if wn, err = wr.Write(tEsc1); err != nil {
		return -1, wn, false, err
	}
	for i := np; i < len(tmpl); i++ {
		b := tmpl[i]
		switch b {
		case tmplEsc:
			name := tmpl[np:i]
			if tn, err = io.WriteString(wr, name); err != nil {
				return -1, wn + tn, false, err
			}
			wn += tn
			if tn, err = wr.Write(nmSepStr); err != nil {
				return -1, wn + tn, false, err
			}
			wn += tn
			if tn, err = writeArg(ValueEsc{wr}, idx, name); err != nil {
				return -1, wn + tn, false, err
			}
			wn += tn
			if tn, err = wr.Write(tEsc1); err != nil {
				return i + 1, wn + tn, false, err
			}
			wn += tn
			return i + 1, wn, false, nil
		case nmSep:
			name := tmpl[np:i]
			if i, idx, err = argIdx(tmpl, i+1); err != nil {
				return -1, wn, true, err
			}
			if tn, err = io.WriteString(wr, name); err != nil {
				return -1, wn + tn, true, err
			}
			wn += tn
			if tn, err = wr.Write(nmSepStr); err != nil {
				return -1, wn + tn, true, err
			}
			wn += tn
			if tn, err = writeArg(ValueEsc{wr}, idx, name); err != nil {
				return -1, wn + tn, true, err
			}
			wn += tn
			if tn, err = wr.Write(tEsc1); err != nil {
				return -1, wn + tn, true, err
			}
			wn += tn
			return i, wn, true, nil
		}
	}
	return -1, wn, false, SyntaxError{tmpl, len(tmpl), "unterminated argument"}
}

// Expands uses Expand to return the expanded temaplate as a string.
func Expands(tmpl string, writeArg ParamWriter) (string, error) {
	var buf bytes.Buffer
	_, err := Expand(&buf, tmpl, writeArg)
	return buf.String(), err
}
