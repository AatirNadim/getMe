import os

# Base connectivity
DEFAULT_SOCKET_PATH = "/tmp/getMeStore/sockDir/getMe.sock"
DEFAULT_BASE_URL = "http://unix"


def socket_path() -> str:
    return os.environ.get("GETME_SOCKET_PATH", DEFAULT_SOCKET_PATH)


def base_url() -> str:
    # The host part is irrelevant for UDS, but must be a valid URL.
    return os.environ.get("GETME_BASE_URL", DEFAULT_BASE_URL)


# Safety & Scope config
def is_read_only() -> bool:
    return os.environ.get("GETME_READ_ONLY", "false").lower() == "true"


def allow_clear() -> bool:
    return os.environ.get("GETME_ALLOW_CLEAR", "false").lower() == "true"


def key_prefix() -> str:
    return os.environ.get("GETME_KEY_PREFIX", "")


# Validation constraints
def max_value_size_bytes() -> int:
    return int(
        os.environ.get("GETME_MAX_VALUE_SIZE_BYTES", 5 * 1024 * 1024)
    )  # default 5MB


def max_batch_items() -> int:
    return int(os.environ.get("GETME_MAX_BATCH_ITEMS", 100))


def max_key_length() -> int:
    return int(os.environ.get("GETME_MAX_KEY_LENGTH", 1024))
