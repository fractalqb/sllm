package anon

import (
	"fmt"
	"os"

	"git.fractalqb.de/fractalqb/sllm"
)

func ExampleByName() {
	anon := ByName{
		Clear: func(wr sllm.ValueEsc, idx int, name string) (int, error) {
			return fmt.Printf("<param #%d %s>", idx, name)
		},
		Anon: map[string]sllm.ParamWriter{
			"bar": Replace("XXXXX"),
		},
	}
	sllm.Expand(os.Stdout, "msg is `foo`, `bar`, `baz`", anon.Param)
	// Output:
	// msg is `foo:<param #0 foo>`, `bar:XXXXX`, `baz:<param #2 baz>`
}
