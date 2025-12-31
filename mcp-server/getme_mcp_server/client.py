from __future__ import annotations

import json
from dataclasses import dataclass
from typing import Any

import httpx

from . import config


class GetMeError(RuntimeError):
    pass


@dataclass(frozen=True)
class GetMeClient:
    socket_path: str = config.DEFAULT_SOCKET_PATH
    base_url: str = config.DEFAULT_BASE_URL
    timeout_s: float = 10.0

    def _client(self) -> httpx.Client:
        transport = httpx.HTTPTransport(uds=self.socket_path)
        return httpx.Client(
            base_url=self.base_url, transport=transport, timeout=self.timeout_s
        )

    def _ensure_ok(self, resp: httpx.Response) -> str:
        if resp.status_code != 200:
            raise GetMeError(f"core returned {resp.status_code}: {resp.text}")
        return resp.text

    def get(self, key: str) -> str:
        if not key:
            raise ValueError("key must be non-empty")
        with self._client() as client:
            resp = client.get("/get", params={"key": key})
        return self._ensure_ok(resp)

    def get_json(self, key: str) -> Any:
        raw = self.get(key)
        try:
            return json.loads(raw)
        except json.JSONDecodeError as e:
            raise GetMeError(f"value for key '{key}' is not valid JSON") from e

    def put(self, key: str, value: str) -> str:
        if not key:
            raise ValueError("key must be non-empty")
        if value is None or value == "":
            raise ValueError("value must be non-empty")
        with self._client() as client:
            resp = client.post("/put", json={"key": key, "value": value})
        return self._ensure_ok(resp)

    def put_json(self, key: str, json_value: Any) -> str:
        if isinstance(json_value, str):
            # If a string is provided, require it to be valid JSON and compact it.
            try:
                parsed = json.loads(json_value)
            except json.JSONDecodeError as e:
                raise ValueError("json_value string must be valid JSON") from e
            compact = json.dumps(parsed, separators=(",", ":"))
        else:
            compact = json.dumps(json_value, separators=(",", ":"))
        return self.put(key, compact)

    def delete(self, key: str) -> str:
        if not key:
            raise ValueError("key must be non-empty")
        with self._client() as client:
            resp = client.delete("/delete", params={"key": key})
        return self._ensure_ok(resp)

    def clear(self) -> str:
        with self._client() as client:
            resp = client.delete("/clearStore")
        return self._ensure_ok(resp)

    def batch_put(self, pairs: dict[str, str]) -> str:
        if not isinstance(pairs, dict) or not pairs:
            raise ValueError("pairs must be a non-empty object/map")
        with self._client() as client:
            resp = client.post("/batch-put", json=pairs)
        return self._ensure_ok(resp)
