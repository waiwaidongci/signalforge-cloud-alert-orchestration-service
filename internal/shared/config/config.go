package config

// Config is the root configuration object loaded from YAML and environment.
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Database  DatabaseConfig  `yaml:"database"`
	Log       LogConfig       `yaml:"log"`
	Scheduler SchedulerConfig `yaml:"scheduler"`
	Auth      AuthConfig      `yaml:"auth"`
	Metrics   MetricsConfig   `yaml:"metrics"`
}

type ServerConfig struct {
	Addr            string   `yaml:"addr"`
	ReadTimeout     Duration `yaml:"read_timeout"`
	WriteTimeout    Duration `yaml:"write_timeout"`
	IdleTimeout     Duration `yaml:"idle_timeout"`
	ShutdownTimeout Duration `yaml:"shutdown_timeout"`
	RequestTimeout  Duration `yaml:"request_timeout"`
	MaxBodyBytes    int64    `yaml:"max_body_bytes"`
}

type DatabaseConfig struct {
	Driver          string   `yaml:"driver"`
	DSN             string   `yaml:"dsn"`
	MaxOpenConns    int      `yaml:"max_open_conns"`
	MaxIdleConns    int      `yaml:"max_idle_conns"`
	ConnMaxLifetime Duration `yaml:"conn_max_lifetime"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type SchedulerConfig struct {
	Enabled  bool     `yaml:"enabled"`
	Interval Duration `yaml:"interval"`
	Jitter   Duration `yaml:"jitter"`
}

type AuthConfig struct {
	PlaceholderToken string `yaml:"placeholder_token"`
}

type MetricsConfig struct {
	Enabled bool `yaml:"enabled"`
}

func Default() Config {
	return Config{
		Server: ServerConfig{
			Addr:            ":8080",
			ReadTimeout:     Duration(5_000_000_000),
			WriteTimeout:    Duration(15_000_000_000),
			IdleTimeout:     Duration(60_000_000_000),
			ShutdownTimeout: Duration(10_000_000_000),
			RequestTimeout:  Duration(12_000_000_000),
			MaxBodyBytes:    1 << 20,
		},
		Database: DatabaseConfig{
			Driver:          "sqlite",
			DSN:             "./data/signalforge.db",
			MaxOpenConns:    8,
			MaxIdleConns:    4,
			ConnMaxLifetime: Duration(30 * 60 * 1_000_000_000),
		},
		Log: LogConfig{Level: "info", Format: "json"},
		Scheduler: SchedulerConfig{
			Enabled:  true,
			Interval: Duration(10_000_000_000),
			Jitter:   Duration(2_000_000_000),
		},
		Metrics: MetricsConfig{Enabled: true},
	}
}
