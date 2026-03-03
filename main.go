package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const channelSize = 2

func main() {
	filename := "my-file.txt"
	file, _ := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	buffer := bufio.NewWriter(file)
	defer func() {
		buffer.Flush()
		file.Close()
	}()

	ch := make(chan []byte, channelSize)

	go func() {
		for i := 0; i < 300; i++ {
			data := fmt.Sprintf("\nDevansh Singhal, %d", i)
			byteData := []byte(data)
			fmt.Println("time: ", time.Now())
			ch <- byteData
		}
		close(ch)
	}()

	for data := range ch {
		fmt.Println("consumed ")
		dataSize := len(data)
		buffer.Write(data)
		bufferSize := buffer.Buffered()
		fmt.Println(fmt.Sprintf("bufferSize: %d bytes, written bytes: %d", bufferSize, dataSize))
		time.Sleep(1 * time.Second)
	}

}
