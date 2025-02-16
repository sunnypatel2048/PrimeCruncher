package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	output := ""

	M := flag.Int("M", 1, "number of worker threads")
	N := flag.Int("N", 64000, "partition size in bytes")
	C := flag.Int("C", 1000, "chunk size in bytes")
	flag.Parse()
	pathName := flag.Arg(0)

	if *N%8 != 0 || *C%8 != 0 {
		fmt.Println("N and C must be multiples of 8")
		os.Exit(1)
	}

	if pathName == "" {
		fmt.Println("Please provide a file name")
		os.Exit(1)
	}

	output += fmt.Sprintf("M: %d N: %d C: %d path: %s", *M, *N, *C, pathName)
	fmt.Println(output)
}
