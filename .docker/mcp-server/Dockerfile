FROM ghcr.io/astral-sh/uv:python3.14-alpine3.23 AS dev

WORKDIR /app

COPY ../../mcp-server /app
RUN uv sync --locked

CMD ["uv", "run", "uvicorn", "src.presentation.http.main:app", "--host", "0.0.0.0", "--port", "8080", "--reload"]

FROM ghcr.io/astral-sh/uv:python3.14-alpine3.23 AS runner

WORKDIR /app

COPY ../../mcp-server /app
RUN uv sync --locked

CMD ["uv", "run", "uvicorn", "src.presentation.http.main:app", "--host", "0.0.0.0", "--port", "8080"]