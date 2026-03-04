package main

import (
	"fileIO/writer"
	"fmt"
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
		jsonEncoder.Encode(rec) //nolint:errcheck
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

// TestAllocsPerOperation measures exact heap allocations for each operation.
// Run with: go test -v -run TestAllocsPerOperation
func TestAllocsPerOperation(t *testing.T) {
	rec := Record{
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

	// Measure allocs for Encode only
	encodeAllocs := testing.AllocsPerRun(100, func() {
		enc := _jsonPOOL.Get().(*JSONEncoder)
		enc.Encode(rec)
		_jsonPOOL.Put(enc)
	})
	fmt.Printf("Encode (pool):          %.0f allocs/op\n", encodeAllocs)

	// Measure allocs for Encode WITHOUT pool
	encodeNoPoolAllocs := testing.AllocsPerRun(100, func() {
		enc := NewJSONEncoder()
		enc.Encode(rec)
	})
	fmt.Printf("Encode (no pool):       %.0f allocs/op\n", encodeNoPoolAllocs)

	// Measure allocs for multiWriter.Write
	fileWriter := writer.NewFileWriter("bench.log")
	multiWriter := writer.NewMultiWriter(fileWriter)
	enc := _jsonPOOL.Get().(*JSONEncoder)
	data, _ := enc.Encode(rec)
	_jsonPOOL.Put(enc)

	writeAllocs := testing.AllocsPerRun(100, func() {
		multiWriter.Write(data)
	})
	fmt.Printf("MultiWriter.Write:      %.0f allocs/op\n", writeAllocs)

	fileWriteAllocs := testing.AllocsPerRun(100, func() {
		fileWriter.Write(data)
	})
	fmt.Printf("FileWriter.Write:       %.0f allocs/op\n", fileWriteAllocs)

	// Measure allocs for the full pipeline (Encode + Write)
	fullAllocs := testing.AllocsPerRun(100, func() {
		e := _jsonPOOL.Get().(*JSONEncoder)
		d, _ := e.Encode(rec)
		multiWriter.Write(d)
		_jsonPOOL.Put(e)
	})
	fmt.Printf("Full pipeline (pool):   %.0f allocs/op\n", fullAllocs)

	multiWriter.Close()
}
