from __future__ import annotations

import sys

from .server import build_mcp

from .utills import _install_graceful_shutdown_handlers


def main() -> int:
    _install_graceful_shutdown_handlers()
    mcp = build_mcp()

    try:
        print("Starting MCP server...")
        mcp.run()
        return 0
    except (KeyboardInterrupt, SystemExit):
        # quite shutdown in case of interrupts
        return 0
    except BrokenPipeError:
        # treating broken pipe as a normal exit
        return 0
    finally:
        try:
            sys.stdout.flush()
        except Exception:
            pass
        try:
            sys.stderr.flush()
        except Exception:
            pass


if __name__ == "__main__":
    raise SystemExit(main())
