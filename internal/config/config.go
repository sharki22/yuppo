package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	WatchPath       string        `yaml:"watch_path"`
	CommitMessage   string        `yaml:"commit_message"`
	AutoPush        bool          `yaml:"auto_push"`
	DebounceTimeout time.Duration `yaml:"-"`
	RawDebounce     int           `yaml:"debounce_seconds"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config %s: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config %s: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	cfg.DebounceTimeout = time.Duration(cfg.RawDebounce) * time.Second
	return &cfg, nil
}

func (c *Config) validate() error {
	if c.WatchPath == "" {
		return fmt.Errorf("config: watch_path is required")
	}

	info, err := os.Stat(c.WatchPath)
	if err != nil {
		return fmt.Errorf("config: watch_path %q: %w", c.WatchPath, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("config: watch_path %q is not a directory", c.WatchPath)
	}

	if c.CommitMessage == "" {
		c.CommitMessage = "auto: update"
	}

	if c.RawDebounce <= 0 {
		c.RawDebounce = 5
	}

	return nil
}
