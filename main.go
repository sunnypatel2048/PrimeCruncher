package main

import (
	"flag"
	"log/slog"
	"os"
	"sync"
)

func main() {
	M := flag.Int64("M", 1, "number of worker threads")
	N := flag.Int64("N", 64000, "partition size in bytes")
	C := flag.Int64("C", 1000, "chunk size in bytes")
	flag.Parse()
	pathName := flag.Arg(0)

	if *N%8 != 0 || *C%8 != 0 {
		slog.Error("N and C must be multiples of 8")
		os.Exit(1)
	}

	if pathName == "" {
		slog.Error("Please provide a file name")
		os.Exit(1)
	}

	jobQueue := make(chan Job, *M)
	resultQueue := make(chan Result, *M)
	done := make(chan int)

	var wg sync.WaitGroup

	go func() {
		err := Dispatcher(jobQueue, pathName, *N)
		if err != nil {
			slog.Error("Dispatcher failed", "error", err)
			os.Exit(1)
		}
	}()

	for i := int64(1); i <= *M; i++ {
		wg.Add(1)
		go Worker(i, jobQueue, resultQueue, &wg, *C)
	}

	fileSize, err := GetFileSize(pathName)
	if err != nil {
		os.Exit(1)
	}
	numJobs := (fileSize + *N - 1) / *N

	go Consolidator(resultQueue, numJobs, done)

	wg.Wait()
	close(resultQueue)

	totalPrimes := <-done

	slog.Info("Total primes found", "totalPrimes", totalPrimes)
}
