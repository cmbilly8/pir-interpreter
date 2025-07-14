package writer

import (
	"bytes"
	"fmt"
	"sync"
)

var (
	outputBuffer bytes.Buffer
	outputWriter = &outputBuffer
	writerMutex  sync.Mutex
)

func WriteOutput(format string, args ...interface{}) {
	writerMutex.Lock()
	defer writerMutex.Unlock()
	fmt.Fprintf(outputWriter, format, args...)
}

func GetOutput() string {
	writerMutex.Lock()
	defer writerMutex.Unlock()
	return outputBuffer.String()
}

func ClearOutput() {
	writerMutex.Lock()
	defer writerMutex.Unlock()
	outputBuffer.Reset()
}
