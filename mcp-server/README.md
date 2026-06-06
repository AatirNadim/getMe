# getMe MCP Server

An MCP server that exposes `getMe` key-value operations as Model Context Protocol (MCP) tools, enabling LLMs to talk to the getMe core database over **HTTP via a Unix Domain Socket (UDS)**.

<!-- mcp-name: io.github.AatirNadim/getme-mcp-server -->

## 🔗 Quick Links

- **GitHub Repository**: [AatirNadim/getMe](https://github.com/AatirNadim/getMe)
- **PyPI Package**: [getme-mcp-server](https://pypi.org/project/getme-mcp-server/)
- **Core Database Docs**: [Architecture & Deep Dive](https://github.com/AatirNadim/getMe/tree/main/server)
- **DockerHub Image**: [getme](https://hub.docker.com/r/aatir0docking/getme)
- **Blog Part I - Building getMe**: [Read here!](https://techtom.hashnode.dev/building-getme-i)
- **Blog Part II - Building getMe**: [Read here!](https://techtom.hashnode.dev/building-getme-ii)

## Using the MCP Server

The `getme-mcp-server` package is published on PyPI. For integration with MCP-compatible clients like Claude Desktop, Cursor, or the MCP Inspector, you can seamlessly run it using `uvx` or `pipx`.

### Claude Desktop Configuration

Add the following to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "getme": {
      "command": "uvx",
      "args": ["getme-mcp-server"],
      "env": {
        "GETME_SOCKET_PATH": "/tmp/getMeStore/sockDir/getMe.sock",
        "GETME_KEY_PREFIX": "agent-context:",
        "GETME_READ_ONLY": "false"
      }
    }
  }
}
```

## Features & Safety

- **Read-Only Mode:** Ensure LLMs don't mutate data via `GETME_READ_ONLY=true`.
- **Destructive Operation Guards:** The `clear` tool is disabled by default, ensuring agents can't accidentally drop the store.
- **Payload & Connection Limits:** Defends against massive payloads/batches and re-uses HTTP connections for performance.
- **Prefix Isolation:** Restrict the LLM to a specific namespace in the DB with `GETME_KEY_PREFIX`.
- **Masked Logging:** Tool execution metadata is logged, but payloads are masked to prevent leaking secrets.

## Configuration (Environment Variables)

Control safety rules and performance via environment variables.

- `GETME_SOCKET_PATH` (default: `/tmp/getMeStore/sockDir/getMe.sock`): Path to the UNIX socket used by the core getMe daemon.
- `GETME_READ_ONLY` (default: `false`): If `true`, only `get` and `get_json` tools are registered.
- `GETME_ALLOW_CLEAR` (default: `false`): If `true`, enables the dangerous `clear` tool.
- `GETME_KEY_PREFIX` (default: `""`): String prefix prepended to all keys. Sandboxes the LLM (e.g., `agent1:`).
- `GETME_MAX_VALUE_SIZE_BYTES` (default: `5242880` / 5MB): Limit on single value payload sizes.
- `GETME_MAX_BATCH_ITEMS` (default: `100`): Maximum keys processed in one `batch_put`.

## Best Practices

1. **Sandbox with Prefixes**: Always use `GETME_KEY_PREFIX` when handing the database over to an LLM. For instance, using `GETME_KEY_PREFIX="claude:"` ensures the model can't overwrite core application data.
2. **Restrict Privileges**: If the agent only needs context (e.g., fetching RAG documents or reading app state), enforce `GETME_READ_ONLY=true`.
3. **Never Allow Clear**: Keep `GETME_ALLOW_CLEAR=false` (the default) in production environments to prevent catastrophic drops caused by hallucinated commands.
4. **Volume Mounts**: Make sure the environment running the MCP server has adequate read/write permissions to the `GETME_SOCKET_PATH`.

## Tools

Once connected, the LLM has access to the following operations:

- `get(key) -> str`
- `get_json(key) -> object`
- `put(key, value) -> str`
- `put_json(key, json_value) -> str`
- `delete(key) -> str`
- `clear() -> str` _(Requires `GETME_ALLOW_CLEAR=true`)_
- `batch_put(pairs: object) -> object`
- `batch_get(keys: list) -> object`
- `batch_delete(keys: list) -> object`

## Development & Installation

### Prereqs

- `uv` installed
- `getMe` core server running and listening on the Unix socket

### Install from Source

```bash
cd mcp-server
uv sync
```

### Fixing VS Code "package not installed" warnings

If Pylance shows warnings like `Package "httpx" is not installed in the selected environment`, VS Code is using a different Python interpreter than the `uv` virtualenv.

- Open Command Palette → `Python: Select Interpreter`
- Select: `getMe/mcp-server/.venv/bin/python`

This repo also includes a workspace setting that points the interpreter at `mcp-server/.venv`.

### Run Locally

```bash
cd mcp-server
# optional override
export GETME_SOCKET_PATH=/tmp/getMeStore/sockDir/getMe.sock
uv run getme-mcp-server
```

### Docker deployment

This MCP server speaks **stdio** (JSON-RPC), so when running in Docker you must avoid allocating a TTY.

Build the image:

```bash
cd mcp-server
docker compose build
```

Run it over stdio (recommended for MCP clients):

```bash
cd mcp-server
docker compose run --rm -T getme-mcp-server
```

The compose file bind-mounts the host UDS directory:

- host: `/tmp/getMeStore/sockDir`
- container: `/tmp/getMeStore/sockDir`

So the MCP server can reach the core getMe server via `GETME_SOCKET_PATH=/tmp/getMeStore/sockDir/getMe.sock`.
