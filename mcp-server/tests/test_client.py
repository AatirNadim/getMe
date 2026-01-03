import json

import httpx
import pytest

from getme_mcp_server.client import GetMeClient, GetMeError


class _Transport(httpx.BaseTransport):
    def __init__(self, handler):
        self._handler = handler

    def handle_request(self, request: httpx.Request) -> httpx.Response:
        return self._handler(request)


def test_get_happy_path(monkeypatch):
    def handler(req: httpx.Request) -> httpx.Response:
        assert req.method == "GET"
        assert req.url.path == "/get"
        assert req.url.params.get("key") == "a"
        return httpx.Response(200, text="hello")

    client = GetMeClient(socket_path="/tmp/x.sock")

    def _client(self):
        return httpx.Client(base_url=self.base_url, transport=_Transport(handler))

    monkeypatch.setattr(GetMeClient, "_client", _client)

    assert client.get("a") == "hello"


def test_get_json_invalid(monkeypatch):
    def handler(req: httpx.Request) -> httpx.Response:
        return httpx.Response(200, text="not-json")

    client = GetMeClient(socket_path="/tmp/x.sock")

    def _client(self):
        return httpx.Client(base_url=self.base_url, transport=_Transport(handler))

    monkeypatch.setattr(GetMeClient, "_client", _client)

    with pytest.raises(GetMeError):
        client.get_json("a")


def test_put_json_compacts(monkeypatch):
    captured = {}

    def handler(req: httpx.Request) -> httpx.Response:
        assert req.method == "POST"
        assert req.url.path == "/put"
        payload = json.loads(req.content.decode("utf-8"))
        captured["payload"] = payload
        return httpx.Response(200, text="ok")

    client = GetMeClient(socket_path="/tmp/x.sock")

    def _client(self):
        return httpx.Client(base_url=self.base_url, transport=_Transport(handler))

    monkeypatch.setattr(GetMeClient, "_client", _client)

    out = client.put_json("k", {"a": 1, "b": [2, 3]})
    assert out == "ok"
    assert captured["payload"]["key"] == "k"
    assert captured["payload"]["value"] == '{"a":1,"b":[2,3]}'
