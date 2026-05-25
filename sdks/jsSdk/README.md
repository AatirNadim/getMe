# getMe JavaScript/TypeScript SDK

The **getMe JS/TS SDK** is the official Node.js client for the **getMe Key-Value Store**. It provides a simple, strongly-typed interface to interact with the getMe engine, allowing you to easily store and retrieve data from your JavaScript and TypeScript applications.

## About the getMe Project

**getMe** is a fast, lightweight, and robust key-value store. It is designed to be highly reliable and operates using Unix Domain Sockets for ultra-high-performance inter-process communication (IPC). 

For a comprehensive overview of the architecture, design decisions, and core workings, please refer to the [Main getMe README](https://github.com/AatirNadim/getMe/blob/main/README.md).

## Ecosystem & Other SDKs

This SDK is part of a broader ecosystem. To provide flexibility across different technology stacks, we officially support multiple SDKs. You can find out more about the other available languages (like Go, Java, and Python) by visiting the [SDKs Overview README](https://github.com/AatirNadim/getMe/tree/main/sdks).

## Why a Standalone Core Engine?

Unlike embedded database libraries (like SQLite), the getMe SDK **does not run the database engine itself.** The core engine must be running independently for the SDK to function. 

Managing the core engine separately is highly advantageous because it:
- **Provides Better Resource Isolation:** Your Node.js application and the database engine do not compete for the same memory heap and garbage collector resources.
- **Enables True Microservices:** Multiple discrete applications, written in completely different languages using varying SDKs, can safely connect to the exact same data store simultaneously.
- **Simplifies Scaling and Deployment:** You can scale, monitor, restart, or update your application containers without taking down your underlying database.

## Running the Core Engine

To use this SDK, you must first have the getMe core engine running. For ease of use, we package the engine as a lightweight, secure Docker container. 

**Standalone Server Image:** [aatir0docking/getme.server](https://hub.docker.com/r/aatir0docking/getme.server/)

### Deployment Requirements

The getMe engine operates by securely listening on a native Unix Domain Socket rather than an open TCP port. For host machines (or other containers) to communicate with it, this socket directory **must be mounted as a Docker Volume**.

Run the core engine container using the following command:

```bash
# Create a local directory for the socket mapping
mkdir -p /tmp/my-getme-sockets

# Run the container and mount the socket directory
docker run -d \
  --name getme-server \
  -v /tmp/my-getme-sockets:/tmp/getMeStore/sockDir \
  aatir0docking/getme.server:latest
```

Once running, the core engine will listen for connections at `/tmp/my-getme-sockets/getMe.sock` on your host machine.

## Installation & Usage

Install the SDK via npm, yarn, or pnpm:

```bash
npm install getme-js-sdk
# or
pnpm add getme-js-sdk
```

Next, ensure you tell the SDK where to find your locally mounted Unix socket by setting the `GETME_SOCKET_PATH` environment variable:

```bash
export GETME_SOCKET_PATH=/tmp/my-getme-sockets/getMe.sock
```

### Basic Example

```typescript
import { GetMeClient } from 'getme-js-sdk';

const client = new GetMeClient(); // Connects automatically using GETME_SOCKET_PATH

async function run() {
  // Store a value
  await client.put('myKey', 'myValue');

  // Retrieve the value
  const value = await client.get('myKey');
  console.log('Value:', value);
}

run();
```

## Useful Links
- [Main getMe Repository README](https://github.com/AatirNadim/getMe/blob/main/README.md)
- [getMe SDKs Documentation](https://github.com/AatirNadim/getMe/tree/main/sdks)
- [getMe Server Docker Hub Image](https://hub.docker.com/r/aatir0docking/getme.server/)
- [Blog Part I - Building getMe](https://techtom.hashnode.dev/building-getme-i)
- [Blog Part II - Building getMe](https://techtom.hashnode.dev/building-getme-ii)
