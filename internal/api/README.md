# internal/api

HTTP API server for cronwatch.

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/healthz` | Liveness check |
| GET | `/status` | Current job states |
| GET | `/history` | Recent execution history |

## Authentication

Optional API key authentication is supported via the `X-API-Key` request header.

Set `api.key` in your `config.yaml` to enable it:

```yaml
api:
  addr: ":8080"
  key: "your-secret-key"
```

When a key is configured:
- Requests without the `X-API-Key` header receive **401 Unauthorized**.
- Requests with an incorrect key receive **403 Forbidden**.
- Requests with the correct key are passed through normally.

If `api.key` is empty or omitted, all requests are allowed without authentication.

## Middleware

- **loggingMiddleware** — logs method, path, status, and duration for every request.
- **methodGuard** — restricts each route to its allowed HTTP method.
- **apiKeyAuth** — enforces static API key authentication when configured.
