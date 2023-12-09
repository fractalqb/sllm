package sllm

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Appender interface {
	AppendSllm([]byte) []byte
}

func IdxArgs(args ...any) func([]byte, int, string) ([]byte, error) {
	return func(buf []byte, i int, n string) ([]byte, error) {
		if i < 0 || i >= len(args) {
			return buf, MissingArg{Index: i, Name: n}
		}
		return AppendArg(buf, args[i]), nil
	}
}

func IdxArgsDefault(d any, args ...any) func([]byte, int, string) ([]byte, error) {
	return func(buf []byte, i int, n string) ([]byte, error) {
		if i < 0 || i >= len(args) {
			return AppendArg(buf, d), nil
		}
		return AppendArg(buf, args[i]), nil
	}
}

func NmArgs(args map[string]any) func([]byte, int, string) ([]byte, error) {
	return func(buf []byte, i int, n string) ([]byte, error) {
		if a, ok := args[n]; ok {
			return AppendArg(buf, a), nil
		}
		return buf, MissingArg{Index: i, Name: n}
	}
}

func NmArgsDefault(d any, args map[string]any) func([]byte, int, string) ([]byte, error) {
	return func(buf []byte, i int, n string) ([]byte, error) {
		if a, ok := args[n]; ok {
			return AppendArg(buf, a), nil
		}
		return AppendArg(buf, d), nil
	}
}

func AppendArg(to []byte, v any) []byte {
	switch a := v.(type) {
	case Appender:
		return a.AppendSllm(to)
	case string:
		return EscString(to, a)
	case int:
		return strconv.AppendInt(to, int64(a), 10)
	case int64:
		return strconv.AppendInt(to, a, 10)
	case bool:
		return strconv.AppendBool(to, a)
	case float64:
		return strconv.AppendFloat(to, a, 'f', -1, 64)
	case float32:
		return strconv.AppendFloat(to, float64(a), 'f', -1, 32)
	case uint:
		return strconv.AppendUint(to, uint64(a), 10)
	case uint64:
		return strconv.AppendUint(to, a, 10)
	case fmt.Stringer:
		return EscString(to, a.String())
	default:
		return EscString(to, fmt.Sprint(a))
	}
}

func EscString(to []byte, val string) []byte {
	for tic := strings.IndexByte(val, '`'); tic >= 0; tic = strings.IndexByte(val, '`') {
		tic++
		to = append(to, val[:tic]...)
		to = append(to, '`')
		val = val[tic:]
	}
	return append(to, val...)
}

func EscBytes(to, val []byte) []byte {
	for tic := bytes.IndexByte(val, '`'); tic >= 0; tic = bytes.IndexByte(val, '`') {
		tic++
		to = append(to, val[:tic]...)
		to = append(to, '`')
		val = val[tic:]
	}
	return append(to, val...)
}
