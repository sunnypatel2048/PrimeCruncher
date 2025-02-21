package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"
)

func main() {
	output := ""

	M := flag.Int64("M", 1, "number of worker threads")
	N := flag.Int64("N", 64000, "partition size in bytes")
	C := flag.Int64("C", 1000, "chunk size in bytes")
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

	go func() {
		err := Dispatcher(jobQueue, pathName, *N)
		if err != nil {
			fmt.Println("Error dispatching jobs:", err)
			os.Exit(1)
		}
	}()

	file, err := os.Open(pathName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		os.Exit(1)
	}
	fileSize := fileInfo.Size()

	// Ensure the file size is a multiple of 8 bytes (64 bits)
	if fileSize%8 != 0 {
		fmt.Println("File size must be a multiple of 8")
		os.Exit(1)
	}

	data := make([]byte, fileSize)
	_, err = file.Read(data)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// Interpret teh data as little-endian 64-bit unsigned integers
	integers := make([]uint64, fileSize/8)
	for i := 0; i < len(integers); i++ {
		integers[i] = binary.LittleEndian.Uint64(data[i*8 : (i+1)*8])
	}

	for i := 0; i < len(integers); i++ {
		fmt.Println(isPrime(integers[i]))
	}

	output += fmt.Sprintf("M: %d N: %d C: %d path: %s", *M, *N, *C, pathName)
	fmt.Println(output)
}
