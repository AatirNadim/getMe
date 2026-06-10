# Testing and Benchmarking Suite

This directory contains the testing and benchmarking suite for the `getMe` key-value store server. The purpose of this suite is to measure the performance of the core storage engine under various workloads, identify potential bottlenecks, and ensure data correctness.

## Purpose

- **Performance Measurement**: To get concrete numbers on the throughput (operations per second) and latency (nanoseconds per operation) of `Put`, `Get`, `Delete`, `BatchPut`, and `BatchGet` operations.
- **Allocation Analysis**: To measure how much memory is allocated (`B/op`) and how many distinct allocations are made (`allocs/op`) for each operation. This is critical for diagnosing garbage collector pressure.
- **Concurrency and Stress Testing**: To test the database's stability and performance under high-concurrency scenarios with many goroutines reading and writing at the same time.
- **Correctness Verification**: To ensure that data remains consistent and is not corrupted during high-volume, concurrent operations, and to verify API edge cases.

## How to Run Tests and Benchmarks

The tests are written using Go's built-in `testing` package. They should be run from the `server/tests` directory (or the root module) using the `go test` command.

### Running All Benchmarks

To run all benchmark tests in the project:

```bash
cd server/tests
go test -bench . ./...
```

- `-bench .`: This flag tells the Go tool to run all benchmark functions (those starting with `Benchmark...`).
- `./...`: This pattern instructs Go to run the tests in the current directory and all subdirectories.

### Running Specific Benchmarks

You can run a specific benchmark by providing a regular expression to the `-bench` flag that matches its name. For example, to run only the `BenchmarkPut` tests:

```bash
cd server/tests
go test -bench=BenchmarkPut ./...
```

To run only the correctness benchmarks:

```bash
cd server/tests
go test -bench=Correctness ./...
```

### Running Unit Tests

To run the standard correctness unit tests (e.g., edge cases and limit checks):

```bash
cd server/tests
go test ./...
```

## Directory Structure

The testing suite is split into specific directories based on their purpose:

### `stressTest/`
Contains benchmarks focused purely on performance and load testing.
- **`index_test.go`**: Contains core performance benchmarks:
  - `BenchmarkGet`
  - `BenchmarkBatchGet`
  - `BenchmarkPut`
  - `BenchmarkBatchPut`
  - `BenchmarkReadWriteMixed_90_10` (90% reads, 10% writes)
  - `BenchmarkReadWriteMixed_80_20` (80% reads, 20% writes)
  - `BenchmarkDelete`

### `correctnessCheck/`
Contains tests and benchmarks focused on data integrity and API validation.
- **`index_test.go`**: Contains benchmarks that also perform data validation:
  - `BenchmarkGet_Correctness`
  - `BenchmarkReadWriteMixed_Correctness`
  - `BenchmarkBatchGet_Correctness`
- **Controller Tests**: Files like `get_test.go`, `put_test.go`, `batchGet_test.go`, and `batchPut_test.go` contain unit tests for verifying API edge cases and limits (e.g., `TestBatchGetController_MaxBodySize`).

The benchmarks are designed to be run in parallel (`b.RunParallel`) to accurately simulate a multi-threaded server environment and uncover race conditions.
