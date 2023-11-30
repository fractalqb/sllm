package sllm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	tmplEsc byte = '`'
	nameSep byte = ':'
	argErr  byte = '!'
)

// ArgsFunc appends the escaped argument i with name n to the buffer buff and
// returns the resulting buffer. Even in case of error, buf ist returned. To
// escape an argument implementers should use [EscString] or [EscBytes].
type ArgsFunc func(buf []byte, i int, n string) ([]byte, error)

type MissingArg struct {
	Index int
	Name  string
}

func (e MissingArg) Is(err error) bool {
	_, ok := err.(MissingArg)
	return ok
}

func (e MissingArg) Error() string {
	return fmt.Sprintf("missing argument %d'%s'", e.Index, e.Name)
}

func Append(to []byte, tmpl string, args ArgsFunc) ([]byte, error) {
	var err error
	argn := 0
	phst := strings.IndexByte(tmpl, tmplEsc)
	for phst >= 0 {
		phst++
		phnd := strings.IndexByte(tmpl[phst:], tmplEsc)
		if phnd < 0 {
			to = append(to, tmpl[:phst]...)
			return to, errors.New("unterminated parameter")
		}
		phnd += phst
		n := tmpl[phst:phnd]
		if n == "" {
			to = append(to, tmpl[:phnd]...)
		} else {
			if colon := strings.IndexByte(n, nameSep); colon >= 0 {
				if colon == 0 {
					return to, fmt.Errorf("empty parameter in '%s'", n)
				}
				idx, err := strconv.Atoi(n[colon+1:])
				if err != nil {
					return to, fmt.Errorf("index in '%s': %w", n, err)
				}
				to = append(to, tmpl[:phnd-len(n)+colon]...)
				to = append(to, nameSep)
				to, err = args(to, idx, n[:colon])
			} else {
				to = append(to, tmpl[:phnd]...)
				to = append(to, nameSep)
				to, err = args(to, argn, n)
				argn++
			}
			if err != nil {
				if !errors.Is(err, MissingArg{}) {
					return to, err
				}
				to[len(to)-1] = argErr
				to = append(to, '(')
				to = append(to, err.Error()...)
				to = append(to, ')')
			}
		}
		to = append(to, tmplEsc)
		tmpl = tmpl[phnd+1:]
		phst = strings.IndexByte(tmpl, tmplEsc)
	}
	to = append(to, tmpl...)
	return to, nil
}
