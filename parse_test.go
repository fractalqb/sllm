package sllm

import (
	"fmt"
	"strings"
)

func ExampleParseMap() {
	var tmpl strings.Builder
	m, err := ParseMap("added `count:7` ⨉ `item:Hat` to shopping cart by `user:John Doe`", &tmpl)
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range m {
		fmt.Printf("%s:[%s]\n", k, v)
	}
	// Unordered output:
	// count:[7]
	// item:[Hat]
	// user:[John Doe]
}
