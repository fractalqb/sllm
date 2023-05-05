package sllm

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"git.fractalqb.de/fractalqb/catch/test"
)

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
		args := test.Err(ParseMap("there is no empty `` arg", &tmpl)).ShouldNot(t)
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
