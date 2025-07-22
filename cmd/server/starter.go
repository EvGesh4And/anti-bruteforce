// Package main implements the entry point for the application.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	pb "github.com/EvGesh4And/anti-bruteforce/api"
	"github.com/EvGesh4And/anti-bruteforce/config"
	"github.com/EvGesh4And/anti-bruteforce/internal/antibruteforce"
	grpcserver "github.com/EvGesh4And/anti-bruteforce/internal/server"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func startGRPCServer(
	ctx context.Context,
	g *errgroup.Group,
	cfg config.GRPCConfig,
	lg *slog.Logger,
	svc *antibruteforce.Service,
) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		g.Go(func() error {
			return fmt.Errorf("gRPC failed to listen on port %s: %w", addr, err)
		})
		return
	}

	server := grpc.NewServer()
	pb.RegisterAntiBruteforceServer(server, grpcserver.New(lg, svc))

	g.Go(func() error {
		log.Printf("gRPC server starting %s...", lis.Addr().String())
		if err := server.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return fmt.Errorf("grpc serve: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		done := make(chan struct{})
		go func() {
			server.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			log.Print("[shutdown] gRPC server stopped gracefully...")
		case <-shutdownCtx.Done():
			log.Print("[shutdown] gRPC graceful shutdown timeout, calling Stop()")
			server.Stop()
		}
		return ctx.Err()
	})
}
