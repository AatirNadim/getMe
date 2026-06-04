# getMe MCP Server

An MCP server that exposes getMe key-value operations as MCP tools, and talks to the getMe core application over **HTTP via a Unix Domain Socket (UDS)**.

## Features & Safety

- **Read-Only Mode:** Ensure LLMs don't mutate data via `GETME_READ_ONLY=true`.
- **Destructive Operation Guards:** The `clear` tool is disabled by default, ensuring agents can't accidentally drop the store.
- **Payload & Connection Limits:** Defends against massive payloads/batches and re-uses HTTP connections for performance.
- **Prefix Isolation:** Restrict the LLM to a specific namespace in the DB with `GETME_KEY_PREFIX`.
- **Masked Logging:** Tool execution metadata is logged, but payloads are masked to prevent leaking secrets.

## Prereqs

- `uv` installed
- getMe core server running and listening on the Unix socket (default: `/tmp/getMeStore/sockDir/getMe.sock`)

## Install

```bash
cd mcp-server
uv sync
```

## Fixing VS Code "package not installed" warnings

If Pylance shows warnings like `Package "httpx" is not installed in the selected environment`, VS Code is using a different Python interpreter than the `uv` virtualenv.

- Open Command Palette → `Python: Select Interpreter`
- Select: `getMe/mcp-server/.venv/bin/python`

This repo also includes a workspace setting that points the interpreter at `mcp-server/.venv`.

## Run

```bash
cd mcp-server
# optional override
export GETME_SOCKET_PATH=/tmp/getMeStore/sockDir/getMe.sock
uv run getme-mcp-server
```

## Docker

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

By default, the MCP server runs over **stdio** (the typical MCP deployment model).

## Configuration

Control safety rules and performance via environment variables:

- `GETME_SOCKET_PATH` (default: `/tmp/getMeStore/sockDir/getMe.sock`): Path to the UNIX socket.
- `GETME_READ_ONLY` (default: `false`): If `true`, only `get` and `get_json` tools are registered.
- `GETME_ALLOW_CLEAR` (default: `false`): If `true`, enables the dangerous `clear` tool.
- `GETME_KEY_PREFIX` (default: `""`): String prefix prepended to all keys. Sandboxes the LLM (e.g. `agent1:`).
- `GETME_MAX_VALUE_SIZE_BYTES` (default: `5242880` / 5MB): Limit on single value payload sizes.
- `GETME_MAX_BATCH_ITEMS` (default: `100`): Maximum keys processed in one `batch_put`.

## Tools

- `get(key) -> str`
- `get_json(key) -> object`
- `put(key, value) -> str`
- `put_json(key, json_value) -> str`
- `delete(key) -> str`
- `clear() -> str`
- `batch_put(pairs: object) -> str`

