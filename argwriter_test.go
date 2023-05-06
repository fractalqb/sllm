package sllm

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func ExampleTimeFormat() {
	t := time.Date(2023, 05, 04, 21, 43, 1, 2003000, time.UTC)
	buf := ArgWriter([]byte{})
	buf.WriteTime(t, Tdefault)
	fmt.Println(string(buf))
	buf = buf[:0]
	buf.WriteTime(t, Tyear|Tmicros)
	fmt.Println(string(buf))
	buf = buf[:0]
	buf.WriteTime(t, Tyear|Tmillis)
	fmt.Println(string(buf))
	buf = buf[:0]
	buf.WriteTime(t, Tyear|Tweekday|Tmillis)
	fmt.Println(string(buf))
	// Output:
	// 05-04 Th 21:43:01
	// 2023-05-04 21:43:01.002003
	// 2023-05-04 21:43:01.002
	// 2023-05-04 Th 21:43:01.002
}

func ExampleArgWriter_Write() {
	var buf []byte
	(*ArgWriter)(&buf).Write([]byte("this `is a funny ``text with `backtik"))
	os.Stdout.Write(buf)
	// Output:
	// this ``is a funny ````text with ``backtik
}

func ExampleArgWriter_WriteString() {
	var buf []byte
	(*ArgWriter)(&buf).WriteString("this `is a funny ``text with `backtik")
	os.Stdout.Write(buf)
	// Output:
	// this ``is a funny ````text with ``backtik
}

func Test_uitoa(t *testing.T) {
	buf := uitoa(nil, 123, 5)
	if s := string(buf); s != "00123" {
		t.Errorf("expect '00123', have '%s'", s)
	}
	buf = uitoa(nil, 123, 4)
	if s := string(buf); s != "0123" {
		t.Errorf("expect '0123', have '%s'", s)
	}
	buf = uitoa(nil, 123, 3)
	if s := string(buf); s != "123" {
		t.Errorf("expect '123', have '%s'", s)
	}
	buf = uitoa(nil, 123, 2)
	if s := string(buf); s != "123" {
		t.Errorf("expect '123', have '%s'", s)
	}
	buf = uitoa(nil, 123, 0)
	if s := string(buf); s != "123" {
		t.Errorf("expect '123', have '%s'", s)
	}
}

func BenchmarkArgWriter(b *testing.B) {
	txt := []byte("this `is a funny ``text with `backtik")
	var pw ArgWriter
	for i := 0; i < b.N; i++ {
		pw = pw[:0]
		pw.Write(txt)
	}
}

func BenchmarkSllmExpand(b *testing.B) {
	var (
		sllmForm = "`service`: Sent `signal` to main `process` (`name`) on client request `at`."
		sllmArgs = Argv("???", "rsyslog", "SIGHUP", 1611, "rsyslogd", Tdefault.Format(time.Now()))
	)
	var buf []byte
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf, _ = Expand(buf[:0], sllmForm, sllmArgs)
	}
}
