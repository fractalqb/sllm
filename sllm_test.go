package sllm

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Example() {
	// Positional arguments
	FprintIdx(os.Stdout, "added `count` ⨉ `item` to shopping cart by `user`\n",
		7, "Hat", "John Doe",
	)
	// Named arguments
	Fprint(os.Stdout, "added `count` ⨉ `item` to shopping cart by `user`\n",
		NmArgs(map[string]any{
			"count": 7,
			"item":  "Hat",
			"user":  "John Doe",
		}),
	)
	// Output:
	// added `count:7` ⨉ `item:Hat` to shopping cart by `user:John Doe`
	// added `count:7` ⨉ `item:Hat` to shopping cart by `user:John Doe`
}

func ExamplePrint() {
	var writeTestArg = func(buf []byte, idx int, name string) ([]byte, error) {
		w := bytes.NewBuffer(buf)
		fmt.Fprintf(w, "#%d/'%s'", idx, name)
		return w.Bytes(), nil
	}
	Fprint(os.Stdout, "want `arg1` here and `arg2` here", writeTestArg)
	fmt.Println()
	Fprint(os.Stdout, "template with backtick '``' and an `arg` here", writeTestArg)
	fmt.Println()
	Fprint(os.Stdout, "touching args: `one``two``three`", IdxArgsDefault("–", 4711, true))
	fmt.Println()
	FprintIdx(os.Stdout, "explicit `index:0` and `same:0`", 4711)
	fmt.Println()
	// Output:
	// want `arg1:#0/'arg1'` here and `arg2:#1/'arg2'` here
	// template with backtick '``' and an `arg:#0/'arg'` here
	// touching args: `one:4711``two:true``three:–`
	// explicit `index:4711` and `same:4711`
}

func ExamplePrint_explicitIndex() {
	var writeTestArg = func(buf []byte, idx int, name string) ([]byte, error) {
		buf = append(buf, '#')
		return strconv.AppendInt(buf, int64(idx), 10), nil
	}
	Fprint(os.Stdout, "`a`, `b:11`, `c`, `d:0`, `e`", writeTestArg)
	// Output:
	// `a:#0`, `b:#11`, `c:#1`, `d:#0`, `e:#2`
}

func ExamplePrint_argError() {
	_, err := FprintIdx(os.Stdout, "`argok` but `notok`\n", 4711)
	fmt.Println(err)
	// Output:
	// `argok:4711` but `notok!(missing argument 1'notok')`
	// <nil>
}

const (
	testTmpl = "`service`: Sent `signal` to main `process` (`name`) on client request `at`."
	testSvc  = "rsyslog"
	testSig  = "SIGHUP"
	testProc = 1611
	testName = "rsyslogd"
)

var (
	testNow  = Tdefault.Fmt(time.Date(2023, 11, 27, 21, 30, 00, 0, time.UTC))
	outBytes int
)

func testArgsN(to []byte, idx int, n string) ([]byte, error) {
	switch idx {
	case 0:
		return EscString(to, testSvc), nil
	case 1:
		return EscString(to, testSig), nil
	case 2:
		return strconv.AppendInt(to, testProc, 10), nil
	case 3:
		return EscString(to, testName), nil
	case 4:
		return testNow.AppendSllm(to), nil
	}
	return to, MissingArg{idx, n}
}

func BenchmarkSprintf(b *testing.B) {
	now := time.Now()
	outBytes = 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out := fmt.Sprintf("service '%s': Sent signal '%s' to main process %d (name %s) on client request at %s.",
			testSvc, testSig, testProc, testName, now,
		)
		outBytes += len(out)
	}
}

func ExampleAppend_testArgsN() {
	var buf []byte
	buf, _ = Append(buf, testTmpl, testArgsN)
	os.Stdout.Write(buf)
	fmt.Println()
	// Output:
	// `service:rsyslog`: Sent `signal:SIGHUP` to main `process:1611` (`name:rsyslogd`) on client request `at:11-27 Mo 21:30:00`.
}

func BenchmarkAppend_testArgsN(b *testing.B) {
	var buf []byte
	outBytes = 0
	for i := 0; i < b.N; i++ {
		buf = buf[:0]
		buf, _ = Append(buf, testTmpl, testArgsN)
		outBytes += len(buf)
	}
}

