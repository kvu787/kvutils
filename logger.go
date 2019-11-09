package util

import (
	"io"
	"log"
	"os"
	"sync"
)

type Options struct {
	UseStdout         bool
	FilePath          string
	AdditionalWriters []io.Writer
}

// Logger can be used concurrently.
func NewLogger(options Options) (*log.Logger, error) {
	writers := make([]io.Writer, len(options.AdditionalWriters))
	copy(writers, options.AdditionalWriters)

	if options.FilePath != "" {
		file, err := os.Create(options.FilePath)
		if err != nil {
			return nil, err
		}
		err = file.Close()
		if err != nil {
			return nil, err
		}
		file, err = os.OpenFile(options.FilePath, os.O_WRONLY, 0)
		if err != nil {
			return nil, err
		}
		writers = append(writers, file)
	}

	if options.UseStdout {
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
