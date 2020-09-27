package sllm

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

func Error(tmpl string, a ...interface{}) error {
	var sb strings.Builder
	_, err := Expand(&sb, tmpl, func(wr io.Writer, idx int, _ string) (int, error) {
		if idx < len(a) {
			return fmt.Fprint(&sb, a[idx])
		}
		return sb.WriteString("<nil>")
	})
	if err != nil {
		fmt.Fprintf(&sb, "[sllm error:%s]", err)
	}
	return errors.New(sb.String())
}
