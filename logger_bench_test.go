package main

import (
	"fileIO/writer"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func BenchmarkEncoderWriter(b *testing.B) {
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

// TestJSONEncoderMethodAllocs measures allocs/op for each individual method of JSONEncoder.
// Run with: go test -v -run TestJSONEncoderMethodAllocs
func TestJSONEncoderMethodAllocs(t *testing.T) {
	enc := NewJSONEncoder()

	printAllocs := func(name string, n float64) {
		fmt.Printf("%-40s %.0f allocs/op\n", name, n)
	}

	fmt.Println("--- JSONEncoder method-level allocs/op ---")

	// addString
	printAllocs("addString", testing.AllocsPerRun(100, func() {
		enc.addString("hello-world")
		enc.reset()
	}))

	// addInt
	printAllocs("addInt", testing.AllocsPerRun(100, func() {
		enc.addInt(12345)
		enc.reset()
	}))

	// addCaller  (runtime.Caller + string concat)
	printAllocs("addRawCaller", testing.AllocsPerRun(100, func() {
		enc.addRawCaller()
		enc.reset()
	}))

	// time.Now().UTC().Format  (isolated — the timestamp line in Encode)
	printAllocs("time.Now().UTC().AppendFormat(enc.b, time.RFC3339Nano)", testing.AllocsPerRun(100, func() {
		enc.b = time.Now().UTC().AppendFormat(enc.b, time.RFC3339Nano)
		enc.reset()
	}))

	// addKeyValue with a string Value
	printAllocs("addKeyValue (string)", testing.AllocsPerRun(100, func() {
		enc.addKeyValue(KV{Key: "k", Value: &Value{val: "v", valType: reflect.String}})
		enc.reset()
	}))

	// addKeyValue with an int64 Value
	printAllocs("addKeyValue (int64)", testing.AllocsPerRun(100, func() {
		enc.addKeyValue(KV{Key: "k", Value: &Value{val: int64(42), valType: reflect.Int64}})
		enc.reset()
	}))

	// addStruct (flat struct — no nested struct)
	type FlatStruct struct {
		Name string
		Age  int64
	}
	printAllocs("addStruct (flat)", testing.AllocsPerRun(100, func() {
		enc.addStruct(FlatStruct{Name: "Ayush", Age: 22})
		enc.reset()
	}))

	// addStruct (nested struct — like MyStruct with MyInfo inside)
	printAllocs("addStruct (nested)", testing.AllocsPerRun(100, func() {
		enc.addStruct(MyStruct{Name: "Ayush", Age: 22, MyInfo: MyInfo{Gender: "Male"}})
		enc.reset()
	}))

	// Full Encode — with pool
	rec := Record{
		Message: "Ayush Singhal",
		Level:   Warn,
		KVs: []KV{
			AddString("my-key", "my-value"),
			AddString("my-key-2", "my-value-2"),
			AddInt("my-int-key", 34),
			AddStruct("my-struct-key", MyStruct{Name: "Ayush", Age: 22, MyInfo: MyInfo{Gender: "Male"}}),
		},
	}
	fmt.Println()
	printAllocs("Encode (with pool)", testing.AllocsPerRun(100, func() {
		e := _jsonPOOL.Get().(*JSONEncoder)
		e.Encode(rec) //nolint:errcheck
		_jsonPOOL.Put(e)
	}))
	printAllocs("Encode (no pool)", testing.AllocsPerRun(100, func() {
		e := NewJSONEncoder()
		e.Encode(rec) //nolint:errcheck
	}))
}
