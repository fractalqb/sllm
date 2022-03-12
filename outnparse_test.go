package sllm

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func Test_OutAndParse(t *testing.T) {
	undef := []byte("<?>")
	test := func(t *testing.T, tmpl string, args ArgMap) {
		var buf strings.Builder
		_, err := Expand(&buf, tmpl, Map(undef, args))
		if err != nil {
			t.Fatal(err)
		}
		var ptmpl bytes.Buffer
		pargs, err := ParseMap(buf.String(), &ptmpl)
		if err != nil {
			t.Fatal(err)
		}
		if ts := ptmpl.String(); ts != tmpl {
			t.Errorf("Template [%s] changed to [%s]", tmpl, ts)
		}
		for k, v := range args {
			pv, ok := pargs[k]
			switch {
			case !ok:
				t.Errorf("Missing arg '%s'", k)
			case len(pv) == 0:
				t.Errorf("Arg '%s' found but has no value", k)
			case len(pv) > 1:
				t.Errorf("Arg '%s' has multiple values: %s", k, pv)
			}
			vs := fmt.Sprint(v)
			if vs != pv[0] {
				t.Errorf("Arg '%s' changed from [%s] to [%s]", k, vs, pv[0])
			}
		}
	}
	type testCase struct {
		tmpl string
		args ArgMap
	}
	tests := map[string]testCase{
		"no args": testCase{
			tmpl: "foo bar bar",
		},
		"single arg only": testCase{
			tmpl: "`arg`",
			args: ArgMap{"arg": 4711},
		},
		"arg with \\r": testCase{
			tmpl: "with `arg` carriage return",
			args: ArgMap{"arg": "hide \r this"},
		},
		"arg with \\n": testCase{
			tmpl: "with `arg` new line",
			args: ArgMap{"arg": "break \n this"},
		},
		"arg with 0-byte": testCase{
			tmpl: "with `arg` new line",
			args: ArgMap{"arg": "Zero \x00 byte"},
		},
	}
	for n, c := range tests {
		t.Run(n, func(t *testing.T) { test(t, c.tmpl, c.args) })
	}
}
