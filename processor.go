package main

type Job struct {
	FilePath string
	Start    int64
	Lenght   int64
}

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
			Lenght:   length,
		}
		start += N
	}
	close(jobQueue)
	return nil
}
