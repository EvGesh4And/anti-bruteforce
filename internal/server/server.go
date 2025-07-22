package grpcserver

import (
	"context"
	"log/slog"

	pb "github.com/EvGesh4And/anti-bruteforce/api"
	"github.com/EvGesh4And/anti-bruteforce/internal/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AntiBruteforceService interface {
	Check(ctx context.Context, login, password, ip string) (bool, error)
	Reset(ctx context.Context, login, ip string)
	AddToBlacklist(ctx context.Context, network string) error
	RemoveFromBlacklist(ctx context.Context, network string) error
	AddToWhitelist(ctx context.Context, network string) error
	RemoveFromWhitelist(ctx context.Context, network string) error
}

type Server struct {
	pb.UnimplementedAntiBruteforceServer
	service AntiBruteforceService
	logger  *slog.Logger
}

func New(logger *slog.Logger, service AntiBruteforceService) *Server {
	return &Server{logger: logger, service: service}
}

func (s *Server) setLogCompMeth(ctx context.Context, method string) context.Context {
	ctx = logger.WithLogComponent(ctx, "server.grpc")
	return logger.WithLogMethod(ctx, method)
}

func (s *Server) Check(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {
	ctx = s.setLogCompMeth(ctx, "Check")
	ctx = logger.WithLogLogin(ctx, req.Login)
	ctx = logger.WithLogIP(ctx, req.Ip)

	s.logger.DebugContext(ctx, "starting Check request")

	ok, err := s.service.Check(ctx, req.Login, req.Password, req.Ip)
	if err != nil {
		s.logger.ErrorContext(ctx, "Check failed: "+err.Error())
		return nil, err
	}

	s.logger.InfoContext(ctx, "Check completed successfully")
	return &pb.CheckResponse{Ok: ok}, nil
}

func (s *Server) Reset(ctx context.Context, req *pb.ResetRequest) (*emptypb.Empty, error) {
	ctx = s.setLogCompMeth(ctx, "Reset")
	ctx = logger.WithLogLogin(ctx, req.Login)
	ctx = logger.WithLogIP(ctx, req.Ip)

	s.logger.DebugContext(ctx, "Reset request initiated")

	s.service.Reset(ctx, req.Login, req.Ip)

	s.logger.InfoContext(ctx, "Rate limits reset successfully")
	return &emptypb.Empty{}, nil
}

func (s *Server) AddToBlacklist(ctx context.Context, req *pb.NetworkRequest) (*emptypb.Empty, error) {
	ctx = s.setLogCompMeth(ctx, "AddToBlacklist")
	ctx = logger.WithLogNetwork(ctx, req.Network)

	s.logger.DebugContext(ctx, "Adding network to blacklist")

	if err := s.service.AddToBlacklist(ctx, req.Network); err != nil {
		s.logger.ErrorContext(ctx, "Failed to add to blacklist: "+err.Error())
		return nil, err
	}

	s.logger.InfoContext(ctx, "Network added to blacklist successfully")
	return &emptypb.Empty{}, nil
}

func (s *Server) RemoveFromBlacklist(ctx context.Context, req *pb.NetworkRequest) (*emptypb.Empty, error) {
	ctx = s.setLogCompMeth(ctx, "RemoveFromBlacklist")
	ctx = logger.WithLogNetwork(ctx, req.Network)

	s.logger.DebugContext(ctx, "Removing network from blacklist")

	if err := s.service.RemoveFromBlacklist(ctx, req.Network); err != nil {
		s.logger.ErrorContext(ctx, "Failed to remove from blacklist: "+err.Error())
		return nil, err
	}

	s.logger.InfoContext(ctx, "Network removed from blacklist successfully")
	return &emptypb.Empty{}, nil
}

func (s *Server) AddToWhitelist(ctx context.Context, req *pb.NetworkRequest) (*emptypb.Empty, error) {
	ctx = s.setLogCompMeth(ctx, "AddToWhitelist")
	ctx = logger.WithLogNetwork(ctx, req.Network)

	s.logger.DebugContext(ctx, "Adding network to whitelist")

	if err := s.service.AddToWhitelist(ctx, req.Network); err != nil {
		s.logger.ErrorContext(ctx, "Failed to add to whitelist: "+err.Error())
		return nil, err
	}

	s.logger.InfoContext(ctx, "Network added to whitelist successfully")
	return &emptypb.Empty{}, nil
}

func (s *Server) RemoveFromWhitelist(ctx context.Context, req *pb.NetworkRequest) (*emptypb.Empty, error) {
	ctx = s.setLogCompMeth(ctx, "RemoveFromWhitelist")
	ctx = logger.WithLogNetwork(ctx, req.Network)

	s.logger.DebugContext(ctx, "Removing network from whitelist")

	if err := s.service.RemoveFromWhitelist(ctx, req.Network); err != nil {
		s.logger.ErrorContext(ctx, "Failed to remove from whitelist: "+err.Error())
		return nil, err
	}

	s.logger.InfoContext(ctx, "Network removed from whitelist successfully")
	return &emptypb.Empty{}, nil
}
