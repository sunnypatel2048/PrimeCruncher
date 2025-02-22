package main

import (
	"flag"
	"log/slog"
	"os"
	"sort"
	"sync"
	"time"
)

func main() {
	M := flag.Int64("M", 1, "number of worker threads")
	N := flag.Int64("N", 65536, "partition size in bytes")
	C := flag.Int64("C", 1024, "chunk size in bytes")
	flag.Parse()
	pathName := flag.Arg(0)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

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

	jobsCompleted := make([]int64, *M)

	startTime := time.Now()

	go func() {
		err := Dispatcher(jobQueue, pathName, *N)
		if err != nil {
			slog.Error("Dispatcher failed", "error", err)
			os.Exit(1)
		}
	}()

	for i := int64(1); i <= *M; i++ {
		wg.Add(1)
		go Worker(i, jobQueue, resultQueue, &wg, *C, &jobsCompleted)
	}

	fileSize, err := GetFileSize(pathName)
	if err != nil {
		os.Exit(1)
	}
	numJobs := (fileSize + *N - 1) / *N

	go Consolidator(resultQueue, numJobs, done)

	go func() {
		wg.Wait()
		close(resultQueue)
	}()

	totalPrimes := <-done

	elapsedTime := time.Since(startTime)
	slog.Info("Elapsed time", "elapsedTime", elapsedTime.String())

	slog.Info("Total primes found", "totalPrimes", totalPrimes)

	if len(jobsCompleted) > 0 {
		sort.Slice(jobsCompleted,
			func(i, j int) bool { return jobsCompleted[i] < jobsCompleted[j] })

		min := jobsCompleted[0]
		max := jobsCompleted[len(jobsCompleted)-1]
		sum := int64(0)
		for _, count := range jobsCompleted {
			sum += count
		}
		avg := float64(sum) / float64(len(jobsCompleted))
		median := jobsCompleted[len(jobsCompleted)/2]
		if len(jobsCompleted)%2 == 0 {
			median = (jobsCompleted[len(jobsCompleted)/2-1] + jobsCompleted[len(jobsCompleted)/2]) / 2
		}

		slog.Info("Jobs completed statistics",
			"min", min,
			"max", max,
			"average", avg,
			"median", median,
		)
	}

	println("Total primes found:", totalPrimes)
}
