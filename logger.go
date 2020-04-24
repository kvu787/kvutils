package util

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LoggerOptions struct {
	UseStdout         bool
	FilePath          string
	AdditionalWriters []io.Writer
}

func NewPrefixLogger(prefix string) (*log.Logger, error) {
	name := fmt.Sprintf("%v_log_%v.txt", prefix, time.Now().Unix())
	return NewLogger(LoggerOptions{true, name, []io.Writer{}})
}

func NewDefaultLogger() (*log.Logger, error) {
	return NewPrefixLogger("default")
}

// Logger can be used concurrently.
func NewLogger(loggerOptions LoggerOptions) (*log.Logger, error) {
	writers := make([]io.Writer, len(loggerOptions.AdditionalWriters))
	copy(writers, loggerOptions.AdditionalWriters)

	if loggerOptions.FilePath != "" {
		doesFileExist, err := DoesFileExist(loggerOptions.FilePath)
		if err != nil {
			return nil, err
		}
		if doesFileExist {
			absoluteFilePath, err := filepath.Abs(loggerOptions.FilePath)
			if err != nil {
				return nil, err
			}
			return nil, errors.New(fmt.Sprintf("File already exists at %v", absoluteFilePath))
		}
		file, err := os.Create(loggerOptions.FilePath)
		if err != nil {
			return nil, err
		}
		// Close, so we can reopen as write-only
		err = file.Close()
		if err != nil {
			return nil, err
		}
		file, err = os.OpenFile(loggerOptions.FilePath, os.O_WRONLY, 0)
		if err != nil {
			return nil, err
		}
		writers = append(writers, file)
	}

	if loggerOptions.UseStdout {
		writers = append(writers, os.Stdout)
	}

	writer := newSyncWriter(io.MultiWriter(writers...))
	logger := log.New(writer, "", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile)
	return logger, nil
}

type syncWriter struct {
	mutex  *sync.Mutex
	writer io.Writer
}

func newSyncWriter(writer io.Writer) io.Writer {
	return &syncWriter{&sync.Mutex{}, writer}
}

func (sw *syncWriter) Write(p []byte) (n int, err error) {
	sw.mutex.Lock()
	n, err = sw.writer.Write(p)
	sw.mutex.Unlock()
	return n, err
}
