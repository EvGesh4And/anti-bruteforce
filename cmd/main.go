// Package main implements the entry point for the application.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/EvGesh4And/anti-bruteforce/config"
	"github.com/EvGesh4And/anti-bruteforce/internal/logger"
)

var pathConfigFile string

func init() {
	flag.StringVar(&pathConfigFile, "config", "configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	log.SetOutput(os.Stdout)

	cfg := config.Config{}
	fmt.Println(cfg)
	err := config.LoadConfig(pathConfigFile, &cfg)
	if err != nil {
		log.Printf("error initializing config: %v", err)
		return
	}

	_, closer, err := logger.NewSlogLogger(cfg.Logger)
	if err != nil {
		log.Printf("error initializing logger: %v", err)
		return
	}
	if closer != nil {
		defer func() {
			if err := closer.Close(); err != nil {
				log.Printf("error closing logger: %v", err)
			}
		}()
	}

	fmt.Println(cfg)
}
