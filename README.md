<div align="center">
  <img src="./getme-landing/app/icon.png" width="150" alt="getMe Icon" style="vertical-align: middle; margin-right: 20px;"/>
  <span style="font-family: 'Mona Sans', sans-serif; font-size: 6em; font-weight: 800; margin-bottom: 0;vertical-align: middle;">getme</span>
  <div style="margin-top: 10px;"><strong>A High-Performance Key-Value Store</strong></div>
</div>

<!-- --- -->

## 📑 Index

- [Overview](#-overview)
- [Project Structure](#-project-structure)
- [Core Architecture](#-core-architecture)
- [Getting Started](#-getting-started)
  - [Running the Server](#running-the-server)
  - [Using the CLI](#using-the-cli)
  - [HTTP Proxy](#http-proxy)
  - [MCP Server](#mcp-server)
- [Running Benchmarks & Tests](#-running-benchmarks--tests)
- [SDKs](#-sdks)
- [License](#-license)

<!-- --- -->

## 📖 Overview

`getMe` is a persistent, embeddable key-value store written in Go. It is inspired by the design of Bitcask and is optimized for high write throughput and low-latency reads.

It uses a log-structured storage approach, ensuring that all data is appended sequentially. It uses Unix Domain Sockets (UDS) for incredibly fast local inter-process communication, alongside several interfaces like an HTTP proxy, a CLI, and a Model Context Protocol (MCP) server for LLMs.

<!-- --- -->

## 🏗 Project Structure

This project is a monorepo containing the core storage server, multiple client interfaces, and tools.

- **[`server/`](./server/)**: The core storage daemon and engine. Implements the log-structured hash table for persistent storage. See [`server/README.md`](./server/README.md) for architectural deep-dives.
- **[`cli/`](./cli/)**: A command-line interface for interacting with the `getMe` server for testing and debugging.
- **[`sdks/`](./sdks/)**: Client libraries (`goSdk`, `javaSdk`, `jsSdk`, `pythonSdk`) to integrate `getMe` into your applications.
- **[`http-proxy-go/`](./http-proxy-go/)**: An HTTP server built using the `goSdk` that exposes the core engine's Unix Domain Socket connection over standard HTTP routes.
- **[`mcp-server/`](./mcp-server/)**: A Model Context Protocol (MCP) server that exposes the `getMe` database as tools to Large Language Models (like Claude or Cursor).
- **[`commons/`](./commons/)**: Shared code, socket paths, types, and constants used across the monorepo to ensure consistency.
- **[`utils/`](./utils/)**: Shared utility packages, including logging stack configurations (Loki + Alloy + Grafana).

> **Spotlight:** The curated inner docs are the quickest way to understand the system end-to-end. Start with [`server/README.md`](./server/README.md) for architecture fundamentals, then explore the `cli` and `mcp-server` modules for integrations.

<!-- --- -->

## 🧠 Core Architecture

The storage engine relies on a few core principles:

- **Log-Structured Storage**: All data is written to an append-only log file. This makes writes extremely fast as it avoids slow, random disk I/O.
- **In-Memory Hash Index**: A hash table is kept in memory, mapping each key to the exact location of its value on disk. This allows for very fast read operations (typically one disk seek).
- **Compaction**: A background process that periodically cleans up old, stale data from the log files to reclaim disk space.
- **Fast Local Transport**: Communication is done predominantly via Unix Domain Sockets, avoiding standard TCP overhead locally.

<!-- --- -->

## 🚀 Getting Started

### Running the Server

The repository ships with helper scripts to bootstrap the environment.

#### Option A: Local binaries + logging stack

Switch to the server module and run the local init script:

```bash
cd server
./init-server-local.sh
```

This script builds the Go binary into `server/dist/`, prepares data/log/socket directories, and starts the Loki + Alloy + Grafana logging stack via Docker Compose before launching the server in the foreground.

> **Warning:** **Do not prefix this script with `sudo`**. It will invoke elevated privileges internally where needed. Using `sudo` at the top level causes permission errors for local development.

#### Option B: Full Docker Compose stack

From the same `server` directory run:

```bash
cd server
./init-server-docker.sh
```

This ensures host directories exist, exports your UID/GID, and invokes `docker compose up --build` to run everything in containers.

### Using the CLI

Interact directly with the local server:

```bash
cd cli
go run . put mykey "hello world"
go run . get mykey
go run . delete mykey
```

### HTTP Proxy

If you want standard HTTP REST endpoints instead of Unix Sockets, run the Go HTTP proxy:

```bash
cd http-proxy-go
go run main.go -port 8080
```

This will allow you to run `curl http://localhost:8080/get?key=mykey`.

### MCP Server

`getMe` can be used by LLM clients (like Claude Desktop) through the Model Context Protocol.

```bash
cd mcp-server
uv run getme-mcp-server
```

(See [`mcp-server/README.md`](./mcp-server/README.md) for configuration and integration instructions).

<!-- --- -->

## 📊 Running Benchmarks & Tests

To ensure no performance regressions or to stress test the database:

1. Navigate to the specific module (e.g., `server`).
2. Run standard tests:
   ```bash
   go test ./...
   ```
3. Run benchmarks:
   ```bash
   go test -bench . ./...
   ```
   (Note: For heavier stress/correctness testing, look into `server/tests/`).

<!-- --- -->

## 📦 SDKs

We provide SDKs across different languages. Find them in the `sdks/` directory:

- [**Go SDK**](./sdks/goSdk/)
- [**JavaScript / TypeScript SDK**](./sdks/jsSdk/)
- [**Python SDK**](./sdks/pythonSdk/)
- [**Java SDK**](./sdks/javaSdk/)

All SDKs interface directly with the Unix Domain Socket to provide optimal latency.

<!-- --- -->

## 📄 License

This project is licensed under the GNU Affero General Public License v3.0 (AGPLv3) - see the [LICENSE](LICENSE) file for details.
