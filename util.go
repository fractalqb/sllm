package sllm

import (
	"errors"
	"fmt"
)

func Error(tmpl string, args ...any) error {
	s, err := Sprint(tmpl, Argv("<?>", args...))
	if err != nil {
		s = fmt.Sprintf("[sllm:%s]", err)
	}
	return errors.New(s)
}
