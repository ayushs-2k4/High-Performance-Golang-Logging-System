package main

import (
	"os"
)

func main() {
	fileName := "my-file.txt"
	buffer, err := NewBuffer(fileName)
	if err != nil {
		panic(err)
	}
	defer buffer.Sync()
	for i := 0; i < 10; i++ {
		buffer.Write([]byte("\nDevansh Singhal"))
	}
}

func createFileIfNotExists(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if file == nil || err != nil {
		return nil, err
	}

	return file, nil
}
