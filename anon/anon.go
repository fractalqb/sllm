package anon

import (
	"git.fractalqb.de/fractalqb/sllm"
)

func Replace(s string) sllm.ParamWriter {
	return func(wr sllm.ValueEsc, idx int, name string) (int, error) {
		return wr.Write([]byte(s))
	}
}

type ByName struct {
	Clear sllm.ParamWriter
	Anon  map[string]sllm.ParamWriter
}

func (a ByName) Param(wr sllm.ValueEsc, idx int, name string) (int, error) {
	if pw, ok := a.Anon[name]; ok {
		return pw(wr, idx, name)
	}
	return a.Clear(wr, idx, name)
}
