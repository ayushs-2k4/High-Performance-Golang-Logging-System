.PHONY: all test test-allocs bench bench-encoder bench-writer bench-compare bench-save bench-profile-cpu bench-profile-mem bench-stat install-benchstat clean

# Default: run tests + benchmarks
all: test bench

# ── Tests ────────────────────────────────────────────────────────────────────

# Run all unit tests
test:
	go test -v ./...

# Run the method-level allocs test specifically
test-allocs:
	go test -v -run TestJSONEncoderMethodAllocs .

# ── Benchmarks ───────────────────────────────────────────────────────────────

# Run all benchmarks with memory stats
bench:
	go test -bench=. -benchmem -count=3 .

# Run only the raw encoder benchmark (no I/O)
bench-encoder:
	go test -bench=BenchmarkEncoder$$ -benchmem -count=3 .

# Run full pipeline benchmark (encode + write to file)
bench-writer:
	go test -bench=BenchmarkEncoderWriter -benchmem -count=3 .

# Run custom logger vs Zap side-by-side
bench-compare:
	go test -bench='BenchmarkMyLogger10Fields|BenchmarkZap10Fields' -benchmem -count=3 .

# Run benchmarks and save results to a file (useful for benchstat)
bench-save:
	go test -bench=. -benchmem -count=5 . | tee bench_results.txt

# ── Profiling ─────────────────────────────────────────────────────────────────

# CPU profiling — opens pprof interactively after the run
bench-profile-cpu:
	go test -bench=BenchmarkMyLogger10Fields -benchmem -cpuprofile=cpu.out .
	go tool pprof cpu.out

# Memory profiling — opens pprof interactively after the run
bench-profile-mem:
	go test -bench=BenchmarkMyLogger10Fields -benchmem -memprofile=mem.out .
	go tool pprof mem.out

# ── Benchstat ────────────────────────────────────────────────────────────────

# Install benchstat tool for comparing benchmark results
install-benchstat:
	go install golang.org/x/perf/cmd/benchstat@latest

# Compare saved benchmark results (run bench-save first)
bench-stat:
	benchstat bench_results.txt

# ── Cleanup ───────────────────────────────────────────────────────────────────

clean:
	rm -f cpu.out mem.out bench.log bench_results.txt
