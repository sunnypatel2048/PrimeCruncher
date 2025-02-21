package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

func main() {
	M := flag.Int64("M", 1, "number of worker threads")
	N := flag.Int64("N", 64000, "partition size in bytes")
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

	jobQueue := make(chan Job, *M)
	resultQueue := make(chan Result, *M)

	var wg sync.WaitGroup

	go func() {
		err := Dispatcher(jobQueue, pathName, *N)
		if err != nil {
			fmt.Println("Error dispatching jobs:", err)
			os.Exit(1)
		}
	}()

	for i := 1; i <= int(*M); i++ {
		wg.Add(1)
		go Worker(i, jobQueue, resultQueue, &wg, *C)
	}
}
