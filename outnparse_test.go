package sllm

import (
	"bytes"
	"fmt"
	"testing"
)

func Test_OutAndParse(t *testing.T) {
	undef := "<?>"
	test := func(t *testing.T, tmpl string, args map[string]any) {
		var buf bytes.Buffer
		_, err := Bprint(&buf, tmpl, Named(undef, args))
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
		args map[string]any
	}
	tests := map[string]testCase{
		"no args": {
			tmpl: "foo bar bar",
		},
		"single arg only": {
			tmpl: "`arg`",
			args: map[string]any{"arg": 4711},
		},
		"arg with \\r": {
			tmpl: "with `arg` carriage return",
			args: map[string]any{"arg": "hide \r this"},
		},
		"arg with \\n": {
			tmpl: "with `arg` new line",
			args: map[string]any{"arg": "break \n this"},
		},
		"arg with 0-byte": {
			tmpl: "with `arg` new line",
			args: map[string]any{"arg": "Zero \x00 byte"},
		},
	}
	for n, c := range tests {
		t.Run(n, func(t *testing.T) { test(t, c.tmpl, c.args) })
	}
}
