<div align="center">
  <img src="https://raw.githubusercontent.com/AatirNadim/getMe/main/getme-landing/public/extended-logo-rounded.png" alt="getMe Logo" style="width: 400px; max-width: 100%; margin-bottom: 20px;"/>
  
  <div style="margin-top: 10px;"><strong>Container Images for a High-Performance Key-Value Store</strong></div>
</div>

<!-- --- -->

## 📑 Index

- [Overview](#-overview)
- [Available Docker Images](#-available-docker-images)
- [All-in-One Container (HTTP Proxy, CLI, Server)](#-all-in-one-container-http-proxy-cli-server)
  - [Quick Start](#quick-start)
  - [Using the Built-in CLI](#using-the-built-in-cli)
  - [Docker Compose Example](#docker-compose-example)
- [MCP Server Container](#-mcp-server-container)
- [Core Server Only Container](#-core-server-only-container)
- [Volumes and Persistence](#-volumes-and-persistence)
- [Security & Customization](#-security--customization)
- [Links](#-links)

<!-- --- -->

## 📖 Overview

`getme` is a persistent, embeddable, log-structured key-value store optimized for high write throughput and low-latency reads. 

This document covers the different containerized environments available for `getme`, allowing you to run the database seamlessly across various setups.

<!-- --- -->

## 📦 Available Docker Images

`getme` provides multiple Docker architectures depending on your use case:

1. **`ContainerFile` (All-in-one)**: Packages the Core Server, HTTP Proxy, and CLI into a single lightweight image.
2. **`ContainerFile.server` (Core Engine)**: Runs *only* the core storage engine over Unix Domain Sockets (UDS).
3. **`mcp-server/Dockerfile`**: A specialized Python-based container exposing the `getme` database as tools for Large Language Models using the Model Context Protocol.

<!-- --- -->

## 🚀 All-in-One Container (HTTP Proxy, CLI, Server)

Built from the `ContainerFile`, this is the standard image. It runs the storage engine internally and exposes it over port `8080` via the HTTP proxy, alongside pre-installing the CLI.

### Quick Start

Run the container in the background, exposing the HTTP proxy port and mounting a volume for data persistence:

```bash
docker run -d \
  --name getme-store \
  -p 8080:8080 \
  -v getme_data:/var/lib/getMeStore \
  -v getme_tmp:/tmp/getMeStore \
  your-dockerhub-username/getme:latest
```

### Using the Built-in CLI

The image comes with the `getme` CLI pre-installed. The image is configured to automatically load an alias (`getme-cli`), allowing you to interact with the database directly from your host using `docker exec`:

```bash
# Set a value
docker exec -it getme-store sh -ic "getme-cli set mykey 'hello world'"

# Get a value
docker exec -it getme-store sh -ic "getme-cli get mykey"
```

_(Note: the `-ic` flags are required to invoke an interactive shell that loads the alias configuration inside the container)._

### Docker Compose Example

For an easier deployment, use a `docker-compose.yml`:

```yaml
version: "3.8"

services:
  getme:
    image: your-dockerhub-username/getme:latest
    container_name: getme-store
    ports:
      - "8080:8080"
    volumes:
      - getme_data:/var/lib/getMeStore
      - getme_tmp:/tmp/getMeStore
    restart: unless-stopped

volumes:
  getme_data:
  getme_tmp:
```

<!-- --- -->

## 🤖 MCP Server Container

If you are using `getme` with an LLM agent (like Claude Desktop or Cursor), deploy the MCP server container defined in `mcp-server/Dockerfile`.

This container requires access to the Unix Domain Socket file of the running `getme` server.

```bash
cd mcp-server
docker compose run --rm -T getme-mcp-server
```

This communicates over standard I/O (JSON-RPC) avoiding TTY allocation (`-T`). Note that it mounts `/tmp/getMeStore/sockDir` from the host.

<!-- --- -->

## ⚡ Core Server Only Container

If you only need the raw storage engine and are interacting directly via UDS (without HTTP), you can build `ContainerFile.server`.

```bash
docker build -t getme.server -f ContainerFile.server .
docker run -d \
  --name getme-engine \
  -v getme_data:/var/lib/getMeStore \
  -v /tmp/getMeStore/sockDir:/tmp/getMeStore/sockDir \
  getme.server
```
**Crucial:** `/tmp/getMeStore/sockDir` MUST be bind-mounted so external processes (SDKs) can reach the `getMe.sock`.

<!-- --- -->

## 💾 Volumes and Persistence

To ensure your data survives container restarts, you must mount volumes to the following directories inside the container:

- `/var/lib/getMeStore/dataDir`: The primary directory where the log-structured segments (database files) are stored.
- `/tmp/getMeStore/sockDir`: Used for internal Unix Domain Socket communication.
- `/tmp/getMeStore/dumpDir`: Used by the internal application logger.

<!-- --- -->

## 🛡️ Security & Customization

Security is built-in by design. The containers **do not run as root**.

During the multi-stage build processes, an unprivileged user named `appuser` (along with `appgroup`) is created. All binaries are executed under this user profile, and the ownership of all critical data directories is automatically assigned to `appuser:appgroup`.

You can override the default UID and GID during the build phase using `build-args` if your environment requires specific user ID mappings to avoid local permission issues:

```bash
docker build --build-arg UID=2000 --build-arg GID=2000 -t getme -f ContainerFile .
```

<!-- --- -->

## 🔗 Links

- **GitHub Repository**: [**Visit here!**](https://github.com/AatirNadim/getMe)
- **Blog Part I - Building getme**: [**Read here!**](https://techtom.hashnode.dev/building-getme-i)
- **Blog Part II - Building getme**: [**Read here!**](https://techtom.hashnode.dev/building-getme-ii)
- **SDKs Available**: [Go](https://github.com/AatirNadim/getMe/releases?q=gosdk&expanded=true), [Java](https://central.sonatype.com/artifact/io.github.aatirnadim/getme-javasdk), [JavaScript](https://www.npmjs.com/package/getme-js-sdk), [Python](https://pypi.org/p/getme-python-sdk)
