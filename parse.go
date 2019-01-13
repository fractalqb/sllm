package sllm

import (
	"bytes"
	"strings"
)

func Parse(msg string, onArg func(name, value string) error) error {
	for len(msg) > 0 {
		idx := strings.IndexByte(msg, tmplEsc)
		if idx < 0 {
			return nil
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
		msg = msg[idx+1:]
	}
	return nil
}

type DuplicateArg struct {
	Msg  string
	Arg  string
	Vals [2]string
}

func (err DuplicateArg) Error() string {
	var sb bytes.Buffer
	ExpandArgs(&sb,
		"duplcate `arg` with values `old` / `new` in `message`", nil,
		err.Arg, err.Vals[0], err.Vals[1], err.Msg)
	return sb.String()
}

func ParseMap(msg string) (map[string]string, error) {
	res := make(map[string]string)
	err := Parse(msg, func(nm, val string) error {
		if vs, ok := res[nm]; ok {
			return DuplicateArg{msg, nm, [2]string{vs, val}}
		}
		res[nm] = val
		return nil
	})
	return res, err
}
