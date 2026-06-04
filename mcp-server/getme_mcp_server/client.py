from __future__ import annotations

import json
import logging
from typing import Any

import httpx

from . import config

logger = logging.getLogger("getme.client")


class GetMeError(RuntimeError):
    pass


class GetMeClient:
    def __init__(
        self,
        socket_path: str = config.DEFAULT_SOCKET_PATH,
        base_url: str = config.DEFAULT_BASE_URL,
        timeout_s: float = 10.0,
    ) -> None:
        self.socket_path = socket_path
        self.base_url = base_url
        self.timeout_s = timeout_s
        self.prefix = config.key_prefix()
        self.max_key_length = config.max_key_length()
        self.max_value_size = config.max_value_size_bytes()
        self.max_batch_items = config.max_batch_items()

        # AsyncClient initialized lazily
        self._client: httpx.AsyncClient | None = None

    def get_client(self) -> httpx.AsyncClient:
        if self._client is None:
            transport = httpx.AsyncHTTPTransport(uds=self.socket_path)
            self._client = httpx.AsyncClient(
                base_url=self.base_url, transport=transport, timeout=self.timeout_s
            )
        return self._client

    async def close(self) -> None:
        if self._client is not None:
            await self._client.aclose()
            self._client = None

    def _ensure_ok(self, resp: httpx.Response) -> str:
        if resp.status_code != 200:
            raise GetMeError(f"core returned {resp.status_code}: {resp.text}")
        return resp.text

    def _check_key(self, key: str) -> str:
        if not key:
            raise ValueError("key must be non-empty")
        if len(key) > self.max_key_length:
            raise ValueError(f"key exceeds max length of {self.max_key_length}")
        return f"{self.prefix}{key}"

    def _check_value(self, value: str) -> None:
        if value is None or value == "":
            raise ValueError("value must be non-empty")
        if len(value.encode("utf-8")) > self.max_value_size:
            raise ValueError(f"value exceeds max size of {self.max_value_size} bytes")

    async def get(self, key: str) -> str:
        full_key = self._check_key(key)
        client = self.get_client()
        try:
            resp = await client.get("/get", params={"key": full_key})
            return self._ensure_ok(resp)
        except httpx.RequestError as e:
            raise GetMeError(f"Database backend unavailable: {e}") from e

    async def get_json(self, key: str) -> Any:
        raw = await self.get(key)
        try:
            return json.loads(raw)
        except json.JSONDecodeError as e:
            raise GetMeError(f"value for key '{key}' is not valid JSON") from e

    async def put(self, key: str, value: str) -> str:
        full_key = self._check_key(key)
        self._check_value(value)
        client = self.get_client()
        try:
            resp = await client.post("/put", json={"key": full_key, "value": value})
            return self._ensure_ok(resp)
        except httpx.RequestError as e:
            raise GetMeError(f"Database backend unavailable: {e}") from e

    async def put_json(self, key: str, json_value: Any) -> str:
        if isinstance(json_value, str):
            # If a string is provided, require it to be valid JSON and compact it.
            try:
                parsed = json.loads(json_value)
            except json.JSONDecodeError as e:
                raise ValueError("json_value string must be valid JSON") from e
            compact = json.dumps(parsed, separators=(",", ":"))
        else:
            compact = json.dumps(json_value, separators=(",", ":"))
        return await self.put(key, compact)

    async def delete(self, key: str) -> str:
        full_key = self._check_key(key)
        client = self.get_client()
        try:
            resp = await client.delete("/delete", params={"key": full_key})
            return self._ensure_ok(resp)
        except httpx.RequestError as e:
            raise GetMeError(f"Database backend unavailable: {e}") from e

    async def clear(self) -> str:
        client = self.get_client()
        try:
            resp = await client.delete("/clearStore")
            return self._ensure_ok(resp)
        except httpx.RequestError as e:
            raise GetMeError(f"Database backend unavailable: {e}") from e

    async def batch_put(self, pairs: dict[str, str]) -> str:
        if not isinstance(pairs, dict) or not pairs:
            raise ValueError("pairs must be a non-empty object/map")
        if len(pairs) > self.max_batch_items:
            raise ValueError(f"batch exceeds max items of {self.max_batch_items}")

        processed_pairs = {}
        for k, v in pairs.items():
            full_key = self._check_key(k)
            self._check_value(v)
            processed_pairs[full_key] = v

        client = self.get_client()
        try:
            resp = await client.post("/batch-put", json=processed_pairs)
            return self._ensure_ok(resp)
        except httpx.RequestError as e:
            raise GetMeError(f"Database backend unavailable: {e}") from e
