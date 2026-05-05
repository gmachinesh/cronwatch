package monitor

import (
	"testing"
	"time"

	"cronwatch/internal/config"
)

type captureAlerts struct {
	alerts []Alert
}

func (c *captureAlerts) Send(a Alert) {
	c.alerts = append(c.alerts, a)
}

func makeConfig(jobs []config.Job) *config.Config {
	return &config.Config{Jobs: jobs}
}

func TestRecordSuccess(t *testing.T) {
	cap := &captureAlerts{}
	cfg := makeConfig([]config.Job{{Name: "backup", IntervalMinutes: 60}})
	m := New(cfg, cap)

	m.RecordSuccess("backup")

	m.mu.Lock()
	defer m.mu.Unlock()
	state := m.jobs["backup"]
	if state.Failed {
		t.Error("expected Failed=false after success")
	}
	if state.LastSeen.IsZero() {
		t.Error("expected LastSeen to be set")
	}
	if len(cap.alerts) != 0 {
		t.Errorf("expected no alerts, got %d", len(cap.alerts))
	}
}

func TestRecordFailure(t *testing.T) {
	cap := &captureAlerts{}
	cfg := makeConfig([]config.Job{{Name: "sync", IntervalMinutes: 30}})
	m := New(cfg, cap)

	m.RecordFailure("sync", "exit code 1")

	m.mu.Lock()
	defer m.mu.Unlock()
	state := m.jobs["sync"]
	if !state.Failed {
		t.Error("expected Failed=true after failure")
	}
	if len(cap.alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(cap.alerts))
	}
	if cap.alerts[0].Kind != KindFailure {
		t.Errorf("expected KindFailure, got %s", cap.alerts[0].Kind)
	}
}

func TestCheckDrift(t *testing.T) {
	cap := &captureAlerts{}
	cfg := makeConfig([]config.Job{{Name: "report", IntervalMinutes: 1}})
	m := New(cfg, cap)

	// Simulate a job that ran 5 minutes ago (well past 2× interval of 1 min)
	m.mu.Lock()
	m.jobs["report"].LastSeen = time.Now().Add(-5 * time.Minute)
	m.mu.Unlock()

	m.CheckDrift()

	if len(cap.alerts) != 1 {
		t.Fatalf("expected 1 drift alert, got %d", len(cap.alerts))
	}
	if cap.alerts[0].Kind != KindDrift {
		t.Errorf("expected KindDrift, got %s", cap.alerts[0].Kind)
	}
}

func TestCheckDrift_NoDriftWhenFresh(t *testing.T) {
	cap := &captureAlerts{}
	cfg := makeConfig([]config.Job{{Name: "cleanup", IntervalMinutes: 60}})
	m := New(cfg, cap)

	m.mu.Lock()
	m.jobs["cleanup"].LastSeen = time.Now()
	m.mu.Unlock()

	m.CheckDrift()

	if len(cap.alerts) != 0 {
		t.Errorf("expected no alerts for fresh job, got %d", len(cap.alerts))
	}
}
