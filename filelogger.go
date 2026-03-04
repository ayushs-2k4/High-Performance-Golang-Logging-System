package main

import (
	"bufio"
	"os"
)

type FileLogger struct {
	file   *os.File
	writer *bufio.Writer
	ch     chan []byte
	done   chan struct{}
}

func NewFileLogger(filename string) *FileLogger {
	file, _ := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	writer := bufio.NewWriter(file)

	ch := make(chan []byte, channelSize)

	fileLogger := &FileLogger{
		file:   file,
		writer: writer,
		ch:     ch,
		done:   make(chan struct{}),
	}

	go fileLogger.run()

	return fileLogger
}

func (f *FileLogger) run() {
	for msg := range f.ch {
		f.writer.Write(msg)
		//time.Sleep(100 * time.Millisecond)
	}

	// channel closed → flush remaining data
	f.writer.Flush()
	f.file.Close()

	close(f.done)
}

func (f *FileLogger) Log(b []byte) {
	f.ch <- b
}

func (f *FileLogger) Close() {
	close(f.ch) // signal no more logs
	<-f.done
}
