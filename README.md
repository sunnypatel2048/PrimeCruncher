# Prime Cruncher

## Introduction

Prime Cruncher is a multi-threaded GoLang application designed to compute the number of prime 64-bit unsigned integers stored in a binary file. The program efficiently processes large binary files by dividing the workload among multiple worker threads, ensuring high concurrency and optimal performance. It uses a dispatcher-worker-consolidator architecture implemented with goroutines and channels for efficient thread communication.

## Usage

To run the program, use the following command:

```bash
go run main.go -M <no-of-threads> -N <segment-size> -C <chunk-size> <path-name> 
```

* **no-of-threads:** The number of worker threads (default: 1). Must be a positive integer.
* **segment-size:** The segment size in bytes for each job (default: 64KB). Must be divisible by 8.
* **chunk-size:** The chunk-size in bytes for reading the file (default: 1KB). Must be divisible by 8.
* **path-name:** The relative or absolute path to the input binary file.

## Generate Data

To test Prime Cruncher, you can generate a random binary data file containing. Use the following command to create a file with random binary data:

```bash
head -c <no-of-bytes> /dev/urandom > <file-name>
```

* **no-of-bytes:** The size of the file in bytes. Since the program processes 64-bit (8-byte) integers, this should ideally be a multiple of 8 for consistent results.
* **file-name:** The name of the output file (e.g., ```data.dat```).

## Design and Implementation

Prime Cruncher employs a multi-threaded architecture consisting of one dispatcher, M worker threads, and one consolidator thread, all implemented as goroutines in GoLang. The program uses channels to manage job queues and result queues, ensuring safe and efficient communication between threads.

### Dispatcher

The dispatcher is responsible for partitioning the input file into segments of N bytes each (except possibly the last segment, which may be smaller). For each segment, it creates a job descriptor containing:

* The pathname to the datafile.
* The start position of the segment in bytes.
* The length of the segment in bytes.

The dispatcher places each job descriptor into the job queue (a GoLang channel) and signals completion when all jobs have been created.

### Workers

Each of the M worker threads performs the following steps:

1. Sleeps for a random duration between 400 and 600 milliseconds before starting.
2. Repeatedly fetches jobs from the job queue (blocking until a job becomes available).
3. For each job:
    * Reads the specified segment from the file in chunks of C bytes (or less for the last chunk).
    * Uses the GoLang encoding/binary package to decode 64-bit unsigned integers in little-endian byte order.
    * Counts the number of 64-bit unsigned integers in the segment that are prime.
    * Places a result descriptor into the results queue, including the job descriptor and the prime count.

Terminates when there are no more jobs in the job queue.

### Consolidator

The consolidator thread:

1. Takes results from the results queue (a GoLang channel).
2. Accumulates the prime counts from each worker's results.
3. Once all workers have terminated (detected via a termination signal), it provides the total prime count to the main function and terminates.

## Logging

Each worker logs the completion of every job using GoLang's ```slog.Info()```, recording:

* The job descriptor (pathname, start position, segment length).
* The number of primes found in that segment.

This logging helps monitor worker activity and debug performance issues.

## Efficiency and Concurrency

The implementation is designed to achieve maximum concurrency with minimal synchronization delays:

* **Goroutines and Channels:** Using GoLang's lightweight goroutines and channels minimizes overhead and ensures thread-safe communication without explicit locks.
Chunked Reading: Reading the file in C-byte chunks optimizes I/O operations, balancing memory usage and performance.
* **Minimal Communication:** Job and result descriptors are compact, minimizing the amount of data communicated between threads.
* **Random Worker Delays:** Initial random sleep delays (400-600 ms) help stagger worker startup, reducing contention on the job queue.

The implementation is correct (produces accurate prime counts) and free of race conditions due to the use of channels for all inter-thread communication.
