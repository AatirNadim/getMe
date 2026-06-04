import pytest
from getme_mcp_server.client import GetMeClient


@pytest.mark.asyncio
async def test_key_length_validation():
    client = GetMeClient(socket_path="/tmp/x.sock")
    client.max_key_length = 5
    with pytest.raises(ValueError, match="key exceeds max length"):
        await client.get("123456")


@pytest.mark.asyncio
async def test_value_size_validation():
    client = GetMeClient(socket_path="/tmp/x.sock")
    client.max_value_size = 10
    with pytest.raises(ValueError, match="value exceeds max size"):
        await client.put("k", "12345678901")


@pytest.mark.asyncio
async def test_batch_items_validation():
    client = GetMeClient(socket_path="/tmp/x.sock")
    client.max_batch_items = 2
    with pytest.raises(ValueError, match="batch exceeds max items"):
        await client.batch_put({"k1": "v1", "k2": "v2", "k3": "v3"})
