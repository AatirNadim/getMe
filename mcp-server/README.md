# getMe MCP Server

An MCP server that exposes getMe key-value operations as MCP tools, and talks to the getMe core application over **HTTP via a Unix Domain Socket (UDS)**.

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

By default, the MCP server runs over **stdio** (the typical MCP deployment model).

## Tools

- `get(key) -> str`
- `get_json(key) -> object`
- `put(key, value) -> str`
- `put_json(key, json_value) -> str`
- `delete(key) -> str`
- `clear() -> str`
- `batch_put(pairs: object) -> str`

