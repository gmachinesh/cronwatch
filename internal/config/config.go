package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Job describes a single cron job to monitor.
type Job struct {
	Name            string `yaml:"name"`
	Schedule        string `yaml:"schedule"`
	IntervalMinutes int    `yaml:"interval_minutes"`
	TimeoutSeconds  int    `yaml:"timeout_seconds"`
}

// Config holds the full cronwatch configuration.
type Config struct {
	LogLevel string `yaml:"log_level"`
	Alert    struct {
		Email   string `yaml:"email"`
		Webhook string `yaml:"webhook"`
	} `yaml:"alert"`
	Jobs []Job `yaml:"jobs"`
}

// Load reads and parses the YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validate(cfg *Config) error {
	if len(cfg.Jobs) == 0 {
		return errors.New("config must define at least one job")
	}
	for i, j := range cfg.Jobs {
		if j.Name == "" {
			return fmt.Errorf("job[%d] missing required field: name", i)
		}
	}
	return nil
}
