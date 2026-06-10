FROM ghcr.io/astral-sh/uv:python3.12-bookworm AS builder

WORKDIR /app

# The project uses hatch-vcs which requires .git to resolve the dynamic version.
# Since .git is not copied into the Docker context, we must provide a pretend version
# so that building the wheel and syncing does not fail with setuptools-scm errors.
ENV SETUPTOOLS_SCM_PRETEND_VERSION=1.0.0

COPY pyproject.toml uv.lock README.md ./

RUN uv sync --frozen --no-dev --no-install-project

COPY getme_mcp_server ./getme_mcp_server

RUN uv pip install --python /app/.venv/bin/python --no-deps .


# slim Python image, no uv
FROM python:3.12-slim-bookworm AS runtime

WORKDIR /app

ENV PYTHONUNBUFFERED=1 \
  GETME_BASE_URL=http://unix \
  GETME_SOCKET_PATH=/tmp/getMeStore/sockDir/getMe.sock \
  PATH=/app/.venv/bin:$PATH

COPY --from=builder /app/.venv /app/.venv

COPY --from=builder /app/README.md /app/README.md

ENTRYPOINT ["getme-mcp-server"]
