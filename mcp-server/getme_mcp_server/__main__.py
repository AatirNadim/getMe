from __future__ import annotations

from .server import build_mcp


def main() -> None:
    mcp = build_mcp()
    mcp.run()


if __name__ == "__main__":
    main()
