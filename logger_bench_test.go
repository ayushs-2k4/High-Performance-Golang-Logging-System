package main

import (
	"fileIO/writer"
	"testing"
)

func BenchmarkEncoderLogger(b *testing.B) {
	fileWriter := writer.NewFileWriter("bench.log")
	//consoleWriter := writer.NewConsoleWriter()
	multiWriter := writer.NewMultiWriter(fileWriter)

	record := Record{
		Message: "Ayush Singhal",
		Level:   Warn,
		KVs: []KV{
			AddString("my-key", "my-value"),
			AddString("my-key-2", "my-value-2"),
			AddInt("my-int-key", 34),
			AddStruct("my-struct-key", MyStruct{
				Name:   "Ayush",
				Age:    22,
				MyInfo: MyInfo{Gender: "Male"},
			}),
		},
	}

	b.ResetTimer()

	record.Message = "Ayush Singhal"
	for i := 0; i < b.N; i++ {
		jsonEncoder := _jsonPOOL.Get().(*JSONEncoder)
		encodedData, _ := jsonEncoder.Encode(record)

		multiWriter.Write(encodedData)

		_jsonPOOL.Put(jsonEncoder)
	}

	b.StopTimer()
	multiWriter.Close()
}

func BenchmarkEncoder(b *testing.B) {
	rec := Record{
		Message: "Ayush Singhal",
		Level:   Warn,
		KVs: []KV{
			AddString("my-key", "my-value"),
			AddInt("my-int-key", 34),
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		jsonEncoder := _jsonPOOL.Get().(*JSONEncoder)
		jsonEncoder.Encode(rec)
		_jsonPOOL.Put(jsonEncoder)
	}
}

func BenchmarkFileWriter(b *testing.B) {
	fileWriter := writer.NewFileWriter("bench.log")

	record := Record{
		Message: "Ayush Singhal",
		Level:   Warn,
		KVs: []KV{
			AddString("my-key", "my-value"),
			AddString("my-key-2", "my-value-2"),
			AddInt("my-int-key", 34),
			AddStruct("my-struct-key", MyStruct{
				Name:   "Ayush",
				Age:    22,
				MyInfo: MyInfo{Gender: "Male"},
			}),
		},
	}

	b.ResetTimer()

	record.Message = "Ayush Singhal"
	jsonEncoder := _jsonPOOL.Get().(*JSONEncoder)
	encodedData, _ := jsonEncoder.Encode(record)
	for i := 0; i < b.N; i++ {

		fileWriter.Write(encodedData)

		_jsonPOOL.Put(jsonEncoder)
	}

	b.StopTimer()
	fileWriter.Close()
}
