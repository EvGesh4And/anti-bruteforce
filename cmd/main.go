// Package main implements the entry point for the application.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/EvGesh4And/anti-bruteforce/config"
	"github.com/EvGesh4And/anti-bruteforce/internal/antibruteforce"
	"github.com/EvGesh4And/anti-bruteforce/internal/logger"
	"golang.org/x/sync/errgroup"
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

	lg, closer, err := logger.NewSlogLogger(cfg.Logger)
	if err != nil {
		log.Printf("error initializing logger: %v", err)
		return
	}
	if closer != nil {
		defer closer.Close()
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	svc := antibruteforce.NewService(lg, cfg.Security)
	defer svc.Shutdown()

	g, ctx := errgroup.WithContext(ctx)

	startGRPCServer(ctx, g, cfg.GRPC, lg, svc)

	if err := g.Wait(); err != nil {
		log.Printf("service stopped with error: %v", err)
		return
	}
}
