package sllm

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func ExampleTimeFormat() {
	t := time.Date(2023, 05, 04, 21, 43, 1, 2003000, time.UTC)

	os.Stdout.Write(Tdefault.Fmt(t).AppendSllm(nil))
	fmt.Println()
	os.Stdout.Write((Tyear | Tmicros).Fmt(t).AppendSllm(nil))
	fmt.Println()
	os.Stdout.Write((Tyear | Tmillis).Fmt(t).AppendSllm(nil))
	fmt.Println()
	os.Stdout.Write((Tyear | Tweekday | Tmillis).Fmt(t).AppendSllm(nil))
	fmt.Println()

	// Output:
	// 05-04 Th 21:43:01
	// 2023-05-04 21:43:01.002003
	// 2023-05-04 21:43:01.002
	// 2023-05-04 Th 21:43:01.002
}

func BenchmarkExpand_StdTime(b *testing.B) {
	t := time.Date(2023, 11, 27, 21, 30, 00, 0, time.UTC)
	buf := make([]byte, 64)
	for i := 0; i < b.N; i++ {
		buf, _ = Append(buf[:0], "`its`", IdxArgs(t))
	}
}

func BenchmarkExpand_TimeFormat(b *testing.B) {
	t := Tdefault.Fmt(time.Date(2023, 11, 27, 21, 30, 00, 0, time.UTC))
	buf := make([]byte, 64)
	for i := 0; i < b.N; i++ {
		buf, _ = Append(buf[:0], "`its`", IdxArgs(t))
	}
}
