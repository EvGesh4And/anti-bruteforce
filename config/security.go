package config

import "time"

// SecurityConfig holds the configuration for security.
type SecurityConfig struct {
	LoginRate       int           `toml:"loginrate"`
	PassRate        int           `toml:"passrate"`
	IPRate          int           `toml:"iprate"`
	CleanupInterval time.Duration `toml:"cleanupInterval"`
	BucketMaxIdle   time.Duration `toml:"bucketMaxIdle"`
}
