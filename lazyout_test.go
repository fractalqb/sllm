package sllm

import (
	"bytes"
	"fmt"
	"log"
	"testing"
)

func ExampleArgs() {
	fmt.Println(Args("just an `what`", "example"))
	// Output:
	// just an `what:example`
}

func ExampleUArgs() {
	fmt.Println(UArgs("just an `what`: `miss`", []byte("<undef>"), "example"))
	// Output:
	// just an `what:example`: `miss:<undef>`
}

func ExampleMap() {
	fmt.Println(Map("just an `what`", ArgMap{
		"what":  "example",
		"dummy": false,
	}))
	// Output:
	// just an `what:example`
}

func ExampleUMap() {
	fmt.Println(UMap("just an `what`: `miss`", []byte("<undef>"),
		ArgMap{
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
		logr.Print(UArgs("just an `what`: `miss`", []byte("<undef>"), "example"))
	}
}
