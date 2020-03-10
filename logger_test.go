package util

import (
	"sync"
	"testing"
)

var DefaultOptions LoggerOptions = LoggerOptions{true, "test.log", nil}

func TestLog(t *testing.T) {
	logger, err := NewLogger(DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}
	logger.Println("here")
}

func TestDefaultLogger(t *testing.T) {
	logger, err := NewDefaultLogger()
	if err != nil {
		t.Fatal(err)
	}
	logger.Println("here")
}

func TestLogConcurrent(t *testing.T) {
	logger, err := NewLogger(DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}
	logger.Println("here")
	numThreads := 10
	wg := &sync.WaitGroup{}
	wg.Add(numThreads)
	for i := 0; i < numThreads; i++ {
		go func(n int) {
			logger.Println(n)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
