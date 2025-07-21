package config

// LoggerConfig holds the configuration for the logger.
type LoggerConfig struct {
	Mod   string `toml:"mod"`
	Path  string `toml:"path"`
	JSON  bool   `toml:"json"`
	Level string `toml:"level"`
}
