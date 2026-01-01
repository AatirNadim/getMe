import threading
import os
import signal


def _install_graceful_shutdown_handlers() -> None:
    """
    this is intentionally defensive because MCP servers often run under
    supervisors that send SIGTERM, and stdio transport can be closed abruptly.
    """

    lock = threading.Lock()
    state = {"seen": False}

    def _handler(signum: int, _frame) -> None:
        with lock:
            if state["seen"]:
                os._exit(1)
            state["seen"] = True
        raise SystemExit(0)

    for sig in (signal.SIGINT, signal.SIGTERM):
        try:
            signal.signal(sig, _handler)
        except Exception:
            pass
