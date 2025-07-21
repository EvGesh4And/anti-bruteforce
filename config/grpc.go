package config

// GRPCConfig holds the configuration for gRPC.
type GRPCConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}
