from __future__ import annotations

import logging
from typing import Any

from mcp.server.fastmcp import FastMCP

from .client import GetMeClient
from . import config


def setup_logging():
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    )


def build_mcp() -> FastMCP:
    setup_logging()
    logger = logging.getLogger("getme.server")

    mcp = FastMCP("getMe")
    client = GetMeClient(socket_path=config.socket_path(), base_url=config.base_url())

    @mcp.tool()
    async def get(key: str) -> str:
        """Get a value by key."""
        logger.info("Executing get for key: %s", key)
        return await client.get(key)

    @mcp.tool()
    async def get_json(key: str) -> Any:
        """Get a value by key and parse it as JSON."""
        logger.info("Executing get_json for key: %s", key)
        return await client.get_json(key)

    @mcp.tool()
    async def batch_get(keys: list[str]) -> Any:
        """Batch get multiple values by keys."""
        logger.info("Executing batch_get for %d keys", len(keys))
        return await client.batch_get(keys)

    # Only register write tools if the store is not read-only. This allows the same server code to be used for both read-only and read-write stores, with the appropriate tools exposed based on configuration.
    if not config.is_read_only():

        @mcp.tool()
        async def put(key: str, value: str) -> str:
            """Put a (key, value) pair."""
            # Mask value in logs
            logger.info("Executing put for key: %s (value length: %d)", key, len(value))
            return await client.put(key, value)

        @mcp.tool()
        async def put_json(key: str, json_value: Any) -> str:
            """Put a key with a JSON value (object/array/string). Stored compacted."""
            logger.info("Executing put_json for key: %s", key)
            return await client.put_json(key, json_value)

        @mcp.tool()
        async def delete(key: str) -> str:
            """Delete a key."""
            logger.info("Executing delete for key: %s", key)
            return await client.delete(key)

        @mcp.tool()
        async def batch_delete(keys: list[str]) -> Any:
            """Batch delete multiple keys."""
            logger.info("Executing batch_delete for %d keys", len(keys))
            return await client.batch_delete(keys)

        @mcp.tool()
        async def batch_put(pairs: dict[str, str]) -> Any:
            """Batch put from a map of key -> value."""
            logger.info("Executing batch_put for %d keys", len(pairs))
            return await client.batch_put(pairs)

        if config.allow_clear():

            @mcp.tool()
            async def clear() -> str:
                """Clear the entire store."""
                logger.warning("Executing clear store")
                return await client.clear()

    return mcp
