package config

import (
	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v10"
)

// LoadConfig загружает конфиг из TOML-файла и env-переменных в любую структуру.
func LoadConfig[T any](path string, cfg *T) error {
	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return err
	}
	return env.Parse(cfg)
}
