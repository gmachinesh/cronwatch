package scheduler

import (
	"testing"
	"time"

	"cronwatch/internal/config"
	"cronwatch/internal/monitor"
	"cronwatch/internal/notifier"
)

func makeTestConfig(schedule, command string) *config.Config {
	return &config.Config{
		Jobs: []config.Job{
			{
				Name:     "test-job",
				Schedule: schedule,
				Command:  command,
			},
		},
	}
}

func makeTestMonitor(cfg *config.Config) *monitor.Monitor {
	n := notifier.New(nil)
	return monitor.New(cfg, n)
}

func TestScheduler_StartStop(t *testing.T) {
	cfg := makeTestConfig("@every 1m", "echo ok")
	mon := makeTestMonitor(cfg)
	s := New(cfg, mon)

	if err := s.Start(); err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	s.Stop()
}

func TestScheduler_InvalidSchedule(t *testing.T) {
	cfg := makeTestConfig("not-a-cron", "echo ok")
	mon := makeTestMonitor(cfg)
	s := New(cfg, mon)

	if err := s.Start(); err == nil {
		t.Fatal("expected error for invalid schedule, got nil")
		s.Stop()
	}
}

func TestScheduler_NextRun(t *testing.T) {
	cfg := makeTestConfig("@every 1m", "echo ok")
	mon := makeTestMonitor(cfg)
	s := New(cfg, mon)

	if err := s.Start(); err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	defer s.Stop()

	next, ok := s.NextRun("test-job")
	if !ok {
		t.Fatal("NextRun: job not found")
	}
	if next.Before(time.Now()) {
		t.Errorf("NextRun: expected future time, got %v", next)
	}
}

func TestScheduler_NextRun_Unknown(t *testing.T) {
	cfg := makeTestConfig("@every 1m", "echo ok")
	mon := makeTestMonitor(cfg)
	s := New(cfg, mon)

	_, ok := s.NextRun("ghost-job")
	if ok {
		t.Error("expected false for unknown job")
	}
}
