package sllm

import (
	"fmt"
	"io"
	"strconv"
)

type IllegalArgIndex int

func (err IllegalArgIndex) Error() string {
	return strconv.Itoa(int(err))
}

type wrArgs struct {
	undef []byte
	argv  []interface{}
}

func (ta wrArgs) wr(wr io.Writer, idx int, name string) (n int, err error) {
	if idx >= 0 && idx < len(ta.argv) {
		return fmt.Fprint(wr, ta.argv[idx])
	}
	if ta.undef == nil {
		return 0, IllegalArgIndex(idx)
	}
	return wr.Write(ta.undef)
}

func Args(u []byte, av ...interface{}) ParamWriter {
	return wrArgs{undef: u, argv: av}.wr
}

type UndefinedArg string

func (err UndefinedArg) Error() string {
	return string(err)
}

type ArgMap = map[string]interface{}

func Map(u []byte, m ArgMap) ParamWriter {
	return wrMap{undef: u, args: m}.wr
}

type wrMap struct {
	undef []byte
	args  map[string]interface{}
}

func (m wrMap) wr(wr io.Writer, idx int, name string) (n int, err error) {
	val, ok := m.args[name]
	if ok {
		return fmt.Fprint(wr, val)
	}
	if m.undef == nil {
		return 0, UndefinedArg(name)
	}
	return wr.Write(m.undef)
}
