# history

The `history` package provides persistent storage for cron job execution records.

## Overview

Each time a job runs, a `history.Entry` is appended to a JSON file on disk. This allows cronwatch to:

- Display recent execution history per job
- Detect patterns of repeated failures
- Support future drift analysis across restarts

## Usage

```go
store, err := history.New("/var/lib/cronwatch/history.json")
if err != nil {
    log.Fatal(err)
}

// Record a completed job run
store.Append(history.Entry{
    JobName:   "daily-backup",
    StartedAt: time.Now(),
    Duration:  5 * time.Second,
    Success:   true,
})

// Retrieve the last 10 runs for a job
entries := store.Recent("daily-backup", 10)
```

## Entry Fields

| Field       | Type            | Description                        |
|-------------|-----------------|------------------------------------|
| `job_name`  | string          | Name of the cron job               |
| `started_at`| time.Time       | When the job started               |
| `duration`  | time.Duration   | How long the job ran               |
| `success`   | bool            | Whether the job exited cleanly     |
| `output`    | string          | Captured stdout/stderr (optional)  |
| `error`     | string          | Error message on failure (optional)|

## Storage

Entries are stored as a JSON array. The file is created automatically if it does not exist.
