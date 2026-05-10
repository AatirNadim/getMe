# getMe - High-Performance Key-Value Store

`getMe` is a persistent, embeddable, log-structured key-value store optimized for high write throughput and low-latency reads. 

This Docker image provides a complete, containerized environment for `getMe`, packaging the core storage server, the HTTP proxy, and the built-in command-line interface (CLI) into a single, lightweight image.

## Image Architecture & Multi-Stage Build

This image is built using a highly optimized **multi-stage Docker build** process defined in our `ContainerFile`:

1. **Build Stage (`golang:1.23.1-alpine`)**: The image compiles the core system (Server, CLI, and HTTP Proxy) from source. It gathers all essential local modules and builds statically linked Go binaries (`getMe-server`, `getMe-cli`, `getMe-proxy`).
2. **Final Stage (`alpine:latest`)**: The built binaries and entrypoint scripts are copied into a clean, minimal Alpine Linux base image. This keeps the final image size incredibly small and reduces the attack surface.

## Quick Start

Run the container in the background, exposing the HTTP proxy port and mounting a volume for data persistence:

```bash
docker run -d \
  --name getme-store \
  -p 8080:8080 \
  -v getme_data:/var/lib/getMeStore \
  your-dockerhub-username/getme:latest
```

## Using the Built-in CLI

The image comes with the `getMe` CLI pre-installed. We have configured the image to automatically load an alias (`getme-cli`), allowing you to interact with the database directly from your host using `docker exec`:

```bash
# Set a value
docker exec -it getme-store sh -ic "getme-cli set mykey 'hello world'"

# Get a value
docker exec -it getme-store sh -ic "getme-cli get mykey"
```
*(Note: the `-ic` flags are required to invoke an interactive shell that loads the alias configuration inside the container).*

## Volumes and Persistence

To ensure your data survives container restarts, you must mount volumes to the following directories inside the container:

* `/var/lib/getMeStore/dataDir`: The primary directory where the log-structured segments (database files) are stored.
* `/tmp/getMeStore/sockDir`: Used for internal Unix socket communication.
* `/tmp/getMeStore/dumpDir`: Used by the internal application logger.

## Security

Security is built-in by design. The container **does not run as root**. 

During the build process, an unprivileged user named `appuser` (along with `appgroup`) is created. All binaries are executed under this user profile, and the ownership of all critical data directories (`/var/lib/getMeStore` and `/tmp/getMeStore`) is automatically assigned to `appuser:appgroup`.

You can override the default UID and GID during the build phase using `build-args` if your environment requires specific user ID mappings:
```bash
docker build --build-arg UID=2000 --build-arg GID=2000 -t getme -f ContainerFile .
```

## Exposed Ports

* **`8080`**: The default port exposed by the HTTP proxy to handle incoming REST requests.

## Docker Compose Example

For an easier deployment, you can use `docker-compose.yml`:

```yaml
version: '3.8'

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

## 🔗 Links

* **GitHub Repository**: [**Visit here!**](https://github.com/AatirNadim/getMe)
* **SDKs Available**: Go, Java, JavaScript, Python
