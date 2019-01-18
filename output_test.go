package sllm

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
)

func Example_valEsc_Write() {
	ew := ValueEsc{os.Stdout}
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

func writeTestArg(wr ValueEsc, idx int, name string) error {
	_, err := fmt.Fprintf(wr, "#%d/'%s'", idx, name)
	return err
}

func ExampleExpand() {
	Expand(os.Stdout, "want `arg1` here and `arg2` here", writeTestArg)
	fmt.Println()
	Expand(os.Stdout, "template with backtick '``' and an `arg` here", writeTestArg)
	fmt.Println()
	ExpandArgs(os.Stdout, "touching args: `one``two`", []byte("–"), 4711, true)
	fmt.Println()
	// Output:
	// want `arg1:#0/'arg1'` here and `arg2:#1/'arg2'` here
	// template with backtick '``' and an `arg:#0/'arg'` here
	// touching args: `one:4711``two:true`
}

func TestExpand_unterminated1(t *testing.T) {
	var out bytes.Buffer
	err := Expand(&out, "foo `bar without end", nil)
	if err == nil {
		t.Fatal("expected Expand error, got none")
	}
	if se, ok := err.(SyntaxError); !ok {
		t.Fatal("received wrong error type:", reflect.TypeOf(err).Name())
	} else if se.Err != "unterminated argument" {
		t.Fatal("received wrong error:", err)
	}
}

func TestExpand_unterminated2(t *testing.T) {
	var out bytes.Buffer
	err := Expand(&out, "without end `", nil)
	if err == nil {
		t.Fatal("expected Expand error, got none")
	}
	if se, ok := err.(SyntaxError); !ok {
		t.Fatal("received wrong error type:", reflect.TypeOf(err).Name())
	} else if se.Err != "unterminated argument" {
		t.Fatal("received wrong error:", err)
	}
}

func TestExpand_name_with_colon(t *testing.T) {
	var out bytes.Buffer
	err := Expand(&out, "foo `ba:r` baz", nil)
	if err == nil {
		t.Fatal("expected Expand error, got none")
	}
	if se, ok := err.(SyntaxError); !ok {
		t.Fatal("received wrong error type:", reflect.TypeOf(err).Name())
	} else if se.Err != "name contains ':'" {
		t.Fatal("received wrong error:", err)
	}
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
		ExpandArgs(&out,
			"just an `what`: `miss`", nil,
			"example",
			"<undef>")
	}
}

// BenchmarkExpandMap shall give an indication of the voverhad for map creation
// compared to the ExpandArgs function.
func BenchmarkExpandMap(b *testing.B) {
	var out bytes.Buffer
	for i := 0; i < b.N; i++ {
		out.Reset()
		ExpandMap(&out,
			"just an `what`: `miss`", nil,
			map[string]interface{}{
				"what": "example",
				"miss": "<undef>",
			})
	}
}

func ExampleForDocGo() {
	const (
		count = 7
		item  = "Hat"
		user  = "John Doe"
	)
	logr := log.New(os.Stdout, "", log.LstdFlags)
	logr.Printf("added %d ⨉ %s to shopping cart by %s", count, item, user)
	logr.Print(Map("added `count` ⨉ `item` to shopping cart by `user`",
		ArgMap{
			"count": count,
			"item":  item,
			"user":  user,
		}))
}
