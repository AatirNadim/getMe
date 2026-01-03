import os


DEFAULT_SOCKET_PATH = "/tmp/getMeStore/sockDir/getMe.sock"
DEFAULT_BASE_URL = "http://unix"


def socket_path() -> str:
    return os.environ.get("GETME_SOCKET_PATH", DEFAULT_SOCKET_PATH)


def base_url() -> str:
    # The host part is irrelevant for UDS, but must be a valid URL.
    return os.environ.get("GETME_BASE_URL", DEFAULT_BASE_URL)
