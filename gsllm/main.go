package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"git.fractalqb.de/fractalqb/sllm"
)

func listIds(rd io.Reader, prefix string) {
	args := make(map[string]int)
	scn := bufio.NewScanner(rd)
	for scn.Scan() {
		line := scn.Text()
		params, err := sllm.ParseMap(line, nil)
		if err != nil {
			log.Fatal(err)
		}
		for arg, _ := range params {
			args[arg] = 1
		}
	}
	for arg, _ := range args {
		fmt.Printf("%s: '%s'\n", prefix, arg)
	}
}

func listFileIds(files []string) {
	if len(files) == 0 {
		listIds(os.Stdin, "<stdin>")
		return
	}
	var rd *os.File
	defer rd.Close()
	var err error
	for _, name := range files {
		rd, err = os.Open(name)
		if err != nil {
			log.Fatal(err)
		}
		listIds(rd, name)
		rd.Close()
	}
}

func filter(rd io.Reader, ids []string) {
	match := errors.New("match")
	scn := bufio.NewScanner(rd)
	for scn.Scan() {
		show := sllm.Parse(scn.Text(), nil, func(n, v string) error {
			for _, id := range ids {
				if id == n {
					return match
				}
			}
			return nil
		})
		if show != nil {
			os.Stdout.Write(scn.Bytes())
			fmt.Println()
		}
	}
}

func filterFiles(files []string, ids []string) {
	var rd *os.File
	var err error
	defer rd.Close()
	for _, name := range files {
		rd, err = os.Open(name)
		if err != nil {
			log.Fatal(err)
		}
		filter(rd, ids)
		rd.Close()
	}
}

var (
	fList bool
	fIds  string
)

func main() {
	flag.BoolVar(&fList, "l", false, "List identifiers from input")
	flag.StringVar(&fIds, "i", "", "Filter lines by id")
	flag.Parse()
	if fList {
		listFileIds(flag.Args())
	} else {
		ids := strings.Split(fIds, ",")
		filterFiles(flag.Args(), ids)
	}
}
