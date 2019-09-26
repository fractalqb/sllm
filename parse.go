package sllm

import (
	"bytes"
	"io/ioutil"
	"strings"
)

// Parse parses a sllm message create by Expand and calls onArg for every
// `name:value` parameter it finds in the message. When a non-nil buffer is
// passed as tmpl Parse will also reconstruct the original template into the
// buffer. Note that the template is appended to tmpl's content.
func Parse(msg string, tmpl *bytes.Buffer, onArg func(name, value string) error) error {
	for len(msg) > 0 {
		idx := strings.IndexByte(msg, tmplEsc)
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
		idx = strings.IndexByte(msg, nmSep)
		if idx < 0 {
			return nil
		}
		name := msg[:idx]
		msg = msg[idx+1:]
		idx = strings.IndexByte(msg, tmplEsc)
		for {
			if idx < 0 {
				return nil
			}
			if idx+1 >= len(msg) || msg[idx+1] != tmplEsc {
				break
			}
			nidx := strings.IndexByte(msg[idx+2:], tmplEsc)
			if nidx < 0 {
				return nil
			}
			idx += nidx + 2
		}
		value := msg[:idx]
		err := onArg(name, value)
		if err != nil {
			return err
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

// ParseMap uses Parse to create a map with all parameters assigned to an
// argument in the passed message msg. ParseMap can also reconstruct the
// template when passing a Buffer to tmpl.
func ParseMap(msg string, tmpl *bytes.Buffer) map[string][]string {
	res := make(map[string][]string)
	Parse(msg, tmpl, func(nm, val string) error {
		vls := res[nm]
		res[nm] = append(vls, val)
		return nil
	})
	return res
}

// ExtractParams extracs the parameter names from template tmpl and appends them
// to appendTo.
func ExtractParams(appendTo []string, tmpl string) ([]string, error) {
	_, err := Expand(ioutil.Discard, tmpl,
		func(wr ValueEsc, idx int, name string) (int, error) {
			appendTo = append(appendTo, name)
			return len(name), nil
		},
	)
	return appendTo, err
}
