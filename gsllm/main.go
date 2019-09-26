package main

import (
	"bufio"
	"fmt"
	"os"

	"git.fractalqb.de/fractalqb/sllm"
)

func main() {
	args := make(map[string]int)
	scn := bufio.NewScanner(os.Stdin)
	for scn.Scan() {
		line := scn.Text()
		params := sllm.ParseMap(line, nil)
		for arg, _ := range params {
			args[arg] = 1
		}
	}
	for arg, _ := range args {
		fmt.Println(arg)
	}
}
