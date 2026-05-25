# getMe Python SDK

The official Python client library for the **getMe** key-value store. This SDK abstracts the underlying HTTP API calls, providing Python developers with a simple, idiomatic interface to interact with a `getMe` runtime.

## About getMe

`getMe` is a persistent, embeddable key-value store written in Go. It is heavily inspired by the design of Bitcask and is optimized for high write throughput and low-latency reads. 

The engine runs on a few core principles:
- **Log-Structured Storage**: All data is written to an append-only log file, maximizing write speed by avoiding slow, random disk I/O.
- **In-Memory Hash Index**: A hash table is maintained in memory mapping each key to its exact on-disk location, enabling single disk seek lookups.
- **Compaction**: A background process that periodically removes stale data to reclaim disk space.
- **Batch Operations**: A batch API amortizes the cost of writes, allowing for extremely high ingestion rates.

For a deep dive into the architecture, please refer to the [Root README](https://github.com/AatirNadim/getMe/blob/main/README.md) and the [Server README](https://github.com/AatirNadim/getMe/blob/main/server/README.md).

## Server Dependency

**Important:** This SDK acts merely as a client. The `getMe` core server engine **must** be running independently for this library to function. 

### Why manage the core engine separately?

Decoupling the database engine from the application logic provides several critical advantages:
- **Performance & Isolation**: The core engine heavily utilizes memory for its hash index and performs continuous background disk I/O for compaction. Running it separately prevents the storage engine from starving your Python application of memory or compute resources (and vice-versa).
- **Independent Scaling**: You can scale your Python applications (e.g., across multiple workers or containers) dynamically while connecting to a centralized, standalone `getMe` instance.
- **Polyglot Ecosystems**: A decoupled server can serve multiple applications simultaneously, even if those applications are built across entirely different technology stacks.

## Other Available SDKs

While this is the Python client, `getMe` supports multiple languages. Detailed information on all available SDKs can be found in the [SDKs Hub README](https://github.com/AatirNadim/getMe/blob/main/sdks/README.md). 

The available official SDKs are:
- [Go SDK](https://github.com/AatirNadim/getMe/tree/main/sdks/goSdk)
- [Java SDK](https://github.com/AatirNadim/getMe/tree/main/sdks/javaSdk)
- [JavaScript/TypeScript SDK](https://github.com/AatirNadim/getMe/tree/main/sdks/jsSdk)
- [Python SDK (This package)](https://github.com/AatirNadim/getMe/tree/main/sdks/pythonSdk)

## Links and Resources

Below is a comprehensive list of resources and documentation linked throughout the `getMe` ecosystem:

- **[getMe Root Repository](https://github.com/AatirNadim/getMe)**
- **[Root README](https://github.com/AatirNadim/getMe/blob/main/README.md)**
- **[Server Engine Source](https://github.com/AatirNadim/getMe/tree/main/server)** | **[Server README](https://github.com/AatirNadim/getMe/blob/main/server/README.md)**
- **[CLI Tool](https://github.com/AatirNadim/getMe/tree/main/cli)** | **[CLI README](https://github.com/AatirNadim/getMe/blob/main/cli/README.md)**
- **[SDKs Hub](https://github.com/AatirNadim/getMe/tree/main/sdks)** | **[SDKs README](https://github.com/AatirNadim/getMe/blob/main/sdks/README.md)**
- **[Benchmarking Suite](https://github.com/AatirNadim/getMe/tree/main/benchmarking)** | **[Benchmarking README](https://github.com/AatirNadim/getMe/blob/main/benchmarking/README.md)**
- **[Shared Utils Package](https://github.com/AatirNadim/getMe/tree/main/utils)**
- **[AGPLv3 License](https://github.com/AatirNadim/getMe/blob/main/LICENSE)**

### Engineering Blog Series
If you are interested in the internal workings, design decisions, and the journey of building the `getMe` storage engine from scratch, check out our engineering blog series:
- **[Building getMe - Part I](https://techtom.hashnode.dev/building-getme-i)**
- **[Building getMe - Part II](https://techtom.hashnode.dev/building-getme-ii)**

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPLv3) - see the [LICENSE](https://github.com/AatirNadim/getMe/blob/main/LICENSE) file for details.
