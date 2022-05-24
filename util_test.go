package sllm

import "fmt"

func ExampleError() {
	fmt.Println(Error("this is just `an` message with missing `param`", "error"))
	fmt.Println(Error("this will `fail", "dummy"))
	// Output:
	// this is just `an:error` message with missing `param:<?>`
	// [sllm:syntax error in `tmpl:this will ``fail`:`pos:11`:`desc:unterminated argument`]
}
