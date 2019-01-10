package sllm

import (
	"bytes"
)

func Parse(msg string, onArg func(name, value string) error) error {
	// TODO
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
