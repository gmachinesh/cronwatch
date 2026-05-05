package scheduler

import (
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"cronwatch/internal/config"
	"cronwatch/internal/monitor"
)

// Scheduler wraps cron scheduling and ties jobs to the monitor.
type Scheduler struct {
	cron    *cron.Cron
	monitor *monitor.Monitor
	cfg     *config.Config
	mu      sync.Mutex
	entries map[string]cron.EntryID
}

// New creates a new Scheduler.
func New(cfg *config.Config, mon *monitor.Monitor) *Scheduler {
	return &Scheduler{
		cron:    cron.New(cron.WithSeconds()),
		monitor: mon,
		cfg:     cfg,
		entries: make(map[string]cron.EntryID),
	}
}

// Start registers all configured jobs and starts the cron engine.
func (s *Scheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, job := range s.cfg.Jobs {
		j := job // capture loop variable
		id, err := s.cron.AddFunc(j.Schedule, func() {
			s.runJob(j)
		})
		if err != nil {
			return err
		}
		s.entries[j.Name] = id
		log.Printf("scheduler: registered job %q (%s)", j.Name, j.Schedule)
	}

	s.cron.Start()
	return nil
}

// Stop gracefully shuts down the scheduler.
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("scheduler: stopped")
}

// NextRun returns the next scheduled time for a named job.
func (s *Scheduler) NextRun(name string) (time.Time, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, ok := s.entries[name]
	if !ok {
		return time.Time{}, false
	}
	entry := s.cron.Entry(id)
	return entry.Next, true
}

func (s *Scheduler) runJob(job config.Job) {
	log.Printf("scheduler: executing job %q", job.Name)
	start := time.Now()

	err := execCommand(job.Command)
	duration := time.Since(start)

	if err != nil {
		log.Printf("scheduler: job %q failed after %s: %v", job.Name, duration, err)
		s.monitor.RecordFailure(job.Name, err)
	} else {
		log.Printf("scheduler: job %q succeeded in %s", job.Name, duration)
		s.monitor.RecordSuccess(job.Name)
	}
}
