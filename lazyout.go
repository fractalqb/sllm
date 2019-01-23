package sllm

import (
	"bytes"
	"io"
)

type message struct {
	Tmpl  string
	Undef []byte
}

type argsMsg struct {
	message
	args []interface{}
}

func Args(tmpl string, args ...interface{}) argsMsg {
	return argsMsg{
		message: message{Tmpl: tmpl},
		args:    args,
	}
}

func UArgs(tmpl string, undef []byte, args ...interface{}) argsMsg {
	return argsMsg{
		message: message{Tmpl: tmpl, Undef: undef},
		args:    args,
	}
}

func (msg argsMsg) WriteTo(w io.Writer) (n int64, err error) {
	wn, err := ExpandArgs(w, msg.Tmpl, msg.Undef, msg.args...)
	return int64(wn), err
}

func (msg argsMsg) String() string {
	var buf bytes.Buffer
	ExpandArgs(&buf, msg.Tmpl, msg.Undef, msg.args...) // TODO error
	return buf.String()
}

type mapMsg struct {
	message
	Map ArgMap
}

func Map(tmpl string, args ArgMap) mapMsg {
	return mapMsg{
		message: message{Tmpl: tmpl},
		Map:     args,
	}
}

func UMap(tmpl string, undef []byte, args ArgMap) mapMsg {
	return mapMsg{
		message: message{Tmpl: tmpl, Undef: undef},
		Map:     args,
	}
}

func (msg mapMsg) WriteTo(w io.Writer) (n int64, err error) {
	wn, err := ExpandMap(w, msg.Tmpl, msg.Undef, msg.Map)
	return int64(wn), err
}

func (msg mapMsg) String() string {
	var buf bytes.Buffer
	ExpandMap(&buf, msg.Tmpl, msg.Undef, msg.Map) // TODO error
	return buf.String()
}
