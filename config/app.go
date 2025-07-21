package config

// Config holds the application configuration.
type Config struct {
	Logger   LoggerConfig   `toml:"logger" env-prefix:"LOGGER_"`
	GRPC     GRPCConfig     `toml:"grpc" env-prefix:"GRPC_"`
	Security SecurityConfig `toml:"security" env-prefix:"SECURITY_"`
}
