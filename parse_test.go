package sllm

import (
	"fmt"
)

func ExampleParseMap() {
	m, err := ParseMap("added `count:7` ⨉ `item:Hat` to shopping cart by `user:John Doe`")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(m)
	// Output:
	// map[count:7 item:Hat user:John Doe]
}
