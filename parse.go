package sllm

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// Parameters extracs the parameter names from template tmpl and appends them
// to a.
func Parameters(tmpl string, a []string) ([]string, error) {
	_, err := Append(nil, tmpl, func(to []byte, idx int, n string) ([]byte, error) {
		a = append(a, n)
		return to[:0], nil
	})
	return a, err
}

// Parse parses a sllm message create by Append and calls onArg for every
// `name:value` parameter it finds in the message. When a non-nil buffer is
// passed as tmpl Parse will also reconstruct the original template into the
// buffer. Note that the template is appended to tmpl's content.
func Parse(msg string, tmpl *bytes.Buffer, onArg func(name, value string, argError bool) error) error {
	for len(msg) > 0 {
		idx := strings.IndexByte(msg, tmplEscChar)
		if idx < 0 {
			if tmpl != nil {
				tmpl.WriteString(msg)
			}
			return nil
		}
		if tmpl != nil {
			tmpl.WriteString(msg[:idx])
		}
		msg = msg[idx+1:]
		switch {
		case msg == "":
			return errors.New("empty arg")
		case msg[0] == tmplEscChar:
			msg = msg[1:]
			continue
		}
		idx = strings.IndexAny(msg, nameEnd)
		if idx < 0 {
			return fmt.Errorf("unterminated arg name '%s'", msg)
		}
		name := msg[:idx]
		isErr := msg[idx] == '!'
		if isErr {
			switch {
			case idx+1 >= len(msg):
				return fmt.Errorf("no error marker for arg '%s'", name)
			case msg[idx+1] != '(':
				return fmt.Errorf("invalid error start marker '%c'", msg[idx+1])
			}
			msg = msg[idx+2:]
		} else {
			msg = msg[idx+1:]
		}
		idx = strings.IndexByte(msg, tmplEscChar)
		for {
			if idx < 0 {
				return fmt.Errorf("unterminated arg '%s'", name)
			}
			if idx+1 >= len(msg) || msg[idx+1] != tmplEscChar {
				break
			}
			nidx := strings.IndexByte(msg[idx+2:], tmplEscChar)
			if nidx < 0 {
				return nil
			}
			idx += nidx + 2
		}
		var value string
		if isErr {
			if r := msg[idx-1]; r != ')' {
				return fmt.Errorf("invalid error end marker '%c'", r)
			}
			value = msg[:idx-1]
		} else {
			value = msg[:idx]
		}
		err := onArg(name, value, isErr)
		if err != nil {
			if isErr {
				return fmt.Errorf("error arg '%s': %w", name, err)
			}
			return fmt.Errorf("arg '%s': %w", name, err)
		}
		if tmpl != nil {
			tmpl.WriteRune('`')
			tmpl.WriteString(name)
			tmpl.WriteRune('`')
		}
		msg = msg[idx+1:]
	}
	return nil
}

var nameEnd = string([]byte{nameSepChar, argErrChar})

// ParseMap uses Parse to create a map with all parameters assigned to an
// argument in the passed message msg. ParseMap can also reconstruct the
// template when passing a Buffer to tmpl.
func ParseMap(msg string, tmpl *bytes.Buffer) (map[string][]any, error) {
	res := make(map[string][]any)
	err := Parse(msg, tmpl, func(nm, val string, isErr bool) error {
		vls := res[nm]
		if isErr {
			res[nm] = append(vls, errors.New(val))
		} else {
			res[nm] = append(vls, val)
		}
		return nil
	})
	return res, err
}
