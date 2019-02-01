package sllm

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func ExampleArgs() {
	Expand(os.Stdout, "just an `what`", Args(nil, "example"))
	// Output:
	// just an `what:example`
}

func ExampleArgs_undef() {
	Expand(os.Stdout, "just an `what`: `miss`", Args([]byte("<undef>"), "example"))
	// Output:
	// just an `what:example`: `miss:<undef>`
}

func ExampleMap() {
	Expand(os.Stdout, "just an `what`", Map(nil, ArgMap{
		"what":  "example",
		"dummy": false,
	}))
	// Output:
	// just an `what:example`
}

func ExampleMap_undef() {
	Expand(os.Stdout, "just an `what`: `miss`", Map([]byte("<undef>"), ArgMap{
		"what":  "example",
		"dummy": false,
	}))
	// Output:
	// just an `what:example`: `miss:<undef>`
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
		logr.Print(Expands("just an `what`: `miss`", Args([]byte("<undef>"), "example")))
	}
}
