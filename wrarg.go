package sllm

import (
	"fmt"
	"strconv"
)

type IllegalArgIndex int

func (err IllegalArgIndex) Error() string {
	return strconv.Itoa(int(err))
}

type wrArgv struct {
	undef string
	argv  []any
}

func (ta wrArgv) wr(wr *ArgWriter, idx int, name string) (n int, err error) {
	if idx >= 0 && idx < len(ta.argv) {
		switch arg := ta.argv[idx].(type) {
		case string:
			return wr.WriteString(arg)
		case int:
			return wr.WriteInt(arg)
		case timeFormatter:
			return wr.WriteTime(arg.t, arg.fmt)
		case int64:
			return wr.WriteInt64(int64(arg))
		}
		return fmt.Fprint(wr, ta.argv[idx])
	}
	if ta.undef == "" {
		return 0, IllegalArgIndex(idx)
	}
	return wr.Write([]byte(ta.undef))
}

func Argv(u string, av ...any) ArgPrintFunc {
	return wrArgv{undef: u, argv: av}.wr
}

type UndefinedArg string

func (err UndefinedArg) Error() string {
	return string(err)
}

func Named(u string, m map[string]any) ArgPrintFunc {
	return wrNamed{undef: u, args: m}.wr
}

type wrNamed struct {
	undef string
	args  map[string]any
}

func (m wrNamed) wr(wr *ArgWriter, idx int, name string) (n int, err error) {
	val, ok := m.args[name]
	if ok {
		switch arg := val.(type) {
		case string:
			return wr.WriteString(arg)
		case int:
			return wr.WriteInt(arg)
		case timeFormatter:
			return wr.WriteTime(arg.t, arg.fmt)
		case int64:
			return wr.WriteInt64(int64(arg))
		}
		return fmt.Fprint(wr, val)
	}
	if m.undef == "" {
		return 0, UndefinedArg(name)
	}
	return wr.Write([]byte(m.undef))
}
