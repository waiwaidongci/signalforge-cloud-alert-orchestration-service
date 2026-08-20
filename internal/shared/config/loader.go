package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Load reads a YAML file, applies defaults for missing values, and then
// applies SIGNALFORGE_* environment variable overrides.
func Load(path string) (Config, error) {
	cfg := Default()
	if path != "" {
		raw, err := os.ReadFile(path)
		if err != nil {
			return Config{}, fmt.Errorf("read config: %w", err)
		}
		if err := yaml.Unmarshal(raw, &cfg); err != nil {
			return Config{}, fmt.Errorf("decode config: %w", err)
		}
	}
	if err := applyEnv(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func applyEnv(cfg *Config) error {
	applyString := func(key string, target *string) {
		if value, ok := os.LookupEnv(key); ok {
			*target = value
		}
	}
	applyString("SIGNALFORGE_SERVER_ADDR", &cfg.Server.Addr)
	applyString("SIGNALFORGE_DATABASE_DRIVER", &cfg.Database.Driver)
	applyString("SIGNALFORGE_DATABASE_DSN", &cfg.Database.DSN)
	applyString("SIGNALFORGE_LOG_LEVEL", &cfg.Log.Level)
	applyString("SIGNALFORGE_LOG_FORMAT", &cfg.Log.Format)
	applyString("SIGNALFORGE_AUTH_PLACEHOLDER_TOKEN", &cfg.Auth.PlaceholderToken)

	if err := envDuration("SIGNALFORGE_READ_TIMEOUT", &cfg.Server.ReadTimeout); err != nil {
		return err
	}
	if err := envDuration("SIGNALFORGE_WRITE_TIMEOUT", &cfg.Server.WriteTimeout); err != nil {
		return err
	}
	if err := envDuration("SIGNALFORGE_REQUEST_TIMEOUT", &cfg.Server.RequestTimeout); err != nil {
		return err
	}

	if value := os.Getenv("SIGNALFORGE_DATABASE_MAX_OPEN_CONNS"); value != "" {
		if err := setInt(value, &cfg.Database.MaxOpenConns); err != nil {
			return fmt.Errorf("SIGNALFORGE_DATABASE_MAX_OPEN_CONNS: %w", err)
		}
	}
	if value := os.Getenv("SIGNALFORGE_DATABASE_MAX_IDLE_CONNS"); value != "" {
		if err := setInt(value, &cfg.Database.MaxIdleConns); err != nil {
			return fmt.Errorf("SIGNALFORGE_DATABASE_MAX_IDLE_CONNS: %w", err)
		}
	}
	if value := os.Getenv("SIGNALFORGE_SCHEDULER_ENABLED"); value != "" {
		cfg.Scheduler.Enabled = strings.EqualFold(value, "true") || value == "1"
	}
	if value := os.Getenv("SIGNALFORGE_METRICS_ENABLED"); value != "" {
		cfg.Metrics.Enabled = strings.EqualFold(value, "true") || value == "1"
	}
	return nil
}
