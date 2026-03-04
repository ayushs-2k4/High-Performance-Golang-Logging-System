package main

import (
	"fmt"
	"time"
)

const channelSize = 100

func main() {
	filename := "my-file.txt"
	fileLogger := NewFileLogger(filename)

	for i := 0; i < 300; i++ {
		data := fmt.Sprintf("Devansh Singhal, %d", i)
		jsonEncoder := NewJSONEncoder()
		encodedData, _ := jsonEncoder.Encode(data)
		fmt.Println(fmt.Sprintf("producer 1: time: %s, i: %d", time.Now(), i))
		fileLogger.Log(encodedData)
	}

	time.Sleep(5 * time.Second)

	for i := 0; i < 300; i++ {
		data := fmt.Sprintf("Devansh Singhal, %d", i)
		jsonEncoder := NewJSONEncoder()
		encodedData, _ := jsonEncoder.Encode(data)
		fmt.Println(fmt.Sprintf("producer 2: time: %s, i: %d", time.Now(), i))
		fileLogger.Log(encodedData)
	}

	fileLogger.Close()
}

type JSONEncoder struct {
}

func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{}
}

func (j *JSONEncoder) Encode(msg string) ([]byte, error) {
	res := make([]byte, 0, len(msg)+2)

	res = append(res, '{')
	res = append(res, []byte(msg)...)
	res = append(res, '}')

	return res, nil
}
