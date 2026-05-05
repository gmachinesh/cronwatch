package monitor

import (
	"log"
	"sync"
	"time"

	"cronwatch/internal/config"
)

// JobState tracks the last execution time and status of a cron job.
type JobState struct {
	Name        string
	LastSeen    time.Time
	MissedCount int
	Failed      bool
}

// Monitor watches registered jobs and detects drift or failure.
type Monitor struct {
	mu     sync.Mutex
	jobs   map[string]*JobState
	cfg    *config.Config
	alerts AlertSender
}

// New creates a new Monitor from the given config and alert sender.
func New(cfg *config.Config, alerts AlertSender) *Monitor {
	m := &Monitor{
		jobs:   make(map[string]*JobState),
		cfg:    cfg,
		alerts: alerts,
	}
	for _, j := range cfg.Jobs {
		m.jobs[j.Name] = &JobState{Name: j.Name}
	}
	return m
}

// RecordSuccess marks a job as having completed successfully.
func (m *Monitor) RecordSuccess(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	state, ok := m.jobs[name]
	if !ok {
		log.Printf("[monitor] unknown job: %s", name)
		return
	}
	state.LastSeen = time.Now()
	state.MissedCount = 0
	state.Failed = false
	log.Printf("[monitor] job %q succeeded", name)
}

// RecordFailure marks a job as failed and sends an alert.
func (m *Monitor) RecordFailure(name string, reason string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	state, ok := m.jobs[name]
	if !ok {
		log.Printf("[monitor] unknown job: %s", name)
		return
	}
	state.LastSeen = time.Now()
	state.Failed = true
	log.Printf("[monitor] job %q failed: %s", name, reason)
	m.alerts.Send(Alert{Job: name, Kind: KindFailure, Message: reason})
}

// CheckDrift inspects all jobs for schedule drift and alerts if overdue.
func (m *Monitor) CheckDrift() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, j := range m.cfg.Jobs {
		state := m.jobs[j.Name]
		if state.LastSeen.IsZero() {
			continue
		}
		deadline := state.LastSeen.Add(time.Duration(j.IntervalMinutes) * time.Minute * 2)
		if time.Now().After(deadline) {
			state.MissedCount++
			log.Printf("[monitor] job %q drift detected (missed: %d)", j.Name, state.MissedCount)
			m.alerts.Send(Alert{Job: j.Name, Kind: KindDrift, Message: "job overdue"})
		}
	}
}
