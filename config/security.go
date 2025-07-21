package config

// SecurityConfig holds the configuration for security.
type SecurityConfig struct {
	N int `toml:"n"`
	M int `toml:"m"`
	K int `toml:"k"`
}
