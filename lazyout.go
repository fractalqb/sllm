package sllm

import (
	"bytes"
)

type message struct {
	Tmpl  string
	Undef []byte
}

type argsMsg struct {
	message
	Args []interface{}
}

func Args(tmpl string, args ...interface{}) argsMsg {
	return argsMsg{
		message: message{Tmpl: tmpl},
		Args:    args,
	}
}

func UArgs(tmpl string, undef []byte, args ...interface{}) argsMsg {
	return argsMsg{
		message: message{Tmpl: tmpl, Undef: undef},
		Args:    args,
	}
}

func (msg argsMsg) String() string {
	var buf bytes.Buffer
	ExpandArgs(&buf, msg.Tmpl, msg.Undef, msg.Args...) // TODO error
	return buf.String()
}
