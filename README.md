# cronwatch

Lightweight daemon that monitors cron job execution and sends alerts on failure or drift.

## Installation

```bash
go install github.com/yourname/cronwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/cronwatch.git && cd cronwatch && go build ./...
```

## Usage

Define your monitored jobs in a `cronwatch.yaml` config file:

```yaml
jobs:
  - name: daily-backup
    schedule: "0 2 * * *"
    timeout: 30m
    alert_on_failure: true
    drift_threshold: 5m

alerts:
  email: ops@example.com
```

Then start the daemon:

```bash
cronwatch --config cronwatch.yaml
```

Wrap your existing cron commands to report status:

```bash
# In your crontab
0 2 * * * cronwatch exec --job daily-backup -- /usr/local/bin/backup.sh
```

cronwatch will send alerts if a job fails, exceeds its timeout, or runs significantly outside its expected schedule window.

## Configuration

| Field | Description |
|---|---|
| `schedule` | Standard cron expression for expected run time |
| `timeout` | Maximum allowed runtime before alerting |
| `drift_threshold` | How late a job can start before triggering a drift alert |

## License

MIT