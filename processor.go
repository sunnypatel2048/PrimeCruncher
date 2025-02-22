package main

import (
	"encoding/binary"
	"log/slog"
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
func Worker(id int64, jobQueue <-chan Job, resultQueue chan<- Result, wg *sync.WaitGroup, chunkSize int64, jobsCompleted *[]int64) {
	defer wg.Done()

	sleepDuration := time.Duration(rand.Intn(201)+400) * time.Millisecond
	time.Sleep(sleepDuration)

	count := int64(0)
	for job := range jobQueue {
		segment, err := ReadSegment(job.FilePath, job.Start, job.Length)
		if err != nil {
			slog.Error("Failed to read segment", "workerID", id, "error", err)
			continue
		}

		numPrimes := 0
		for i := int64(0); i < int64(len(segment)); i += chunkSize {
			end := i + chunkSize
			if end > int64(len(segment)) {
				end = int64(len(segment))
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

		count++
	}

	*jobsCompleted = append(*jobsCompleted, count)
	slog.Info("Worker completed jobs", "workerID", id, "numJobs", count)
}

// Consolidator collects results from the resultQueue and calculates the total number of primes
func Consolidator(resultQueue <-chan Result, numJobs int64, done chan<- int) {
	totalPrimes := 0

	for i := int64(0); i < numJobs; i++ {
		result := <-resultQueue
		totalPrimes += result.PrimeCount
		slog.Info("Consolidator received result",
			"numPrimes", result.PrimeCount,
			"start", result.Job.Start,
		)
	}

	done <- totalPrimes
}
