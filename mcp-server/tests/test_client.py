import json

import httpx
import pytest

from getme_mcp_server.client import GetMeClient, GetMeError


class _AsyncTransport(httpx.AsyncBaseTransport):
    def __init__(self, handler):
        self._handler = handler

    async def handle_async_request(self, request: httpx.Request) -> httpx.Response:
        return self._handler(request)


@pytest.mark.asyncio
async def test_get_happy_path(monkeypatch):
    def handler(req: httpx.Request) -> httpx.Response:
        assert req.method == "GET"
        assert req.url.path == "/get"
        assert req.url.params.get("key") == "a"
        return httpx.Response(200, text="hello")

    client = GetMeClient(socket_path="/tmp/x.sock")

    def get_client(self):
        return httpx.AsyncClient(
            base_url=self.base_url, transport=_AsyncTransport(handler)
        )

    monkeypatch.setattr(GetMeClient, "get_client", get_client)

    assert await client.get("a") == "hello"


@pytest.mark.asyncio
async def test_get_json_invalid(monkeypatch):
    def handler(req: httpx.Request) -> httpx.Response:
        return httpx.Response(200, text="not-json")

    client = GetMeClient(socket_path="/tmp/x.sock")

    def get_client(self):
        return httpx.AsyncClient(
            base_url=self.base_url, transport=_AsyncTransport(handler)
        )

    monkeypatch.setattr(GetMeClient, "get_client", get_client)

    with pytest.raises(GetMeError):
        await client.get_json("a")


@pytest.mark.asyncio
async def test_put_json_compacts(monkeypatch):
    captured = {}

    def handler(req: httpx.Request) -> httpx.Response:
        assert req.method == "POST"
        assert req.url.path == "/put"
        payload = json.loads(req.content.decode("utf-8"))
        captured["payload"] = payload
        return httpx.Response(200, text="ok")

    client = GetMeClient(socket_path="/tmp/x.sock")

    def get_client(self):
        return httpx.AsyncClient(
            base_url=self.base_url, transport=_AsyncTransport(handler)
        )

    monkeypatch.setattr(GetMeClient, "get_client", get_client)

    out = await client.put_json("k", {"a": 1, "b": [2, 3]})
    assert out == "ok"
    assert captured["payload"]["key"] == "k"
    assert captured["payload"]["value"] == '{"a":1,"b":[2,3]}'


@pytest.mark.asyncio
async def test_batch_get_happy_path(monkeypatch):
    def handler(req: httpx.Request) -> httpx.Response:
        assert req.method == "POST"
        assert req.url.path == "/batch-get"
        payload = json.loads(req.content.decode("utf-8"))
        assert payload["keys"] == ["a"]
        resp_body = {"found": {"a": "hello"}, "notFound": [], "errors": {}}
        return httpx.Response(200, text=json.dumps(resp_body))

    client = GetMeClient(socket_path="/tmp/x.sock")

    def get_client(self):
        return httpx.AsyncClient(
            base_url=self.base_url, transport=_AsyncTransport(handler)
        )

    monkeypatch.setattr(GetMeClient, "get_client", get_client)

    res = await client.batch_get(["a"])
    assert res["found"]["a"] == "hello"


@pytest.mark.asyncio
async def test_batch_put_unprefix(monkeypatch):
    def handler(req: httpx.Request) -> httpx.Response:
        assert req.method == "POST"
        assert req.url.path == "/batch-put"
        resp_body = {"successful": 0, "failed": {"test:a": "err"}}
        return httpx.Response(200, text=json.dumps(resp_body))

    client = GetMeClient(socket_path="/tmp/x.sock")
    client.prefix = "test:"

    def get_client(self):
        return httpx.AsyncClient(
            base_url=self.base_url, transport=_AsyncTransport(handler)
        )

    monkeypatch.setattr(GetMeClient, "get_client", get_client)

    res = await client.batch_put({"a": "hello"})
    assert res["failed"]["a"] == "err"


@pytest.mark.asyncio
async def test_batch_get_unprefix(monkeypatch):
    def handler(req: httpx.Request) -> httpx.Response:
        assert req.method == "POST"
        assert req.url.path == "/batch-get"
        resp_body = {
            "found": {"test:a": "hello"},
            "notFound": ["test:b"],
            "errors": {"test:c": "err"},
        }
        return httpx.Response(200, text=json.dumps(resp_body))

    client = GetMeClient(socket_path="/tmp/x.sock")
    client.prefix = "test:"

    def get_client(self):
        return httpx.AsyncClient(
            base_url=self.base_url, transport=_AsyncTransport(handler)
        )

    monkeypatch.setattr(GetMeClient, "get_client", get_client)

    res = await client.batch_get(["a", "b", "c"])
    assert res["found"]["a"] == "hello"
    assert res["notFound"][0] == "b"
    assert res["errors"]["c"] == "err"
