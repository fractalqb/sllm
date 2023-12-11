package sllm

import (
	"fmt"
	"os"
)

func ExamplePrint_argError() {
	_, err := FprintIdx(os.Stdout, "`argok` but `notok`\n", 4711)
	fmt.Println(err)
	// Output:
	// `argok:4711` but `notok!(missing argument 1 'notok')`
	// <nil>
}
