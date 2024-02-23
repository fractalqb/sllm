package sllm

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParameters(t *testing.T) {
	ptest := func(t *testing.T, msg string, expect ...string) {
		ps, err := Parameters(msg, nil)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(ps, expect) {
			t.Fatal("wrong params:", ps)
		}
	}
	t.Run("single param", func(t *testing.T) {
		ptest(t, "`foo`", "foo")
	})
	t.Run("one param mid", func(t *testing.T) {
		ptest(t, "foo `bar`baz", "bar")
	})
	t.Run("two params", func(t *testing.T) {
		ptest(t, "this is `foo` and `bar`", "foo", "bar")
	})
}

func ExampleParseMap() {
	var tmpl bytes.Buffer
	args, _ := ParseMap(
		"added `count:7` ⨉ `item:Hat` to shopping cart by `user:John Doe`",
		&tmpl,
	)
	fmt.Println(tmpl.String())
	for k, v := range args {
		fmt.Printf("%s:[%s]\n", k, v)
	}
	// Unordered output:
	// added `count` ⨉ `item` to shopping cart by `user`
	// count:[[7]]
	// item:[[Hat]]
	// user:[[John Doe]]
}

func TestParse(t *testing.T) {
	t.Run("empty arg", func(t *testing.T) {
		var tmpl bytes.Buffer
		args, err := ParseMap("there is no empty `` arg", &tmpl)
		if err != nil {
			t.Fatal(err)
		}
		if len(args) != 0 {
			t.Errorf("found args: %v", args)
		}
	})
	t.Run("start arg at end", func(t *testing.T) {
		var tmpl bytes.Buffer
		_, err := ParseMap("there is no arg `", &tmpl)
		switch {
		case err == nil:
			t.Error("error not detected")
		case err.Error() != "empty arg":
			t.Errorf("unexpected error: %s", err)
		}
	})
	t.Run("unterminated arg name", func(t *testing.T) {
		var tmpl bytes.Buffer
		_, err := ParseMap("there is no `arg`", &tmpl)
		switch {
		case err == nil:
			t.Error("error not detected")
		case !strings.HasPrefix(err.Error(), "unterminated arg name '"):
			t.Errorf("unexpected error: %s", err)
		}
	})
	t.Run("no arg-error start marker", func(t *testing.T) {
		var tmpl bytes.Buffer
		_, err := ParseMap("there is no `arg!", &tmpl)
		switch {
		case err == nil:
			t.Error("error not detected")
		case !strings.HasPrefix(err.Error(), "no error marker for arg '"):
			t.Errorf("unexpected error: %s", err)
		}
	})
	t.Run("wrong arg-error start marker", func(t *testing.T) {
		var tmpl bytes.Buffer
		_, err := ParseMap("there is no `arg!<bla>`", &tmpl)
		switch {
		case err == nil:
			t.Error("error not detected")
		case !strings.HasPrefix(err.Error(), "invalid error start marker '"):
			t.Errorf("unexpected error: %s", err)
		}
	})
	t.Run("wrong arg-error end marker", func(t *testing.T) {
		var tmpl bytes.Buffer
		_, err := ParseMap("there is no `arg!(bla>`", &tmpl)
		switch {
		case err == nil:
			t.Error("error not detected")
		case !strings.HasPrefix(err.Error(), "invalid error end marker '"):
			t.Errorf("unexpected error: %s", err)
		}
	})
	t.Run("unterminated arg", func(t *testing.T) {
		var tmpl bytes.Buffer
		_, err := ParseMap("there is no `arg:4711", &tmpl)
		switch {
		case err == nil:
			t.Error("error not detected")
		case !strings.HasPrefix(err.Error(), "unterminated arg '"):
			t.Errorf("unexpected error: %s", err)
		}
	})
}

func Test_OutAndParse(t *testing.T) {
	undef := "<?>"
	test := func(t *testing.T, tmpl string, args map[string]any) {
		var buf bytes.Buffer
		_, err := Fprint(&buf, tmpl, NmArgsDefault(undef, args))
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
