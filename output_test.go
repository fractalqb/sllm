package sllm

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func Example_valEsc_Write() {
	ew := valEsc{os.Stdout}
	ew.Write([]byte("foo"))
	fmt.Fprintln(os.Stdout)
	ew.Write([]byte("`bar`"))
	fmt.Fprintln(os.Stdout)
	ew.Write([]byte("b`az"))
	fmt.Fprintln(os.Stdout)
	// Output:
	// foo
	// ``bar``
	// b``az
}

func ExampleExpand() {
	var writeTestArg = func(wr io.Writer, idx int, name string) (int, error) {
		return fmt.Fprintf(wr, "#%d/'%s'", idx, name)
	}
	Expand(os.Stdout, "want `arg1` here and `arg2` here", writeTestArg)
	fmt.Println()
	Expand(os.Stdout, "template with backtick '``' and an `arg` here", writeTestArg)
	fmt.Println()
	Expand(os.Stdout, "touching args: `one``two``three`", Args([]byte("–"), 4711, true))
	fmt.Println()
	Expand(os.Stdout, "explicit `index:0` and `same:0`", Args(nil, 4711))
	fmt.Println()
	// Output:
	// want `arg1:#0/'arg1'` here and `arg2:#1/'arg2'` here
	// template with backtick '``' and an `arg:#0/'arg'` here
	// touching args: `one:4711``two:true``three:–`
	// explicit `index:4711` and `same:4711`
}

func ExampleExpand_explicitIndex() {
	var writeTestArg = func(wr io.Writer, idx int, name string) (int, error) {
		return fmt.Fprint(wr, idx)
	}
	Expand(os.Stdout, "`a`, `b:11`, `c`, `d:0`, `e`", writeTestArg)
	// Output:
	// `a:0`, `b:11`, `c:1`, `d:0`, `e:2`
}

func TestExpand_syntaxerror(t *testing.T) {
	test := func(t *testing.T, tmpl string, epos int, emsg string) {
		var out bytes.Buffer
		_, err := Expand(&out, tmpl, nil)
		if err == nil {
			t.Fatal("expected Expand error, got none")
		}
		if se, ok := err.(SyntaxError); !ok {
			t.Error("received wrong error type:", reflect.TypeOf(err).Name())
		} else {
			errArgs, err := ParseMap(se.Error(), nil)
			if err != nil {
				t.Fatal(err)
			}
			if pos := errArgs["pos"][0]; pos != strconv.Itoa(epos) {
				t.Errorf("wrong error position %s, expected %d", pos, epos)
			}
			if errArgs["desc"][0] != emsg {
				t.Fatal("received wrong error:", err)
			}
		}
	}
	t.Run("unterminated mid", func(t *testing.T) {
		test(t, "foo `bar without end", 5, "unterminated argument")
	})
	t.Run("unterminated end", func(t *testing.T) {
		test(t, "without end `", 13, "unterminated argument")
	})
	t.Run("empty explicit index", func(t *testing.T) {
		test(t, "foo `ba:` baz", 8, "empty explicit index")
	})
	t.Run("non-numeric explicit index", func(t *testing.T) {
		test(t, "foo `ba:1x2` baz", 9, "not a digit in explicit arg index")
	})
}

func TestExtractParams(t *testing.T) {
	ptest := func(t *testing.T, msg string, expect ...string) {
		ps, err := ExtractParams(nil, msg)
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

func BenchmarkPrintf(b *testing.B) {
	var out bytes.Buffer
	for i := 0; i < b.N; i++ {
		out.Reset()
		fmt.Fprintf(&out, "just an `what:%s`: `miss:%s`",
			"example",
			"<undef>")
	}
}

func BenchmarkExpandArgs(b *testing.B) {
	var out bytes.Buffer
	for i := 0; i < b.N; i++ {
		out.Reset()
		Expand(&out,
			"just an `what`: `miss`",
			Args(nil, "example", "<undef>"))
	}
}

// BenchmarkExpandMap shall give an indication of the voverhad for map creation
// compared to the ExpandArgs function.
func BenchmarkExpandMap(b *testing.B) {
	var out bytes.Buffer
	for i := 0; i < b.N; i++ {
		out.Reset()
		Expand(&out,
			"just an `what`: `miss`",
			Map(nil, map[string]interface{}{
				"what": "example",
				"miss": "<undef>",
			}))
	}
}

// func Example_forDocGo() {
// 	const (
// 		count = 7
// 		item  = "Hat"
// 		user  = "John Doe"
// 	)
// 	logr := log.New(os.Stdout, "", log.LstdFlags)
// 	logr.Printf("added %d ⨉ %s to shopping cart by %s", count, item, user)
// 	logr.Print(Map("added `count` ⨉ `item` to shopping cart by `user`",
// 		ArgMap{
// 			"count": count,
// 			"item":  item,
// 			"user":  user,
// 		}))
// }
