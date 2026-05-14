# Acknowledge API

The acknowledge API allows operators to acknowledge failing or drifting cron jobs,
suppressing repeat alerts while an issue is being investigated.

## Endpoints

### `POST /api/v1/jobs/acknowledge?job=<name>`

Acknowledge a job alert. Optionally provide a duration after which the acknowledgement expires.

**Request body (JSON, optional):**
```json
{
  "acked_by": "alice",
  "note": "investigating disk issue",
  "duration": "2h"
}
```

**Response:**
```json
{"acknowledged": "backup", "acked_at": "2024-01-15T10:30:00Z"}
```

### `POST /api/v1/jobs/unacknowledge?job=<name>`

Remove an acknowledgement for a job.

**Response:**
```json
{"unacknowledged": "backup"}
```

### `GET /api/v1/jobs/acknowledged`

List all currently active acknowledgements.

**Response:**
```json
[
  {
    "job_name": "backup",
    "acked_at": "2024-01-15T10:30:00Z",
    "acked_by": "alice",
    "note": "investigating disk issue",
    "expires_at": "2024-01-15T12:30:00Z"
  }
]
```

## Notes

- Acknowledgements without a `duration` do not expire automatically.
- Expired acknowledgements are excluded from the list and from `IsAcknowledged` checks.
- The `IsAcknowledged(jobName string) bool` helper can be used internally to gate alert delivery.
