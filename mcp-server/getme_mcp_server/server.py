from __future__ import annotations

from typing import Any

from mcp.server.fastmcp import FastMCP

from .client import GetMeClient
from .config import base_url, socket_path


def build_mcp() -> FastMCP:
    mcp = FastMCP("getMe")
    client = GetMeClient(socket_path=socket_path(), base_url=base_url())

    @mcp.tool()
    def get(key: str) -> str:
        """Get a value by key."""
        return client.get(key)

    @mcp.tool()
    def get_json(key: str) -> Any:
        """Get a value by key and parse it as JSON."""
        return client.get_json(key)

    @mcp.tool()
    def put(key: str, value: str) -> str:
        """Put a (key, value) pair."""
        return client.put(key, value)

    @mcp.tool()
    def put_json(key: str, json_value: Any) -> str:
        """Put a key with a JSON value (object/array/string). Stored compacted."""
        return client.put_json(key, json_value)

    @mcp.tool()
    def delete(key: str) -> str:
        """Delete a key."""
        return client.delete(key)

    @mcp.tool()
    def clear() -> str:
        """Clear the entire store."""
        return client.clear()

    @mcp.tool()
    def batch_put(pairs: dict[str, str]) -> str:
        """Batch put from a map of key -> value."""
        return client.batch_put(pairs)

    return mcp
