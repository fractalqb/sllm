package sllm

import (
	"fmt"
	"testing"
	"time"
)

func ExampleTimeFormat() {
	t := time.Date(2023, 05, 04, 21, 43, 1, 2003000, time.UTC)
	buf := ArgWriter([]byte{})
	buf.WriteTime(t, Tdefault)
	fmt.Println(string(buf))
	buf = buf[:0]
	buf.WriteTime(t, Tdate|Tyear|Tmicros)
	fmt.Println(string(buf))
	buf = buf[:0]
	buf.WriteTime(t, Tdate|Tyear|Tmillis)
	fmt.Println(string(buf))
	// Output:
	// May 04 21:43:01
	// 2023 May 04 21:43:01.002003
	// 2023 May 04 21:43:01.002
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
	_ = buf[1]
}
