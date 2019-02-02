package sllm

import (
	"fmt"
	"strconv"
)

type IllegalArgIndex int

func (err IllegalArgIndex) Error() string {
	return strconv.Itoa(int(err))
}

func Args(undef []byte, argv ...interface{}) ParamWriter {
	return func(wr ValueEsc, idx int, name string) (n int, err error) {
		if idx < 0 || idx >= len(argv) {
			if undef == nil {
				return 0, IllegalArgIndex(idx)
			} else {
				n, err = wr.Write(undef)
			}
		} else {
			n, err = writeVal(wr, argv[idx])
		}
		return n, err
	}
}

type UndefinedArg string

func (err UndefinedArg) Error() string {
	return string(err)
}

type ArgMap = map[string]interface{}

func Map(undef []byte, m ArgMap) ParamWriter {
	return func(wr ValueEsc, idx int, name string) (n int, err error) {
		if val, ok := m[name]; !ok {
			if undef == nil {
				return 0, UndefinedArg(name)
			} else {
				n, err = wr.Write(undef)
			}
		} else {
			n, err = writeVal(wr, val)
		}
		return n, err
	}
}

var writeVal = fmt.Fprint
