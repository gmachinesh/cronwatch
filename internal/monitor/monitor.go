package monitor

import (
	"fmt"
	"sync"
	"time"

	"github.com/cronwatch/internal/config"
	"github.com/cronwatch/internal/notifier"
)

// State holds runtime information for a single job.
type State struct {
	LastRun   time.Time
	LastError error
	Drifted   bool
}

// Monitor tracks job execution state and triggers alerts.
type Monitor struct {
	mu       sync.RWMutex
	states   map[string]*State
	cfg      *config.Config
	notifier *notifier.Notifier
}

// New creates a Monitor for the given config and optional notifier.
func New(cfg *config.Config, n *notifier.Notifier) *Monitor {
	states := make(map[string]*State, len(cfg.Jobs))
	for _, j := range cfg.Jobs {
		states[j.Name] = &State{}
	}
	return &Monitor{cfg: cfg, notifier: n, states: states}
}

// RecordSuccess marks a job as having completed successfully.
func (m *Monitor) RecordSuccess(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	st := m.stateFor(name)
	st.LastRun = time.Now()
	st.LastError = nil
	st.Drifted = false
}

// RecordFailure marks a job as failed and sends an alert.
func (m *Monitor) RecordFailure(name string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	st := m.stateFor(name)
	st.LastRun = time.Now()
	st.LastError = err
	m.sendAlert(fmt.Sprintf("job %q failed: %v", name, err))
}

// CheckDrift inspects whether a job has exceeded its expected interval.
func (m *Monitor) CheckDrift(name string, maxAge time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	st := m.stateFor(name)
	if st.LastRun.IsZero() {
		return
	}
	if time.Since(st.LastRun) > maxAge {
		st.Drifted = true
		m.sendAlert(fmt.Sprintf("job %q has not run in %v", name, maxAge))
	}
}

// States returns a snapshot of all job states.
func (m *Monitor) States() map[string]State {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[string]State, len(m.states))
	for k, v := range m.states {
		out[k] = *v
	}
	return out
}

func (m *Monitor) stateFor(name string) *State {
	if _, ok := m.states[name]; !ok {
		m.states[name] = &State{}
	}
	return m.states[name]
}

func (m *Monitor) sendAlert(msg string) {
	if m.notifier == nil {
		return
	}
	_ = m.notifier.Send(msg)
}
