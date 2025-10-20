# getMe - A High-Performance Key-Value Store

`getMe` is a persistent, embeddable key-value store written in Go. It is inspired by the design of Bitcask and is optimized for high write throughput and low-latency reads.

## Project Structure

This project is a monorepo containing the core storage server, a command-line interface (CLI), client SDKs for various languages, and a benchmarking suite.

- **[`server/`](./server/)**: The core storage engine and HTTP server. This is the heart of the project, implementing the log-structured hash table for persistent storage.
- **[`cli/`](./cli/)**: A command-line interface for interacting with the `getMe` server. Useful for manual testing, debugging, and scripting.
- **[`sdks/`](./sdks/)**: Client libraries (SDKs) for different programming languages to make it easy to integrate `getMe` into your applications.
  - `goSdk/`
  - `javaSdk/`
  - `jsSdk/`
  - `pythonSdk/`
- **[`benchmarking/`](./benchmarking/)**: A comprehensive suite of benchmarks for measuring performance, analyzing memory allocations, and stress-testing the database under concurrent loads.
- **[`utils/`](./utils/)**: Shared utility packages, such as a logger and global constants, used across the project.

## Core Concepts

The storage engine is built on a few key principles:

- **Log-Structured Storage**: All data is written to an append-only log file. This makes writes extremely fast as it avoids slow, random disk I/O.
- **In-Memory Hash Index**: A hash table is kept in memory, mapping each key to the exact location of its value on disk. This allows for very fast read operations, typically requiring only a single disk seek.
- **Compaction**: A background process that periodically cleans up old, stale data from the log files to reclaim disk space.
- **Batch Operations**: A `BatchPut` API is provided to amortize the cost of writes, allowing for very high throughput when ingesting large amounts of data.

## Getting Started

### Running the Server

1.  Navigate to the `server` directory:
    ```bash
    cd server
    ```
2.  Run the server:
    ```bash
    go run .
    ```

### Using the CLI

1.  Navigate to the `cli` directory:
    ```bash
    cd cli
    ```
2.  Use the `set`, `get`, or `delete` commands:
    ```bash
    go run . set mykey "hello world"
    go run . get mykey
    ```

### Running Benchmarks

1.  Navigate to the project root.
2.  Run the tests with the `-bench` flag:
    ```bash
    go test -bench . ./...
    ```

For more detailed information, please refer to the `README.md` file within each respective directory.
