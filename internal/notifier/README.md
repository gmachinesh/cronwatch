# notifier

The `notifier` package provides a pluggable alert-dispatch layer for cronwatch.

## Overview

When the monitor detects a job failure or schedule drift it constructs an `Alert`
and passes it to `Notifier.Send`. The notifier fans the alert out to every
configured backend.

## Supported backends

| Backend | Config field | Notes |
|---------|-------------|-------|
| Slack   | `SlackWebhookURL` | Uses Incoming Webhooks |

## Usage

```go
import "cronwatch/internal/notifier"

n := notifier.New(notifier.Config{
    SlackWebhookURL: "https://hooks.slack.com/services/…",
})

err := n.Send(notifier.Alert{
    JobName: "db-backup",
    Level:   notifier.AlertFailure,
    Message: "process exited with code 2",
})
```

## Adding a new backend

1. Add the relevant config fields to `Config`.
2. Implement a private `sendXxx(a Alert) error` method on `*Notifier`.
3. Call it inside `Send` when the config field is non-empty.
4. Add tests in `notifier_test.go`.
