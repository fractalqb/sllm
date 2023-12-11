package sllm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	tmplEscChar byte = '`'
	nameSepChar byte = ':'
	argErrChar  byte = '!'
)

// ArgsFunc appends the escaped argument i with name n to the buffer buff and
// returns the resulting buffer. Even in case of error, buf must be returned. To
// escape an argument implementers should use [EscString] or [EscBytes].
type ArgsFunc func(buf []byte, i int, n string) ([]byte, error)

type ArgError struct {
	Index int
	Name  string
	Err   error
}

func (e ArgError) Error() string {
	return fmt.Sprintf("argument %d '%s': %s", e.Index, e.Name, e.Err.Error())
}

func (e ArgError) Unwrap() error { return e.Err }

type ArgErrors []error

func (e ArgErrors) Is(err error) bool {
	_, ok := err.(ArgErrors)
	return ok
}

func (e ArgErrors) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	}
	var sb strings.Builder
	for _, err := range e {
		fmt.Fprintln(&sb, err.Error())
	}
	return sb.String()
}

func (e ArgErrors) Unwrap() []error { return e }

func Append(to []byte, tmpl string, args ArgsFunc) ([]byte, error) {
	var argErrs ArgErrors
	argErr := func(i int, n string, err error) {
		argErrs = append(argErrs, ArgError{Index: i, Name: n, Err: err})
		to[len(to)-1] = argErrChar
		to = append(to, '(')
		to = append(to, err.Error()...)
		to = append(to, ')')
	}
	argn := 0
	phst := strings.IndexByte(tmpl, tmplEscChar)
	for phst >= 0 {
		phst++
		phnd := strings.IndexByte(tmpl[phst:], tmplEscChar)
		if phnd < 0 {
			to = append(to, tmpl[:phst]...)
			return to, errors.New("unterminated parameter")
		}
		phnd += phst
		n := tmpl[phst:phnd]
		if n == "" {
			to = append(to, tmpl[:phnd]...)
		} else {
			var err error
			if colon := strings.IndexByte(n, nameSepChar); colon >= 0 {
				if colon == 0 {
					return to, fmt.Errorf("empty parameter in '%s'", n)
				}
				idx, err := strconv.Atoi(n[colon+1:])
				if err != nil {
					return to, fmt.Errorf("index in '%s': %w", n, err)
				}
				to = append(to, tmpl[:phnd-len(n)+colon]...)
				to = append(to, nameSepChar)
				if to, err = args(to, idx, n[:colon]); err != nil {
					argErr(argn, n, err)
				}
			} else {
				to = append(to, tmpl[:phnd]...)
				to = append(to, nameSepChar)
				if to, err = args(to, argn, n); err != nil {
					argErr(argn, n, err)
				}
				argn++
			}
		}
		to = append(to, tmplEscChar)
		tmpl = tmpl[phnd+1:]
		phst = strings.IndexByte(tmpl, tmplEscChar)
	}
	to = append(to, tmpl...)
	if len(argErrs) > 0 {
		return to, argErrs
	}
	return to, nil
}
