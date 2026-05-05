# cronwatch HTTP API

The `api` package exposes a lightweight read-only HTTP interface for inspecting
cronwatch's runtime state.

## Endpoints

### `GET /healthz`

Returns `{"status": "ok"}` when the daemon is running.

```json
{"status": "ok"}
```

### `GET /status`

Returns the current state of every monitored job.

```json
[
  {
    "job": "backup",
    "last_run_at": "2024-01-15T03:00:00Z",
    "last_result": "success",
    "drifted": false
  }
]
```

### `GET /history?job=<name>`

Returns the most recent 50 history records for the given job.
Omit `job` to retrieve records for all jobs.

## Configuration

Enable the API in `config.yaml`:

```yaml
api:
  enabled: true
  addr: ":8080"
```

Defaults to `disabled`. When disabled the HTTP server is not started.
