package consumer

import (
	"log"
	"os"
	"sync"
)

type FileWriter interface {
	Write(content string)
}

// FileWriter represents a writer that writes to a file concurrently.
type FileWriterImpl struct {
	filename string
	dataCh   chan string
	wg       sync.WaitGroup
	logger   *log.Logger
}

// NewFileWriter creates a new instance of FileWriter.
func NewFileWriter(filename string, logger *log.Logger) *FileWriterImpl {
	return &FileWriterImpl{
		filename: filename,
		dataCh:   make(chan string),
		logger:   logger,
	}
}

// Write writes data to the file.
func (fw *FileWriterImpl) Write(data string) {
	// Writing to the file might be considered a "long io" operation.
	// We may even delay it "intentionally", so we can run our "delivery processing" in thread pool of say cpuCores*8 threads.
	// Technically such approach will not give much advantage.
	// Neither speed - because of single bottleneck file in the end where in general we don't write even pseudo-simultaneously.
	// Nor usability advantage, because we work via message broker and don't need to say listen to client connections directly.
	// But as far as I understand - need to demonstrate explicit parallelism somewhere :)
	// Intentional delay for emulating "slow io operation"
	// <-time.After(time.Second * 10)
	fw.dataCh <- data
}

// Start starts the FileWriter to handle writing data to the file.
func (fw *FileWriterImpl) Start() {
	fw.wg.Add(1)
	go func() {
		defer fw.wg.Done()
		file, err := os.OpenFile(fw.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fw.logger.Printf("Error opening file: %v\n", err)
			return
		}
		defer file.Close()

		for data := range fw.dataCh {
			_, err := file.WriteString(data)
			if err != nil {
				fw.logger.Printf("Error writing to file: %v\n", err)
			}
		}
	}()
}

// Close closes the FileWriter and waits for all pending writes to complete.
func (fw *FileWriterImpl) Close() {
	close(fw.dataCh)
	fw.wg.Wait()
}
