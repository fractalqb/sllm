package sllm

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func ExampleArgs() {
	Fprint(os.Stdout, "just an `what`", Argv("", "example"))
	// Output:
	// just an `what:example`
}

func ExampleArgs_undef() {
	Fprint(os.Stdout, "just an `what`: `miss`", Argv("<undef>", "example"))
	// Output:
	// just an `what:example`: `miss:<undef>`
}

func ExampleMap() {
	Fprint(os.Stdout, "just an `what`", Named("", map[string]any{
		"what":  "example",
		"dummy": false,
	}))
	// Output:
	// just an `what:example`
}

func ExampleMap_undef() {
	Fprint(os.Stdout, "just an `what`: `miss`", Named("<undef>", map[string]any{
		"what":  "example",
		"dummy": false,
	}))
	// Output:
	// just an `what:example`: `miss:<undef>`
}

func fuzzArgs[T any](t *testing.T, arg T) {
	const tmpl = "Just `a` template"
	var sb bytes.Buffer
	n, err := Bprint(&sb, tmpl, Argv("", arg))
	if err != nil {
		t.Error(err)
	}
	if n != sb.Len() {
		t.Errorf("Output len %d; written %d byte", sb.Len(), n)
	}
	msg := sb.String()
	if !strings.HasPrefix(msg, "Just `a:") {
		t.Errorf("Message does not start properly")
	}
	if !strings.HasSuffix(msg, "` template") {
		t.Errorf("Message does not end properly")
	}
}

func FuzzExpand_stringArgs(f *testing.F) {
	f.Add("thing")
	f.Add("")
	f.Add(" łäü")
	f.Add("just a rather long piece of content")
	f.Fuzz(func(t *testing.T, arg string) { fuzzArgs(t, arg) })
}

func FuzzExpand_intArgs(f *testing.F) {
	f.Add(0)
	f.Add(-1)
	f.Add(1)
	f.Add(4711)
	f.Fuzz(func(t *testing.T, arg int) { fuzzArgs(t, arg) })
}

func BenchmarkStdLog(b *testing.B) {
	var out bytes.Buffer
	logr := log.New(&out, "stdlog", log.LstdFlags)
	for i := 0; i < b.N; i++ {
		out.Reset()
		logr.Printf("just an `what:%s`: `miss:%s`",
			"example",
			"<undef>")
	}
}

func BenchmarkSllmLog(b *testing.B) {
	var out bytes.Buffer
	logr := log.New(&out, "stdlog", log.LstdFlags)
	for i := 0; i < b.N; i++ {
		out.Reset()
		logr.Print(Sprint("just an `what`: `miss`", Argv("<undef>", "example")))
	}
}
