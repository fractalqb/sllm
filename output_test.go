package sllm

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
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

func writeTestArg(wr io.Writer, idx int, name string) error {
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
