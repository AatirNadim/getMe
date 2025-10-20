# Benchmarking Suite

This directory contains the benchmarking suite for the `getMe` key-value store. The purpose of these benchmarks is to measure the performance of the core storage engine under various workloads and to identify potential bottlenecks.

## Purpose

- **Performance Measurement**: To get concrete numbers on the throughput (operations per second) and latency (nanoseconds per operation) of `Put`, `Get`, `Delete`, and `BatchPut` operations.
- **Allocation Analysis**: To measure how much memory is allocated (`B/op`) and how many distinct allocations are made (`allocs/op`) for each operation. This is critical for diagnosing garbage collector pressure.
- **Concurrency and Stress Testing**: To test the database's stability and performance under high-concurrency scenarios with many goroutines reading and writing at the same time.
- **Correctness Verification**: To ensure that data remains consistent and is not corrupted during high-volume, concurrent operations.

## How to Run Benchmarks

The benchmarks are written using Go's built-in `testing` package. They should be run from the root of the `getMe` project using the `go test` command.

### Running All Benchmarks

To run all benchmark tests in the project:

```bash
go test -bench . ./...
```

- `-bench .`: This flag tells the Go tool to run all benchmark functions (those starting with `Benchmark...`).
- `./...`: This pattern instructs Go to run the tests in the current directory and all subdirectories.

### Running Specific Benchmarks

You can run a specific benchmark by providing a regular expression to the `-bench` flag that matches its name. For example, to run only the `BenchmarkPut` tests:

```bash
go test -bench=Put ./...
```

To run only the correctness checks:

```bash
go test -bench=Correctness ./...
```

## Key Benchmark Files

- **`main_test.go`**: This is the primary file containing all the benchmark logic.
  - **`BenchmarkPut`**: Measures the performance of concurrent single-key writes.
  - **`BenchmarkGet_Correctness`**: Measures the performance of concurrent reads while also verifying that the data read is correct.
  - **`BenchmarkBatchPut`**: Measures the performance of bulk-writing batches of key-value pairs.
  - **`BenchmarkDelete`**: Measures the performance of concurrent single-key deletions.
  - **`BenchmarkReadWriteMixed_Correctness`**: Simulates a more realistic workload with a mix of read and write operations, and includes correctness checks.

The benchmarks are designed to be run in parallel (`b.RunParallel`) to accurately simulate a multi-threaded server environment and uncover race conditions.
