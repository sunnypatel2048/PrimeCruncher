package main

import (
	"fmt"
	"log/slog"
	"os"
)

// OpenFile opens a file and returns a file descriptor
func OpenFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		slog.Error("error opening file", "error", err)
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	return file, nil
}

// GetFileSize returns the size of a file in bytes
func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		slog.Error("error getting file info", "error", err)
		return 0, fmt.Errorf("error getting file info: %v", err)
	}
	return fileInfo.Size(), nil
}

// ReadSegment reads a segment of a file and returns it as a byte slice
func ReadSegment(filePath string, start, length int64) ([]byte, error) {
	file, err := OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	segment := make([]byte, length)
	_, err = file.ReadAt(segment, start)
	if err != nil {
		slog.Error("error reading segment", "error", err)
		return nil, fmt.Errorf("error reading segment: %v", err)
	}
	return segment, nil
}
