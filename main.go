package main

import (
	buffer2 "fileIO/buffer"
	"fmt"
	"sync"
)

const channelSize = 2

func main() {
	fileName := "my-file.txt"
	var wg sync.WaitGroup
	buffer, err := buffer2.NewBuffer(fileName)
	if err != nil {
		panic(err)
	}
	defer buffer.Sync()

	ch := make(chan []byte, channelSize)

	for i := 0; i < 300; i++ {
		data := fmt.Sprintf("\nDevansh Singhal, %d", i)
		byteData := []byte(data)
		wg.Add(1)
		go func(b []byte) {
			ch <- b
			wg.Done()
		}(byteData)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for data := range ch {
		buffer.Write(data)
	}

}
