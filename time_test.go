package sllm

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func ExampleTimeFormat() {
	tz := time.FixedZone("West", -int((3 * time.Hour).Seconds()))
	t := time.Date(2023, 05, 04, 21, 43, 1, 2003000, tz)

	os.Stdout.Write(TimeFormat(0).Fmt(t).AppendSllm(nil))
	fmt.Print("\n\n")

	os.Stdout.Write(TUTC.Fmt(t).AppendSllm(nil))
	fmt.Println()
	os.Stdout.Write(TNoDate.Fmt(t).AppendSllm(nil))
	fmt.Println()
	os.Stdout.Write(TNoWeekday.Fmt(t).AppendSllm(nil))
	fmt.Println()
	os.Stdout.Write(TYear.Fmt(t).AppendSllm(nil))
	fmt.Println()
	os.Stdout.Write(TNoClock.Fmt(t).AppendSllm(nil))
	fmt.Println()
	os.Stdout.Write(TMillis.Fmt(t).AppendSllm(nil))
	fmt.Println()
	os.Stdout.Write(TMicros.Fmt(t).AppendSllm(nil))
	fmt.Println()
	// Output:
	// 05-04 Th 21:43:01-03
	//
	// 05-05 Fr 00:43:01
	// 21:43:01-03
	// 05-04 21:43:01-03
	// 2023-05-04 Th 21:43:01-03
	// 05-04 Th -03
	// 05-04 Th 21:43:01.002-03
	// 05-04 Th 21:43:01.002003-03
}

func TestTimeFormat_Append(t *testing.T) {
	tz := time.FixedZone("West", -int((3 * time.Hour).Seconds()))
	ts := time.Date(2023, 05, 04, 21, 43, 1, 2003000, tz)
	out := TMicros.Append(nil, ts)
	if s := string(out); s != "05-04 Th 21:43:01.002003-03" {
		t.Fatalf("unexpected time string '%s'", s)
	}
}

func BenchmarkAppend_StdTime(b *testing.B) {
	t := time.Date(2023, 11, 27, 21, 30, 00, 0, time.UTC)
	buf := make([]byte, 64)
	for i := 0; i < b.N; i++ {
		buf, _ = Append(buf[:0], "`its`", IdxArgs(t))
	}
}

func BenchmarkAppend_TimeFormat(b *testing.B) {
	t := TDefault.Fmt(time.Date(2023, 11, 27, 21, 30, 00, 0, time.UTC))
	buf := make([]byte, 64)
	for i := 0; i < b.N; i++ {
		buf, _ = Append(buf[:0], "`its`", IdxArgs(t))
	}
}
