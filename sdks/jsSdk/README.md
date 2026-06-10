<div align="center">
  <img src="https://raw.githubusercontent.com/AatirNadim/getMe/main/getme-landing/public/extended-logo-rounded.png" alt="getMe Logo" style="width: 400px; max-width: 100%; margin-bottom: 20px;"/>
  
  <div style="margin-top: 10px;"><strong>JavaScript/TypeScript SDK for a High-Performance Key-Value Store</strong></div>
</div>

<!-- --- -->

## 📑 Index

- [Overview](#-overview)
- [About getMe Core Engine](#-about-getme-core-engine)
- [Server Dependency & Docker](#-server-dependency--docker)
  - [Running the Core Engine via Docker](#running-the-core-engine-via-docker)
  - [Connecting the SDK to the Core Engine](#connecting-the-sdk-to-the-core-engine)
  - [Why manage the core engine separately?](#why-manage-the-core-engine-separately)
- [Installation & Usage](#-installation--usage)
- [Other Available SDKs](#-other-available-sdks)
- [Links and Resources](#-links-and-resources)
- [License](#-license)

<!-- --- -->

## 📖 Overview

The official JavaScript/TypeScript client library for the **getMe** key-value store. This SDK abstracts the underlying API calls, providing Node.js developers with a simple, strongly-typed, and idiomatic interface to interact directly with the `getMe` runtime.

<!-- --- -->

## 🧠 About getMe Core Engine

`getMe` is a persistent, embeddable key-value store written in Go. It is heavily inspired by the design of Bitcask and is optimized for high write throughput and low-latency reads.

The engine runs on a few core principles:

- **Log-Structured Storage**: All data is written to an append-only log file, maximizing write speed by avoiding slow, random disk I/O.
- **In-Memory Hash Index**: A hash table is maintained in memory mapping each key to its exact on-disk location, enabling single disk seek lookups.
- **Compaction**: A background process that periodically removes stale data to reclaim disk space.
- **Batch Operations**: A batch API amortizes the cost of writes, allowing for extremely high ingestion rates.

For a deep dive into the architecture, please refer to the [Root README](https://github.com/AatirNadim/getMe/blob/main/README.md) and the [Server README](https://github.com/AatirNadim/getMe/blob/main/server/README.md).

<!-- --- -->

## 🚀 Server Dependency & Docker

**Important:** This SDK acts merely as a client. The `getMe` core server engine **must** be running independently for this library to function. The SDK interacts with the Core Engine through a Unix Domain Socket (UDS).

### Running the Core Engine via Docker

You can use the official Docker images mentioned in our [DockerHub README](https://github.com/AatirNadim/getMe/blob/main/DOCKERHUB.md).

For maximum performance, using the Core Server Only container (`ContainerFile.server`) is recommended, which isolates the database engine:

```bash
docker build -t getme-core -f ../../ContainerFile.server ../../
docker run -d \
  --name getme-engine \
  -v getme_data:/var/lib/getMeStore \
  -v /tmp/getMeStore/sockDir:/tmp/getMeStore/sockDir \
  getme-core
```

### Connecting the SDK to the Core Engine

The SDK communicates with the engine via the Unix Domain Socket (`getMe.sock`). Because the engine creates the socket in `/tmp/getMeStore/sockDir`, it is **critical** that this exact directory is bind-mounted when running the container so that your local Node.js application can reach it.

To connect your JS/TS SDK to the core engine:

1. Ensure the `getMe` server is running (either locally or via Docker with the `/tmp/getMeStore/sockDir` volume mounted).
2. Point your SDK initialization to the socket file location by setting the `GETME_SOCKET_PATH` environment variable (default: `/tmp/getMeStore/sockDir/getMe.sock`).

### Why manage the core engine separately?

Decoupling the database engine from the JavaScript logic provides several critical advantages:

- **Performance & Isolation**: The core engine heavily utilizes memory for its hash index and performs continuous background disk I/O for compaction. Running it separately prevents the storage engine from starving your Node.js application (and V8 garbage collector) of memory or compute resources (and vice-versa).
- **Independent Scaling**: You can scale your Node.js applications dynamically while connecting to a centralized, standalone `getMe` instance.
- **Polyglot Ecosystems**: A decoupled server can serve multiple applications simultaneously, even if those applications are built across entirely different technology stacks.

<!-- --- -->

## 💻 Installation & Usage

Install the SDK via npm, yarn, or pnpm:

```bash
npm install getme-js-sdk
# or
pnpm add getme-js-sdk
```

Next, ensure you tell the SDK where to find your locally mounted Unix socket by setting the `GETME_SOCKET_PATH` environment variable:

```bash
export GETME_SOCKET_PATH=/tmp/getMeStore/sockDir/getMe.sock
```

### Basic Example

```typescript
import { GetMeClient } from "getme-js-sdk";

const client = new GetMeClient(); // Connects automatically using GETME_SOCKET_PATH

async function run() {
  // Store a value
  await client.put("myKey", "myValue");

  // Retrieve the value
  const value = await client.get("myKey");
  console.log("Value:", value);
}

run();
```

<!-- --- -->

## 📦 Other Available SDKs

While this is the JavaScript/TypeScript client, `getMe` supports multiple languages. Detailed information on all available SDKs can be found in the [Root README](https://github.com/AatirNadim/getMe/blob/main/README.md).

- [Go SDK](https://github.com/AatirNadim/getMe/tree/main/sdks/goSdk)
- [Java SDK](https://github.com/AatirNadim/getMe/tree/main/sdks/javaSdk)
- [JavaScript/TypeScript SDK (This package)](https://github.com/AatirNadim/getMe/tree/main/sdks/jsSdk)
- [Python SDK](https://github.com/AatirNadim/getMe/tree/main/sdks/pythonSdk)

<!-- --- -->

## 🔗 Links and Resources

Below is a comprehensive list of resources and documentation linked throughout the `getMe` ecosystem:

- **[getMe Root Repository](https://github.com/AatirNadim/getMe)**
- **[Root README](https://github.com/AatirNadim/getMe/blob/main/README.md)**
- **[DockerHub README](https://github.com/AatirNadim/getMe/blob/main/DOCKERHUB.md)**
- **[Server Engine Source](https://github.com/AatirNadim/getMe/tree/main/server)**
- **[AGPLv3 License](https://github.com/AatirNadim/getMe/blob/main/LICENSE)**

### Engineering Blog Series

- **[Building getMe - Part I](https://techtom.hashnode.dev/building-getme-i)**
- **[Building getMe - Part II](https://techtom.hashnode.dev/building-getme-ii)**

<!-- --- -->

## 📄 License

This project is licensed under the GNU Affero General Public License v3.0 (AGPLv3) - see the [LICENSE](https://github.com/AatirNadim/getMe/blob/main/LICENSE) file for details.
