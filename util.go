package sllm

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

func FprintIdx(w io.Writer, tmpl string, args ...any) (int, error) {
	return Fprint(w, tmpl, IdxArgs(args...))
}

func Fprint(w io.Writer, tmpl string, args ArgsFunc) (int, error) {
	if buf, ok := w.(*bytes.Buffer); ok {
		tmp, err := Append(buf.Bytes(), tmpl, args)
		if err != nil {
			return 0, err
		}
		buf.Reset()
		return buf.Write(tmp)
	}
	tmp, err := Append(nil, tmpl, args)
	if err != nil {
		return 0, err
	}
	return w.Write(tmp)
}

func StringIdx(tmpl string, args ...any) (string, error) {
	return String(tmpl, IdxArgs(args...))
}

func String(tmpl string, args ArgsFunc) (string, error) {
	buf, err := Append(nil, tmpl, args)
	return string(buf), err
}

func ErrorIdx(tmpl string, args ...any) error {
	return ErrorIdx(tmpl, IdxArgs(args...))
}

func Error(tmpl string, args ArgsFunc) error {
	s, err := String(tmpl, args)
	if err != nil {
		return fmt.Errorf("appending '%s': %w", s, err)
	}
	return errors.New(s)
}
