# High-Performance Structured Logger in Go

A zero-dependency, allocation-conscious structured logging library written in Go. It features a custom JSON encoder, a pooled encoder architecture, and a suite of pluggable writers — all benchmarked against [Uber Zap](https://github.com/uber-go/zap).

---

## Features

- **Zero-alloc JSON encoding** — custom `JSONEncoder` built on `[]byte` appends, avoiding `encoding/json` for hot paths
- **`sync.Pool` encoder pooling** — reuses encoder instances to minimise GC pressure
- **Pluggable writers** — `FileWriter`, `ConsoleWriter`, `MultiWriter`, and `DiscardWriter` all implement a common `Writer` interface
- **Async file writer** — buffered, channel-backed `FileWriter` for non-blocking log writes
- **Structured key-value fields** — strongly-typed helpers for strings, ints, floats, bools, structs, arrays, and custom `ArrayMarshal` types
- **Log levels** — `Debug`, `Info`, `Warn`, `Error`
- **Optional caller info** — file/line/function capture via `runtime.Caller`
- **Prettified output** — optional indented JSON for human-readable logs (toggle via `shouldPrettify`)

---

## Project Structure

```
root/
├── main.go                   # Example usage: concurrent logging with MultiWriter
├── go.mod / go.sum
├── Makefile                  # Build, test, bench, and profiling targets
│
├── logger/
│   ├── encoder.go            # Encoder interface
│   ├── jsonencoder.go        # Core JSON encoder + sync.Pool
│   ├── record.go             # Record, KV, Value types and Add* helpers
│   ├── level.go              # Log level constants and String()
│   └── logger_bench_test.go  # Benchmarks and alloc tests (vs. Zap)
│
├── writer/
│   ├── writer.go             # Writer interface
│   ├── filewriter.go         # Async, buffered file writer
│   ├── consolewriter.go      # stdout writer
│   ├── multiwriter.go        # Fan-out to multiple writers
│   └── DiscardWriter.go      # No-op writer (for benchmarks)
│
├── buffer/
│   ├── buffer.go             # Manual byte-buffer with flush-on-full semantics
│   └── helper.go
│
└── models/
    └── models.go             # Example domain models (Person, Address, etc.)
```

---

## Getting Started

### Prerequisites

- Go 1.21+

### Install dependencies

```zsh
go mod download
```

### Run the example

```zsh
go run main.go
```

This spawns 200 goroutines, each encoding a deeply nested `Person` struct as JSON and writing it to both a file (`my-file.txt`) and stdout via `MultiWriter`.

---

## Usage

### Creating a Record and encoding it

```go
record := logger.Record{
    Message: "user signed in",
    Level:   logger.Info,
    KVs: []logger.KV{
        logger.AddString("user_id", "u-123"),
        logger.AddInt64("attempt", 1),
        logger.AddBool("success", true),
    },
}

enc := logger.GetJSONEncoder()
data, _ := enc.Encode(record)
logger.PutJSONEncoder(enc) // return to pool
```

### Writing output

```go
// File only
fileWriter := writer.NewFileWriter("app.log")
fileWriter.Write(data)
fileWriter.Close()

// Console only
consoleWriter := writer.NewConsoleWriter()
consoleWriter.Write(data)

// Both simultaneously
multi := writer.NewMultiWriter(fileWriter, consoleWriter)
multi.Write(data)
multi.Close()
```

### Available field types

| Helper | Type |
|---|---|
| `AddString(key, value)` | `string` |
| `AddInt(key, value)` | `int` |
| `AddInt32(key, value)` | `int32` |
| `AddInt64(key, value)` | `int64` |
| `AddFloat32(key, value)` | `float32` |
| `AddFloat64(key, value)` | `float64` |
| `AddBool(key, value)` | `bool` |
| `AddStruct(key, value)` | any struct (via reflection) |
| `AddArray(key, value)` | any slice/array (via `encoding/json`) |
| `AddArrayMarshal(key, value)` | `ArrayMarshal` — zero-alloc custom array serialisation |

### Zero-alloc custom array serialisation

Implement the `ArrayMarshal` interface to bypass `encoding/json` entirely for slice fields:

```go
type Tags []string

func (t Tags) MarshalArray(b []byte) ([]byte, error) {
    b = append(b, '[')
    for i, tag := range t {
        if i > 0 {
            b = append(b, ',')
        }
        b = append(b, '"')
        b = append(b, tag...)
        b = append(b, '"')
    }
    return append(b, ']'), nil
}
```

---

## Running Tests & Benchmarks

All targets are available via `make`:

| Command | Description |
|---|---|
| `make test` | Run all unit tests |
| `make test-allocs` | Print per-method allocs/op for `JSONEncoder` |
| `make bench` | Run all benchmarks with memory stats |
| `make bench-encoder` | Raw encoder only (no I/O) |
| `make bench-writer` | Full pipeline: encode + write to file |
| `make bench-compare` | Custom logger vs Zap side-by-side |
| `make bench-save` | Save benchmark results to `bench_results.txt` |
| `make bench-profile-cpu` | CPU profile with pprof |
| `make bench-profile-mem` | Memory profile with pprof |
| `make install-benchstat` | Install `benchstat` for statistical comparison |
| `make bench-stat` | Run `benchstat` on saved results |
| `make clean` | Remove generated files |

### Example benchmark comparison (10 fields, discard writer)

```
BenchmarkMyLogger10Fields-10         ~0 allocs/op
BenchmarkZap10Fields-10              ~1 allocs/op
```

---

## Configuration

Two compile-time flags in `logger/jsonencoder.go` control output behaviour:

| Constant | Default | Effect |
|---|---|---|
| `shouldPrettify` | `true` | Indented, human-readable JSON output |
| `shouldAddCallerInfo` | `false` | Attach `file:line funcName` to every log entry |

---

## Dependencies

| Package | Purpose |
|---|---|
| `go.uber.org/zap` | Benchmark baseline only |
| `go.uber.org/multierr` | Transitive dependency of zap |

The core logger and writer packages have **no external dependencies**.
