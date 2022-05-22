package sllm

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func Example() {
	// Positional arguments
	Fprint(os.Stdout, "added `count` ⨉ `item` to shopping cart by `user`\n",
		Argv("", 7, "Hat", "John Doe"),
	)
	// Named arguments
	Fprint(os.Stdout, "added `count` ⨉ `item` to shopping cart by `user`\n",
		Named("", map[string]any{
			"count": 7,
			"item":  "Hat",
			"user":  "John Doe",
		}),
	)
	// Output:
	// added `count:7` ⨉ `item:Hat` to shopping cart by `user:John Doe`
	// added `count:7` ⨉ `item:Hat` to shopping cart by `user:John Doe`
}

func ExampleArgWriter() {
	var ew ArgWriter
	ew.Write([]byte("foo"))
	fmt.Println(string(ew))
	ew = ew[:0]
	ew.Write([]byte("`bar`"))
	fmt.Println(string(ew))
	ew = ew[:0]
	ew.Write([]byte("b`az"))
	fmt.Println(string(ew))
	// Output:
	// foo
	// ``bar``
	// b``az
}

func BenchmarkArgWriter(b *testing.B) {
	txt := []byte("this `is a funny ``text with `backtik")
	var pw ArgWriter
	for i := 0; i < b.N; i++ {
		pw = pw[:0]
		pw.Write(txt)
	}
}

func ExamplePrint() {
	var writeTestArg = func(wr *ArgWriter, idx int, name string) (int, error) {
		return fmt.Fprintf(wr, "#%d/'%s'", idx, name)
	}
	Fprint(os.Stdout, "want `arg1` here and `arg2` here", writeTestArg)
	fmt.Println()
	Fprint(os.Stdout, "template with backtick '``' and an `arg` here", writeTestArg)
	fmt.Println()
	Fprint(os.Stdout, "touching args: `one``two``three`", Argv("–", 4711, true))
	fmt.Println()
	Fprint(os.Stdout, "explicit `index:0` and `same:0`", Argv("", 4711))
	fmt.Println()
	// Output:
	// want `arg1:#0/'arg1'` here and `arg2:#1/'arg2'` here
	// template with backtick '``' and an `arg:#0/'arg'` here
	// touching args: `one:4711``two:true``three:–`
	// explicit `index:4711` and `same:4711`
}

func ExamplePrint_explicitIndex() {
	var writeTestArg = func(wr *ArgWriter, idx int, name string) (int, error) {
		return fmt.Fprint(wr, idx)
	}
	Fprint(os.Stdout, "`a`, `b:11`, `c`, `d:0`, `e`", writeTestArg)
	// Output:
	// `a:0`, `b:11`, `c:1`, `d:0`, `e:2`
}

func TestBprint_syntaxerror(t *testing.T) {
	test := func(t *testing.T, tmpl string, epos int, emsg string) {
		var out bytes.Buffer
		_, err := Bprint(&out, tmpl, nil)
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
		fmt.Fprintf(&out,
			"added `count:%d` x `item:%s` to shopping cart by `user:%s`",
			7,
			"`hat`",
			"John Doe")
	}
}

func BenchmarkExpandArgs(b *testing.B) {
	var out bytes.Buffer
	for i := 0; i < b.N; i++ {
		out.Reset()
		Fprint(&out,
			"added `count` x `item` to shopping cart by `user`",
			Argv("", 7, "`hat`", "John Doe"),
		)
	}
}

// BenchmarkExpandMap shall give an indication of the voverhad for map creation
// compared to the ExpandArgs function.
func BenchmarkExpandMap(b *testing.B) {
	var out bytes.Buffer
	for i := 0; i < b.N; i++ {
		out.Reset()
		Bprint(&out,
			"added `count` x `item` to shopping cart by `user`",
			Named("", map[string]any{
				"count": 7,
				"item":  "`hat`",
				"user":  "John Doe",
			}),
		)
	}
}
