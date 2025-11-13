# getMe - A High-Performance Key-Value Store

`getMe` is a persistent, embeddable key-value store written in Go. It is inspired by the design of Bitcask and is optimized for high write throughput and low-latency reads.

## Project Structure

This project is a monorepo containing the core storage server, a command-line interface (CLI), client SDKs for various languages, and a benchmarking suite.

- **[`server/`](./server/)**: The core storage engine and HTTP server. This is the heart of the project, implementing the log-structured hash table for persistent storage. <br>**Important:** See [`server/README.md`](./server/README.md) for a diagram-rich tour of the storage layers, HTTP controllers, and background workers.
- **[`cli/`](./cli/)**: A command-line interface for interacting with the `getMe` server. Useful for manual testing, debugging, and scripting.
- **[`sdks/`](./sdks/)**: Client libraries (SDKs) for different programming languages to make it easy to integrate `getMe` into your applications.
  - `goSdk/`
  - `javaSdk/`
  - `jsSdk/`
  - `pythonSdk/`
- **[`benchmarking/`](./benchmarking/)**: A comprehensive suite of benchmarks for measuring performance, analyzing memory allocations, and stress-testing the database under concurrent loads.
- **[`utils/`](./utils/)**: Shared utility packages, such as a logger and global constants, used across the project.

> **Spotlight:** The curated inner docs are the quickest way to understand the system end-to-end. Start with [`server/README.md`](./server/README.md) for architecture fundamentals, then explore [`benchmarking/README.md`](./benchmarking/README.md) to see how we measure performance and [`cli/README.md`](./cli/README.md) for a walkthrough of the local tooling.

## Core Concepts

The storage engine is built on a few key principles:

- **Log-Structured Storage**: All data is written to an append-only log file. This makes writes extremely fast as it avoids slow, random disk I/O.
- **In-Memory Hash Index**: A hash table is kept in memory, mapping each key to the exact location of its value on disk. This allows for very fast read operations, typically requiring only a single disk seek.
- **Compaction**: A background process that periodically cleans up old, stale data from the log files to reclaim disk space.
- **Batch Operations**: A `BatchPut` API is provided to amortize the cost of writes, allowing for very high throughput when ingesting large amounts of data.

## Getting Started

### Running the Server

The repository ships with helper scripts that bootstrap everything you need for a local or containerised deployment of the store.

#### Option A: Local binaries + logging stack

1. Switch to the server module and run the local init script:

    ```bash
    cd server
    ./init-server-local.sh
    ```

    This script builds the Go binary into `server/dist/`, prepares data/log/socket directories, and starts the Loki + Alloy + Grafana logging stack via Docker Compose before launching the server in the foreground.

    **Note: Do not prefix this script with `sudo`**—the script already calls the individual setup helpers with elevated privileges where needed. Running the top-level script as root would cause all generated folders and files to be owned by `root`, making subsequent local development much harder to manage.

#### Option B: Full Docker Compose stack

1. From the same `server` directory run:

    ```bash
    cd server
    ./init-server-docker.sh
    ```

    The script ensures host bind-mount directories exist, exports your UID/GID for correct ownership, and then invokes `docker compose up --build` to start the containerised server alongside its logging dependencies.

> Prefer Option A when iterating on Go code locally; use Option B to validate the container stack or share an environment with teammates.

### Using the CLI

1. Navigate to the `cli` directory:

    ```bash
    cd cli
    ```

2. Use the `set`, `get`, or `delete` commands:

    ```bash
    go run . set mykey "hello world"
    go run . get mykey
    ```

### Running Benchmarks

1. Navigate to the project root.
2. Run the tests with the `-bench` flag:

    ```bash
    go test -bench . ./...
    ```

For more detailed information, please refer to the `README.md` file within each respective directory.
