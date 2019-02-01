package sllm

import (
	"bytes"
	"fmt"
)

func ExampleParseMap() {
	var tmpl bytes.Buffer
	m, err := ParseMap(
		"added `count:7` ⨉ `item:Hat` to shopping cart by `user:John Doe`",
		&tmpl,
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tmpl.String())
	for k, v := range m {
		fmt.Printf("%s:[%s]\n", k, v)
	}
	// Unordered output:
	// added `count` ⨉ `item` to shopping cart by `user`
	// count:[7]
	// item:[Hat]
	// user:[John Doe]
}