func ExampleAppend_IdxArgs() {
	args := IdxArgsDefault("???", testSvc, testSig, testProc, testName, testNow)
	var buf []byte
	buf, _ = Append(buf, testTmpl, args)
	os.Stdout.Write(buf)
	fmt.Println()
	// Output:
	// `service:rsyslog`: Sent `signal:SIGHUP` to main `process:1611` (`name:rsyslogd`) on client request `at:11-27 Mo 21:30:00`.
}

func BenchmarkAppend_IdxArgs(b *testing.B) {
	args := IdxArgsDefault("???", testSvc, testSig, testProc, testName, testNow)
	var buf []byte
	outBytes = 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf = buf[:0]
		buf, _ = Append(buf, testTmpl, args)
		outBytes += len(buf)
	}
}

func BenchmarkWrite_IdxArgs(b *testing.B) {
	args := IdxArgsDefault("???", testSvc, testSig, testProc, testName, testNow)
	var sb bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sb.Reset()
		Fprint(&sb, testTmpl, args)
		outBytes += sb.Len()
	}
}

func ExampleAppend_NmArgs() {
	args := NmArgs(map[string]any{
		"service": testSvc,
		"signal":  testSig,
		"process": testProc,
		"name":    testName,
		"at":      testNow,
	})
	var buf []byte
	buf, _ = Append(buf, testTmpl, args)
	os.Stdout.Write(buf)
	fmt.Println()
	// Output:
	// `service:rsyslog`: Sent `signal:SIGHUP` to main `process:1611` (`name:rsyslogd`) on client request `at:11-27 Mo 21:30:00`.
}

func BenchmarkAppend_NmArgs(b *testing.B) {
	args := NmArgs(map[string]any{
		"service": testSvc,
		"signal":  testSig,
		"process": testProc,
		"name":    testName,
		"at":      testNow,
	})
	var buf []byte
	outBytes = 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf = buf[:0]
		buf, _ = Append(buf, testTmpl, args)
		outBytes += len(buf)
	}
}

func ExampleAppend_time() {
	t := time.Date(2023, 11, 27, 21, 30, 00, 0, time.UTC)
	buf, _ := Append(nil, "`its`", IdxArgs(t))
	os.Stdout.Write(buf)
	fmt.Println()
	// Output:
	// `its:2023-11-27 21:30:00 +0000 UTC`
}

func TestAppend_errors(t *testing.T) {
	t.Run("unterminated mid", func(t *testing.T) {
		_, err := Append(nil, "this `is unterminated", IdxArgs())
		if err == nil {
			t.Fatal("no error")
		}
		if err.Error() != "unterminated parameter" {
			t.Fatalf("unexpected error: %s", err)
		}
	})
	t.Run("unterminated end", func(t *testing.T) {
		_, err := Append(nil, "without end `", IdxArgs())
		if err == nil {
			t.Fatal("no error")
		}
		if err.Error() != "unterminated parameter" {
			t.Fatalf("unexpected error: %s", err)
		}
	})
	t.Run("empty explicit index", func(t *testing.T) {
		_, err := Append(nil, "illegal `place:`", IdxArgs())
		if err == nil {
			t.Fatal("no error")
		}
		if err.Error() != "index in 'place:': strconv.Atoi: parsing \"\": invalid syntax" {
			t.Fatalf("unexpected error: %s", err)
		}
	})
	t.Run("non-numeric explicit index", func(t *testing.T) {
		_, err := Append(nil, "illegal `place:holder`", IdxArgs())
		if err == nil {
			t.Fatal("no error")
		}
		if err.Error() != "index in 'place:holder': strconv.Atoi: parsing \"holder\": invalid syntax" {
			t.Fatalf("unexpected error: %s", err)
		}
	})
	t.Run("empty parameter", func(t *testing.T) {
		_, err := Append(nil, "foo `:7` baz", IdxArgs())
		if err == nil {
			t.Fatal("no error")
		}
		if err.Error() != "empty parameter in ':7'" {
			t.Fatalf("unexpected error: %s", err)
		}
	})
}

func fuzzArgs[T any](t *testing.T, arg T) {
	const tmpl = "Just `a` template"
	var sb bytes.Buffer
	n, err := Fprint(&sb, tmpl, IdxArgsDefault("", arg))
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
