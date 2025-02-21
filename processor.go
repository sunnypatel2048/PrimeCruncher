package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	FilePath string
	Start    int64
	Length   int64
}

type Result struct {
	Job        Job
	PrimeCount int
}

// Dispatcher divides the file into N-byte segments and sends them to the jobQueue
func Dispatcher(jobQueue chan<- Job, filePath string, N int64) error {
	fileSize, err := GetFileSize(filePath)
	if err != nil {
		return err
	}

	var start int64 = 0
	for start < fileSize {
		length := N
		if start+N > fileSize {
			length = fileSize - start
		}
		jobQueue <- Job{
			FilePath: filePath,
			Start:    start,
			Length:   length,
		}
		start += N
	}
	close(jobQueue)
	return nil
}

// Worker processes jobs from the jobQueue and sends results to the resultQueue
func Worker(id int, jobQueue <-chan Job, resultQueue chan<- Result, wg *sync.WaitGroup, chunkSize int) {
	defer wg.Done()

	sleepDuration := time.Duration(rand.Intn(201)+400) * time.Millisecond
	time.Sleep(sleepDuration)

	for job := range jobQueue {
		segment, err := ReadSegment(job.FilePath, job.Start, job.Length)
		if err != nil {
			fmt.Printf("Worker %d: error reading segment: %v\n", id, err)
			continue
		}

		numPrimes := 0
		for i := 0; i < len(segment); i += chunkSize {
			end := i + chunkSize
			if end > len(segment) {
				end = len(segment)
			}
			chunk := segment[i:end]

			for j := 0; j < len(chunk)-8; j += 8 {
				value := binary.LittleEndian.Uint64(chunk[j : j+8])

				if isPrime(value) {
					numPrimes++
				}
			}
		}

		resultQueue <- Result{
			Job:        job,
			PrimeCount: numPrimes,
		}

		fmt.Printf("Worker %d processed segment starting at %d, length %d, found %d primes\n", id, job.Start, job.Length, numPrimes)
	}
}
