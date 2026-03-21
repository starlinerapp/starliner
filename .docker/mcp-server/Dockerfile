FROM ghcr.io/astral-sh/uv:python3.14-alpine3.23 AS dev

WORKDIR /app

CMD ["sh", "-c", "uv sync && uv run fastmcp run src/presentation/http/main.py:mcp --transport streamable-http --host 0.0.0.0 --port 8080 --reload"]

FROM ghcr.io/astral-sh/uv:python3.14-alpine3.23 AS runner

WORKDIR /app

COPY ../../mcp-server /app
RUN uv sync --locked

CMD ["uv", "run", "fastmcp", "run", "src/presentation/http/main.py:mcp", "--transport", "streamable-http", "--host", "0.0.0.0", "--port", "8080"]